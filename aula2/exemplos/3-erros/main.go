// Slide "Dois jeitos de dar errado".
// Rode com o servidor ligado e depois desligado para ver os dois casos.
package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

const api = "http://localhost:8080"

var client = &http.Client{Timeout: 5 * time.Second}

func buscar(url string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return nil
}

func main() {
	fmt.Println("rota certa: ", buscar(api+"/exemplo"))
	fmt.Println("rota errada:", buscar(api+"/nao-existe"))
}
