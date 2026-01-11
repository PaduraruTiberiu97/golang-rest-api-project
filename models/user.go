package models

import (
	"apiproject/db"
	"fmt"

	"apiproject/utils"
)

type User struct {
	Id       int64  `json:"id"`
	Email    string `binding:"required" json:"email"`
	Password string `binding:"required" json:"password"`
}

func (u User) Save() error {
	query := `INSERT INTO users (email, password) VALUES (?, ?)`
	statement, err := db.DB.Prepare(query)

	if err != nil {
		fmt.Println("Error preparing statement: " + err.Error())
		return err
	}

	defer statement.Close()

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		fmt.Println("Error hashing password: " + err.Error())
		return err
	}

	result, err := statement.Exec(u.Email, hashedPassword)
	if err != nil {
		fmt.Println("Error executing statement: " + err.Error())
		return err
	}

	userId, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error getting last insert id: " + err.Error())
		return err
	}

	u.Id = userId

	return nil
}

func (u *User) ValidateCredentials() error {
	query := `SELECT id, password FROM users WHERE email = ?`

	row := db.DB.QueryRow(query, u.Email)

	var storedHashedPassword string
	err := row.Scan(&u.Id, &storedHashedPassword)
	if err != nil {
		return err
	}

	passwordIsValid := utils.CheckPasswordHash(u.Password, storedHashedPassword)
	if !passwordIsValid {
		return fmt.Errorf("invalid credentials")
	}

	return nil
}
