package ctrl

import "errors"

var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")
var ErrBadExtReq = errors.New("bad external request")
var ExtSrvErr = errors.New("external service error")
var ErrExtUnreachable = errors.New("unreachable")
