package flags

import (
	"database/sql"
	"fmt"
	"github.com/charmbracelet/log"
	_ "github.com/jackc/pgx/stdlib"
	"go.nhat.io/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

type PostgresFlags struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Port     string `mapstructure:"port"`
	DBName   string `mapstructure:"dbname"`
}

func (p *PostgresFlags) InitDB(logger *log.Logger) (*sql.DB, error) {
	logger.Debug("POSTGRES! Start init postgreSQL", "user", p.User, "DBName", p.DBName,
		"host", p.Host, "port", p.Port)

	dsnPGConn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		p.User, p.DBName, p.Password,
		p.Host, p.Port)

	driverName, err := otelsql.Register("pgx",
		otelsql.AllowRoot(),
		otelsql.TraceQueryWithoutArgs(),
		otelsql.TraceRowsClose(),
		otelsql.TraceRowsAffected(),
		otelsql.WithDatabaseName("postgres"),                         // Optional.
		otelsql.WithSystem(semconv.ServiceNameKey.String("backend")), // Optional.
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	db, err := sql.Open(driverName, dsnPGConn)
	if err != nil {
		logger.Fatal("POSTGRES! Error in method open")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.Fatal("POSTGRES! Error in method ping")
		return nil, err
	}

	db.SetMaxOpenConns(10)

	if err := otelsql.RecordStats(db); err != nil {
		return nil, err
	}

	logger.Info("POSTGRES! Successfully init postgreSQL")
	return db, nil
}
