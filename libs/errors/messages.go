package errors

func Message(errStatus string) string {

	var errMessage string

	if errStatus == "ACCESS_DENIED" {
		errMessage = "You need to be authorized to provide this operation"
	}
	if errStatus == "USER_NOT_FOUND" {
		errMessage = "User not found"
	}
	if errStatus == "INTERNAL_SERVER_ERROR"{
		errMessage = "oops, error occurred: Please provide bug to github: https://github.com/lastbackend/lastbackend/issues"
	}

	return errMessage
}
