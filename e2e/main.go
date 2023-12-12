package main

import (
	"consoleApp/client_controller"
	"fmt"
	"net/http"
)

func main() {
	client := &http.Client{}

	err := client_controller.ClientTest(client)
	if err != nil {
		fmt.Println(err)
	}
}
