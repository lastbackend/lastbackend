package errors

var Project project

type project struct {
	*Err
}

func (project) NotFound(e ...error) *Err {
	return &Err{
		Code:   StatusNotFound,
		origin: getError("Project: not found", e...),
		http:   HTTP.getNotFound("project"),
	}
}

func (project) BadParameter(attr string, e ...error) *Err {
	return &Err{
		Code:   StatusBadParameter,
		Attr:   attr,
		origin: getError("Project: bad parameter", e...),
		http:   HTTP.getBadParameter(attr),
	}
}

func (project) IncorrectJSON(e ...error) *Err {
	return &Err{
		Code:   StatusIncorrectJson,
		origin: getError("Project: incorrect json", e...),
		http:   HTTP.getIncorrectJSON(),
	}
}

func (project) Unknown(e ...error) *Err {
	return &Err{
		Code:   StatusUnknown,
		origin: getError("Project: unknown error", e...),
		http:   HTTP.getUnknown(),
	}
}
