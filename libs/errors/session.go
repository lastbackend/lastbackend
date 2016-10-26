package errors

var Session session

type session struct {
	*Err
}

func (session) NotFound(e ...error) *Err {
	return &Err{
		Code:   StatusNotFound,
		origin: getError("Session: not found", e...),
		http:   HTTP.getNotFound("session"),
	}
}

func (session) BadParameter(attr string, e ...error) *Err {
	return &Err{
		Code:   StatusBadParameter,
		Attr:   attr,
		origin: getError("Session: bad parameter", e...),
		http:   HTTP.getBadParameter(attr),
	}
}

func (session) IncorrectJSON(e ...error) *Err {
	return &Err{
		Code:   StatusIncorrectJson,
		origin: getError("Session: incorrect json", e...),
		http:   HTTP.getIncorrectJSON(),
	}
}

func (session) Unknown(e ...error) *Err {
	return &Err{
		Code:   StatusUnknown,
		origin: getError("Session: unknown error", e...),
		http:   HTTP.getUnknown(),
	}
}
