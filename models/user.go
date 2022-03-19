package models

import (
	"intra/db"

	"go.uber.org/zap"
)

const (
	getUsersQuery = `
		SELECT
			id, email, first_name, last_name
		FROM users;
	`

	getUserByEmailQuery = `
		SELECT
			id, email, first_name, last_name
		FROM users
		WHERE email = ?;
	`

	insertUserQuery = `
		INSERT INTO users
			(email, first_name, last_name)
		VALUES
			(:email, :first_name, :last_name);
	`
)

type User struct {
	Id
	Email     string `json:"email" db:"email"`
	FirstName string `json:"firstname" db:"first_name"`
	LastName  string `json:"lastname" db:"last_name"`
	// Groups    []Group `json:"groups" gorm:"foreignKey:ID"`
}

func GetUsers() ([]*User, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer Commit(tx, err)

	var users []*User
	err = tx.Select(&users, getUsersQuery)

	zap.S().Info("Retrieved all users")
	return users, err
}

func GetUserByEmail() (*User, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer Commit(tx, err)

	var u *User
	err = tx.Select(u, getUserByEmailQuery)
	return u, err
}

func (u *User) Insert() error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	defer Commit(tx, err)

	_, err = tx.NamedExec(insertUserQuery, u)
	if err != nil {
		return err
	}
	zap.S().Info("Created user ", u.FirstName, " ", u.LastName, " ", u.Email, ".")

	return nil
}
