package purehttp_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ryanc414/purehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	f := func(req *http.Request) (*purehttp.Response, error) {
		return nil, errors.New("not implemented")
	}

	h := purehttp.NewHandler(f)
	ts := httptest.NewServer(h)
	t.Cleanup(ts.Close)

	rsp, err := http.Get(ts.URL)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, rsp.StatusCode)
	data, err := ioutil.ReadAll(rsp.Body)
	rsp.Body.Close()
	require.NoError(t, err)

	assert.Equal(t, "not implemented\n", string(data))
}

func TestHappy(t *testing.T) {
	f := func(req *http.Request) (*purehttp.Response, error) {
		return &purehttp.Response{
			Body:       []byte("{\"foo\":\"bar\"}"),
			StatusCode: http.StatusAccepted,
			JSON:       true,
			Header:     map[string][]string{"Server": {"Go"}},
		}, nil
	}

	h := purehttp.NewHandler(f)
	ts := httptest.NewServer(h)
	t.Cleanup(ts.Close)

	rsp, err := http.Get(ts.URL)
	require.NoError(t, err)

	assert.Equal(t, "application/json", rsp.Header.Get("Content-Type"))
	assert.Equal(t, "Go", rsp.Header.Get("Server"))

	assert.Equal(t, http.StatusAccepted, rsp.StatusCode)
	data, err := ioutil.ReadAll(rsp.Body)
	rsp.Body.Close()
	require.NoError(t, err)

	var res struct{ Foo string }
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)

	assert.Equal(t, "bar", res.Foo)
}

func TestDefault(t *testing.T) {
	f := func(req *http.Request) (*purehttp.Response, error) {
		return &purehttp.Response{}, nil
	}

	h := purehttp.NewHandler(f)
	ts := httptest.NewServer(h)
	t.Cleanup(ts.Close)

	rsp, err := http.Get(ts.URL)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rsp.StatusCode)
	data, err := ioutil.ReadAll(rsp.Body)
	rsp.Body.Close()
	require.NoError(t, err)

	assert.Empty(t, data)
}
