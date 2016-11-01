package errors

var Image image

type image struct {
	*Err
}

func (image) NotFound(e ...error) *Err {
	return &Err{
		Code:   StatusNotFound,
		origin: getError("Image: not found", e...),
		http:   HTTP.getNotFound("image"),
	}
}

func (image) BadParameter(attr string, e ...error) *Err {
	return &Err{
		Code:   StatusBadParameter,
		Attr:   attr,
		origin: getError("Image: bad parameter", e...),
		http:   HTTP.getBadParameter(attr),
	}
}

func (image) IncorrectJSON(e ...error) *Err {
	return &Err{
		Code:   StatusIncorrectJson,
		origin: getError("Image: incorrect json", e...),
		http:   HTTP.getIncorrectJSON(),
	}
}

func (image) Unknown(e ...error) *Err {
	return &Err{
		Code:   StatusUnknown,
		origin: getError("Image: unknown error", e...),
		http:   HTTP.getUnknown(),
	}
}
