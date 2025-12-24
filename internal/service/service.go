package service

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

	// Применяем вашу функцию к каждой строке
	var results []string
	for _, line := range data {
		masked := s.maskaAfterURL(line)
		results = append(results, masked)
	}

	return s.pres.Present(results)
}

func (s *Service) maskaAfterURL(text string) string {

	result := []byte(text) // ← ЗДЕСЬ определяем переменную result
	target := "https://"

	for i := 0; i <= len(text)-len(target); i++ {
		if string(result[i:i+len(target)]) == target {

			for j := i + len(target); j < len(result); j++ {

				if result[j] == ' ' {
					break
				}
				result[j] = '*'
			}
			i += len(target)
		}
	}
	return string(result)
}
