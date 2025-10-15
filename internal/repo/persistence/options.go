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
