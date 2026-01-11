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

	tasks := make(chan string, len(data))

	results := make(chan string, len(data))

	var wg sync.WaitGroup

	workers := 10
	if len(data) < workers {
		workers = len(data)
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for task := range tasks {

				masked := s.maskURL(task)

				results <- masked
			}
		}()
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
