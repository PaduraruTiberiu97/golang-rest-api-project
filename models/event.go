package models

import (
	"apiproject/db"
	"fmt"
	"time"
)

type Event struct {
	Id          int64     `json:"id"`
	Name        string    `binding:"required" json:"name"`
	Description string    `binding:"required" json:"description"`
	Location    string    `binding:"required" json:"location"`
	Date        time.Time `binding:"required" json:"date"`
	UserId      int       `json:"user_id"`
}

var events []Event = []Event{}

func (e Event) Save() error {
	query := `INSERT INTO events (name, description, location, date, user_id) VALUES (?, ?, ?, ?, ?)`

	statement, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println("Error preparing statement: " + err.Error())
		return err
	}

	defer statement.Close()

	result, err := statement.Exec(e.Name, e.Description, e.Location, e.Date, e.UserId)
	if err != nil {
		fmt.Println("Error executing statement: " + err.Error())
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error getting last insert id: " + err.Error())
		return err
	}
	e.Id = id

	return nil
}

func GetAllEvents() ([]Event, error) {
	query := `SELECT * FROM events`

	rows, err := db.DB.Query(query)
	if err != nil {
		fmt.Println("Error querying events: " + err.Error())
		return nil, err
	}

	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.Id, &event.Name, &event.Description, &event.Location, &event.Date, &event.UserId)
		if err != nil {
			fmt.Println("Error scanning event with ID " + fmt.Sprint(event.Id) + ": " + err.Error())
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func GetEventById(id int64) (*Event, error) {
	query := `SELECT * FROM events WHERE id = ?`

	var event Event

	err := db.DB.QueryRow(query, id).Scan(&event.Id, &event.Name, &event.Description, &event.Location, &event.Date, &event.UserId)
	if err != nil {
		fmt.Println("Error querying event by ID: " + err.Error())
		return nil, err
	}
	return &event, nil
}

func (event Event) Update() error {
	query := `UPDATE events SET name = ?, description = ?, location = ?, date = ? WHERE id = ?`
	statement, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println("Error preparing update statement: " + err.Error())
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(event.Name, event.Description, event.Location, event.Date, event.Id)
	if err != nil {
		fmt.Println("Error executing update statement: " + err.Error())
		return err
	}
	return nil
}

func (event Event) Delete() error {
	query := `DELETE FROM events WHERE id = ?`

	statement, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println("Error preparing delete statement: " + err.Error())
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(event.Id)
	if err != nil {
		fmt.Println("Error executing delete statement: " + err.Error())
		return err
	}

	return nil
}
