package models

import (
	"github.com/Atelier-Epita/intra-atelier/db"

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

	getUserGroupsQuery = `
		SELECT
  			groups.id,
  			groups.name
		FROM
  		groups
  		JOIN users_group ON groups.id = users_group.group_id
  		JOIN users ON users.id = users_group.user_id
		WHERE
  			users.id = ?;
	`

	getUserFilesQuery = `
		SELECT
			id, permission, owner_id, group_id, equipment_id, file_name, file_hash
		FROM files
		JOIN users ON users.id = files.owner_id
		WHERE users.id = ?;
	`

	insertUserQuery = `
		INSERT INTO users
			(email, first_name, last_name)
		VALUES
			(:email, :first_name, :last_name);
	`

	AddGroupQuery = `
		INSERT INTO users_group
			(user_id, group_id)
		VALUES
			(:userID, :groupID);
	`
)

type User struct {
	Id        uint64 `json:"-" db:"id"`
	Email     string `json:"email" db:"email"`
	FirstName string `json:"firstname" db:"first_name"`
	LastName  string `json:"lastname" db:"last_name"`
	// In database: Groups    []Group
}

type UserGroup struct {
	Id      uint64 `db:"id"`
	UserId  uint64 `db:"userID"`
	GroupId uint64 `db:"groupID"`
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

func GetUserByMail(email string) (*User, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer Commit(tx, err)

	var user User
	err = tx.Get(&user, getUserByEmailQuery, email)
	return &user, err
}

func (u *User) GetGroups() ([]Group, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer Commit(tx, err)

	var groups []Group
	err = tx.Select(&groups, getUserGroupsQuery, u.Id)
	return groups, err
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

func (u *User) AddGroup(group *Group) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	defer Commit(tx, err)

	_, err = tx.NamedExec(AddGroupQuery, UserGroup{UserId: u.Id, GroupId: group.Id})
	if err != nil {
		return err
	}
	zap.S().Info("Binded group ", group.Name, " to user ", u.Email, ".")

	return nil
}

func GetUserFiles(email string) ([]*File, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer Commit(tx, err)

	var files []*File
	err = tx.Select(&files, getUserFilesQuery)
	return files, err
}
