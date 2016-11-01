package errors

var Build build

type build struct {
	*Err
}

func (build) NotFound(e ...error) *Err {
	return &Err{
		Code:   StatusNotFound,
		origin: getError("Build: not found", e...),
		http:   HTTP.getNotFound("build"),
	}
}

func (build) BadParameter(attr string, e ...error) *Err {
	return &Err{
		Code:   StatusBadParameter,
		Attr:   attr,
		origin: getError("Build: bad parameter", e...),
		http:   HTTP.getBadParameter(attr),
	}
}

func (build) IncorrectJSON(e ...error) *Err {
	return &Err{
		Code:   StatusIncorrectJson,
		origin: getError("Build: incorrect json", e...),
		http:   HTTP.getIncorrectJSON(),
	}
}

func (build) Unknown(e ...error) *Err {
	return &Err{
		Code:   StatusUnknown,
		origin: getError("Build: unknown error", e...),
		http:   HTTP.getUnknown(),
	}
}
