package httper

import (
	"net/http"
)

type Resp struct {
	ByteBody []byte
	*http.Response
}

func newResp(resp *http.Response, body []byte) *Resp {
	return &Resp{
		body,
		resp,
	}
}
