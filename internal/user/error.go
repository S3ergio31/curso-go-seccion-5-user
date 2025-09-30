package user

import (
	"errors"
	"fmt"
)

var ErrorFirstNameRequired = errors.New("first name is required")

var ErrorLastNameRequired = errors.New("last name is required")

type ErrorUserNotFound struct {
	UserID string
}

func (e ErrorUserNotFound) Error() string {
	return fmt.Sprintf("user '%s' does not found", e.UserID)
}
