package db

import (
	"fmt"
	"strings"
	"time"
)

type Event struct {
	ID          int
	Name        string
	Description string
	CreatedAt   time.Time
}

type RoleEvent struct {
	ID  int `json:"id"`
	Max int `json:"max"`
}

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func AddEvent(name string, description string, date time.Time, roles []RoleEvent) error {
	var eventID int

	err := db.QueryRow("INSERT INTO events (name, description, data) VALUES ($1, $2, $3) RETURNING id", name, description, date).Scan(&eventID)
	if err != nil {
		return err
	}

	valueStrings := make([]string, len(roles))
	valueArgs := make([]interface{}, 0, len(roles)*3)

	for i, r := range roles {
		valueStrings[i] = fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3)
		valueArgs = append(valueArgs, eventID, r.ID, r.Max)
	}

	query := fmt.Sprintf(
		"INSERT INTO event_roles (id_evento, id_ruolo, max) VALUES %s",
		strings.Join(valueStrings, ", "),
	)

	_, err = db.Exec(query, valueArgs...)
	return err
}

func GetEvents() ([]Event, error) {
	rows, err := db.Query("SELECT id, name, description, created_at FROM events")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err = rows.Scan(&event.ID, &event.Name, &event.Description, &event.CreatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func GetRoles() ([]Role, error) {
	rows, err := db.Query("SELECT id, nome FROM Roles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles = make([]Role, 0)
	for rows.Next() {
		var role Role
		err = rows.Scan(&role.ID, &role.Name)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}
