package models

import "github.com/google/uuid"

type Category struct {
	ID   uuid.UUID
	Name string
}

type CategoryBody struct {
	Name string
}
