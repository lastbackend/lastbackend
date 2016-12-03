package errors

var Volume volume

type volume struct {
	*Err
}

func (volume) NotFound(e ...error) *Err {
	return &Err{
		Code:   StatusNotFound,
		origin: getError("Volume: not found", e...),
		http:   HTTP.getNotFound("volume"),
	}
}

func (volume) Unknown(e ...error) *Err {
	return &Err{
		Code:   StatusUnknown,
		origin: getError("Volume: unknown error", e...),
		http:   HTTP.getUnknown(),
	}
}
