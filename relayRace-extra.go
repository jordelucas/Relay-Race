package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const LANES = 6
const RUNNERS = 4
const MIN_TIME = 5

var waitGroup sync.WaitGroup

func main() {
	classificationChan := make(chan string)
	var classification []string
	count := 0

	fmt.Printf("\n--- RACE STARTED ---\n\n")

	waitGroup.Add(LANES)

	for i := 1; i <= LANES; i++ {
		go startTeam(fmt.Sprintf("Team %d", i), classificationChan)
	}

	for !(count == LANES) {
		select {
		case value := <-classificationChan:
			count = count + 1
			classification = append(classification, value)
		}
	}

	waitGroup.Wait()

	fmt.Printf("\n--- RACE ENDED ---\n")

	showResult(classification)

	time.Sleep(1e9)
}

func startTeam(teamName string, classificationChan chan string) {
	runnerChan := make(chan int, 1)
	timeChan := make(chan int)
	endChan := make(chan bool)

	finished := false
	totalTime := 0
	result := ""

	go runner(1, teamName, runnerChan, timeChan, endChan)

	for !finished {
		select {
		case order := <-runnerChan:
			go runner(order, teamName, runnerChan, timeChan, endChan)
		case value := <-timeChan:
			totalTime = totalTime + value
		case value := <-endChan:
			result = fmt.Sprintf("Duration of %s: %ds\n", teamName, totalTime)
			finished = value
		}
	}

	classificationChan <- result

	waitGroup.Done()
}

func runner(order int, teamName string, runnerChan chan int, timeChan chan<- int, endChan chan<- bool) {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	duration := random.Intn(MIN_TIME) + MIN_TIME //duração da corrida

	fmt.Printf("Runner %d of %s is starting \n", order, teamName)
	time.Sleep(time.Second * time.Duration(duration))
	fmt.Printf("Runner %d of %s finished after %ds \n", order, teamName, duration)

	timeChan <- duration

	if order == RUNNERS {
		endChan <- true
	} else {
		runnerChan <- order + 1
	}
}

func showResult(classification []string) {
	fmt.Println("\nResult: ")

	for index, value := range classification {
		if index == 0 {
			fmt.Printf("(winner) ")
		}
		fmt.Printf("%s", value)
	}
}
