package models

import (
	"intra/db"

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
	Id
	Event_id uint64 `json:"eventId" db:"event_id"`
	User_id  uint64 `json:"userId" db:"user_id"`
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

func (p *Participant) insert() error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(insertParticipantQuery, p)
	if err != nil {
		return err
	}

	zap.S().Info("Participant ", p.User_id, " just created for event ", p.Event_id, ".")
	return nil
}
