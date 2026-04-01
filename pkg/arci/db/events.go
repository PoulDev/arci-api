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

type Volunteer struct {
	Name string `json:"name"`
}

type Event struct {
	ID           int                    `json:"id"`
	Title        string                 `json:"titolo"`
	Description  string                 `json:"descrizione"`
	Date         time.Time              `json:"data"`
	Roles        []EventRole            `json:"roles"`
	SelectedRole *int                   `json:"selected-role"`
	Volunteers   map[string][]Volunteer `json:"volunteers"`
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
		var roleName string
		err := db.QueryRow("SELECT nome FROM Roles WHERE id = ?", r.ID).Scan(&roleName)
		if err != nil {
			return fmt.Errorf("Role %d not found: %w", r.ID, err)
		}
		valueStrings[i] = "(?, ?, ?)"
		valueArgs = append(valueArgs, eventID, roleName, r.Max)
	}

	query := fmt.Sprintf(
		"INSERT INTO EventRoles (id_evento, nome_ruolo, max) VALUES %s",
		strings.Join(valueStrings, ", "),
	)
	_, err = db.Exec(query, valueArgs...)
	return err
}

func GetEvents(memberID int, lastID *int) ([]Event, error) {
	var (
		idQuery string
		idArgs  []interface{}
	)
	if lastID != nil {
		idQuery = `SELECT id FROM Events WHERE id < ? ORDER BY data DESC, id DESC LIMIT 10`
		idArgs = []interface{}{*lastID}
	} else {
		idQuery = `SELECT id FROM Events ORDER BY data DESC, id DESC LIMIT 10`
	}

	idRows, err := db.Query(idQuery, idArgs...)
	if err != nil {
		return nil, err
	}
	defer idRows.Close()

	var eventOrder []int
	for idRows.Next() {
		var id int
		if err = idRows.Scan(&id); err != nil {
			return nil, err
		}
		eventOrder = append(eventOrder, id)
	}
	if err = idRows.Err(); err != nil {
		return nil, err
	}
	if len(eventOrder) == 0 {
		return []Event{}, nil
	}

	placeholders := strings.Repeat("?,", len(eventOrder))
	placeholders = placeholders[:len(placeholders)-1]
	pageArgs := make([]interface{}, len(eventOrder))
	for i, id := range eventOrder {
		pageArgs[i] = id
	}

	eventRows, err := db.Query(fmt.Sprintf(`
		SELECT
			e.id, e.titolo, e.descrizione, e.data,
			er.id, r.id, r.nome, er.max
		FROM Events e
		LEFT JOIN EventRoles er ON er.id_evento = e.id
		LEFT JOIN Roles r      ON r.nome = er.nome_ruolo
		WHERE e.id IN (%s)
		ORDER BY e.data DESC, e.id DESC
	`, placeholders), pageArgs...)
	if err != nil {
		return nil, err
	}
	defer eventRows.Close()

	eventMap := make(map[int]*Event, len(eventOrder))
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
		if err = eventRows.Scan(&eventID, &title, &description, &date, &erID, &roleID, &roleName, &roleMax); err != nil {
			return nil, err
		}
		if _, exists := eventMap[eventID]; !exists {
			eventMap[eventID] = &Event{
				ID:          eventID,
				Title:       title,
				Description: description,
				Date:        date,
				Roles:       []EventRole{},
				Volunteers:  map[string][]Volunteer{},
			}
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

	volRows, err := db.Query(fmt.Sprintf(`
		SELECT p.id_evento, p.ruolo, m.name, p.id_partecipante
		FROM Partecipation p
		JOIN Members m ON m.id = p.id_partecipante
		WHERE p.id_evento IN (%s)
	`, placeholders), pageArgs...)
	if err != nil {
		return nil, err
	}
	defer volRows.Close()

	for volRows.Next() {
		var (
			eventID       int
			roleName      string
			memberName    string
			participantID int
		)
		if err = volRows.Scan(&eventID, &roleName, &memberName, &participantID); err != nil {
			return nil, err
		}
		ev, exists := eventMap[eventID]
		if !exists {
			continue
		}

		ev.Volunteers[roleName] = append(ev.Volunteers[roleName], Volunteer{Name: memberName})

		if participantID == memberID {
			for _, r := range ev.Roles {
				if r.Name == roleName {
					id := r.ID
					ev.SelectedRole = &id
					break
				}
			}
		}
	}
	if err = volRows.Err(); err != nil {
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
