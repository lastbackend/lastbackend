package errors

func Message(errStatus string) string {

	var errMessage string

	if errStatus == "ACCESS_DENIED" {
		errMessage = "Access denied"
	}

	if errStatus == "USER_NOT_FOUND" {
		errMessage = "Incorrect login or password"
	}

	if errStatus == "INTERNAL_SERVER_ERROR" {
		errMessage = "oops, error occurred: Please provide bug to github: https://github.com/lastbackend/lastbackend/issues"
	}

	if errStatus == "INCORRECT_JSON" {
		errMessage = "Incorrect json"
	}

	if errStatus == "BAD_PARAMETER_USERNAME" {
		errMessage = "Bad parameter username"
	}

	if errStatus == "BAD_PARAMETER_EMAIL" {
		errMessage = "Bad parameter e-mail"
	}

	if errStatus == "BAD_PARAMETER_PASSWORD" {
		errMessage = "Bad parameter password"
	}

	if errStatus == "USERNAME_NOT_UNIQUE" {
		errMessage = "Username not unique"
	}

	if errStatus == "EMAIL_NOT_UNIQUE" {
		errMessage = "E-mail not unique"
	}

	if errStatus == "BAD_PARAMETER_LOGIN" {
		errMessage = "Bad parameter login"
	}

	if errStatus == "BAD_PARAMETER_PASSWORD" {
		errMessage = "Bad parameter password"
	}

	if errStatus == "BAD_PARAMETER_NAME" {
		errMessage = "Bad project name parameter"
	}

	if errStatus == "NAME_NOT_UNIQUE" {
		errMessage = "Project name not unique"
	}

	if errStatus == "PROJECT_NOT_FOUND" {
		errMessage = "Project not found"
	}

	return errMessage
}
