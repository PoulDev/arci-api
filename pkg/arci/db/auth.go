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
	
	var memberID int
	err := db.QueryRow(
		"INSERT INTO \"Members\" (\"email\", \"showname\", \"is_admin\") VALUES ($1, $2, $3) RETURNING id",
		hashedEmail, showname, isAdmin,
	).Scan(&memberID)
	
	if err != nil {
		return nil, err
	}
	
	return &Member{
		ID:       memberID,
		Email:    email,
		ShowName: showname,
		IsAdmin:  isAdmin,
	}, nil
}

func Login(email string) (*Member, error) {
	hashedEmail := hash(email)
	
	var member Member
	err := db.QueryRow(
		"SELECT id, showname, is_admin FROM \"Members\" WHERE email = $1",
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
		"SELECT id, email, showname, is_admin FROM Members WHERE id = $1",
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
