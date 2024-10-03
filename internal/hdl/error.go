package hdl

import "errors"

var ErrInternal = errors.New("internal error")
var ErrDecodeRequest = errors.New("failed to decode request")
var ErrMissingSongID = errors.New("missing song ID")
var ErrMethodNotAllowed = errors.New("method not allowed")
