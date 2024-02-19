package main

import (
	"fmt"
	"time"

	"github.com/eiannone/keyboard"
)

type Item struct {
	Name string
	Date time.Time
	Tags string
	Link string
}

func main() {
	var links []Item

	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}

		switch char {
		case 'a': // Добавить новую ссылку
			var newItem Item
			fmt.Println("Enter name:")
			fmt.Scanln(&newItem.Name)
			newItem.Date = time.Now()
			fmt.Println("Enter tags:")
			fmt.Scanln(&newItem.Tags)
			fmt.Println("Enter link:")
			fmt.Scanln(&newItem.Link)

			links = append(links, newItem)
			fmt.Println("Link added successfully!")

		case 'd': // Удалить ссылку
			if len(links) == 0 {
				fmt.Println("No links to delete")
			} else {
				fmt.Println("Links:")
				for i, link := range links {
					fmt.Printf("%d. %s\n", i+1, link.Name)
				}
				fmt.Println("Enter the number of the link you want to delete:")
				var num int
				fmt.Scanln(&num)

				if num <= 0 || num > len(links) {
					fmt.Println("Invalid link number")
				} else {
					links = append(links[:num-1], links[num:]...)
					fmt.Println("Link deleted successfully!")
				}
			}

		case 'l': // Вывести список ссылок
			if len(links) == 0 {
				fmt.Println("No links to display")
			} else {
				fmt.Println("Links:")
				for i, link := range links {
					fmt.Printf("%d. Name: %s, Tags: %s, Link: %s, Date Added: %s\n", i+1, link.Name, link.Tags, link.Link, link.Date.Format(time.RFC3339))
				}
			}

		case 'q': // Выйти из программы
			fmt.Println("Quitting program...")
			return

		default:
			fmt.Println("Invalid input. Press 'a' to add a new link, 'd' to delete a link, 'l' to list all links, or 'q' to quit")
		}
	}
}
