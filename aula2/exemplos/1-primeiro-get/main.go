// Slide "Primeiro GET".
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const api = "http://localhost:8080" // ou o endereço que o servidor do professor mostrou

func main() {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(api + "/exemplo")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
}
