package storage

import (
    "context"
    "database/sql"
    "auth/db"
    "time"
)

const userColumns = "id, email, password_hash, created_at, updated_at, last_login, is_active, failed_attempts, locked_until"

type User struct {
    ID             int       `json:"id"`
    Email          string    `json:"email"`
    PasswordHash   string    `json:"password_hash"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
    LastLogin      time.Time `json:"last_login,omitempty"`
    IsActive       bool      `json:"is_active"`
    FailedAttempts int       `json:"failed_attempts"`
    LockedUntil    *time.Time `json:"locked_until,omitempty"`
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

func (s *Storage) CreateUser(ctx context.Context, u User) error {
    query := `INSERT INTO users (` + userColumns + `)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
    _, err := s.db.DB.ExecContext(ctx, query, u.ID, u.Email, u.PasswordHash, u.CreatedAt, u.UpdatedAt, u.LastLogin, u.IsActive, u.FailedAttempts, u.LockedUntil)
    if err != nil {
        return err
    }

    return nil
}

func (s *Storage) GetUser(ctx context.Context, userID int) (*User, error) {
    var u User
    query := "SELECT " + userColumns + " FROM users WHERE id = $1"

    row := s.db.DB.QueryRowContext(ctx, query, userID)
    err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt, &u.LastLogin, &u.IsActive, &u.FailedAttempts, &u.LockedUntil)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }

    return &u, nil
}

func (s *Storage) UpdateUserByID(ctx context.Context, u User) error {
    query := `UPDATE users SET email = $1, password_hash = $2, updated_at = $3, last_login = $4, is_active = $5, failed_attempts = $6, locked_until = $7
              WHERE id = $8`
    _, err := s.db.DB.ExecContext(ctx, query, u.Email, u.PasswordHash, u.UpdatedAt, u.LastLogin, u.IsActive, u.FailedAttempts, u.LockedUntil, u.ID)
    if err != nil {
        return err
    }

    return nil
}

func (s *Storage) DeleteUser(ctx context.Context, userID int) error {
    query := `DELETE FROM users WHERE id = $1`
    _, err := s.db.DB.ExecContext(ctx, query, userID)
    if err != nil {
        return err
    }

    return nil
}