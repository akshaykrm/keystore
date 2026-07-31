package user

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrEmailConflict = errors.New("email already exists")
