// Slide "Lendo o JSON da resposta" (continua o primeiro GET).
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const api = "http://localhost:8080"

type Mensagem struct {
	Texto string `json:"texto"`
}

func main() {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(api + "/exemplo")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var m Mensagem
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(m.Texto)
}
