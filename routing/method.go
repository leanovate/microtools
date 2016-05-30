package routing

import "net/http"

func Method(method string, handler http.Handler) Matcher {
	return func(remainingPath string, resp http.ResponseWriter, req *http.Request) bool {
		if req.Method == method {
			handler.ServeHTTP(resp, req)
			return true
		}
		return false
	}
}

func GET(handler http.Handler) Matcher {
	return Method("GET", handler)
}

func GETFunc(handler func(http.ResponseWriter, *http.Request)) Matcher {
	return GET(http.HandlerFunc(handler))
}

func POST(handler http.Handler) Matcher {
	return Method("POST", handler)
}

func POSTFunc(handler func(http.ResponseWriter, *http.Request)) Matcher {
	return POST(http.HandlerFunc(handler))
}

func PUT(handler http.Handler) Matcher {
	return Method("PUT", handler)
}

func PUTFunc(handler func(http.ResponseWriter, *http.Request)) Matcher {
	return PUT(http.HandlerFunc(handler))
}

func PATCH(handler http.Handler) Matcher {
	return Method("PATCH", handler)
}

func PATCHFunc(handler func(http.ResponseWriter, *http.Request)) Matcher {
	return PATCH(http.HandlerFunc(handler))
}

func DELETE(handler http.Handler) Matcher {
	return Method("DELETE", handler)
}

func DELETEFunc(handler func(http.ResponseWriter, *http.Request)) Matcher {
	return DELETE(http.HandlerFunc(handler))
}

func MethodNotAllowed(remainingPath string, resp http.ResponseWriter, req *http.Request) bool {
	resp.WriteHeader(405)
	return true
}
