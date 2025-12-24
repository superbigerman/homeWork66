package internal

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// Мок Presenter
type MockPresenter struct {
	mock.Mock
}

func (m *MockPresenter) Present(id string, data interface{}) error {
	args := m.Called(id, data)
	return args.Error(0)
}

// Мок Producer
type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) Produce(msg string) error {
	args := m.Called(msg)
	return args.Error(0)
}

// Тест Presenter
func TestPresenter_Present(t *testing.T) {
	presenter := new(MockPresenter)
	testData := map[string]string{"key": "value"}

	presenter.On("Present", "test-id", testData).Return(nil)

	err := presenter.Present("test-id", testData)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	presenter.AssertExpectations(t)
}

// Тест Producer
func TestProducer_Produce(t *testing.T) {
	producer := new(MockProducer)

	producer.On("Produce", "test message").Return(nil)

	err := producer.Produce("test message")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	producer.AssertExpectations(t)
}
