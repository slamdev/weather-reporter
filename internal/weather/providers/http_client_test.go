package providers

import (
	"bytes"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"net/http"
)

func NewClientStub(data string, statusCode int, err error, requestAsserts ...func(*http.Request)) http.Client {
	handler := func(req *http.Request) (*http.Response, error) {
		for _, assertFunc := range requestAsserts {
			assertFunc(req)
		}
		response := &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(data)),
			StatusCode: statusCode,
			Header:     make(http.Header),
		}
		return response, err
	}
	return http.Client{Transport: promhttp.RoundTripperFunc(handler)}
}
