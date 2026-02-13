package models

import (
	"apiproject/db"
	"apiproject/utils"
	"database/sql"
	"errors"
	"fmt"

	sqlite3 "github.com/mattn/go-sqlite3"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailTaken         = errors.New("email already in use")
)

// User stores account credentials.
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (u *User) Save() error {
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	query := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
	result, err := db.DB.Exec(query, u.Username, u.Email, hashedPassword)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code == sqlite3.ErrConstraint {
			return ErrEmailTaken
		}
		return fmt.Errorf("insert user: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("read inserted user id: %w", err)
	}

	u.ID = userID
	return nil
}

func (u *User) ValidateCredentials() error {
	query := `SELECT id, password FROM users WHERE email = ?`

	var storedHashedPassword string
	if err := db.DB.QueryRow(query, u.Email).Scan(&u.ID, &storedHashedPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidCredentials
		}
		return fmt.Errorf("lookup credentials: %w", err)
	}

	if !utils.CheckPasswordHash(u.Password, storedHashedPassword) {
		return ErrInvalidCredentials
	}

	return nil
}
