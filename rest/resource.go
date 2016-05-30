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

	SubResource(name string) (routing.Matcher, error)
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

func (ResourceBase) SubResource(name string) (routing.Matcher, error) {
	return nil, NotFound
}

func ResourceMatcher(resource Resource) routing.Matcher {
	return routing.Sequence(
		routing.StringPart(func(name string) routing.Matcher {
			subResource, err := resource.SubResource(name)
			if err != nil {
				return HttpErrorMatcher(WrapError(err))
			}
			return subResource
		}),
		routing.EndSeq(
			routing.GET(restHandler(resource.Get)),
			routing.PUT(restHandler(resource.Update)),
			routing.PATCH(restHandler(resource.Patch)),
			routing.DELETE(restHandler(resource.Delete)),
			HttpErrorMatcher(MethodNotAllowed),
		),
	)
}
