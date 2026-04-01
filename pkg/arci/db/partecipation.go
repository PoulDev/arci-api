package db

import "errors"

func Partecipate(eventID int, memberID int, roleName string) error {
	// controllo esistenza ruolo
	var roleExists bool
	err := db.QueryRow("SELECT COUNT(*) > 0 FROM Roles WHERE nome = ?", roleName).Scan(&roleExists)
	if err != nil {
		return err
	}
	if !roleExists {
		return errors.New("role not found")
	}

	// controlla massimo di partecipanti
	var maxSlots, currentCount int
	err = db.QueryRow(
		"SELECT max FROM EventRoles WHERE id_evento = ? AND nome_ruolo = ?",
		eventID, roleName,
	).Scan(&maxSlots)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return errors.New("role not available for this event")
		}
		return err
	}

	err = db.QueryRow(
		"SELECT COUNT(*) FROM Partecipation WHERE id_evento = ? AND ruolo = ?",
		eventID, roleName,
	).Scan(&currentCount)
	if err != nil {
		return err
	}
	if currentCount >= maxSlots {
		return errors.New("role is full")
	}

	// se l'utente è già iscritto con un altro ruolo, lo toglie da quello precedente
	// e lo aggiunge a quello nuovo
	var existingRole *string
	err = db.QueryRow(
		"SELECT ruolo FROM Partecipation WHERE id_evento = ? AND id_partecipante = ?",
		eventID, memberID,
	).Scan(&existingRole)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return err
	}
 
	if existingRole != nil {
		if *existingRole == roleName {
			return errors.New("already signed up for this role")
		}
		_, err = db.Exec(
			"UPDATE Partecipation SET ruolo = ? WHERE id_evento = ? AND id_partecipante = ?",
			roleName, eventID, memberID,
		)
		return err
	}
 
	_, err = db.Exec(
		"INSERT INTO Partecipation (id_evento, id_partecipante, ruolo) VALUES (?, ?, ?)",
		eventID, memberID, roleName,
	)
	return err
}

func CancelPartecipation(eventID int, memberID int, roleName string) error {
	result, err := db.Exec(
		"DELETE FROM Partecipation WHERE id_evento = ? AND id_partecipante = ? AND ruolo = ?",
		eventID, memberID, roleName,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("partecipation not found")
	}

	return nil
}
