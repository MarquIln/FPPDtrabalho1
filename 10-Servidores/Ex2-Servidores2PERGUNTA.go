// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// servidor com criacao dinamica de thread de servico
// Problema:
//   considere um servidor que recebe pedidos por um canal (representando uma conexao)
//   ao receber o pedido, sabe-se através de qual canal (conexao) responder ao cliente.
//   Abaixo uma solucao sequencial para o servidor.
// Exercicio
//   deseja-se tratar os clientes concorrentemente, e nao sequencialmente.
//   como ficaria a solucao ?
// Veja abaixo a resposta ...
//   quantos clientes podem estar sendo tratados concorrentemente ?
//
// Exercicio:
//   agora suponha que o seu servidor pode estar tratando no maximo 10 clientes concorrentemente.
//   como voce faria ?
//

package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	NCL  = 10
	Pool = 10
)

type Request struct {
	v      int
	ch_ret chan int
}

// cliente
func cliente(i int, req chan Request) {
	var v, r int
	my_ch := make(chan int)
	for {
		v = rand.Intn(1000)
		req <- Request{v, my_ch}
		r = <-my_ch
		fmt.Println("cli:", i, " req:", v, "  resp:", r)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500))) // simulando um tempo de resposta variável
	}
}

// servidor
func trataReq(id int, req Request) {
	fmt.Println("                                 trataReq", id)
	req.ch_ret <- req.v * 2
}

// servidor que dispara threads de servico
func servidorConc(in chan Request) {
	pool := make(chan struct{}, Pool) // canal para controlar o número de clientes concorrentes
	for {
		req := <-in
		pool <- struct{}{} // adicionar marca ao canal para ocupar um espaço no pool
		go func(req Request) {
			defer func() { <-pool }() // remover a marca do canal quando a goroutine terminar
			trataReq(rand.Int(), req)
		}(req)
	}
}

func main() {
	fmt.Println("------ Servidores - criacao dinamica -------")
	serv_chan := make(chan Request) // canal para os pedidos
	go servidorConc(serv_chan)
	for i := 0; i < NCL; i++ {
		go cliente(i, serv_chan)
	}

	// Mantém o programa em execução indefinidamente
	select {}
}
