// Dev Challenge: SOLUÇÃO (só para o professor).
//
// 1. Troque api pelo endereço do servidor do professor e dupla pelo nome da dupla.
// 2. Escreva nivel1, nivel2 e nivel3 (lá embaixo).
// 3. go run .   (pode rodar quantas vezes quiser: nível certo não perde ponto)
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const api = "http://localhost:8080" // troque pelo endereço do professor
const dupla = "dupla-azul"          // troque pelo nome da sua dupla

var client = &http.Client{Timeout: 5 * time.Second}

type Pedido struct {
	Cliente       string `json:"cliente"`
	ValorCentavos int    `json:"valor_centavos"`
	Pago          bool   `json:"pago"`
}

func main() {
	token := cadastrar(dupla)
	pedidos := buscarPedidos(token)

	responder(token, 1, nivel1(pedidos))
	responder(token, 2, nivel2(pedidos))
	responder(token, 3, nivel3(pedidos))
}

// ====================== é aqui que a dupla trabalha ======================

// Nível 1 · O fiado: quanto falta receber? Soma dos pedidos NÃO pagos, em centavos.
func nivel1(pedidos []Pedido) int {
	fiado := 0
	for _, p := range pedidos {
		if !p.Pago {
			fiado += p.ValorCentavos
		}
	}
	return fiado
}

// Nível 2 · Clientes: quantos clientes diferentes aparecem nos pedidos?
func nivel2(pedidos []Pedido) int {
	vistos := map[string]bool{}
	for _, p := range pedidos {
		vistos[p.Cliente] = true
	}
	return len(vistos)
}

// Nível 3 · Melhor cliente: quem mais gastou, contando só os pedidos pagos?
// Empate: vence o nome que vem primeiro em ordem alfabética.
func nivel3(pedidos []Pedido) string {
	gasto := map[string]int{}
	for _, p := range pedidos {
		if p.Pago {
			gasto[p.Cliente] += p.ValorCentavos
		}
	}
	melhor, maior := "", -1
	for cliente, total := range gasto {
		if total > maior || (total == maior && cliente < melhor) {
			melhor, maior = cliente, total
		}
	}
	return melhor
}

// ============ daqui pra baixo é o HTTP que vimos na aula (não precisa mexer) ============

// POST /equipes: cadastra a dupla e devolve o token. Cadastrar de novo devolve o mesmo token.
func cadastrar(nome string) string {
	var r struct {
		Token string `json:"token"`
	}
	resp := enviar("POST", "/equipes", "", map[string]string{"nome": nome})
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Fatal("cadastro deu ", resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Fatal(err)
	}
	return r.Token
}

// GET /desafio: os pedidos da sua dupla.
func buscarPedidos(token string) []Pedido {
	var dados struct {
		Pedidos []Pedido `json:"pedidos"`
	}
	resp := enviar("GET", "/desafio", token, nil)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal("GET /desafio deu ", resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(&dados); err != nil {
		log.Fatal(err)
	}
	return dados.Pedidos
}

// POST /desafio/responder: manda a resposta de um nível e mostra o que o servidor disse.
func responder(token string, nivel int, resposta any) {
	corpo := map[string]any{"nivel": nivel, "resposta": resposta}
	resp := enviar("POST", "/desafio/responder", token, corpo)
	defer resp.Body.Close()
	var r struct {
		Mensagem string `json:"mensagem"`
		Erro     string `json:"erro"`
	}
	json.NewDecoder(resp.Body).Decode(&r)
	fmt.Printf("nível %d · enviei %v · %s %s\n", nivel, resposta, resp.Status, r.Mensagem+r.Erro)
}

// monta e envia o pedido HTTP (com JSON no corpo e o token no header, quando tiver)
func enviar(metodo, rota, token string, corpo any) *http.Response {
	var b []byte
	if corpo != nil {
		var err error
		b, err = json.Marshal(corpo)
		if err != nil {
			log.Fatal(err)
		}
	}
	req, err := http.NewRequest(metodo, api+rota, bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}
