package entity

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
    Id          int64     `db:"id" json:"id"`
    ServiceName string    `db:"service_name" json:"service_name"`
    Price       uint64    `db:"price" json:"price"`
    UserID      uuid.UUID `db:"user_id" json:"user_id"`
    StartDate   time.Time `db:"start_date" json:"start_date"`
	EndDate  time.Time `db:"end_date" json:"end_date"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

