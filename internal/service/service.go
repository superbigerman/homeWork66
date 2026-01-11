package service

import (
	"sync"
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

	maskedData := s.maskConcurrently(data)
	return s.pres.Present(maskedData)
}

func (s *Service) maskConcurrently(data []string) []string {
	if len(data) == 0 {
		return []string{}
	}

	// 1️⃣ ПЕРВЫЙ канал: для задач
	tasks := make(chan string, len(data))

	// 2️⃣ ВТОРОЙ канал: для результатов (все воркеры пишут сюда)
	results := make(chan string, len(data))

	var wg sync.WaitGroup

	// Определяем количество воркеров (максимум 10)
	workers := 10
	if len(data) < workers {
		workers = len(data)
	}

	// Запускаем воркеров
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Каждый воркер читает задачи из ОБЩЕГО канала tasks
			for task := range tasks {
				// Обрабатываем задачу
				masked := s.maskURL(task)

				// ⚠️ ВСЕ воркеры пишут в ОДИН канал results!
				// Может быть небольшая конкуренция, но для быстрых операций это нормально
				results <- masked
			}
		}()
	}

	// Отправляем все задачи в канал tasks
	go func() {
		for _, line := range data {
			tasks <- line
		}
		close(tasks) // Закрываем канал задач - сигнал воркерам что задач больше нет
	}()

	// Закрываем канал результатов когда все воркеры завершились
	go func() {
		wg.Wait()      // Ждем завершения всех воркеров
		close(results) // Закрываем канал результатов
	}()

	// Собираем все результаты из канала results
	var masked []string
	for result := range results {
		masked = append(masked, result)
	}

	return masked
}

func (s *Service) maskURL(text string) string {
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
