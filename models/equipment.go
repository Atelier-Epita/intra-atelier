package models

import (
	"intra/db"

	"go.uber.org/zap"
)

const (
	getEquipmentsQuery = `
		SELECT
			id, name
		FROM equipments;
	`

	insertEquipmentQuery = `
		INSERT INTO equipments 
			(id, name)
		VALUES
			(:id, :name);
	`
)

type Equipment struct {
	Id
	Name string `json:"name" db:"name"`
	// Files []File `json:"files"`
}

func GetEquipments() ([]*Equipment, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer Commit(tx, err)

	var equipments []*Equipment
	err = tx.Select(&equipments, getEquipmentsQuery)
	return equipments, err
}

func (e *Equipment) Insert() error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(insertEquipmentQuery, e)
	if err != nil {
		return err
	}

	zap.S().Info("Equipment ", e.Name, " just created.")
	return nil
}
