package models

import (
	"errors"

	"github.com/google/uuid"
)

var ErrCategoryNotFound = errors.New("Category not found")

type Category struct {
	ID   uuid.UUID
	Name string
}

type CategoryBody struct {
	Name string
}
