package service_test

import (
	"testing"

	"clean-project/internal/service"
	"clean-project/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Run_WithMocks(t *testing.T) {
	// Тест 1: Проверяем маскировку с моками
	t.Run("маскирует данные через моки", func(t *testing.T) {
		// 1. СОЗДАЁМ МОКИ (ваши сгенерированные)
		mockProd := new(mocks.Producer)
		mockPres := new(mocks.Presenter)

		// 2. Данные, которые ВЕРНЁТ мок Producer
		inputFromMock := []string{
			"Сайт http://example.com",
			"Иди на http://test.org сейчас",
			"Просто текст",
		}

		// 3. Данные, которые ДОЛЖЕН ПОЛУЧИТЬ мок Presenter
		//    (после маскировки вашей функцией)
		expectedToPresenter := []string{
			"Сайт http://***********",
			"Иди на http://******** сейчас",
			"Просто текст",
		}

		// 4. НАСТРАИВАЕМ МОКИ
		//    "Когда у mockProd вызовут Produce(), верни inputFromMock"
		mockProd.On("Produce").Return(inputFromMock, nil)

		//    "Когда у mockPres вызовут Present(), он должен получить expectedToPresenter"
		mockPres.On("Present", expectedToPresenter).Return(nil)

		// 5. СОЗДАЁМ СЕРВИС С МОКАМИ
		svc := service.NewService(mockProd, mockPres)

		// 6. ЗАПУСКАЕМ ТЕСТИРУЕМЫЙ МЕТОД
		err := svc.Run()

		// 7. ПРОВЕРЯЕМ
		assert.NoError(t, err, "Run() не должен возвращать ошибку")

		// Проверяем, что моки вызвались как ожидалось
		mockProd.AssertExpectations(t)
		mockPres.AssertExpectations(t)

		t.Logf("✅ Успех! Producer вернул %d строк, Presenter получил замаскированные данные", len(inputFromMock))
	})

	// Тест 2: Проверяем обработку ошибок с моками
	t.Run("обрабатывает ошибку от мока Producer", func(t *testing.T) {
		mockProd := new(mocks.Producer)
		mockPres := new(mocks.Presenter)

		// Мок Producer возвращает ошибку
		mockProd.On("Produce").Return([]string{}, assert.AnError)

		svc := service.NewService(mockProd, mockPres)
		err := svc.Run()

		// Должна быть ошибка
		assert.Error(t, err, "Должна быть ошибка от Producer")

		// Presenter не должен вызываться при ошибке Producer
		mockPres.AssertNotCalled(t, "Present")
		mockProd.AssertExpectations(t)
	})

	t.Run("работает с пустыми данными от мока", func(t *testing.T) {
		// Настройка мока Producer
		prodMock := &mocks.Producer{}
		prodMock.On("Produce").Return([]string{}, nil) // Пустой слайс

		// Настройка мока Presenter
		presMock := &mocks.Presenter{}
		// Принимаем либо пустой слайс, либо nil
		presMock.On("Present", mock.AnythingOfType("[]string")).Return(nil).Run(func(args mock.Arguments) {
			data := args.Get(0).([]string)
			assert.True(t, len(data) == 0, "Должен быть пустой слайс")
		})

		// Создание сервиса
		s := service.NewService(prodMock, presMock)

		// Запуск
		err := s.Run()

		// Проверки
		assert.NoError(t, err, "Не ожидалась ошибка")

		// Проверяем вызовы моков
		prodMock.AssertExpectations(t)
		presMock.AssertExpectations(t)
	})
}

// Тест для конкретных случаев маскировки
func TestService_MaskingScenarios(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string // что вернёт мок Producer
		expected []string // что должен получить мок Presenter
	}{
		{
			name:     "один URL",
			input:    []string{"Перейдите http://site.com"},
			expected: []string{"Перейдите http://********"}, // site.com = 8 символов
		},
		{
			name:  "несколько URL",
			input: []string{"http://one.com и http://two.org"},
			// one.com = 7 символов, two.org = 8 символов
			expected: []string{"http://******* и http://*******"},
		},
		{
			name:     "https не маскируется",
			input:    []string{"Безопасно: https://bank.ru"},
			expected: []string{"Безопасно: https://bank.ru"},
		},
		{
			name:  "смесь http и https",
			input: []string{"http://insecure.com https://secure.com"},
			// insecure.com = 12 символов
			expected: []string{"http://************ https://secure.com"},
		},
		{
			name:  "URL в конце строки",
			input: []string{"Ссылка: http://end"},
			// end = 3 символа
			expected: []string{"Ссылка: http://***"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Для каждого теста создаём новые моки
			mockProd := new(mocks.Producer)
			mockPres := new(mocks.Presenter)

			mockProd.On("Produce").Return(tc.input, nil)
			mockPres.On("Present", tc.expected).Return(nil)

			svc := service.NewService(mockProd, mockPres)
			err := svc.Run()

			assert.NoError(t, err)

			// Проверяем, что моки вызвались с правильными параметрами
			mockProd.AssertCalled(t, "Produce")
			mockPres.AssertCalled(t, "Present", tc.expected)
		})
	}
}
