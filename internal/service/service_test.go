package service

import "testing"

type TestProducer struct {
	data []string
	err  error
}

func (tp *TestProducer) Produce() ([]string, error) { return tp.data, tp.err }

type TestPresenter struct {
	received []string
	err      error
}

func (tp *TestPresenter) Present(data []string) error { tp.received = data; return tp.err }

func TestService(t *testing.T) {
	// Тест maskaAfterURL
	s := &Service{}

	if got := s.maskaAfterURL("http://site.com"); got != "http://********" {
		t.Errorf("Не работает маскировка: %s", got)
	}

	// Тест NewService
	p := &TestProducer{}
	pr := &TestPresenter{}
	svc := NewService(p, pr)

	if svc == nil {
		t.Error("NewService вернул nil")
	}

	// Тест Run
	p.data = []string{"Visit http://example.com"}
	if err := svc.Run(); err != nil {
		t.Errorf("Run() ошибка: %v", err)
	}

	if pr.received == nil || pr.received[0] != "Visit http://***********" {
		t.Errorf("Неправильный результат: %v", pr.received)
	}
}
