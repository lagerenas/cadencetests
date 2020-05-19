package internal

import (
	"errors"
	"strconv"
	"time"
)

type EventDB struct {
	events map[string]Event
	nextID int
}

func NewEventDB() *EventDB {
	return &EventDB{
		events: map[string]Event{},
		nextID: 1,
	}
}

func (edb *EventDB) GetEvent(id string) (Event, error) {
	e, ok := edb.events[id]
	if !ok {
		return Event{}, errors.New("Event not found")
	}
	return e, nil
}

func (edb *EventDB) AddEvent(start time.Time, end time.Time, description string) {
	e := Event{
		ID:          strconv.Itoa(edb.nextID),
		Start:       start,
		End:         end,
		Cancelled:   false,
		Description: description,
	}
	edb.events[e.ID] = e

	edb.nextID += 1
}
