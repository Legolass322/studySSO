package db

import (
	"database/sql"
	"fmt"
	// _ "github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"log/slog"
)

type Database struct {
	db  *sql.DB
	log *slog.Logger
}

type DBCfg struct {
	DriverName string
	Username   string
	Password   string
	Host       string
	Port       int
	DBName     string
}

func (dbcfg *DBCfg) GetConnStr() string {
	return fmt.Sprintf("%v://%v:%v@%v:%v/%v?sslmode=disable",
		dbcfg.DriverName,
		dbcfg.Username,
		dbcfg.Password,
		dbcfg.Host,
		dbcfg.Port,
		dbcfg.DBName,
	)
}

func MustNew(logger *slog.Logger, dbcfg *DBCfg) *Database {
	db, err := New(logger, dbcfg)
	if err != nil {
		panic("Failed to connect to db")
	}

	return db
}

func New(logger *slog.Logger, dbcfg *DBCfg) (*Database, error) {
	const op = "db.New"

	log := logger.With("op", op)

	log.Info("Connecting to db")

	db, err := sql.Open(dbcfg.DriverName, dbcfg.GetConnStr())
	if err != nil {
		log.Error("Cannot connect to db", slog.String("err", err.Error()))
		return &Database{}, err
	}
	if err :=db.Ping(); err != nil {
		log.Error("Cannot connect to db", slog.String("err", err.Error()))
		return &Database{}, err
	}

	log.Info("Connection success")

	return &Database{db, logger}, nil
}

func (s *Database) Stop() error {
	const op = "db.Stop"

	log := s.log.With(slog.String("op", op))

	log.Info("Closing...")

	if err := s.db.Close(); err != nil {
		log.Error("Failed to close connection")
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Closed")

	return nil
}
