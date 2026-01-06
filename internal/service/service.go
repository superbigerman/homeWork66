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
		masked := s.MaskaAfterURL(line)
		results = append(results, masked)
	}

	return s.pres.Present(results)
}

func (s *Service) MaskaAfterURL(text string) string {
	result := []byte(text)
	target := "http://"
	targetLen := len(target)

	i := 0
	for i <= len(text)-targetLen {
		if string(result[i:i+targetLen]) == target {
			// Начинаем маскировать после "http://"
			start := i + targetLen

			// Маскируем до пробела или конца строки
			for j := start; j < len(result); j++ {
				if result[j] == ' ' {
					break
				}
				result[j] = '*'
			}

			// Переходим к следующему символу после обработанного
			i = start + 1
		} else {
			i++
		}
	}
	return string(result)
}
