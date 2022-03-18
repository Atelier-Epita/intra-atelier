package models

import (
	"intra/db"
	"time"

	"go.uber.org/zap"
)

const (
	getEventsQuery = `
		SELECT
			id, title, description, start_date, end_date, owner, image
		FROM events;
	`

	insertEventQuery = `
		INSERT INTO events
			(title, description, start_date, end_date, owner, image)
		VALUES
			(:title, :description, :start_date, :end_date, :owner, :image);
	`
)

type Event struct {
	Id
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Start_date  time.Time `json:"start_date" db:"start_date"`
	End_date    time.Time `json:"end_date" db:"end_date"`
	Owner       uint64    `json:"owner" db:"owner"`
	Image       uint64    `json:"image" db:"image"`
}

func getEvents() ([]*Event, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer Commit(tx, err)

	var events []*Event
	err = tx.Select(&events, getEventsQuery)
	return events, err
}

func (e *Event) insert() error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(insertEventQuery, e)
	if err != nil {
		return err
	}

	zap.S().Info("Event ", e.Title, " just created.")
	return nil
}