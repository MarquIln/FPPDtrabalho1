// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// PROBLEMA:
//   o dorminhoco especificado no arquivo Ex1-ExplanacaoDoDorminhoco.pdf nesta pasta
// ESTE ARQUIVO
//   Um template para criar um anel generico.
//   Adapte para o problema do dorminhoco.
//   Nada está dito soordem de processos qubre como funciona a e batem.
//   O ultimo leva a rolhada ...
//   ESTE  PROGRAMA NAO FUNCIONA.    É UM RASCUNHO COM DICAS.


package main

import (
	"fmt"
)

const NJ = 5  // numero de jogadores
const M = 4  // numero de cartas
const J = 1 // carta coringa

type carta string      // carta é uma string
var ch [NJ]chan carta  // NJ canais de itens tipo carta que é string 

func jogador(id int, in chan carta, out chan carta, cartasIniciais []carta) {
	mao := cartasIniciais    // estado local - as cartas na mao do jogador
	nroDeCartas := M
  	cartaRecebida := " "

	for {
		if nroDeCartas > 0 {
			fmt.Println(id, " joga") // escreve seu identificador
			cartaParaSair := mao[0]
			// guarda carta que entrou 
			mao[M+1] = cartaRecebida
			// manda carta escolhida o proximo
			out <- cartaParaSair        
		} else {  
			// ...
			cartaRecebida := <-in   // recebe carta na entrada
			//  ...
			// e se alguem bate ?
			fmt.Println(id, " bateu") // escreve seu identificador que bateu
		}
	}
}

func main() {
	for i := 0; i < NJ; i++ {
		ch[i] = make(chan struct{})
	}
	for i := 0; i < NJ; i++ {
		go jogador(i, ch[i], ch[(i+1)%N], cartasEscolhidas , ...)
	}
	
	<-make(chan struct{})
}


