package models

import (
	"errors"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID       uuid.UUID
	Name     string
	Password string
}

type UserAuth struct {
	Name     string
	Password string
}
