package postgres_repo

import (
	"backend/internal/models"
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"testing"
)

var testClientPostgresRepositoryGormCreateSuccess = []struct {
	TestName  string
	InputData struct {
		client *models.Client
	}
	CheckOutput func(t *testing.T, err error)
}{
	{
		TestName: "create success test",
		InputData: struct {
			client *models.Client
		}{&models.Client{Login: "ChicagoTest", Password: "12345"}},

		CheckOutput: func(t *testing.T, err error) {
			require.NoError(t, err)
		},
	},
}

func TestClientPostgresRepositoryGormCreate(t *testing.T) {
	dbContainer, db, err := SetupTestDatabaseGorm()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())

	if err != nil {
		return
	}

	for _, tt := range testClientPostgresRepositoryGormCreateSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			fields := PostgresRepositoryFieldsGorm{DB: db}

			clientRepository := CreateClientPostgresRepositoryGorm(&fields)

			err := clientRepository.Create(tt.InputData.client)
			tt.CheckOutput(t, err)

			client, err := clientRepository.GetClientByLogin("ChicagoTest")
			fmt.Println("!!!!", client.ClientId, client.Login, client.Password)
			tt.CheckOutput(t, err)
		})
	}
}
