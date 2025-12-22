package infrastructure

import (
	"os"
	"strings"

	"clean-project/internal/service"
)

// FilePresenter записывает данные в файл
type FilePresenter struct {
	filePath string
}

// NewFilePresenter создает новый FilePresenter
func NewFilePresenter(filePath string) *FilePresenter {
	return &FilePresenter{
		filePath: filePath,
	}
}

// Present записывает данные в файл
func (fp *FilePresenter) Present(data []string) error {
	content := strings.Join(data, "\n")

	file, err := os.Create(fp.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

// Проверка реализации интерфейса
var _ service.Presenter = (*FilePresenter)(nil)
