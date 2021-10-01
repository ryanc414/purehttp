package purehttp

import "net/http"

type Response struct {
	Body       []byte
	Header     http.Header
	JSON       bool
	StatusCode int
}

type HandlerFunc func(req *http.Request) (*Response, error)

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

	if rsp.JSON && len(rsp.Body) > 0 {
		header.Set("Content-Type", "application/json")
	}

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
