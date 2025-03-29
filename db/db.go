package db

import (
	"auth/config"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var (
	ErrAlreadyInTX      = errors.New("storage already running in a tx")
	ErrNoTXProvided     = errors.New("no tx provided")
	ErrDBNoTInitiated   = errors.New("db not initiated")
	ErrDBNoRowsEffected = errors.New("db no rows effected")
	ErrMustBeInTx       = errors.New("must be in tx")
)

type DB struct {
	DB *sql.DB
}

func NewDB() (*DB, error) {
	config.LoadConfig()
	C := config.GetConfig().Database
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		C.User, C.Password, C.Host, C.Port, C.DbName,
	)
	fmt.Printf("Connecting to database: %s\n", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Successfully connected to database")

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(30 * time.Minute)

	return &DB{DB: db}, nil
}

func (d *DB) Close() {
	d.DB.Close()
}

func (d *DB) Ping(c context.Context) error {
	return d.DB.PingContext(c)
}
