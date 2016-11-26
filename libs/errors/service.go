package errors

var Service service

type service struct {
	*Err
}

func (service) NotFound(e ...error) *Err {
	return &Err{
		Code:   StatusNotFound,
		origin: getError("Service: not found", e...),
		http:   HTTP.getNotFound("service"),
	}
}

func (service) BadParameter(attr string, e ...error) *Err {
	return &Err{
		Code:   StatusBadParameter,
		Attr:   attr,
		origin: getError("Service: bad parameter", e...),
		http:   HTTP.getBadParameter(attr),
	}
}

func (service) IncorrectJSON(e ...error) *Err {
	return &Err{
		Code:   StatusIncorrectJson,
		origin: getError("Service: incorrect json", e...),
		http:   HTTP.getIncorrectJSON(),
	}
}

func (service) Unknown(e ...error) *Err {
	return &Err{
		Code:   StatusUnknown,
		origin: getError("Service: unknown error", e...),
		http:   HTTP.getUnknown(),
	}
}
