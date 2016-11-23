package http

const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"
)

func (r *RawReq) POST(pathURL string) *RawReq {
	r.method = MethodPost
	r.rawURL = r.host + pathURL
	return r
}

func (r *RawReq) GET(pathURL string) *RawReq {
	r.method = MethodGet
	r.rawURL = r.host + pathURL
	return r
}

func (r *RawReq) PUT(pathURL string) *RawReq {
	r.method = MethodPut
	r.rawURL = r.host + pathURL
	return r
}

func (r *RawReq) DELETE(pathURL string) *RawReq {
	r.method = MethodDelete
	r.rawURL = r.host + pathURL
	return r
}
