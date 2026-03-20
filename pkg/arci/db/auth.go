package db

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

type Member struct {
	ID       int
	Email    string
	ShowName string
	IsAdmin  bool
}

func hash(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

func Register(email string, showname string, isAdmin bool) (*Member, error) {
	hashedEmail := hash(email)

	result, err := db.Exec(
		"INSERT INTO Members (email, name, is_admin) VALUES (?, ?, ?)",
		hashedEmail, showname, isAdmin,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Member{
		ID:       int(id),
		Email:    email,
		ShowName: showname,
		IsAdmin:  isAdmin,
	}, nil
}

func Login(email string) (*Member, error) {
	hashedEmail := hash(email)

	var member Member
	err := db.QueryRow(
		"SELECT id, name, is_admin FROM Members WHERE email = ?",
		hashedEmail,
	).Scan(&member.ID, &member.ShowName, &member.IsAdmin)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("email not found")
		}
		return nil, err
	}

	return &member, nil
}

func GetMemberByID(id int) (*Member, error) {
	var member Member
	var hashedEmail string

	err := db.QueryRow(
		"SELECT id, email, name, is_admin FROM Members WHERE id = ?",
		id,
	).Scan(&member.ID, &hashedEmail, &member.ShowName, &member.IsAdmin)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("member not found")
		}
		return nil, err
	}

	member.Email = ""
	return &member, nil
}
