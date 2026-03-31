package db

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type EventRole struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Max  int    `json:"max"`
}

type Event struct {
	ID           int         `json:"id"`
	Title        string      `json:"titolo"`
	Description  string      `json:"descrizione"`
	Date         time.Time   `json:"data"`
	Roles        []EventRole `json:"roles"`
	SelectedRole *int        `json:"selected-role"`
}

type RoleEventInput struct {
	ID  int `json:"id"`
	Max int `json:"max"`
}

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func AddRole(name string) (*Role, error) {
	result, err := db.Exec("INSERT INTO Roles (nome) VALUES (?)", name)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Role{
		ID:   int(id),
		Name: name,
	}, nil
}

func DeleteRole(id int) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM EventRoles WHERE id_ruolo = ?", id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("role is assigned to one or more events")
	}

	result, err := db.Exec("DELETE FROM Roles WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("role not found")
	}

	return nil
}

func AddEvent(name string, description string, date time.Time, roles []RoleEventInput) error {
	result, err := db.Exec(
		"INSERT INTO Events (titolo, descrizione, data) VALUES (?, ?, ?)",
		name, description, date,
	)
	if err != nil {
		return err
	}

	eventID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	valueStrings := make([]string, len(roles))
	valueArgs := make([]interface{}, 0, len(roles)*3)

	for i, r := range roles {
		valueStrings[i] = "(?, ?, ?)"
		valueArgs = append(valueArgs, eventID, r.ID, r.Max)
	}

	query := fmt.Sprintf(
		"INSERT INTO EventRoles (id_evento, nome_ruolo, max) VALUES %s",
		strings.Join(valueStrings, ", "),
	)

	_, err = db.Exec(query, valueArgs...)
	return err
}

func GetEvents(memberID int) ([]Event, error) {
	eventRows, err := db.Query(`
		SELECT
			e.id, e.titolo, e.descrizione, e.data,
			er.id, r.id, r.nome, er.max
		FROM Events e
		LEFT JOIN EventRoles er ON er.id_evento = e.id
			LEFT JOIN Roles r ON r.id = er.nome_ruolo
		ORDER BY e.id
	`)
	if err != nil {
		return nil, err
	}
	defer eventRows.Close()

	// mappa id evento -> evento
	eventMap := make(map[int]*Event)
	var eventOrder []int

	for eventRows.Next() {
		var (
			eventID     int
			title       string
			description string
			date        time.Time
			erID        *int
			roleID      *int
			roleName    *string
			roleMax     *int
		)

		err = eventRows.Scan(&eventID, &title, &description, &date, &erID, &roleID, &roleName, &roleMax)
		if err != nil {
			return nil, err
		}

		if _, exists := eventMap[eventID]; !exists {
			eventMap[eventID] = &Event{
				ID:          eventID,
				Title:       title,
				Description: description,
				Date:        date,
				Roles:       []EventRole{},
			}
			eventOrder = append(eventOrder, eventID)
		}

		if erID != nil && roleID != nil {
			eventMap[eventID].Roles = append(eventMap[eventID].Roles, EventRole{
				ID:   *roleID,
				Name: *roleName,
				Max:  *roleMax,
			})
		}
	}

	if err = eventRows.Err(); err != nil {
		return nil, err
	}

	// Ruolo selezionato dall'evento
	partRows, err := db.Query(`
		SELECT p.id_evento, r.id
		FROM Partecipation p
		JOIN Roles r ON r.nome = p.ruolo
		WHERE p.id_partecipante = ?
	`, memberID)
	if err != nil {
		return nil, err
	}
	defer partRows.Close()

	for partRows.Next() {
		var eventID, roleID int
		if err = partRows.Scan(&eventID, &roleID); err != nil {
			return nil, err
		}
		if ev, exists := eventMap[eventID]; exists {
			id := roleID
			ev.SelectedRole = &id
		}
	}

	if err = partRows.Err(); err != nil {
		return nil, err
	}

	events := make([]Event, 0, len(eventOrder))
	for _, id := range eventOrder {
		events = append(events, *eventMap[id])
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
