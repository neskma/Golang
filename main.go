package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func bubbleSort(arr []int) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}

func promptInput() []int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Введите массив из 6 цифр через пробел:")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	nums := strings.Split(input, " ")
	arr := make([]int, 0)

	for _, numStr := range nums {
		num, _ := strconv.Atoi(numStr)
		arr = append(arr, num)
	}

	return arr
}

func main() {
	arr := promptInput()

	fmt.Println("Исходный массив:", arr)

	bubbleSort(arr)

	fmt.Println("Отсортированный массив:", arr)
}
