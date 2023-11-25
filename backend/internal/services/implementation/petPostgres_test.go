package servicesImplementation

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/repository/postgres_repo"
	// "backend/internal/pkg/errors/dbErrors"
	"backend/internal/services"
	// "context"
	// "database/sql"
	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/require"
	"github.com/jmoiron/sqlx"
	// "github.com/testcontainers/testcontainers-go"
	// "github.com/DATA-DOG/go-sqlmock"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"os"
	"testing"
)

type petServiceFieldsPostgres struct {
	petRepository    *repository.PetRepository
	clientRepository *repository.ClientRepository
	logger           *log.Logger
}

func createPetServiceFieldsPostgres(dbTest *sqlx.DB) *petServiceFieldsPostgres {
	fields := new(petServiceFieldsPostgres)

	petRepo := postgres_repo.CreatePetMockPostgresRepository(dbTest)
	fields.petRepository = &petRepo

	fields.logger = log.New(os.Stderr)
	fields.logger.SetLevel(log.FatalLevel)

	return fields
}

func createPetServicePostgres(fields *petServiceFieldsPostgres) services.PetService {
	return NewPetServiceImplementation(*fields.petRepository, *fields.clientRepository, fields.logger)
}

func TestPetServiceImplementationCreatePostgres(t *testing.T) {
	db, mock, _ := sqlxmock.Newx()
	defer db.Close()

	var testPetCreatePostgresSuccess = []struct {
		TestName        string
		InputData       struct{}
		Mock            func()
		Prepare         func(fields *petServiceFieldsPostgres)
		CheckOutput     func(t *testing.T, err error)
		CheckOutputHelp func(t *testing.T, err error)
	}{
		{
			TestName:  "pet create success",
			InputData: struct{}{},

			Mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("insert into pets").WithArgs("Havrosha", "cat", uint64(1), uint64(10), uint64(1))
				mock.ExpectCommit()
			},

			CheckOutput: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range testPetCreatePostgresSuccess {
		tt := tt
		t.Run(tt.TestName, func(t *testing.T) {
			tt.Mock()
			fields := createPetServiceFieldsPostgres(db)

			pets := fields.petRepository

			(*pets).Create(&models.Pet{Name: "Havrosha", Type: "cat", Age: 1, Health: 10, ClientId: 1})

		})
	}
}


// func TestPetServiceImplementationCreatePostgres(t *testing.T) {
// 	db, _, err := sqlxmock.Newx()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	for _, tt := range testPetCreatePostgresSuccess {
// 		tt := tt
// 		t.Run(tt.TestName, func(t *testing.T) {
// 			tt.Mock()
// 			fields := createPetServiceFieldsPostgres(db)

// 			pets := fields.petRepository

// 			err = (*pets).Create(&models.Pet{Name: "Havrosha", Type: "cat", Age: 1, Health: 10, ClientId: 1})
// 			tt.CheckOutput(t, err)
// 		})
// 	}

// }
