package main

import (
	"fmt"
	"sync"
)

func writer(text string, wg *sync.WaitGroup) {
	fmt.Printf("Hello %s\n", text)
	wg.Done()
}
func main() {
	var wg sync.WaitGroup
	wg.Add(3)
	go writer("world", &wg)
	go writer("golang", &wg)
	go writer("universe", &wg)
	wg.Wait()
	fmt.Println("All done!")
}
