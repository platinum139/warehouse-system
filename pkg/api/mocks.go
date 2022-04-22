package api

import (
	"net/http"
	"net/http/httptest"
)

type ResponseRecorderMock struct {
	httptest.ResponseRecorder
}

func (m ResponseRecorderMock) Write([]byte) (int, error) {
	return 0, http.ErrHijacked
}

func NewRecorderMock() *ResponseRecorderMock {
	return &ResponseRecorderMock{}
}
