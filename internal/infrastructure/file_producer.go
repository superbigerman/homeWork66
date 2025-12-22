package infrastructure

import (
	"bufio"
	"os"

	"clean-project/internal/service"
)

// FileProducer читает данные из файла
type FileProducer struct {
	filePath string
}

// NewFileProducer создает новый FileProducer
func NewFileProducer(filePath string) *FileProducer {
	return &FileProducer{
		filePath: filePath,
	}
}

// Produce читает строки из файла
func (fp *FileProducer) Produce() ([]string, error) {
	file, err := os.Open(fp.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// Проверка реализации интерфейса
var _ service.Producer = (*FileProducer)(nil)
