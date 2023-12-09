package client_controller

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
