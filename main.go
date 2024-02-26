package main

import (
	"fmt"
	"os"
)

func main() {
	fileName := "log.txt"

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Ошибка при открытии файла: %s\n", err.Error())
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("Ошибка при получении информации о файле: %s\n", err.Error())
		return
	}

	if fileInfo.Size() == 0 {
		fmt.Printf("Файл %s пуст\n", fileName)
		return
	}

	var (
		buffer = make([]byte, fileInfo.Size())
		offset int64
	)

	for {
		n, err := file.ReadAt(buffer, offset)
		if err != nil {
			fmt.Printf("Ошибка при чтении файла: %s\n", err.Error())
			return
		}

		if n == 0 {
			break
		}

		fmt.Print(string(buffer[:n]))

		offset += int64(n)
	}
}
