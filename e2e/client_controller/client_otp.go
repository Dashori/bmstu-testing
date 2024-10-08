package client_controller

import (
	errors "consoleApp/errors"
	handlers "consoleApp/handlers"
	models "consoleApp/models"
	utils "consoleApp/utils"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func ClientTestOTP(client *http.Client) error {
	var token string
	var err error
	fmt.Println("\n\n***********************************************")
	fmt.Println("------------ START TEST WITH OTP -------------")

	token, err = createClientOTP(client)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("\n------- 1/4 Successfully create client -------")

	err = getInfo(client, token)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("\n------- 2/4 Successfully get client info -------")

	err = addPet(client, token)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("\n------- 3/4 Successfully add new pet -------")

	err = getPets(client, token)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("\n------- 4/4 Successfully get pets -------")

	fmt.Println("------------- END TEST WITH OTP ---------------")
	fmt.Println("***********************************************")

	return nil
}

func createClientOTP(client *http.Client) (string, error) {
	rand.Seed(time.Now().UnixNano())

	login := randomString(7)
	password := "12345"
	email := "dashylya@huds.su"

	newClient := models.Client{Login: login, Password: password, Email: email}

	response, err := handlers.CreateClientOTP(client, &newClient)
	if err == errors.ErrorResponseStatus {
		return "", utils.CheckErrorInBody(response)
	} else if err != nil {
		return "", err
	}

	passwordEmail := os.Getenv("PASSWORD_TO")

	otp, _ := getOTPfromEmail(newClient.Email, passwordEmail)
	substrings := strings.Split(otp, " ")
	newClient.OTP = strings.TrimRight(substrings[4], "\r\n")

	response, err = handlers.CreateClientOTP(client, &newClient)
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

func getOTPfromEmail(email string, password string) (string, error) {
	fmt.Println("1")
	c, err := client.DialTLS("huds.su:993", nil)
	if err != nil {
		return "", err
	}
	log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	fmt.Println("2")
	if err := c.Login(email, password); err != nil {
		return "", err
	}
	log.Println("Logged in")

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}

	if err := <-done; err != nil {
		return "", err
	}

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return "", err
	}

	// Get the last message
	if mbox.Messages == 0 {
		return "", fmt.Errorf("No messages")
	}
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(mbox.Messages)

	// Get the whole message body
	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 1)
	go func() {
		if err := c.Fetch(seqSet, items, messages); err != nil {
			log.Fatal(err)
		}
	}()

	msg := <-messages
	if msg == nil {
		return "", fmt.Errorf("Server didn't returned message")
		// log.Fatal("Server didn't returned message")
	}

	r := msg.GetBody(&section)
	if r == nil {
		return "", fmt.Errorf("Server didn't returned message body")
	}

	// Create a new mail reader
	mr, err := mail.CreateReader(r)
	if err != nil {
		return "", err
	}

	// Print some info about the message
	header := mr.Header
	if date, err := header.Date(); err == nil {
		log.Println("Date:", date)
	}
	if from, err := header.AddressList("From"); err == nil {
		log.Println("From:", from)
	}
	if to, err := header.AddressList("To"); err == nil {
		log.Println("To:", to)
	}
	if subject, err := header.Subject(); err == nil {
		log.Println("Subject:", subject)
	}

	// Process each message's part
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			// This is the message's text (can be plain-text or HTML)
			b, _ := ioutil.ReadAll(p.Body)
			log.Println("Got text: ", string(b))
			return string(b), nil
		case *mail.AttachmentHeader:
			// This is an attachment
			filename, _ := h.Filename()
			log.Println("Got attachment: ", filename)
		}
	}

	return "", nil
}
