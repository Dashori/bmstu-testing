package benchmark

import (
	"backend/internal/models"
	"backend/internal/repository"
	. "backend/internal/repository/postgres_repo"
	"context"
	"database/sql"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
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

func benchAddClient(repo repository.ClientRepositoryGorm, n int) func(b *testing.B) {
	// fmt.Println("!")
	return func(b *testing.B) {
		b.N = n
		for i := 0; i < b.N; i++ {
			rand.Seed(time.Now().UnixNano())
			login := randomString(7)
			// fmt.Println(login)
			err := repo.Create(&models.Client{Login: login, Password: "12345"})
			if err != nil {
				panic(err)
			}
		}
	}
}

func benchGetClient(repo repository.ClientRepositoryGorm, n int) func(b *testing.B) {
	return func(b *testing.B) {
		b.N = n
		for i := 0; i < b.N; i++ {
			_, err := repo.GetClientByLogin("Cooper1")
			if err != nil {
				panic(err)
			}
		}
	}
}

func benchAddClientSqlx(repo repository.ClientRepository, n int) func(b *testing.B) {
	fmt.Println("!")
	return func(b *testing.B) {
		b.N = n
		for i := 0; i < b.N; i++ {
			rand.Seed(time.Now().UnixNano())
			login := randomString(7)
			// fmt.Println(login)
			err := repo.Create(&models.Client{Login: login, Password: "12345"})
			if err != nil {
				panic(err)
			}
		}
	}
}

func benchGetClientSqlx(repo repository.ClientRepository, n int) func(b *testing.B) {
	return func(b *testing.B) {
		b.N = n
		for i := 0; i < b.N; i++ {
			_, err := repo.GetClientByLogin("Cooper1")
			if err != nil {
				panic(err)
			}
		}
	}
}

func ClientBench() []string {
	dbContainer, db, err := SetupTestDatabaseGorm()
	defer func(dbContainer testcontainers.Container, ctx context.Context) {
		err := dbContainer.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer, context.Background())
	// fmt.Println("statr!")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fields := PostgresRepositoryFieldsGorm{DB: db}

	// ---------- gorm! -----------

	clientRepository := CreateClientPostgresRepositoryGorm(&fields)

	addClient := benchAddClient(clientRepository, 1000)
	resultsAddUser := testing.Benchmark(addClient)
	// fmt.Println("statr5")

	getClient := benchGetClient(clientRepository, 1000)
	resultsGetUser := testing.Benchmark(getClient)

	var res []string
	res = append(res, fmt.Sprintf("gorm.AddClient -- runs %5d times\tCPU: %5d ns/op\tRAM: %5d allocs/op %5d bytes/op\n",
		resultsAddUser.N, resultsAddUser.NsPerOp(), resultsAddUser.AllocsPerOp(), resultsAddUser.AllocedBytesPerOp()))

	res = append(res, fmt.Sprintf("gorm.GetClient -- runs %5d times\tCPU: %5d ns/op\tRAM: %5d allocs/op %5d bytes/op\n",
		resultsGetUser.N, resultsGetUser.NsPerOp(), resultsGetUser.AllocsPerOp(), resultsGetUser.AllocedBytesPerOp()))

	// ---------- sqlx! -----------

	dbContainer2, db2 := SetupTestDatabaseSqlx()
	defer func(dbContainer2 testcontainers.Container, ctx context.Context) {
		err := dbContainer2.Terminate(ctx)
		if err != nil {
			return
		}
	}(dbContainer2, context.Background())

	fields2 := PostgresRepositoryFields{DB: db2}

	clientRepository2 := CreateClientPostgresRepository(&fields2)

	addClient = benchAddClientSqlx(clientRepository2, 1000)
	resultsAddUser = testing.Benchmark(addClient)
	// fmt.Println("statr5")

	getClient = benchGetClientSqlx(clientRepository2, 1000)
	resultsGetUser = testing.Benchmark(getClient)

	res = append(res, fmt.Sprintf("sqlx.AddClient -- runs %5d times\tCPU: %5d ns/op\tRAM: %5d allocs/op %5d bytes/op\n",
		resultsAddUser.N, resultsAddUser.NsPerOp(), resultsAddUser.AllocsPerOp(), resultsAddUser.AllocedBytesPerOp()))

	res = append(res, fmt.Sprintf("sqlx.GetClient -- runs %5d times\tCPU: %5d ns/op\tRAM: %5d allocs/op %5d bytes/op\n",
		resultsGetUser.N, resultsGetUser.NsPerOp(), resultsGetUser.AllocsPerOp(), resultsGetUser.AllocedBytesPerOp()))

	fmt.Println(res)
	return res
}

const (
	USER     = "dashori"
	PASSWORD = "parasha"
	DBNAME   = "postgres"
)

func SetupTestDatabaseGorm() (testcontainers.Container, *gorm.DB, error) {
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:16",
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

	host, err := dbContainer.Host(context.Background())
	// fmt.Println(host, err)
	port, _ := dbContainer.MappedPort(context.Background(), "5432")

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port.Int(), USER, PASSWORD, DBNAME)
	pureDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// fmt.Println("1, ", err)
		return nil, nil, fmt.Errorf("gorm open: %w", err)
	}

	text, err := os.ReadFile("../../db/postgreSQL/init/init.sql")
	if err != nil {
		// fmt.Println("3, ", err)
		return nil, nil, fmt.Errorf("read file: %w", err)
	}

	if err := pureDB.Exec(string(text)).Error; err != nil {
		// fmt.Println("4, ", err)
		return nil, nil, fmt.Errorf("exec: %w", err)
	}

	// fmt.Println("All is ok!")
	return dbContainer, pureDB, nil
}

func SetupTestDatabaseSqlx() (testcontainers.Container, *sql.DB) {
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:16",
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

	// text, err := os.ReadFile("db/postgreSQL/init/init.sql")
	text, err := os.ReadFile("../../db/postgreSQL/init/init.sql")
	if err != nil {
		return dbContainer, nil
	}

	if _, err := db.Exec(string(text)); err != nil {
		fmt.Println(err)
		return dbContainer, nil
	}

	return dbContainer, db
}

// [gorm.AddUser -- runs   200 times       CPU: 5091739 ns/op      RAM:    69 allocs/op  5143 bytes/op
// [gorm.AddUser -- runs   200 times       CPU: 5216902 ns/op      RAM:    68 allocs/op  4868 bytes/op
