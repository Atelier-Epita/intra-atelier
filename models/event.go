package models

import (
	"errors"
	"time"

	"github.com/Atelier-Epita/intra-atelier/db"

	"go.uber.org/zap"
)

const (
	getEventsQuery = `
		SELECT
			id, title, description, start_date, end_date, owner, image
		FROM events;
	`

	getUpcomingEventsQuery = `
		SELECT
			id, title, description, start_date, end_date, owner, image
		FROM events
		WHERE start_date > ?;
	`

	getCurrentEventsQuery = `
		SELECT
			id, title, description, start_date, end_date, owner, image
		FROM events
		WHERE start_date < ? AND end_date > ?;
	`

	getPastEventsQuery = `
		SELECT
			id, title, description, start_date, end_date, owner, image
		FROM events
		WHERE end_date < ?;
	`

	insertEventQuery = `
		INSERT INTO events
			(title, description, start_date, end_date, owner, image)
		VALUES
			(:title, :description, :start_date, :end_date, :owner, :image);
	`
)

type Event struct {
	Id          uint64    `db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Start_date  time.Time `json:"start_date" db:"start_date"`
	End_date    time.Time `json:"end_date" db:"end_date"`
	OwnerID     uint64    `json:"owner" db:"owner"`
	ImageID     uint64    `json:"image" db:"image"`
	// Subscribed  []User    `json:"subscribed"`
}

// Refactor this

func GetEvents() ([]*Event, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer Commit(tx, err)

	var events []*Event
	err = tx.Select(&events, getEventsQuery)
	return events, err
}

func GetUpcomingEvents() ([]*Event, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer Commit(tx, err)

	var events []*Event
	err = tx.Select(&events, getUpcomingEventsQuery, time.Now())
	return events, err
}

func GetCurrentEvents() ([]*Event, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer Commit(tx, err)

	var events []*Event
	err = tx.Select(&events, getCurrentEventsQuery, time.Now(), time.Now())
	return events, err
}

func GetPastEvents() ([]*Event, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer Commit(tx, err)

	var events []*Event
	err = tx.Select(&events, getPastEventsQuery, time.Now())
	return events, err
}

func (e *Event) Insert() error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	file, err := GetFileById(e.Id)
	if err != nil {
		return err
	}
	if file.OwnerID != e.OwnerID {
		return errors.New("File doesn't belong to owner")
	}

	_, err = tx.NamedExec(insertEventQuery, e)
	if err != nil {
		return err
	}

	zap.S().Info("Event ", e.Title, " just created.")
	return nil
}
