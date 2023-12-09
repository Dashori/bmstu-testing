package handlers

import (
	"bytes"
	errors "consoleApp/errors"
	models "consoleApp/models"
	"fmt"
	"net/http"
)

const port = "8080"
const adress = "docker"

func DoRequest(client *http.Client, request *http.Request) (*http.Response, error) {
	response, err := client.Do(request)
	if err != nil {
		return nil, errors.ErrorExecuteRequest
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return response, errors.ErrorResponseStatus
	}

	return response, nil
}

func LoginClient(client *http.Client, newClient *models.Client) (*http.Response, error) {

	url := "http://" + adress + ":" + port + "/api/client/login"
	params := fmt.Sprintf("{\"Login\": \"%s\", \"Password\": \"%s\"}", newClient.Login, newClient.Password)
	var jsonStr = []byte(params)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, errors.ErrorNewRequest
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := DoRequest(client, request)

	return response, err
}

func CreateClient(client *http.Client, newClient *models.Client) (*http.Response, error) {

	url := "http://" + adress + ":" + port + "/api/client/create"
	params := fmt.Sprintf("{\"Login\": \"%s\", \"Password\": \"%s\"}", newClient.Login, newClient.Password)
	var jsonStr = []byte(params)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, errors.ErrorNewRequest
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := DoRequest(client, request)

	return response, err
}

func CreateClientOTP(client *http.Client, newClient *models.Client) (*http.Response, error) {

	url := "http://" + adress + ":" + port + "/api/client/createOTP"
	params := fmt.Sprintf("{\"Login\": \"%s\", \"Password\": \"%s\", \"Email\": \"%s\", \"OTP\": \"%s\"}", newClient.Login, newClient.Password,
		newClient.Email, newClient.OTP)

	var jsonStr = []byte(params)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, errors.ErrorNewRequest
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := DoRequest(client, request)

	return response, err
}

func GetClientPets(client *http.Client, token string) (*http.Response, error) {
	url := "http://" + adress + ":" + port + "/api/client/pets"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.ErrorNewRequest
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := DoRequest(client, request)

	return response, err
}

func GetClientInfo(client *http.Client, token string) (*http.Response, error) {
	url := "http://" + adress + ":" + port + "/api/client/info"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.ErrorNewRequest
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := DoRequest(client, request)

	return response, err
}

func AddPet(client *http.Client, token string, pet models.Pet) (*http.Response, error) {

	url := "http://" + adress + ":" + port + "/api/client/pet"
	params := fmt.Sprintf("{\"Name\": \"%s\", \"Type\": \"%s\", \"Age\": %d,\"Health\": %d}",
		pet.Name, pet.Type, pet.Age, pet.Health)
	var jsonStr = []byte(params)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, errors.ErrorNewRequest
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := DoRequest(client, request)

	return response, err
}
