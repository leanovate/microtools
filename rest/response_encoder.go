package rest

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"strings"
)

type ResponseEncoder func(output io.Writer, data interface{}) error

func JsonResponseEncoder(output io.Writer, data interface{}) error {
	return json.NewEncoder(output).Encode(data)
}

func XmlResponseEncoder(output io.Writer, data interface{}) error {
	return xml.NewEncoder(output).Encode(data)
}

type ResponseEncoderChooser func(*http.Request) ResponseEncoder

func StdResponseEncoderChooser(request *http.Request) ResponseEncoder {
	accept := request.Header.Get("accept")

	if strings.Contains(accept, "text/xml") || strings.Contains(accept, "application/xml") {
		return XmlResponseEncoder
	}
	return JsonResponseEncoder
}
