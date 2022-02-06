package main

import (
	"fmt"
	"math/rand"
	"time"
)

const RUNNERS = 4  //Quantidade de corredores
const MIN_TIME = 5 //Tempo mínimo para cada corredor

func main() {
	runnerChan := make(chan int, 1) //Canal para indicar qual corredor deve iniciar
	timeChan := make(chan int)      //Canal para receber o tempo decorrido de cada corredor
	endChan := make(chan bool)      //Canal para indicar o fim da corrida

	finished := false
	totalTime := 0

	fmt.Printf("\n--- RACE STARTED ---\n\n")

	go runner(1, runnerChan, timeChan, endChan)

	//Enquanto não for identificado o sinal de fim de corrida e toda vez que runnerChan receber um novo valor, uma
	//nova gorotine referente ao corredor será iniciada.
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

//Cada corredor recebe um número e os canais de corrida, tempo e finalização
func runner(order int, runnerChan chan int, timeChan chan<- int, endChan chan<- bool) {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	duration := random.Intn(MIN_TIME) + MIN_TIME //Gera um valor aleatório entre o tempo mínimo indicado e o seu dobro

	fmt.Printf("Runner %d is starting \n", order)
	time.Sleep(time.Second * time.Duration(duration))
	fmt.Printf("Runner %d finished after %ds \n\n", order, duration)

	timeChan <- duration

	if order == RUNNERS { //Se o número do corredor corresponder ao último, a corrida será finalizada
		endChan <- true
	} else { //Caso contrário, um novo corredor dará continuidade a corrida
		runnerChan <- order + 1
	}
}
