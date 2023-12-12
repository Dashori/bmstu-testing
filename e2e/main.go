package main

import (
	"consoleApp/client_controller"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	client := &http.Client{}

	err := client_controller.ClientTest(client)
	if err != nil {
		fmt.Println(err)
	}

	defer func() {
		if r := recover(); r != nil {
			err, _ := json.Marshal(r)
			_ = os.WriteFile("trace.json", err, os.ModePerm)
		}
	}()

	req, err := http.NewRequest(http.MethodGet, "http://localhost:16686/api/traces?service=backend", nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic("incorrect status " + resp.Status)
	}
	jsonBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("trace.json", jsonBody, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
