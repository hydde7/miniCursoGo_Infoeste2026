// Slides "O corpo do POST" e "Enviando o POST".
// Manda para /treino/responder só para ver o POST funcionando: o servidor responde 422 (850 não é a resposta do treino).
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const api = "http://localhost:8080"

var client = &http.Client{Timeout: 5 * time.Second}

func enviar(url, token string) error {
	corpo := map[string]any{
		"nivel":    1,
		"resposta": 850,
	}

	b, err := json.Marshal(corpo)
	if err != nil {
		return err
	}
	fmt.Println("corpo:", string(b))

	leitor := bytes.NewReader(b)
	req, err := http.NewRequest("POST", url, leitor)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("status:", resp.Status)
	return nil
}

func main() {
	if err := enviar(api+"/treino/responder", "abc123"); err != nil {
		log.Fatal(err)
	}
}
