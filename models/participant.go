package models

import (
	"github.com/Atelier-Epita/intra-atelier/db"

	"go.uber.org/zap"
)

const (
	getParticipantsQuery = `
		SELECT
			id, event_id, user_id
		FROM participant;
	`

	insertParticipantQuery = `
		INSERT INTO participants
			(event_id, user_id)
		VALUES
			(:event_id, :user_id);
	`
)

type Participant struct {
	Id      uint64 `db:"id"`
	EventID uint64 `json:"eventId" db:"event_id"`
	UserID  uint64 `json:"userId" db:"user_id"`
}

func getParticipants() ([]*Participant, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer Commit(tx, err)

	var participants []*Participant
	err = tx.Select(&participants, getParticipantsQuery)
	return participants, err
}

func (p *Participant) Insert() error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(insertParticipantQuery, p)
	if err != nil {
		return err
	}

	zap.S().Info("Participant ", p.UserID, " just created for event ", p.EventID, ".")
	return nil
}
