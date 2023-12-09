package client_controller

import (
	errors "consoleApp/errors"
	handlers "consoleApp/handlers"
	models "consoleApp/models"
	utils "consoleApp/utils"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func ClientTest(client *http.Client) error {
	var token string
	var err error

	fmt.Println("\n\n***********************************************")
	fmt.Println("------------ START TEST WITHOUT OTP -------------\n")

	token, err = createClient(client)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("\n---------- 1/4 Successfully create client ----------")

	err = getInfo(client, token)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("\n---------- 2/4 Successfully get client info ----------")

	err = addPet(client, token)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("\n---------- 3/4 Successfully add new pet ----------")

	err = getPets(client, token)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("\n---------- 4/4 Successfully get pets ----------")

	fmt.Println("\n\n------------- END TEST WITHOUT OTP --------------")
	fmt.Println("***********************************************\n\n")

	return nil
}

func createClient(client *http.Client) (string, error) {
	rand.Seed(time.Now().UnixNano())

	login := randomString(7)
	password := "12345"
	newClient := models.Client{Login: login, Password: password}

	response, err := handlers.CreateClient(client, &newClient)
	if err == errors.ErrorResponseStatus {
		return "", utils.CheckErrorInBody(response)
	} else if err != nil {
		return "", err
	}

	result, err := utils.ParseClientBody(response)
	if err != nil {
		return "", err
	}

	fmt.Println("New client login:", login)
	return result.Token, nil
}
