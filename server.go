package purehttp

import "net/http"

type HTTPResponse struct {
	Body       []byte
	StatusCode int
	Header     http.Header
}

type HandlerFunc func(req *http.Request) (*HTTPResponse, error)

type Handler struct {
	handle HandlerFunc
}

func NewHandler(handle HandlerFunc) Handler {
	return Handler{handle: handle}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rsp, err := h.handle(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	header := w.Header()

	for h := range rsp.Header {
		val := rsp.Header.Get(h)
		header.Set(h, val)
	}

	statusCode := rsp.StatusCode
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	w.WriteHeader(statusCode)
	if len(rsp.Body) > 0 {
		w.Write(rsp.Body)
	}
}
