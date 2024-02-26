package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	file, err := os.Create("log.txt")
	if err != nil {
		fmt.Println("Не удалось создать файл")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(os.Stdin)
	lineNumber := 1

	for {
		fmt.Print("Введите строку (для выхода введите 'exit'): ")
		scanner.Scan()
		input := scanner.Text()

		if input == "exit" {
			break
		}

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		logMessage := fmt.Sprintf("%d %s %s\n", lineNumber, timestamp, input)
		fmt.Fprint(file, logMessage)

		lineNumber++
	}

	fmt.Println("Программа завершена. Лог записан в файл log.txt.")
}
