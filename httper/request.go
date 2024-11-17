package httper

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
)

type Req struct {
	params        *Params
	NeedUnmarshal bool
	*http.Request
}

type Params struct {
	Method method
	Url    string

	Body     interface{}
	ByteBody []byte

	Marshal     bool
	MarshalType contentType

	Unmarshal     bool
	UnmarshalTo   interface{}
	UnmarshalType contentType
}

func NewReq(params *Params) (*Req, error) {
	var body []byte

	var err error

	if params.Marshal {
		body, err = marshal(params)
		if err != nil {
			return nil, err
		}
	} else {
		body = params.ByteBody
	}

	reader := bytes.NewReader(body)

	base, err := http.NewRequest(string(params.Method), params.Url, reader)
	if err != nil {
		return nil, err
	}

	return &Req{
		params,
		params.Unmarshal,
		base,
	}, nil
}

func (r *Req) Unmarshal(body []byte) error {
	switch r.params.UnmarshalType {

	case JsonType:
		err := json.Unmarshal(body, r.params.UnmarshalTo)
		if err != nil {
			return err
		}
	case XmlType:
		err := xml.Unmarshal(body, r.params.UnmarshalTo)
		if err != nil {
			return err
		}
	case TextType:
		if ptr, ok := r.params.UnmarshalTo.(*string); ok {
			*ptr = string(body)
		} else {
			return errors.New("incorrect type")
		}
	case HtmlType:
		if ptr, ok := r.params.UnmarshalTo.(*string); ok {
			*ptr = string(body)
		} else {
			return errors.New("incorrect type")
		}
	default:
		return errors.New("incorrect type")
	}

	return nil
}

func marshal(params *Params) ([]byte, error) {
	var body []byte

	switch params.MarshalType {

	case JsonType:
		enc, err := json.Marshal(params.Body)
		if err != nil {
			return nil, err
		}
		body = enc
	case XmlType:
		enc, err := xml.Marshal(params.Body)
		if err != nil {
			return nil, err
		}
		body = enc
	case TextType:
		body = []byte(fmt.Sprintf("%v", params.Body))
	case HtmlType:
		body = []byte(fmt.Sprintf("%v", params.Body))
	default:
		return nil, errors.New("incorrect type")
	}

	return body, nil
}
