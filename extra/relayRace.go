package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const LANES = 6    //Quantidade de raias
const RUNNERS = 4  //Quantidade de corredores
const MIN_TIME = 5 //Tempo mínimo para cada corredor

var waitGroup sync.WaitGroup

func main() {
	classificationChan := make(chan string)
	var classification []string
	count := 0

	fmt.Printf("\n--- RACE STARTED ---\n\n")

	waitGroup.Add(LANES) //Indica quantas gorotinas/raias devem ser finalizadas para prosseguir com o resultado

	//Inicializa os times/raias de acordo com o valor respectiva constante
	for i := 1; i <= LANES; i++ {
		go startTeam(fmt.Sprintf("Team %d", i), classificationChan)
	}

	//Enquanto a quantidade de equipes que finalizaram for menor ao total, toda vez que o canal de classificação recebe
	//um novo calor, o contador é incrementado e o array de classificação é populado
	for !(count == LANES) {
		select {
		case value := <-classificationChan:
			count = count + 1
			classification = append(classification, value)
		}
	}

	waitGroup.Wait() //Aguarda as gorotinas

	fmt.Printf("\n--- RACE ENDED ---\n")

	showResult(classification) //Mostra a classificação

	time.Sleep(1e9)
}

func startTeam(teamName string, classificationChan chan string) {
	runnerChan := make(chan int, 1) //Canal para indicar qual corredor deve iniciar
	timeChan := make(chan int)      //Canal para receber o tempo decorrido de cada corredor
	endChan := make(chan bool)      //Canal para indicar o fim da corrida

	finished := false
	totalTime := 0
	result := ""

	go runner(1, teamName, runnerChan, timeChan, endChan)

	//Enquanto não for identificado o sinal de fim de corrida da equipe atual e toda vez que runnerChan receber um novo
	//valor, uma nova gorotine referente ao corredor será iniciada. Ao final, o tempo da equipe é guardado.
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

//Cada corredor recebe um número e a identificação do respectivo time, além dos canais de corrida, tempo e finalização
func runner(order int, teamName string, runnerChan chan int, timeChan chan<- int, endChan chan<- bool) {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	duration := random.Intn(MIN_TIME) + MIN_TIME //Gera um valor aleatório entre o tempo mínimo indicado e o seu dobrocorrida

	fmt.Printf("Runner %d of %s is starting \n", order, teamName)
	time.Sleep(time.Second * time.Duration(duration))
	fmt.Printf("Runner %d of %s finished after %ds \n", order, teamName, duration)

	timeChan <- duration

	if order == RUNNERS { //Se o número do corredor corresponder ao último, a corrida da equipe será finalizada
		endChan <- true
	} else { //Caso contrário, um novo corredor dará continuidade a corrida
		runnerChan <- order + 1
	}
}

//Mostra o resulta da corrida, indicando qual foi a equipe vencedora
func showResult(classification []string) {
	fmt.Println("\nResult: ")

	for index, value := range classification {
		if index == 0 {
			fmt.Printf("(winner) ")
		}
		fmt.Printf("%s", value)
	}
}
