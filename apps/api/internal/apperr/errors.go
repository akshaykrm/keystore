package apperr

import "errors"

var ErrNotFound = errors.New("RESOURCE_NOT_FOUND")
var ErrInvalidRequest = errors.New("INVALID_REQUEST")
var ErrAccessDenied = errors.New("ACCESS_DENIED")
var ErrUniqueViolation = errors.New("RESOURCE_CONFLICT")
var ErrInternalServerError = errors.New("INTERNAL_SERVER_ERROR")
