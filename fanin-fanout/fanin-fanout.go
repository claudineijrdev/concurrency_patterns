package main

import "fmt"

func cymbal(drum chan string) {
	drum <- "-tsssss-"
}

func hiHat(drum chan string) {
	drum <- "-chik-chik-"
}

func snare(drum chan string) {
	drum <- "-tak-"
}

func bassDrum(drum chan string) {
	drum <- "-boom-"
}

func main() {
	drumSet := make(chan string)
	var song string

	go hiHat(drumSet)
	go bassDrum(drumSet)
	go snare(drumSet)
	go cymbal(drumSet)

	for i := 0; i < 4; i++ {
		song += <-drumSet
	}
	fmt.Println(song)

}
