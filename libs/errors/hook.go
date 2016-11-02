package errors

var Hook hook

type hook struct {
	*Err
}

func (hook) NotFound(e ...error) *Err {
	return &Err{
		Code:   StatusNotFound,
		origin: getError("Hook: not found", e...),
		http:   HTTP.getNotFound("hook"),
	}
}

func (hook) BadParameter(attr string, e ...error) *Err {
	return &Err{
		Code:   StatusBadParameter,
		Attr:   attr,
		origin: getError("Hook: bad parameter", e...),
		http:   HTTP.getBadParameter(attr),
	}
}

func (hook) IncorrectJSON(e ...error) *Err {
	return &Err{
		Code:   StatusIncorrectJson,
		origin: getError("Hook: incorrect json", e...),
		http:   HTTP.getIncorrectJSON(),
	}
}

func (hook) Unknown(e ...error) *Err {
	return &Err{
		Code:   StatusUnknown,
		origin: getError("Hook: unknown error", e...),
		http:   HTTP.getUnknown(),
	}
}
