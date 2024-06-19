package main

import "fmt"

func writer(text string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			c <- fmt.Sprintf("Hello %s", text)
		}
	}()
	return c
}

func main() {
	c := writer("world")
	for i := 0; i < 5; i++ {
		fmt.Println(<-c)
	}
}
