package main

import (
	"fmt"
	"os"
)

func main() {
	// Запрашиваем название файла от пользователя
	fmt.Print("Введите название файла: ")
	var fileName string
	fmt.Scanln(&fileName)

	// Создаем файл только для чтения
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE|os.O_TRUNC, 0444)
	if err != nil {
		fmt.Printf("Ошибка при создании файла: %s\n", err.Error())
		return
	}
	defer file.Close()

	fmt.Printf("Файл %s успешно создан и защищен от записи.\n", fileName)
}
