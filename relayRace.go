package main

import (
	"fmt"
	"math/rand"
	"time"
)

const RUNNERS = 4
const MIN_TIME = 5

func main() {
	runnerChan := make(chan int, 1)
	timeChan := make(chan int)
	endChan := make(chan bool)

	finished := false
	totalTime := 0

	fmt.Printf("\n--- RACE STARTED ---\n\n")

	go runner(1, runnerChan, timeChan, endChan)

	for !finished {
		select {
		case order := <-runnerChan:
			go runner(order, runnerChan, timeChan, endChan)
		case value := <-timeChan:
			totalTime = totalTime + value
		case finished = <-endChan:
		}
	}

	fmt.Printf("--- RACE ENDED ---\n\n")

	fmt.Printf("Duration: %ds\n", totalTime)
}

func runner(order int, runnerChan chan int, timeChan chan<- int, endChan chan<- bool) {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	duration := random.Intn(MIN_TIME) + MIN_TIME //duração da corrida

	fmt.Printf("Runner %d is starting \n", order)
	time.Sleep(time.Second * time.Duration(duration))
	fmt.Printf("Runner %d finished after %ds \n\n", order, duration)

	timeChan <- duration

	if order == RUNNERS {
		endChan <- true
	} else {
		runnerChan <- order + 1
	}
}
