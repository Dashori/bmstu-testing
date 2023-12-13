package servicesImplementation

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/repoErrors"
	"backend/internal/pkg/errors/servicesErrors"
	"backend/internal/pkg/hasher"
	"backend/internal/repository"
	"backend/internal/services"
	"github.com/charmbracelet/log"
)

type ClientServiceImplementation struct {
	ClientRepository repository.ClientRepository
	hasher           hasher.Hasher
	logger           *log.Logger
}

func NewClientServiceImplementation(
	ClientRepository repository.ClientRepository,
	hasher hasher.Hasher,
	logger *log.Logger,
) services.ClientService {

	return &ClientServiceImplementation{
		ClientRepository: ClientRepository,
		hasher:           hasher,
		logger:           logger,
	}
}

func (c *ClientServiceImplementation) SetRole() error {
	err := c.ClientRepository.SetRole()

	return err
}

func (c *ClientServiceImplementation) GetClientByLogin(login string) (*models.Client, error) {
	client, err := c.ClientRepository.GetClientByLogin(login)

	if err != nil {
		c.logger.Warn("CLIENT! Error in repository GetClientByLogin", "login", login, "error", err)
		return nil, err
	}

	c.logger.Debug("CLIENT! Successfully GetClientByLogin", "login", login)
	return client, nil
}

func (c *ClientServiceImplementation) Create(client *models.Client, password string) (*models.Client, error) {
	c.logger.Debug("CLIENT! Start create client with", "login", client.Login)

	_, err := c.ClientRepository.GetClientByLogin(client.Login)

	if err != nil && err != repoErrors.EntityDoesNotExists {
		c.logger.Warn("CLIENT! Error in repository GetClientByLogin", "login", client.Login, "error", err)
		return nil, err
	} else if err == nil {
		c.logger.Warn("CLIENT! Client already exists", "login", client.Login)
		return nil, serviceErrors.ClientAlreadyExists
	}

	passwordHash, err := c.hasher.GetHash(password)
	if err != nil {
		c.logger.Warn("CLIENT! Error get hash for password", "login", client.Login)
		return nil, serviceErrors.ErrorHash
	}
	client.Password = string(passwordHash)

	err = c.ClientRepository.Create(client)
	if err != nil {
		c.logger.Warn("CLIENT! Error in repository Create", "login", client.Login, "error", err)
		return nil, err
	}

	newClient, err := c.GetClientByLogin(client.Login)
	if err != nil {
		c.logger.Warn("CLIENT! Error in repository method GetClientByLogin", "login", client.Login, "error", err)
		return nil, err
	}

	c.logger.Info("CLIENT! Successfully create client", "login", newClient.Login, "id", newClient.ClientId)
	return newClient, nil
}

func (c *ClientServiceImplementation) CreateOTP(client *models.Client) (*models.Client, error) {
	c.logger.Debug("CLIENT! Start create client with", "login", client.Login)

	_, err := c.ClientRepository.GetClientByLogin(client.Login)

	if err != nil && err != repoErrors.EntityDoesNotExists {
		c.logger.Warn("CLIENT! Error in repository GetClientByLogin", "login", client.Login, "error", err)
		return nil, err
	} else if err == nil {
		c.logger.Warn("CLIENT! Client already exists", "login", client.Login)
		return nil, serviceErrors.ClientAlreadyExists
	}

	if client.Email == "" {
		c.logger.Debug("CLIENT! No email")
		return nil, serviceErrors.ErrorNoEmail
	}

	if client.OTP == "" {
		c.logger.Debug("CLIENT! No otp token")

		otpHash, err := c.hasher.GetHash(client.Email)
		if err != nil {
			c.logger.Warn("CLIENT! Error get hash for email", client.Email)
			return nil, serviceErrors.ErrorHash
		}

		err = sendEmail(client.Email, string(otpHash))
		if err != nil {
			c.logger.Warn("CLIENT! Error send email ", err)
			return nil, serviceErrors.ErrorSendEmail
		}

		return nil, serviceErrors.ErrorNoOTP
	}

	if !c.hasher.CheckUnhashedValue(client.OTP, client.Email) { //== false
		c.logger.Warn("CLIENT! Error otp!")
		return nil, serviceErrors.ErrorBadOTP
	}

	passwordHash, err := c.hasher.GetHash(client.Password)
	if err != nil {
		c.logger.Warn("CLIENT! Error get hash for password", "login", client.Login)
		return nil, serviceErrors.ErrorHash
	}
	client.Password = string(passwordHash)

	err = c.ClientRepository.Create(client)
	if err != nil {
		c.logger.Warn("CLIENT! Error in repository Create", "login", client.Login, "error", err)
		return nil, err
	}

	newClient, err := c.GetClientByLogin(client.Login)
	if err != nil {
		c.logger.Warn("CLIENT! Error in repository method GetClientByLogin", "login", client.Login, "error", err)
		return nil, err
	}

	c.logger.Info("CLIENT! Successfully create client", "login", newClient.Login, "id", newClient.ClientId)
	return newClient, nil
}

func (c *ClientServiceImplementation) Login(login, password string) (*models.Client, error) {
	c.logger.Debug("CLIENT! Start login with", "login", login)
	tempClient, err := c.ClientRepository.GetClientByLogin(login)

	if err != nil && err == repoErrors.EntityDoesNotExists {
		c.logger.Warn("CLIENT! Error, client with this login does not exists", "login", login, "error", err)
		return nil, serviceErrors.ClientDoesNotExists
	} else if err != nil {
		c.logger.Warn("CLIENT! Error in repository method GetClientByLogin", "login", login, "error", err)
		return nil, err
	}

	if !c.hasher.CheckUnhashedValue(tempClient.Password, password) { //== false
		c.logger.Warn("CLIENT! Error client password", "login", login)
		return nil, serviceErrors.InvalidPassword
	}

	c.logger.Info("CLIENT! Success login with", "login", login, "id", tempClient.ClientId)
	return tempClient, nil
}

func (c *ClientServiceImplementation) GetClientById(id uint64) (*models.Client, error) {
	client, err := c.ClientRepository.GetClientById(id)

	if err != nil {
		c.logger.Warn("CLIENT! Error in repository method GetClientById", "id", id, "error", err)
		return nil, err
	}

	c.logger.Debug("CLIENT! Success repository method GetClientById", "id", id)

	return client, nil
}
