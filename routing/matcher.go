package routing

import "net/http"

type Matcher func(remainingPath string, resp http.ResponseWriter, req *http.Request) bool
