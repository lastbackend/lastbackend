package errors

var Template template

type template struct {
	*Err
}

func (template) NotFound(e ...error) *Err {
	return &Err{
		Code:   StatusNotFound,
		origin: getError("Template: not found", e...),
		http:   HTTP.getNotFound("template"),
	}
}

func (template) BadParameter(attr string, e ...error) *Err {
	return &Err{
		Code:   StatusBadParameter,
		Attr:   attr,
		origin: getError("Template: bad parameter", e...),
		http:   HTTP.getBadParameter(attr),
	}
}

func (template) IncorrectJSON(e ...error) *Err {
	return &Err{
		Code:   StatusIncorrectJson,
		origin: getError("Template: incorrect json", e...),
		http:   HTTP.getIncorrectJSON(),
	}
}

func (template) Unknown(e ...error) *Err {
	return &Err{
		Code:   StatusUnknown,
		origin: getError("Template: unknow error", e...),
		http:   HTTP.getUnknown(),
	}
}
