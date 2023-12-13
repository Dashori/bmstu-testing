package client_controller

import (
	handlers "consoleApp/handlers"
	models "consoleApp/models"
	utils "consoleApp/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"os"
	"strings"
	"testing"
)

func Test2FA(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E 2FA test")
}

var _ = Describe("2FA", func() {
	client := &http.Client{}

	login := randomString(7)
	password := "12345"
	email := "dashylya@huds.su"

	newClient := models.Client{Login: login, Password: password, Email: email}

	Context("First request with email and without otp", func() {
		_, err := handlers.CreateClientOTP(client, &newClient)
		It("Check response code", func() {
			Expect(err).To(BeZero())
		})
	})

	Context("Get mail from smtp server", func() {
		passwordEmail := os.Getenv("PASSWORD_TO")
		otp, err := getOTPfromEmail(newClient.Email, passwordEmail)

		It("Check response code from getting mail", func() {
			Expect(err).To(BeZero())
		})

		substrings := strings.Split(otp, " ")
		newClient.OTP = strings.TrimRight(substrings[4], "\r\n")
	})

	Context("Second request with email and otp", func() {
		response, err := handlers.CreateClientOTP(client, &newClient)

		It("Check response code", func() {
			Expect(err).To(BeZero())
		})

		result, err := utils.ParseClientBody(response)
		It("Parse client body", func() {
			Expect(err).To(BeZero())
		})

		It("Check client body", func() {
			Expect(result.Login).To(Equal(newClient.Login))
		})

		newClient.Token = result.Token
	})

	Context("Get client info with new token", func() {
		err := getInfo(client, newClient.Token)

		It("Check response code", func() {
			Expect(err).To(BeZero())
		})
	})
})
