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
	return &Service{prod: prod, pres: pres}
}

func (s *Service) Run() error {
	data, err := s.prod.Produce()
	if err != nil {
		return err
	}

	results := s.mask(data)
	return s.pres.Present(results)
}

func (s *Service) mask(data []string) []string {
	if len(data) == 0 {
		return nil
	}

	// Каналы
	tasks := make(chan string, len(data))
	results := make([]string, len(data))
	var wg sync.WaitGroup

	// 10 горутин максимум
	workers := 10
	if len(data) < workers {
		workers = len(data)
	}

	// Запускаем горутины
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				// Читаем индекс из канала и маскируем
				for i, text := range data {
					if text == task {
						results[i] = s.MaskaAfterURL(text)
						break
					}
				}
			}
		}()
	}

	// Отправляем задачи
	for _, line := range data {
		tasks <- line
	}
	close(tasks)

	wg.Wait()
	return results
}

func (s *Service) MaskaAfterURL(text string) string {
	result := []byte(text)
	target := "http://"
	targetLen := len(target)

	i := 0
	for i <= len(text)-targetLen {
		if string(result[i:i+targetLen]) == target {
			start := i + targetLen
			for j := start; j < len(result); j++ {
				if result[j] == ' ' {
					break
				}
				result[j] = '*'
			}
			i = start + 1
		} else {
			i++
		}
	}
	return string(result)
}
