package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/powerman/check"
)

func TestGreet(tt *testing.T) {
	t := check.T(tt)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	happy(w, r)
	res := w.Result()
	t.Nil(res.Body.Close())
	t.Equal(res.StatusCode, http.StatusOK)
	rawbody, err := ioutil.ReadAll(res.Body)
	t.Nil(err)
	body := string(rawbody)
	t.Equal(body, "smile )")
}
