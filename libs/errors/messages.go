package errors

func Message(errStatus string) string {
	
	var errMessage string

	if errStatus == "Internal Server Error" {
		errMessage = "oops, error occurred: Please provide bug to github: https://github.com/lastbackend/lastbackend/issues"
	}

	return errMessage
}
