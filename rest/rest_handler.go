package rest

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
)

type restHandler func(request *http.Request) (interface{}, error)

func (h restHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	encoder := StdResponseEncoderChooser(req)
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, string(debug.Stack()))
			InternalServerError(fmt.Errorf("Paniced: %v", r)).Send(resp, encoder)
		}
	}()
	var err error
	result, err := h(req)
	if err == nil {
		switch result.(type) {
		case *Result:
			err = result.(*Result).Send(resp, encoder)
		default:
			err = Ok().WithBody(result).Send(resp, encoder)
		}
	}
	if err != nil {
		WrapError(err).Send(resp, encoder)
	}
}

type createHandler func(*http.Request) (Resource, error)

func (h createHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	encoder := StdResponseEncoderChooser(req)
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, string(debug.Stack()))
			InternalServerError(fmt.Errorf("Paniced: %v", r)).Send(resp, encoder)
		}
	}()
	var err error
	var resource Resource
	resource, err = h(req)
	if err == nil {
		if resource != nil {
			var result interface{}
			result, err = resource.Get(req)
			if err == nil {
				switch result.(type) {
				case *Result:
					err = result.(*Result).
						AddHeader("location", resource.Self().Href).
						WithStatus(201).
						Send(resp, encoder)
				default:
					err = Created().
						AddHeader("location", resource.Self().Href).
						WithBody(result).
						Send(resp, encoder)
				}
			}
		}
	}
	if err != nil {
		WrapError(err).Send(resp, encoder)
	}
}