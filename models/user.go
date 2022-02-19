package models

import (
	"intra/db"

	"go.uber.org/zap"
)

type User struct {
	base

	Login string `json:"login" db:"login"`
	FirstName string `json:"firstname" db:"first_name"`
	Name string `json:"name" db:"last_name"`
	Promotion uint16 `json:"promotion" db:"promotion"`
}

func GetUsers() ([]*User, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer Commit(tx, err)

	var users []*User
	err = tx.Select(&users, "SELECT login, first_name, last_name, promotion FROM user")

	zap.S().Info("Retrieved all users")
	return users, err
}

func Insert(login string, name string, firstname string, promotion uint16) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	defer Commit(tx, err)

	_, err = tx.Exec("INSERT INTO user (login, first_name, last_name, promotion) VALUES (?, ?, ?, ?)",
		login,
		name,
		firstname,
		promotion,
	)
	if err != nil {
		return err
	}
	zap.S().Info("Created user ", firstname, " ", name)

	return nil
}

