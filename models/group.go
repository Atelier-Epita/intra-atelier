package models

import (
	"intra/db"

	"go.uber.org/zap"
)

const (
	getGroupsQuery = `
		SELECT
			id, name
		FROM groups;
	`

	insertGroupQuery = `
		INSERT INTO groups
			(name)
		VALUES
			(:name);
	`
)

type Group struct {
	Id
	Name string `json:"name" db:"name"`
}

func GetGroups() ([]*Group, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer Commit(tx, err)

	var groups []*Group
	err = tx.Select(&groups, getGroupsQuery)
	return groups, err
}

func (g *Group) Insert() error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(insertGroupQuery, g)
	if err != nil {
		return err
	}

	zap.S().Info("Event ", g.Name, " just created.")
	return nil
}
