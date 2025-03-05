package storage

import (
	"auth/db"
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

const userColumns = "email, password_hash"

type User struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type Storage struct {
	db *db.DB
}

func (s *Storage) Close() error {
	if s.db != nil {
		s.db.Close()
	}
	return nil
}

func NewStorage() (*Storage, error) {
	dbConn, err := db.NewDB()
	if err != nil {
		return nil, err
	}
	return &Storage{db: dbConn}, nil
}

func (s *Storage) SignUp(ctx context.Context, u User) error {
	var db = s.db.DB
	stmt, err := db.Prepare(`
        INSERT INTO users (email, password_hash, created_at, updated_at)
        VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    `)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, u.Email, u.PasswordHash)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetUser(ctx context.Context, e string) (*User, error) {
	var u User
	query := "SELECT " + userColumns + " FROM users WHERE mail = $1"

	row := s.db.DB.QueryRowContext(ctx, query, e)
	err := row.Scan(&u.Email, &u.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (s *Storage) UpdateUser(ctx context.Context, u User) error {
	query := `UPDATE users SET email = $1, password_hash = $2
              WHERE email = $1`
	_, err := s.db.DB.ExecContext(ctx, query, u.Email, u.PasswordHash)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteUser(ctx context.Context, e string) error {
	query := `DELETE FROM users WHERE email = $1`
	_, err := s.db.DB.ExecContext(ctx, query, e)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CheckPassword(ctx context.Context, e, p string) error {
	u, err := s.GetUser(ctx, e)
	if err != nil {
		return err
	}
	if u == nil {
		return nil
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(p))
	if err != nil {
		return err
	}
	return nil
}
