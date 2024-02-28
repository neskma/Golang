package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	done := make(chan struct{})
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		for {
			select {
			case <-done:
				fmt.Println("Завершение программы...")
				return
			default:
				for i := 1; ; i++ {
					fmt.Println(i * i)
					time.Sleep(1 * time.Second)
				}
			}
		}
	}()

	<-signals
	close(done)
	fmt.Println("\nПрограмма завершает работу...")
}
