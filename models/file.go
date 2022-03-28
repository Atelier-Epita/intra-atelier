package models

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"

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
		WHERE equipment_id = ?;
	`

	getFileByIdQuery = `
		SELECT
			id, permission, owner_id, group_id, equipment_id, file_name, file_hash
		FROM files
		WHERE id = ?;
	`
)

type File struct {
	Id          uint64 `json:"-" db:"id"`
	Permission  uint32 `json:"-" db:"permission"` // 0 is public, 1 is private
	OwnerID     uint64 `json:"owner" db:"owner_id"`
	GroupID     uint64 `json:"-" db:"group_id"` // TODO: Could be many groups
	EquipmentID uint64 `json:"-" db:"equipment_id"`

	// file path is ./files/$(file_hash)$(file_name)
	Filename string `json:"filename" db:"file_name"`
	Filehash []byte `json:"-" db:"file_hash"`
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

func CreateFile(file File, fileContent []byte) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}
	defer Commit(tx, err)

	hasher := sha256.New()
	hasher.Write(fileContent)
	sha := hasher.Sum(nil)
	file.Filehash = sha

	filename := "./files/" + hex.EncodeToString(sha)
	if _, err = os.Stat(filename); err == nil { // file already exist
		return nil
	} else if errors.Is(err, os.ErrNotExist) { // file do not exist
		_, err = tx.NamedExec(insertFileQuery, file)
		if err != nil {
			return err
		}

		err = os.WriteFile(filename, fileContent, 0444)
		if err != nil {
			return err
		}

		zap.S().Info("File ", file.Filename, " just created.")
		return nil
	}
	return err
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

func GetFileById(id uint64) (*File, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer Commit(tx, err)

	var file File
	err = tx.Get(&file, getFileByIdQuery, id)
	return &file, err
}
