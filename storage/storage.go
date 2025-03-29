package storage

import (
	"auth/db"
	"context"
	"database/sql"
	"fmt"
	"log"

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
	query := `
        INSERT INTO users (email, password_hash)
        VALUES ($1, $2)
    `
	log.Printf("Executing query: %s with email: %s", query, u.Email)

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return err
	}

	_, err = stmt.ExecContext(ctx, u.Email, u.PasswordHash)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return err
	}

	log.Printf("Successfully created user with email: %s", u.Email)
	return nil
}

func (s *Storage) GetUser(ctx context.Context, e string) (*User, error) {
	var u User
	query := "SELECT email, password_hash FROM users WHERE email = $1"
	log.Printf("Executing GetUser query: %s with email: %s", query, e)

	row := s.db.DB.QueryRowContext(ctx, query, e)
	err := row.Scan(&u.Email, &u.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user found with email: %s", e)
			return nil, nil
		}
		log.Printf("Error scanning user: %v", err)
		return nil, err
	}

	log.Printf("Successfully retrieved user with email: %s", e)
	return &u, nil
}

func (s *Storage) UpdateUser(ctx context.Context, u User) error {
	query := `UPDATE users SET password_hash = $2 WHERE email = $1`
	log.Printf("Executing UpdateUser query: %s with email: %s", query, u.Email)
	log.Printf("Password hash: %s", u.PasswordHash)

	result, err := s.db.DB.ExecContext(ctx, query, u.Email, u.PasswordHash)
	if err != nil {
		log.Printf("Error executing update: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("No user found with email: %s", u.Email)
		return fmt.Errorf("no user found with email: %s", u.Email)
	}

	log.Printf("Successfully updated user with email: %s", u.Email)
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
		log.Printf("Error getting user: %v", err)
		return err
	}
	if u == nil {
		log.Printf("No user found with email: %s", e)
		return fmt.Errorf("no user found with email: %s", e)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(p))
	if err != nil {
		log.Printf("Invalid password for user: %s", e)
		return err
	}

	log.Printf("Password check successful for user: %s", e)
	return nil
}
