package functions

import (
	database "server/createDb"
	"golang.org/x/crypto/bcrypt"
)

func AuthenticateUser(username, password string) (int, error) {
	
	var (
		userID       int
		passwordHash string
	)

	err := database.DB.QueryRow(
		"SELECT id, password FROM users WHERE username = ?",
		username,
	).Scan(&userID, &passwordHash)

	if err != nil {
		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(password),
	); err != nil {
		return 0, err
	}

	return userID, nil
}
