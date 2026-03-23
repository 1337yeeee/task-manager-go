package myerrors

import "errors"

type InvalidCredentialsError error
type CouldNotCreateTokenError error
type IdentityNotFoundInContextError error
type EntityNotFoundError error
type InvalidTaskStatusError error
type EntityAlreadyExistsError error

func InvalidCredentials() InvalidCredentialsError {
	return InvalidCredentialsError(errors.New("invalid credentials"))
}

func CouldNotCreateToken() CouldNotCreateTokenError {
	return CouldNotCreateTokenError(errors.New("could not create token"))
}

func EntityNotFound(entity string) EntityNotFoundError {
	return EntityNotFoundError(errors.New(entity + " not found"))
}

func InvalidTaskStatus() InvalidTaskStatusError {
	return InvalidTaskStatusError(errors.New("invalid task status"))
}

func EntityAlreadyExists(entity string) EntityAlreadyExistsError {
	return EntityAlreadyExistsError(errors.New(entity + " already exists"))
}
