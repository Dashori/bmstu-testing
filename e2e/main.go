package main

import (
	"consoleApp/client_controller"
	"fmt"
	"net/http"
	// "os"
)

func main() {
	// fmt.Println(os.Getenv("PASSWORD"))
	client := &http.Client{}

	err := client_controller.ClientTest(client)
	if err != nil {
		fmt.Println(err)
	}

	// err = client_controller.ClientTestOTP(client)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
