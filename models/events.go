package models

import (
	"klubRanks/db"
	"time"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	CreatedAt   time.Time `binding:"required"`
	UserID      int64
}

func (e *Event) Save() error {

	query := `
	INSERT INTO events(name, description, location, createdAt, user_id)
	VALUES ( ?, ?, ?, ?, ?)`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.CreatedAt, e.UserID)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	e.ID = id
	return err
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.CreatedAt, &event.UserID)

		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func GetEventByID(id int64) (*Event, error) {
	query := "SELECT * FROM events WHERE id = ?"
	row := db.DB.QueryRow(query, id)

	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.CreatedAt, &event.UserID)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (event Event) UpdateEvent() error {
	query := `
	UPDATE events
	SET name = ?, description = ?, location = ?, createdAt = ?
	WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.Name, event.Description, event.Location, event.CreatedAt, event.ID)
	return err
}

func (event Event) DeleteEvent() error {
	query := "DELETE FROM EVENTS where id = ?"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.ID)
	return err
}

func (e Event) Register(userId int64) error {
	query := "INSERT INTO registrations(event_id, user_id) VALUES (?,?)"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)

	return err
}

func (e Event) CancelRegistration(userId int64) error {
	query := "DELETE FROM registrations WHERE event_id = ? AND user_id = ?"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)
	return err
}

func GetAllRegistrations() (any, error) {
	query := "SELECT * FROM registrations"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []any

	for rows.Next() {
		var registrationId, eventId, userId int64
		rows.Scan(&registrationId, &eventId, &userId)
		event, _ := GetEventByID(eventId)
		data = append(data, map[string]any{
			"registration_id": registrationId,
			"user_id":         userId,
			"event":           event,
		})
	}
	return data, nil
}
