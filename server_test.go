package purehttp_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ryanc414/purehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	f := func(req *http.Request) (*purehttp.HTTPResponse, error) {
		return nil, errors.New("not implemented")
	}

	h := purehttp.NewHandler(f)
	ts := httptest.NewServer(h)
	t.Cleanup(ts.Close)

	rsp, err := http.Get(ts.URL)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, rsp.StatusCode)
}
