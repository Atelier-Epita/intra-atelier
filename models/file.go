package models

import (
	"intra/db"

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
	Id
	Permission   uint   `json:"permission" db:"permission"` // 0 is public, 1 is private
	Owner_id     uint64 `json:"ownerId" db:"owner_id"`
	Group_id     uint64 `json:"groupId" db:"group_id"`
	Equipment_id uint64 `json:"equipmentId" db:"equipment_id"`

	// file path is ./files/$(file_hash)$(file_name)
	File_name string   `json:"filename" db:"file_name"`
	File_hash [32]byte `json:"filehash" db:"file_hash"`
}

func getFiles() ([]*File, error) {
	tx, err := db.DB.Beginx()

	if err != nil {
		return nil, err
	}
	defer Commit(tx, err)

	var files []*File
	err = tx.Select(&files, getFilesQuery)
	return files, err
}

func (f *File) insert() error {
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

func (e *Equipment) getFilesByEquipment() ([]*File, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer Commit(tx, err)

	var files []*File
	err = tx.Select(&files, getFilesByEquipmentQuery)
	return files, err
}
