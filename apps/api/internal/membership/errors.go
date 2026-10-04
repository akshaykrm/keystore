package membership

import "errors"

var ErrMemberShipConflict = errors.New("membership already exists")
var ErrInvalidRole = errors.New("invalid membership role")
