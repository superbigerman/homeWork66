package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"clean-project/internal/infrastructure"
	"clean-project/internal/service"
)

func main() {
	inputFile := flag.String("input", "", "Путь к входному файлу")
	outputFile := flag.String("output", "output.txt", "Путь к выходному файлу")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Использование: program -input <входной_файл> [-output <выходной_файл>]")
		fmt.Println("Пример: program -input data.txt -output result.txt")
		fmt.Println("По умолчанию выходной файл: output.txt")
		os.Exit(1)
	}

	if _, err := os.Stat(*inputFile); os.IsNotExist(err) {
		log.Fatalf("Файл %s не найден", *inputFile)
	}

	producer := infrastructure.NewFileProducer(*inputFile)
	presenter := infrastructure.NewFilePresenter(*outputFile)

	svc := service.NewService(producer, presenter)

	if err := svc.Run(); err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	fmt.Printf("Готово! Результат записан в %s\n", *outputFile)
}
