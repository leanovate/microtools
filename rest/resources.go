package rest

import (
	"net/http"

	"github.com/leanovate/microtools/routing"
)

type Resources interface {
	Self() Link
	Create(request *http.Request) (Resource, error)
	List(request *http.Request) (interface{}, error)

	FindById(id string) (Resource, error)
}

type ResourcesBase struct{}

func (ResourcesBase) Create(*http.Request) (Resource, error) {
	return nil, MethodNotAllowed
}

func (ResourcesBase) List(*http.Request) (interface{}, error) {
	return nil, MethodNotAllowed
}

func (ResourcesBase) FindById(id string) (Resource, error) {
	return nil, NotFound
}

func ResourcesMatcher(prefix string, collection Resources) routing.Matcher {
	return routing.PrefixSeq(prefix,
		routing.StringPart(func(id string) routing.Matcher {
			resource, err := collection.FindById(id)
			if err != nil {
				return HttpErrorMatcher(WrapError(err))
			}
			return ResourceMatcher(resource)
		}),
		routing.EndSeq(
			routing.GET(restHandler(collection.List)),
			routing.POST(createHandler(collection.Create)),
			HttpErrorMatcher(MethodNotAllowed),
		),
	)
}
