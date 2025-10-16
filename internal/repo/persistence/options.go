package persistence

import (
	"time"

	"github.com/google/uuid"
)

type ListOption func(*ListOptions)

type ListOptions struct {
	UserID        *uuid.UUID
	ServiceName   *string
	Price         *uint64
	StartDateFrom *time.Time
	StartDateTo   *time.Time
	Limit         int
	Offset        int
	SortBy        string
	SortOrder     string
}

func WithUserID(id uuid.UUID) ListOption {
	return func(l *ListOptions) {
		l.UserID = &id
	}
}

func WithServiceName(name string) ListOption {
	return func(l *ListOptions) {
		l.ServiceName = &name
	}
}

func WithPrice(price uint64) ListOption {
	return func(l *ListOptions) {
		l.Price = &price
	}
}

func WithStartDateFrom(t time.Time) ListOption {
	return func(l *ListOptions) {
		l.StartDateFrom = &t
	}
}

func WithStartDateTo(t time.Time) ListOption {
	return func(l *ListOptions) {
		l.StartDateTo = &t
	}
}
