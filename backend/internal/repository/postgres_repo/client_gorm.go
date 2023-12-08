package postgres_repo

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/dbErrors"
	"backend/internal/repository"
	"fmt"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type ClientPostgresGorm struct {
	Login    string `db:"login"`
	Password string `db:"password"`
}

type ClientPostgresRepositoryGorm struct {
	db *gorm.DB
}

func NewClientPostgresRepositoryGorm(db *gorm.DB) repository.ClientRepositoryGorm {
	return &ClientPostgresRepositoryGorm{db: db}
}

func (c *ClientPostgresRepositoryGorm) Create(client *models.Client) error {

	clientDB := &ClientPostgresGorm{
		Login:    client.Login,
		Password: client.Password,
	}

	res := c.db.Table("clients").Create(clientDB)
	if res.Error != nil {
		return fmt.Errorf("insert: %w", res.Error)
	}

	return nil
}

func (c *ClientPostgresRepositoryGorm) GetClientByLogin(login string) (*models.Client, error) {
	clientDB := &ClientPostgres{}

	res := c.db.Table("clients").Where("login = ?", login).Take(&clientDB)
	if res.Error != nil {
		return nil, dbErrors.ErrorSelect
	}

	clientModels := &models.Client{}
	err := copier.Copy(clientModels, clientDB)
	if err != nil {
		return nil, dbErrors.ErrorCopy
	}

	return clientModels, nil
}

// -------- gorm --------

type PostgresRepositoryFieldsGorm struct {
	DB *gorm.DB
}

func CreateClientPostgresRepositoryGorm(fields *PostgresRepositoryFieldsGorm) repository.ClientRepositoryGorm {

	return NewClientPostgresRepositoryGorm(fields.DB)
}
