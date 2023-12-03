package menu_controller

import (
	errors "consoleApp/errors"
	handlers "consoleApp/handlers"
	models "consoleApp/models"
	utils "consoleApp/utils"
	view "consoleApp/view"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

func clientMenu(client *http.Client) error {

	var token string
	var err error

	token, err = createClient(client)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("------- 1/4 Successfully create client -------")

	err = getInfo(client, token)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("------- 2/4 Successfully get client info -------")

	err = addPet(client, token)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("------- 3/4 Successfully add new pet -------")

	err = getPets(client, token)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("------- 4/4 Successfully get pets -------")

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

func getPets(client *http.Client, token string) error {
	response, err := handlers.GetClientPets(client, token)
	if err == errors.ErrorResponseStatus {
		return utils.CheckErrorInBody(response)
	} else if err != nil {
		return err
	}

	pets, err := utils.ParsePetsBody(response)
	if err != nil {
		return err
	}

	view.PrintPets(pets)

	return nil
}

func getInfo(client *http.Client, token string) error {
	response, err := handlers.GetClientInfo(client, token)
	if err == errors.ErrorResponseStatus {
		return utils.CheckErrorInBody(response)
	} else if err != nil {
		return err
	}

	result, err := utils.ParseClientBody(response)
	if err != nil {
		return err
	}

	view.PrintClientInfo(result)

	return err
}

func addPet(client *http.Client, token string) error {

	pet := models.Pet{Name: randomString(7), Type: "cat", Age: 2, Health: 8}

	response, err := handlers.AddPet(client, token, pet)
	if err == errors.ErrorResponseStatus {
		return utils.CheckErrorInBody(response)
	} else if err != nil {
		return err
	}

	fmt.Println("\nNew pet: ", pet.Name, pet.Type, pet.Age, pet.Health)

	return nil
}
