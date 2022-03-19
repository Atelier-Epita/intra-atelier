package models

import (
	"intra/db"

	"go.uber.org/zap"
)

const (
	getInventoriesQuery = `
		SELECT
			id, group_id, name, amount
		FROM inventory;
	`

	insertInventoryQuery = `
		INSERT INTO inventory
			(group_id, name, amount)
		VALUES
			(:group_id, :name, :amount);
	`
)

type Inventory struct {
	Id
	GroupID uint64 `json:"groupId" db:"group_id"`
	Name    string `json:"name" db:"name"`
	Amount  uint64 `json:"amount" db:"amount"`
}

func getInventories() ([]*Inventory, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer Commit(tx, err)

	var inventories []*Inventory
	err = tx.Select(&inventories, getEquipmentsQuery)
	return inventories, err
}

func (i *Inventory) insert() error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(insertInventoryQuery, i)
	if err != nil {
		return err
	}

	zap.S().Info("Event ", i.Name, " just created.")
	return nil
}
