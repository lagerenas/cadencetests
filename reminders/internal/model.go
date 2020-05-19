package internal

import (
	"context"
	"time"
)

type Event struct {
	ID          string
	Start       time.Time
	End         time.Time
	Cancelled   bool
	Description string
}

type ReminderSender interface {
	CreateReminder(ctx context.Context, m Event) error
	UpdateReminder(ctx context.Context, m Event) error
	CancelReminder(ctx context.Context, m Event) error
}
