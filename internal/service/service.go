package service

// Producer интерфейс для поставщика данных
type Producer interface {
	Produce() ([]string, error)
}

// Presenter интерфейс для вывода результатов
type Presenter interface {
	Present([]string) error
}

// Service основной сервис
type Service struct {
	prod Producer
	pres Presenter
}

// NewService создает новый сервис
func NewService(prod Producer, pres Presenter) *Service {
	return &Service{
		prod: prod,
		pres: pres,
	}
}

// Run запускает обработку данных
func (s *Service) Run() error {
	data, err := s.prod.Produce()
	if err != nil {
		return err
	}

	result := MaskData(data)
	return s.pres.Present(result)
}
