package servicesImplementation

import (
	"backend/internal/repository"
	mock_repository "backend/internal/repository/mocks"
	"backend/internal/repository/postgres_repo"
	"backend/internal/services"
	"context"
	"database/sql"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"net/smtp"
	"os"
	"time"
)

const (
	USER     = "dashori"
	PASSWORD = "parasha"
	DBNAME   = "postgres"
)

func SetupTestDatabase() (testcontainers.Container, *sql.DB) {
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       DBNAME,
			"POSTGRES_PASSWORD": PASSWORD,
			"POSTGRES_USER":     USER,
		},
	}

	dbContainer, _ := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})

	host, _ := dbContainer.Host(context.Background())
	port, _ := dbContainer.MappedPort(context.Background(), "5432")

	dsnPGConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port.Int(), USER, PASSWORD, DBNAME)
	db, err := sql.Open("pgx", dsnPGConn)
	if err != nil {
		return dbContainer, nil
	}

	err = db.Ping()
	if err != nil {
		return dbContainer, nil
	}
	db.SetMaxOpenConns(10)

	text, err := os.ReadFile("../../../db/postgreSQL/init/init.sql")
	if err != nil {
		return dbContainer, nil
	}

	if _, err := db.Exec(string(text)); err != nil {
		return dbContainer, nil
	}

	fmt.Println(string(text))

	return dbContainer, db
}

type recordServiceFields struct {
	recordRepositoryMock *mock_repository.MockRecordRepository
	doctorRepositoryMock *mock_repository.MockDoctorRepository
	clientRepositoryMock *mock_repository.MockClientRepository
	petRepositoryMock    *mock_repository.MockPetRepository
	logger               *log.Logger
}

type RecordServiceFieldsPostgres struct {
	RecordRepository *repository.RecordRepository
	DoctorRepository *repository.DoctorRepository
	ClientRepository *repository.ClientRepository
	PetRepository    *repository.PetRepository
	logger           *log.Logger
}

func СreateRecordServiceFieldsPostgres(dbTest *sql.DB) *RecordServiceFieldsPostgres {
	fields := new(RecordServiceFieldsPostgres)

	repositoryFields := postgres_repo.PostgresRepositoryFields{DB: dbTest}

	recordRepo := postgres_repo.CreateRecordPostgresRepository(&repositoryFields)
	fields.RecordRepository = &recordRepo

	doctorRepo := postgres_repo.CreateDoctorPostgresRepository(&repositoryFields)
	fields.DoctorRepository = &doctorRepo

	clientRepo := postgres_repo.CreateClientPostgresRepository(&repositoryFields)
	fields.ClientRepository = &clientRepo

	petRepo := postgres_repo.CreatePetPostgresRepository(&repositoryFields)
	fields.PetRepository = &petRepo

	fields.logger = log.New(os.Stderr)
	fields.logger.SetLevel(log.FatalLevel)

	return fields
}

func CreateRecordServicePostgres(fields *RecordServiceFieldsPostgres) services.RecordService {
	return NewRecordServiceImplementation(*fields.RecordRepository, *fields.DoctorRepository,
		*fields.ClientRepository, *fields.PetRepository, fields.logger)
}

func sendEmail(emailTo string, otp string) error {

	from := "dashori@huds.su"
	password := os.Getenv("PASSWORD_FROM")

	smtpHost := "huds.su"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	msg := "From: " + from +
		"\r\nTo: " + emailTo +
		"\r\nDate: " + time.Now().Format("Mon, 02 Jan 2006 15:04:05 -0700") +
		"\r\nSubject: Код двухфакторной аутентификации\n" +
		"\r\n\r\nВаш код для входа: " + otp

	fmt.Println(msg)

	// Отправка письма через SMTP
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{emailTo}, []byte(msg))
	if err != nil {
		fmt.Println("Ошибка при отправке письма: ", err)
		return err
	}
	fmt.Println("Письмо успешно отправлено!")

	return nil
}
