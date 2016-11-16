package errors

var User user

type user struct {
	*Err
}

func (user) NotFound(e ...error) *Err {
	return &Err{
		Code:   StatusNotFound,
		origin: getError("User: not found", e...),
		http:   HTTP.getNotFound("user"),
	}
}

func (user) AccessDenied(e ...error) *Err {
	return &Err{
		Code:   StatusAccessDenied,
		origin: getError("User: access denied", e...),
		http:   HTTP.getAccessDenied(),
	}
}

func (user) BadParameter(attr string, e ...error) *Err {
	return &Err{
		Code:   StatusBadParameter,
		Attr:   attr,
		origin: getError("User: bad parameter", e...),
		http:   HTTP.getBadParameter(attr),
	}
}

func (user) EmailExists(e ...error) *Err {
	return &Err{
		Code:   StatusNotUnique,
		origin: getError("User: email not unique", e...),
		http:   HTTP.getNotUnique("email"),
	}
}

func (user) UsernameExists(e ...error) *Err {
	return &Err{
		Code:   StatusNotUnique,
		origin: getError("User: username not unique", e...),
		http:   HTTP.getNotUnique("username"),
	}
}

func (user) IncorrectJSON(e ...error) *Err {
	return &Err{
		Code:   StatusIncorrectJson,
		origin: getError("User: incorrect json", e...),
		http:   HTTP.getIncorrectJSON(),
	}
}

func (user) Unknown(e ...error) *Err {
	return &Err{
		Code:   StatusUnknown,
		origin: getError("User: unknow error", e...),
		http:   HTTP.getUnknown(),
	}
}
