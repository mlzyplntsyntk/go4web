package middleware

import (
	"fmt"
	"../core"
	"net/http"
)

type ContentHandler struct {
	
}

func (t ContentHandler) RunBeforeHandled() bool {
	return false
}

func (t ContentHandler) RunAfterHandled() bool {
	return false
}

func (t ContentHandler) Pattern() string {
	return "*"
}

func (t ContentHandler) Run(w http.ResponseWriter, r *http.Request, route *core.Page) (result bool, err error) {
	fmt.Fprintf(w, route.Name)
	return
}