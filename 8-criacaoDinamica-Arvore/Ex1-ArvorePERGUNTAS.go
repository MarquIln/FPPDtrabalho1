// por Fernando Dotti - PUCRS
// dado abaixo um exemplo de estrutura em arvore, uma arvore inicializada
// e uma operação de caminhamento, pede-se fazer:
//   1.a) a operação que soma todos elementos da arvore.
//        func soma(r *Nodo) int {...}
//   1.b) uma operação concorrente que soma todos elementos da arvore
//   2.a) a operação de search de um elemento v, dizendo true se encontrou v na árvore, ou falso
//        func search(r* Nodo, v int) bool {}...}
//   2.b) a operação de search concorrente de um elemento, que informa imediatamente
//        por um canal se encontrou o elemento (sem acabar a search), ou informa
//        que nao encontrou ao final da search
//   3.a) a operação que escreve todos pares em um canal de saidaPares e todos
//        impares em um canal saidaImpares, e ao final avisa que acabou em um canal fin
//        func returnOddOrEven(r *Nodo, saidaP chan int, saidaI chan int, fin chan struct{}){...}
//   3.b) a versao concorrente da operação acima, ou seja, os varios nodos sao testados
//        concorrentemente se pares ou impares, escrevendo o valor no canal adequado
//
//  ABAIXO: RESPOSTAS A QUESTOES 1a e b
//  APRESENTE A SOLUÇÃO PARA AS DEMAIS QUESTÕES

package main

import (
	"fmt"
)

type Nodo struct {
	v int
	l *Nodo
	r *Nodo
}

func caminhaERD(r *Nodo) {
	if r != nil {
		caminhaERD(r.l)
		fmt.Print(r.v, ", ")
		caminhaERD(r.r)
	}
}

// -------- SOMA ----------
func soma(r *Nodo) int {
	if r != nil {
		//fmt.Print(r.v, ", ")
		return r.v + soma(r.l) + soma(r.r)
	}
	return 0
}

func somaConc(r *Nodo) int {
	s := make(chan int)
	go somaConcCh(r, s)
	return <-s
}
func somaConcCh(r *Nodo, s chan int) {
	if r != nil {
		s1 := make(chan int)
		go somaConcCh(r.l, s1)
		go somaConcCh(r.r, s1)
		s <- (r.v + <-s1 + <-s1)
	} else {
		s <- 0
	}
}

func search(r *Nodo, v int) bool {
	if r != nil {
		if r.v == v {
			return true
		}
		return search(r.l, v) || search(r.r, v)
	}
	return false
}

func concurrentSearch(r *Nodo, v int) bool {
	s := make(chan bool)
	go chConcurrentSearch(r, v, s)
	return <-s
}

func chConcurrentSearch(r *Nodo, v int, s chan bool) {
	if r != nil {
		if r.v == v {
			s <- true
		} else {
			s1 := make(chan bool)
			go chConcurrentSearch(r.l, v, s1)
			go chConcurrentSearch(r.r, v, s1)
			s <- (<-s1 || <-s1)
		}
	}
	s <- false
}

func returnOddOrEven(r *Nodo, outputEven chan int, outputOdd chan int, done chan struct{}) {
	if r != nil {
		fillOutput(r, outputEven, outputOdd)
	}
	done <- struct{}{}
}

func fillOutput(r *Nodo, outputEven chan int, outputOdd chan int) {
	if r != nil {
		fillOutput(r.l, outputEven, outputOdd)
		if r.v%2 == 0 {
			outputEven <- r.v
		} else {
			outputOdd <- r.v
		}
		fillOutput(r.r, outputEven, outputOdd)
	}
	if r == nil {
		return
	}

}

func returnConcurrentOddOrEven(r *Nodo, outputEven chan int, outputOdd chan int, done chan struct{}) {
	finished := make(chan bool)

	go func() {
		fillOutputConcurrent(r, outputEven, outputOdd, finished)
		finished <- true
	}()
	<-finished
	done <- struct{}{}
}

func fillOutputConcurrent(r *Nodo, outputEven chan int, outputOdd chan int, finished chan bool) {
	if r != nil {
		finishedLeft := make(chan bool)
		finishedRight := make(chan bool)

		go func() {
			fillOutputConcurrent(r.l, outputEven, outputOdd, finishedLeft)
			finishedLeft <- true
		}()
		go func() {
			fillOutputConcurrent(r.r, outputEven, outputOdd, finishedRight)
			finishedRight <- true
		}()

		if r.v%2 == 0 {
			outputEven <- r.v
		} else {
			outputOdd <- r.v
		}
		<-finishedLeft
		<-finishedRight
	}
	finished <- true
}

func main() {
	root := &Nodo{v: 10,
		l: &Nodo{v: 5,
			l: &Nodo{v: 3,
				l: &Nodo{v: 1, l: nil, r: nil},
				r: &Nodo{v: 4, l: nil, r: nil}},
			r: &Nodo{v: 7,
				l: &Nodo{v: 6, l: nil, r: nil},
				r: &Nodo{v: 8, l: nil, r: nil}}},
		r: &Nodo{v: 15,
			l: &Nodo{v: 13,
				l: &Nodo{v: 12, l: nil, r: nil},
				r: &Nodo{v: 14, l: nil, r: nil}},
			r: &Nodo{v: 18,
				l: &Nodo{v: 17, l: nil, r: nil},
				r: &Nodo{v: 19, l: nil, r: nil}}}}

	fmt.Println()
	fmt.Print("Valores na árvore: ")
	caminhaERD(root)
	fmt.Println()
	fmt.Println()

	fmt.Println("Soma: ", soma(root))
	fmt.Println("SomaConc: ", somaConc(root))
	fmt.Println()

	fmt.Println("Busca 13: ", search(root, 13))
	fmt.Println("BuscaConc 13: ", concurrentSearch(root, 13))
	fmt.Println("Busca 20: ", search(root, 20))
	fmt.Println("BuscaConc 20: ", concurrentSearch(root, 20))
	fmt.Println()

	fmt.Println("Pares e Ímpares")
	saidaP := make(chan int)
	saidaI := make(chan int)
	fin := make(chan struct{})
	// go returnOddOrEven(root, saidaP, saidaI, fin)
	go returnConcurrentOddOrEven(root, saidaP, saidaI, fin)
	for {
		select {
		case p := <-saidaP:
			fmt.Println("Par: ", p)
		case i := <-saidaI:
			fmt.Println("Ímpar: ", i)
		case <-fin:
			fmt.Println("Fim")
			return
		}
	}
}
