package db

import "time"

type Event struct {
	ID          int
	Name        string
	Description string
	CreatedAt   time.Time
}

func AddEvent(name string, description string) error {
	_, err := db.Exec("INSERT INTO events (name, description) VALUES ($1, $2)", name, description)
	if err != nil {
		return err
	}

	return nil
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


