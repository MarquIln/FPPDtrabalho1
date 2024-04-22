// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// PROBLEMA:
//   o dorminhoco especificado no arquivo Ex1-ExplanacaoDoDorminhoco.pdf nesta pasta
// ESTE ARQUIVO
//   Um template para criar um anel generico.
//   Adapte para o problema do dorminhoco.
//   Nada está dito sobre como funciona a ordem de processos que batem.
//   O ultimo leva a rolhada ...
//   ESTE  PROGRAMA NAO FUNCIONA.    É UM RASCUNHO COM DICAS.
// AGORA ESTA FUNCIONANDO

package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Card string

const (
	totalPlayers = 5
	totalCards   = 4
)

var (
	playersChannels [totalPlayers]chan Card
	newRound        chan bool
	flapPlayerId    chan int
)

func player(id int, in chan Card, out chan Card, initialCards []Card, gameBegin chan bool) {
	hand := initialCards
	nCards := totalCards
	receivedCard := Card(" ")
	fmt.Printf("Jogador %d tem essas cartas na mão: %v\n", id, hand)

	<-gameBegin
	earlyFlap := false
	for {
		if len(flapPlayerId) != 0 && !earlyFlap {
			fmt.Printf("Jogador %d bateu também!\n", id)
			if len(flapPlayerId) == totalPlayers-1 {
				fmt.Printf("Como, o Jogador %d foi o último a bater e ele perdeu o jogo.\n", id)
			}
			flapPlayerId <- id
			return
		}

		if nCards == 5 {
			fmt.Println(id, " joga")
			aux := rand.Intn(nCards)
			cardToLeave := hand[aux]
			if hasSameSuit(hand) {
				cardToLeave = differentCard(hand)
			}
			fmt.Printf("Jogador %d escolheu a Carta %s\n para passar adiante ", id, cardToLeave)

			out <- cardToLeave
			hand = append(hand[:aux], hand[aux+1:]...)
			fmt.Printf("Nova mão do jogador %d: %v\n", id, hand)
			nCards--

			if hasSameSuit(hand) && nCards == 4 {
				fmt.Printf("Jogador %d bate!\n", id)
				earlyFlap = true
				flapPlayerId <- id
				newRound <- false
			}

			if nCards == 0 {
				newRound <- true
			}
		} else {
			select {
			case receivedCard = <-in:
				fmt.Printf("Jogador %d recebeu a Carta %s\n", id, receivedCard)
				hand = append(hand, receivedCard)
				nCards++
				fmt.Printf("Nova mão do jogador %d: %v\n", id, hand)
			default: // wait
			}
		}
	}
}

func hasSameSuit(hand []Card) bool {
	count := make(map[Card]int)
	for _, n := range hand {
		count[n]++
		if count[n] >= 4 {
			return true
		}
	}
	return false
}

func differentCard(hand []Card) Card {
	count := make(map[Card]int)
	for _, n := range hand {
		count[n]++
		if count[n] == 1 {
			return n
		}
	}
	return " "
}

func main() {
	flapPlayerId = make(chan int, 5)
	for i := 0; i < totalPlayers; i++ {
		playersChannels[i] = make(chan Card)
	}

	gameBegin := make(chan bool)
	newRound = make(chan bool)

	for i := 0; i < totalPlayers; i++ {
		chooseCards := make([]Card, 4)
		for j := 0; j < 4; j++ {
			randomIndex := rand.Intn(4)
			chooseCards[j] = Card(rune('D' + randomIndex))
		}
		go player(i, playersChannels[i], playersChannels[(i+1)%totalPlayers], chooseCards, gameBegin)
	}

	for i := 0; i < totalPlayers; i++ {
		gameBegin <- true
	}

	initialCard := Card("Joker")
	playersChannels[0] <- initialCard

	time.Sleep(1 * time.Second)
}
