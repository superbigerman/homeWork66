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

	maskedData := s.maskConcurrentlyWithFanIn(data)
	return s.pres.Present(maskedData)
}

func (s *Service) maskConcurrentlyWithFanIn(data []string) []string {
	if len(data) == 0 {
		return []string{}
	}

	tasks := make(chan string, len(data))

	workers := 10
	if len(data) < workers {
		workers = len(data)
	}

	workerChannels := make([]chan string, workers)

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		workerChannels[i] = make(chan string, 10)

		wg.Add(1)
		go func(resultChan chan<- string) {
			defer wg.Done()
			defer close(resultChan)

			for task := range tasks {
				masked := s.maskURL(task)
				resultChan <- masked
			}
		}(workerChannels[i])
	}

	go func() {
		for _, line := range data {
			tasks <- line
		}
		close(tasks)
	}()

	fanInChan := make(chan string, len(data))

	var fanInWg sync.WaitGroup

	for _, workerChan := range workerChannels {
		fanInWg.Add(1)
		go func(sourceChan <-chan string) {
			defer fanInWg.Done()

			for result := range sourceChan {
				fanInChan <- result
			}
		}(workerChan)
	}

	go func() {
		fanInWg.Wait()
		close(fanInChan)
	}()

	var masked []string
	for result := range fanInChan {
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
