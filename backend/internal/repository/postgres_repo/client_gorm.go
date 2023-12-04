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

	// fmt.Println("------ ok 0 -------")

	res := c.db.Table("clients").Create(clientDB)
	if res.Error != nil {
		return fmt.Errorf("insert: %w", res.Error)
	}

	// fmt.Println("------ ok 1 -------")

	return nil
}

func (c *ClientPostgresRepositoryGorm) GetClientByLogin(login string) (*models.Client, error) {
	clientDB := &ClientPostgres{}

	res := c.db.Table("clients").Where("login = ?", login).Take(&clientDB)
	// fmt.Println("------ ok 0 -------")
	if res.Error != nil {
		return nil, dbErrors.ErrorSelect
	}

	clientModels := &models.Client{}
	err := copier.Copy(clientModels, clientDB)
	// fmt.Println("------ ok 1 -------")

	if err != nil {
		return nil, dbErrors.ErrorCopy
	}

	// fmt.Println("------ ok 2 -------")
	fmt.Println(clientModels.Login, clientModels.Password, clientModels.ClientId)

	return clientModels, nil
}

// -------- gorm --------

type PostgresRepositoryFieldsGorm struct {
	DB *gorm.DB
}

func CreateClientPostgresRepositoryGorm(fields *PostgresRepositoryFieldsGorm) repository.ClientRepositoryGorm {

	return NewClientPostgresRepositoryGorm(fields.DB)
}
