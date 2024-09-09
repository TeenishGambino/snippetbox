package models

import (
	"errors"
)

// Declaring the variables
// doing this is similar to just doing
// var one at a time
var (
	ErrNoRecord = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: ivalid credentials")
	ErrDuplicateEmail = errors.New("models : duplicate email")
)

