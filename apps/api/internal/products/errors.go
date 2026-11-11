package products

import "errors"

var ErrProductNotFound = errors.New("PRODUCT_NOT_FOUND")
var ErrProductUniqueViolation = errors.New("UNIQUE_CONSTRAINT_VIOLATION")
