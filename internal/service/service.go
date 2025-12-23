package service

import "strings"

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

	result := s.MaskData(data)
	return s.pres.Present(result)
}

func (s *Service) MaskData(data []string) []string {
	var result []string

	for _, line := range data {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		words := strings.Fields(line)
		maskedWords := make([]string, len(words))

		for i, word := range words {
			maskedWords[i] = s.maskWord(word)
		}

		result = append(result, strings.Join(maskedWords, " "))
	}

	return result
}

func (s *Service) maskWord(word string) string {
	if len(word) <= 1 {
		return word
	}
	return string(word[0]) + strings.Repeat("*", len(word)-1)
}
