package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

func main() {
	fileName := "log.txt"
	exitCommand := "exit"

	// Создаем или открываем файл для записи
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Ошибка при открытии файла: %s\n", err.Error())
		return
	}
	defer file.Close()

	// Бесконечный цикл для ввода строк
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Введите строку (или 'exit' для выхода): ")
		scanner.Scan()
		input := scanner.Text()

		// Проверяем наличие команды выхода
		if input == exitCommand {
			break
		}

		// Формируем строку для записи в файл
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		lineNumber := getFileLineCount(fileName) + 1
		logLine := fmt.Sprintf("%s. %s %s\n", strconv.Itoa(lineNumber), timestamp, input)

		// Записываем строку в файл
		_, err := file.WriteString(logLine)
		if err != nil {
			fmt.Printf("Ошибка при записи в файл: %s\n", err.Error())
			return
		}

		fmt.Println("Запись успешно добавлена в файл.")

		// Читаем и выводим строки из файла
		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Файл %s не существует\n", fileName)
			} else {
				fmt.Printf("Ошибка при чтении файла: %s\n", err.Error())
			}
			return
		}

		if len(data) == 0 {
			fmt.Printf("Файл %s пуст\n", fileName)
		} else {
			fmt.Println("Содержимое файла:")
			fmt.Println(string(data))
		}
	}
}

// Функция для получения количества строк в файле
func getFileLineCount(fileName string) int {
	file, err := os.Open(fileName)
	if err != nil {
		return 0
	}

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	return lineCount
}
