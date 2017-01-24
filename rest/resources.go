package rest

import (
	"errors"
	"net/http"

	"github.com/leanovate/microtools/routing"
)

type Resources interface {
	Self() Link
	Create(request *http.Request) (Resource, error)
	List(request *http.Request) (interface{}, error)

	FindById(id string) (interface{}, error)
}

type ResourcesBase struct{}

func (ResourcesBase) Create(*http.Request) (Resource, error) {
	return nil, MethodNotAllowed
}

func (ResourcesBase) List(*http.Request) (interface{}, error) {
	return nil, MethodNotAllowed
}

func (ResourcesBase) FindById(id string) (interface{}, error) {
	return nil, NotFound
}

func ResourcesMatcher(prefix string, collection Resources) routing.Matcher {
	return routing.PrefixSeq(prefix,
		routing.StringPart(func(id string) routing.Matcher {
			result, err := collection.FindById(id)
			if err != nil {
				return HttpErrorMatcher(WrapError(err))
			}
			switch resource := (result).(type) {
			case Resource:
				return ResourceMatcher(resource)
			case Resources:
				return ResourcesMatcher("", resource)
			default:
				return HttpErrorMatcher(InternalServerError(errors.New("Invalid result")))
			}
		}),
		routing.EndSeq(
			routing.GET(restHandler(collection.List)),
			routing.POST(createHandler(collection.Create)),
			HttpErrorMatcher(MethodNotAllowed),
		),
	)
}
