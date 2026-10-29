// Treino guiado: somar os pedidos pagos. É o código que você monta junto com a turma.
// buscar e enviar são o código dos slides de GET, JSON e POST dentro de funções.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const api = "http://localhost:8080" // ou o endereço que o servidor do professor mostrou

var client = &http.Client{Timeout: 5 * time.Second}

type Pedido struct {
	Cliente       string `json:"cliente"`
	ValorCentavos int    `json:"valor_centavos"`
	Pago          bool   `json:"pago"`
}

type Dados struct {
	Pedidos []Pedido `json:"pedidos"`
}

// GET + Decode
func buscar(url string, destino any) {
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal("GET deu ", resp.Status)
	}
	err = json.NewDecoder(resp.Body).Decode(destino)
	if err != nil {
		log.Fatal(err)
	}
}

// Marshal + POST
func enviar(url string, corpo any) {
	b, err := json.Marshal(corpo)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var r struct {
		Mensagem string `json:"mensagem"`
		Erro     string `json:"erro"`
	}
	json.NewDecoder(resp.Body).Decode(&r)
	fmt.Println(resp.Status, r.Mensagem+r.Erro)
}

func faturamento(pedidos []Pedido) int {
	total := 0
	for _, p := range pedidos {
		if p.Pago {
			total += p.ValorCentavos
		}
	}
	return total
}

func main() {
	var dados Dados
	buscar(api+"/treino", &dados)

	total := faturamento(dados.Pedidos)
	fmt.Println("total:", total)

	corpo := map[string]any{"resposta": total}
	enviar(api+"/treino/responder", corpo)
}
