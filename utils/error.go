package utils

import "errors"

var (
	ErrNotFound = errors.New("data not found")
	ErrInternal = errors.New("internal server error")
	//ErrBadRequest       = errors.New("bad request")
	//ErrUnauthorized     = errors.New("unauthorized access")
	//ErrNPMAlreadyExists = errors.New("NPM already exists")
	//ErrKamar            = errors.New("Kamar Penuh")
)
