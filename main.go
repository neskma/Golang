package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

func main() {
	var text string

	if len(os.Args) > 1 {
		text = strings.Join(os.Args[1:], " ")
	} else {
		fmt.Println("Введите предложение в качестве аргумента")
		return
	}

	textWords := make(map[rune]int)
	totalWords := 0

	// Подсчет количества каждой введенной буквы
	for _, word := range text {
		word = unicode.ToLower(word)
		if word >= 'a' && word <= 'z' {
			textWords[word]++
			totalWords++
		}
	}

	// Посчет доли букв в процентах
	fmt.Println("Letter - Count - Percentage")
	for word, count := range textWords {
		percentage := float64(count) / float64(totalWords) * 100
		fmt.Printf("%c - %d - %.2f%%\n", word, count, percentage)
	}
}
