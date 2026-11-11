package apperr

import "errors"

var ErrNotFound = errors.New("RESOURCE_NOT_FOUND")
var ErrUniqueViolation = errors.New("RESOURCE_CONFLICT")
var ErrInternalServerError = "INTERNAL_SERVER_ERROR"
