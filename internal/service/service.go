package service

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Producer interface {
	Produce() ([]string, error)
}

type Presenter interface {
	Present([]string) error
}

type Service struct {
	prod Producer
	pres Presenter
}

func NewService(prod Producer, pres Presenter) *Service {
	return &Service{
		prod: prod,
		pres: pres,
	}
}

func (s *Service) Run() error {
	data, err := s.prod.Produce()
	if err != nil {
		return err
	}

	maskedData := s.maskConcurrentlyWithFanIn(data)
	return s.pres.Present(maskedData)
}

// maskConcurrentlyWithFanIn - использует Fan-In паттерн для сбора результатов
func (s *Service) maskConcurrentlyWithFanIn(data []string) []string {
	if len(data) == 0 {
		return []string{}
	}

	// ШАГ 1: Создаем канал задач
	tasks := make(chan string, len(data))

	// ШАГ 2: Запускаем N рабочих горутин, каждая имеет СВОЙ канал результатов
	workers := 10
	if len(data) < workers {
		workers = len(data)
	}

	// Каждая рабочая горутина будет писать в свой отдельный канал
	workerChannels := make([]chan string, workers)

	var wg sync.WaitGroup

	// Запускаем рабочие горутины
	for i := 0; i < workers; i++ {
		// Каждая горутина получает свой канал для результатов
		workerChannels[i] = make(chan string, 10) // Буферизованный канал

		wg.Add(1)
		go func(workerID int, resultChan chan<- string) {
			defer wg.Done()
			defer close(resultChan) // Закрываем канал когда горутина завершится

			// Рабочая горутина читает задачи из общего канала
			for task := range tasks {
				// 1. Вызываем maskURL (основная работа)
				masked := s.maskURL(task)

				// 2. Отправляем результат в СВОЙ канал
				resultChan <- masked

				// Для отладки
				goID := getGoroutineID()
				fmt.Printf("Воркер %d (горутина %d) обработал: %.20s...\n",
					workerID, goID, task)
			}
		}(i+1, workerChannels[i])
	}

	// ШАГ 3: Запускаем горутину которая отправляет задачи
	go func() {
		for _, line := range data {
			tasks <- line
		}
		close(tasks) // Закрываем канал задач
	}()

	// ШАГ 4: Fan-In - собираем результаты из всех каналов в один
	fanInChan := make(chan string, len(data))

	// Запускаем горутину которая читает из всех каналов и объединяет
	var fanInWg sync.WaitGroup

	for _, workerChan := range workerChannels {
		fanInWg.Add(1)
		go func(sourceChan <-chan string) {
			defer fanInWg.Done()

			// Читаем из канала рабочей горутины
			for result := range sourceChan {
				// Отправляем в объединенный канал
				fanInChan <- result
			}
		}(workerChan)
	}

	// Закрываем fanInChan когда все читатели закончили
	go func() {
		fanInWg.Wait()
		close(fanInChan)
	}()

	// ШАГ 5: Собираем все результаты из объединенного канала
	var masked []string
	for result := range fanInChan {
		masked = append(masked, result)
	}

	return masked
}

// Старая версия (для сравнения)
func (s *Service) maskConcurrently(data []string) []string {
	if len(data) == 0 {
		return []string{}
	}

	tasks := make(chan string, len(data))
	results := make(chan string, len(data))
	var wg sync.WaitGroup

	workers := 10
	if len(data) < workers {
		workers = len(data)
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for task := range tasks {
				results <- s.maskURL(task)
			}
		}(i + 1)
	}

	go func() {
		for _, line := range data {
			tasks <- line
		}
		close(tasks)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var masked []string
	for result := range results {
		masked = append(masked, result)
	}

	return masked
}

func (s *Service) maskURL(text string) string {
	goID := getGoroutineID()
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		if elapsed > time.Millisecond {
			fmt.Printf("[Горутина %d] maskURL выполнилась за %v\n", goID, elapsed)
		}
	}()

	result := []byte(text)
	target := "http://"
	targetLen := len(target)

	i := 0
	for i <= len(text)-targetLen {
		if string(result[i:i+targetLen]) == target {
			start := i + targetLen
			for j := start; j < len(result) && result[j] != ' '; j++ {
				result[j] = '*'
			}
			i = start + 1
		} else {
			i++
		}
	}

	return string(result)
}

func getGoroutineID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	stack := string(buf[:n])
	if strings.HasPrefix(stack, "goroutine ") {
		fields := strings.Fields(stack)
		if len(fields) >= 2 {
			idStr := fields[1]
			id, err := strconv.Atoi(idStr)
			if err == nil {
				return id
			}
		}
	}
	return 0
}
