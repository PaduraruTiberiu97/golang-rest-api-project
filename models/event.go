package models

import (
	"apiproject/db"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sqlite3 "github.com/mattn/go-sqlite3"
)

var ErrAlreadyRegistered = errors.New("user is already registered for this event")

// Event stores an event and its creator.
type Event struct {
	ID          int64     `json:"id"`
	Name        string    `binding:"required" json:"name"`
	Description string    `binding:"required" json:"description"`
	Location    string    `binding:"required" json:"location"`
	Date        time.Time `binding:"required" json:"date"`
	UserID      int64     `json:"user_id"`
}

func (e *Event) Save() error {
	query := `INSERT INTO events (name, description, location, date, user_id) VALUES (?, ?, ?, ?, ?)`

	result, err := db.DB.Exec(query, e.Name, e.Description, e.Location, e.Date, e.UserID)
	if err != nil {
		return fmt.Errorf("insert event: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("read inserted event id: %w", err)
	}

	e.ID = id
	return nil
}

func GetAllEvents() ([]Event, error) {
	query := `SELECT id, name, description, location, date, user_id FROM events ORDER BY date`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query events: %w", err)
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		if err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.Date, &event.UserID); err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate events: %w", err)
	}

	return events, nil
}

func GetEventByID(id int64) (*Event, error) {
	query := `SELECT id, name, description, location, date, user_id FROM events WHERE id = ?`

	var event Event
	if err := db.DB.QueryRow(query, id).Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.Date, &event.UserID); err != nil {
		return nil, fmt.Errorf("query event by id: %w", err)
	}

	return &event, nil
}

func (e *Event) Update() error {
	query := `UPDATE events SET name = ?, description = ?, location = ?, date = ? WHERE id = ?`

	result, err := db.DB.Exec(query, e.Name, e.Description, e.Location, e.Date, e.ID)
	if err != nil {
		return fmt.Errorf("update event: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read updated rows: %w", err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (e *Event) Delete() error {
	query := `DELETE FROM events WHERE id = ?`

	result, err := db.DB.Exec(query, e.ID)
	if err != nil {
		return fmt.Errorf("delete event: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read deleted rows: %w", err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (e *Event) Register(userID int64) error {
	var existing int
	checkQuery := `SELECT 1 FROM registrations WHERE event_id = ? AND user_id = ? LIMIT 1`
	if err := db.DB.QueryRow(checkQuery, e.ID, userID).Scan(&existing); err == nil {
		return ErrAlreadyRegistered
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("check existing event registration: %w", err)
	}

	query := `INSERT INTO registrations (event_id, user_id) VALUES (?, ?)`

	if _, err := db.DB.Exec(query, e.ID, userID); err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code == sqlite3.ErrConstraint {
			return ErrAlreadyRegistered
		}
		return fmt.Errorf("register for event: %w", err)
	}

	return nil
}

func (e *Event) CancelRegistration(userID int64) error {
	query := `DELETE FROM registrations WHERE event_id = ? AND user_id = ?`

	result, err := db.DB.Exec(query, e.ID, userID)
	if err != nil {
		return fmt.Errorf("cancel event registration: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read canceled registration rows: %w", err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
