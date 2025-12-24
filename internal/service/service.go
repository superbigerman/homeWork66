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

	result := s.maskData(data)
	return s.pres.Present(result)
}

// maskData - ТОЧНО ТАКАЯ ЖЕ как maskaAfterURL, маскирует после https://
func (s *Service) maskData(data []string) []string {
	var result []string

	// Для каждой строки в данных
	for _, line := range data {
		// Превращаем строку в массив байт
		text := []byte(line)
		target := []byte("https://")

		// Ищем "https://" в строке
		for i := 0; i <= len(text)-len(target); i++ {
			// Проверяем совпадение
			match := true
			for t := 0; t < len(target); t++ {
				if text[i+t] != target[t] {
					match = false
					break
				}
			}

			if match {
				// Маскируем всё после "https://" до пробела
				for j := i + len(target); j < len(text); j++ {
					if text[j] == ' ' {
						break
					}
					text[j] = '*'
				}
				i += len(target) - 1
			}
		}

		result = append(result, string(text))
	}

	return result
}
