package servicesImplementation

import (
	"backend/internal/models"
	"backend/internal/repository/postgres_repo"
	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"os"
	"testing"
)

func createPetServiceFieldsPostgres2(dbTest *sqlx.DB) *petServiceFieldsPostgres {
	fields := new(petServiceFieldsPostgres)

	petRepo := postgres_repo.CreatePetMockPostgresRepository(dbTest)
	fields.petRepository = &petRepo

	fields.logger = log.New(os.Stderr)
	fields.logger.SetLevel(log.FatalLevel)

	return fields
}

func TestPetServiceImplementationCreatePostgres2(t *testing.T) {
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
			fields := createPetServiceFieldsPostgres2(db)

			pets := fields.petRepository

			err := (*pets).Create(&models.Pet{Name: "Havrosha", Type: "cat", Age: 1, Health: 10, ClientId: 1})

			tt.CheckOutput(t, err)
		})
	}
}
