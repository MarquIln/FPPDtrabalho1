// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// PROBLEMA:
//   o dorminhoco especificado no arquivo Ex1-ExplanacaoDoDorminhoco.pdf nesta pasta
// ESTE ARQUIVO
//   Um template para criar um anel generico.
//   Adapte para o problema do dorminhoco.
//   Nada está dito sobre como funciona a ordem de processos que batem.
//   O ultimo leva a rolhada ...
//   ESTE  PROGRAMA NAO FUNCIONA.    É UM RASCUNHO COM DICAS.

package main

import (
	"fmt"
	"math/rand"
)

const NJ = 5 // numero de jogadores
const M = 4  // numero de cartas

type card string // card é um strirng

var ch [NJ]chan card // NJ canais de itens tipo card

var newRound chan bool
var flap chan int

func jogador(id int, in chan card, out chan card, initalCards []card, gameBegin chan bool) {
	hand := initalCards // estado local - as cartas na mão do jogador
	nCards := M         // quantas cartas ele tem
	reciviedCard := card("")
	fmt.Printf("Jogador %d cartas na mão: %v\n", id, hand)

	<-gameBegin

	bateuAntes := false

	for {
		if len(flap) != 0 && !bateuAntes { //se alguém bateu e não possui o id dentro do canal
			fmt.Printf("Jogador %d bateu também!\n", id)
			if len(flap) == NJ-1 {
				fmt.Printf("Jogador %d foi o último a bater e perdeu o jogo.\n", id)
			}
			flap <- id
			return
		}

		if nCards == 5 {
			fmt.Println(id, " joga")
			aux := rand.Intn(nCards)
			cardToLeave := hand[aux]

			if hasSameNaipe(hand) {
				cardToLeave = differentCard(hand)
			}

			fmt.Printf("Jogador %d escolheu a card %s\n para passar adiante", id, cardToLeave)

			out <- cardToLeave
			hand = append(hand[:aux], hand[aux+1:]...)
			fmt.Printf("Nova mão do jogador (%d): %v\n", id, hand)
			nCards--

			if hasSameNaipe(hand) && nCards == 4 {
				fmt.Printf("Jogador %d bate!\n", id)
				bateuAntes = true
				flap <- id
				newRound <- false
			}

			if nCards == 0 {
				newRound <- true
			}
		} else {
			select {
			case reciviedCard = <-in:
				fmt.Printf("Jogador %d recebeu a card %s\n", id, reciviedCard)
				hand = append(hand, reciviedCard)
				nCards++
				fmt.Printf("Nova mão do jogador %d: %v\n", id, hand)
			default: // wait
			}
		}
	}
}

func hasSameNaipe(hand []card) bool {
	count := make(map[card]int)
	for _, n := range hand {
		count[n]++
		if count[n] >= 4 {
			return true
		}
	}
	return false
}

func differentCard(hand []card) card {
	count := make(map[card]int)
	for _, n := range hand {
		count[n]++
		if count[n] == 1 {
			return n
		}
	}
	return " "
}

func main() {
	flap = make(chan int, 5)
	for i := 0; i < NJ; i++ {
		ch[i] = make(chan card)
	}

	gameBegin := make(chan bool) // Canal para sinalizar o início do jogo
	newRound = make(chan bool)   // Canal para sinalizar o início de uma nova rodada

	for i := 0; i < NJ; i++ {
		chooseCards := make([]card, 4)
		for j := 0; j < 4; j++ {
			randomIndex := rand.Intn(4)
			chooseCards[j] = card(rune('A' + randomIndex))
		}
		go jogador(i, ch[i], ch[(i+1)%NJ], chooseCards, gameBegin)
	}
	fmt.Println()

	for i := 0; i < NJ; i++ {
		gameBegin <- true
	}

	initalCard := card("Joker")
	ch[0] <- initalCard

	for i := 0; i < NJ; i++ {
		<-newRound
	}
}
