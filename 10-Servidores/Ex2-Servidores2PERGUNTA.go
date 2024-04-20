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
	NumClients = 100
	PoolSize   = 10
)

type Request struct {
	value      int
	returnChan chan int
}

func client(id int, req chan Request, canal chan struct{}) {
	var value, response int
	myChannel := make(chan int)
	for {
		value = rand.Intn(1000)
		req <- Request{value, myChannel}
		response = <-myChannel
		fmt.Println("Cliente: ", id, " Requisição: ", value, " Resposta:", response, "Processos: ", len(canal))
		<-canal
		time.Sleep(60 * time.Second)
	}
}

func requestTreatment(id int, req Request) {
	fmt.Println("                                               Tratando Requisição ", id)
	req.returnChan <- req.value * 2
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}

func concurrentServer(input chan Request, semaphore chan struct{}) {
	var j int = 0
	for {
		j++
		request := <-input
		semaphore <- struct{}{}
		go requestTreatment(j, request)
	}
}

// ------------------------------------
// Main
func main() {
	fmt.Println("------ Servidores - criacao dinamica -------")
	servChan := make(chan Request, PoolSize) // CANAL POR ONDE SERVIDOR RECEBE PEDIDOS
	sem := make(chan struct{}, PoolSize)     // Use um canal com tamanho 10
	go concurrentServer(servChan, sem)       // LANÇA PROCESSO SERVIDOR
	for i := 0; i < NumClients; i++ {
		go client(i, servChan, sem)
	}
	<-make(chan int)
}
