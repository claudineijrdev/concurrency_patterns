package main

import "fmt"

// Cria um canal de inteiros e o preenche com os valores do slice de inteiros
func generate(data []int) chan int {
	out := make(chan int)
	go func() {
		for _, v := range data {
			out <- v
		}
		close(out)
	}()
	return out
}

// Recebe um canal de inteiros e retorna um canal de inteiros com os valores pares
func filter(in chan int) chan int {
	out := make(chan int)
	go func() {
		for v := range in {
			if v%2 == 0 {
				out <- v
			}
		}
		close(out)
	}()
	return out
}

// Recebe um canal de inteiros e retorna um canal de inteiros com os valores ao quadrado
func square(in chan int) chan int {
	out := make(chan int)
	go func() {
		for v := range in {
			out <- v * v
		}
		close(out)
	}()
	return out
}

func main() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	c1 := generate(data)
	c2 := filter(c1)
	c3 := square(c2)

	for v := range c3 {
		fmt.Println(v)
	}
}
