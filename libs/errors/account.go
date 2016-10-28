package errors

var Account account

type account struct {
	*Err
}

func (account) NotFound(e ...error) *Err {
	return &Err{
		Code:   StatusNotFound,
		origin: getError("Account: not found", e...),
		http:   HTTP.getNotFound("account"),
	}
}

func (account) AccessDenied(e ...error) *Err {
	return &Err{
		Code:   StatusAccessDenied,
		origin: getError("Account: access denied", e...),
		http:   HTTP.getAccessDenied(),
	}
}

func (account) BadParameter(attr string, e ...error) *Err {
	return &Err{
		Code:   StatusBadParameter,
		Attr:   attr,
		origin: getError("Account: bad parameter", e...),
		http:   HTTP.getBadParameter(attr),
	}
}

func (account) IncorrectJSON(e ...error) *Err {
	return &Err{
		Code:   StatusIncorrectJson,
		origin: getError("Account: incorrect json", e...),
		http:   HTTP.getIncorrectJSON(),
	}
}

func (account) Unknown(e ...error) *Err {
	return &Err{
		Code:   StatusUnknown,
		origin: getError("Account: unknow error", e...),
		http:   HTTP.getUnknown(),
	}
}
