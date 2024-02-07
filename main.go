package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Укажите полный путь до файла вторым аргументом")
	}

	filePth := os.Args[1]

	var fileName, fileExt string

	fileName = strings.TrimSuffix(filepath.Base(filePth), filepath.Ext(filePth))
	fileExt = strings.TrimPrefix(filepath.Ext(filePth), ".")

	fmt.Printf("filename: %s\n", fileName)
	fmt.Printf("extension: %s\n", fileExt)
}
