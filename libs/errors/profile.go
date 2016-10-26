package errors

var Profile profile

type profile struct {
	*Err
}

func (profile) NotFound(e ...error) *Err {
	return &Err{
		Code:   StatusNotFound,
		origin: getError("Profile: not found", e...),
		http:   HTTP.getNotFound("profile"),
	}
}

func (profile) BadParameter(attr string, e ...error) *Err {
	return &Err{
		Code:   StatusBadParameter,
		Attr:   attr,
		origin: getError("Profile: bad parameter", e...),
		http:   HTTP.getBadParameter(attr),
	}
}

func (profile) Unknown(e ...error) *Err {
	return &Err{
		Code:   StatusUnknown,
		origin: getError("Profile: unknow error", e...),
		http:   HTTP.getUnknown(),
	}
}
