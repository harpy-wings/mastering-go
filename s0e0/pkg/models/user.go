package models

import "time"

type User struct {
	UUID         string    `json:"uuid"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	RegisteredAt time.Time `json:"registered_at"`
}
