package http

import (
	"bytes"
	"encoding/json"
	"io"
)

func (r *RawReq) BodyJSON(bodyJSON interface{}) *RawReq {
	r.bodyJSON = bodyJSON
	r.header.Set(contentType, jsonContentType)
	return r
}

func (r *RawReq) getRequestBody() (body io.Reader, err error) {
	if r.bodyJSON != nil && r.header.Get(contentType) == jsonContentType {
		body, err = encodeBodyJSON(r.bodyJSON)
		if err != nil {
			return nil, err
		}
	}
	return body, nil
}

func encodeBodyJSON(bodyJSON interface{}) (io.Reader, error) {

	var buf = new(bytes.Buffer)
	if bodyJSON != nil {
		buf = &bytes.Buffer{}
		err := json.NewEncoder(buf).Encode(bodyJSON)
		if err != nil {
			return nil, err
		}
		//fmt.Fprintf(os.Stdout, "JSON %s", buf.String())
	}
	return buf, nil

}
