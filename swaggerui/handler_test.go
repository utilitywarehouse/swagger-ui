package swaggerui

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHandler(t *testing.T) {
	cases := []struct {
		Request *http.Request
		Code    int
	}{
		{
			Request: &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: "/"},
			},
			Code: http.StatusSeeOther,
		},
		{
			Request: &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: specFile},
			},
			Code: http.StatusOK,
		},
		{
			Request: &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: basePath},
			},
			Code: http.StatusOK,
		},
	}
	f := Handler()
	for _, tc := range cases {
		w := &httptest.ResponseRecorder{}
		f.ServeHTTP(w, tc.Request)
		if w.Code != tc.Code {
			t.Logf("headers: %#v", w.Header())
			t.Fatalf("unexpected code for %s: want %d, have %d",
				tc.Request.URL.Path, tc.Code, w.Code)
		}
	}
}
