package products

import "errors"

var ErrProductNotFound = errors.New("products not found")
var ErrProductUniqueViolation = errors.New("unique constraint violation")
