package myerrors

import "errors"

func MissingAuthorizationHeader() error {
	return errors.New("missing authorization header")
}

func InvalidAuthorizationHeader() error {
	return errors.New("invalid authorization header format")
}
