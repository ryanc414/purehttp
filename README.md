purehttp
--------

Normally when you want to serve HTTP requests in Go you define a type which
implements the http.Handler interface:

```
func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {...}
```

The Handler interface requires mutating the http.ResponseWriter object, which
has its own particular set of rules on which order you must call certain
methods.

This module provides an alternate interface, where you supply a pure function
which returns a response object without requiring side-effects. Here is what it
looks like:

```
func main() {
	srv := http.Server{
		Addr:    "localhost:8080",
		Handler: purehttp.NewHandler(MyHandler),
	}
	srv.ListenAndServe()
}

func MyHandler(req *http.Request) (*purehttp.Response, error) {
	if req.Method != http.MethodPost {
		return &purehttp.Response{StatusCode: http.StatusMethodNotAllowed}, nil
	}

	if err := req.ParseForm(); err != nil {
		return nil, err
	}

	if req.Form.Get("foo") == "" {
		return &purehttp.Response{StatusCode: http.StatusBadRequest}, nil
	}

	response := map[string]interface{}{
		"code":    420,
		"message": "hello world",
	}
	rspData, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return &purehttp.Response{Body: rspData, JSON: true}, nil
}
```

