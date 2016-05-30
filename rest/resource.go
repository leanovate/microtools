package rest

import (
	"net/http"

	"github.com/leanovate/microtools/routing"
)

type Resource interface {
	Self() Link
	Get(request *http.Request) (interface{}, error)
	Patch(request *http.Request) (interface{}, error)
	Update(request *http.Request) (interface{}, error)
	Delete(request *http.Request) (interface{}, error)

	SubResources() routing.Matcher
}

type ResourceBase struct{}

func (ResourceBase) Get(request *http.Request) (interface{}, error) {
	return nil, MethodNotAllowed
}

func (ResourceBase) Patch(request *http.Request) (interface{}, error) {
	return nil, MethodNotAllowed
}

func (ResourceBase) Update(request *http.Request) (interface{}, error) {
	return nil, MethodNotAllowed
}

func (ResourceBase) Delete(request *http.Request) (interface{}, error) {
	return nil, MethodNotAllowed
}

func (ResourceBase) SubResources() routing.Matcher {
	return HttpErrorMatcher(NotFound)
}

func ResourceMatcher(resource Resource) routing.Matcher {
	return routing.Sequence(
		routing.EndSeq(
			routing.GET(restHandler(resource.Get)),
			routing.PUT(restHandler(resource.Update)),
			routing.PATCH(restHandler(resource.Patch)),
			routing.DELETE(restHandler(resource.Delete)),
			HttpErrorMatcher(MethodNotAllowed),
		),
		resource.SubResources(),
	)
}
