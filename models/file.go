package models

import (
	"github.com/Atelier-Epita/intra-atelier/db"

	"go.uber.org/zap"
)

const (
	getFilesQuery = `
		SELECT
			id, permission, owner_id, group_id, equipment_id, file_name, file_hash
		FROM files;
	`

	insertFileQuery = `
		INSERT INTO files
			(permission, owner_id, group_id, equipment_id, file_name, file_hash)
		VALUES
			(:permission, :owner_id, :group_id, :equipment_id, :file_name, :file_hash)
	`

	getFilesByEquipmentQuery = `
		SELECT
			id, permission, owner_id, group_id, equipment_id, file_name, file_hash
		FROM files
		WHERE files.equipment_id = ?;
	`
)

type File struct {
	Id          uint64 `db:"id"`
	Permission  bool   `json:"permission" db:"permission"` // 0 is public, 1 is private
	OwnerID     uint64 `json:"ownerId" db:"owner_id"`
	GroupID     uint64 `json:"groupId" db:"group_id"`
	EquipmentID uint64 `json:"equipmentId" db:"equipment_id"`

	// file path is ./files/$(file_hash)$(file_name)
	File_name string   `json:"filename"`
	File_hash [32]byte `json:"filehash"`
}

func GetFiles() ([]*File, error) {
	tx, err := db.DB.Beginx()

	if err != nil {
		return nil, err
	}
	defer Commit(tx, err)

	var files []*File
	err = tx.Select(&files, getFilesQuery)
	return files, err
}

func (f *File) Insert() error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(insertFileQuery, f)
	if err != nil {
		return err
	}

	zap.S().Info("Event ", f.File_name, " just created.")
	return nil
}

func (e *Equipment) GetFilesByEquipment() ([]*File, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer Commit(tx, err)

	var files []*File
	err = tx.Select(&files, getFilesByEquipmentQuery)
	return files, err
}
