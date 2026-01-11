package models

import (
	"apiproject/db"
	"fmt"
)

type User struct {
	Id       int64  `json:"id"`
	Email    string `binding:"required" json:"email"`
	Password string `binding:"required" json:"password"`
}

func (u User) Save() error {
	query := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
	statement, err := db.DB.Prepare(query)

	if err != nil {
		fmt.Println("Error preparing statement: " + err.Error())
		return err
	}

	defer statement.Close()

	result, err := statement.Exec(u.Email, u.Email, u.Password)
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
