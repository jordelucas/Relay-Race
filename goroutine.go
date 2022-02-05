package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	runnerChan := make(chan int, 1)
	timeChan := make(chan int)
	endChan := make(chan bool)

	fmt.Printf("- - RACE STARTED - - \n\n")

	go runner(1, runnerChan, timeChan, endChan)

	finished := false
	totalTime := 0

	for !finished {
		select {
		case order := <-runnerChan:
			go runner(order, runnerChan, timeChan, endChan)
		case value := <-timeChan:
			totalTime = totalTime + value
		case finished = <-endChan:
		}
	}

	fmt.Printf("- - RACE ENDED - - \n")

	fmt.Printf("Duration: %ds\n", totalTime)
}

func runner(order int, runnerChan chan int, timeChan chan<- int, endChan chan<- bool) {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	duration := random.Intn(5) + 5 //duração da corrida

	fmt.Printf("Runner %d is starting \n", order)
	time.Sleep(time.Second * time.Duration(duration))
	fmt.Printf("Runner %d finished after %ds \n\n", order, duration)

	runnerChan <- order + 1
	timeChan <- duration

	if order == 4 {
		fmt.Println("")
		endChan <- true
	}
}
