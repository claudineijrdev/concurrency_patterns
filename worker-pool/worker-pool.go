package main

import (
	"fmt"
	"time"
)

func start(workersCount int, maxNumber int) (time.Duration, []int) {
	start := time.Now()
	tasks := make(chan int, maxNumber)
	results := make(chan int, maxNumber)
	nums := make([]int, maxNumber)

	for i := 0; i < workersCount; i++ {
		go worker(tasks, results)
	}

	for i := 0; i < maxNumber; i++ {
		tasks <- i
	}

	close(tasks)

	for i := 0; i < maxNumber; i++ {
		results := <-results
		nums[i] = results
	}
	elapsed := time.Since(start)
	return elapsed, nums
}

func worker(tasks <-chan int, results chan<- int) {
	for task := range tasks {
		results <- count(task)
	}
}

func count(n int) int {
	c := 1
	for i := 1; i <= n; i++ {
		c++
	}
	return c
}

func main() {
	maxNumber := 10000
	timeOneWorker, numsOw := start(1, maxNumber)
	timeTwoWorkers, _ := start(2, maxNumber)
	timeThreeWorkers, _ := start(3, maxNumber)
	tomeFourWorkers, numsFw := start(4, maxNumber)

	fmt.Println("1 worker: ", timeOneWorker)
	fmt.Println("2 workers: ", timeTwoWorkers)
	fmt.Println("3 workers: ", timeThreeWorkers)
	fmt.Println("4 workers: ", tomeFourWorkers)

	fmt.Println("-----------------")
	fmt.Println("1 Workers", numsOw[:100])
	fmt.Println("4 Workers", numsFw[:100])

	fmt.Println("end")
}
