package errors

import "errors"

var (
	NotLoggedMessage   = errors.New("You are currently not logged in to the system, to get proper access create a new user or login with an existing user.")
	LoginErrorMessage  = errors.New("Incorrect login or password")
	LogoutErrorMessage = errors.New("Some problems with logout")
	UnknownMessage     = errors.New("Oops, error occurred: Please provide bug to github: https://github.com/lastbackend/lastbackend/issues")
)
