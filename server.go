package purehttp

import "net/http"

type HTTPResponse struct{}

type HandlerFunc func(req *http.Request) (*HTTPResponse, error)

type Handler struct {
	handle HandlerFunc
}

func NewHandler(handle HandlerFunc) Handler {
	return Handler{handle: handle}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	_, err := h.handle(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
