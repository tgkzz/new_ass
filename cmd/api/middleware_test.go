package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"greenlight.bcc/internal/assert"
)

func TestRecoverPanicMiddleware(t *testing.T) {
	app := newTestApplication(t)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something went wrong")
	})

	middleware := app.recoverPanic(handler)

	middleware.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	// expectedBody := `{"error":"the server encountered a problem and could not process your request"}`
	// assert.Equal(t, expectedBody, recorder.Body.String())

	expectedJSON := `{"error":"the server encountered a problem and could not process your request"}`
	actualJSON := strings.TrimSpace(recorder.Body.String())

	if !json.Valid([]byte(actualJSON)) {
		t.Fatalf("invalid JSON response: %s", actualJSON)
	}

	var expected interface{}
	var actual interface{}
	if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(actualJSON), &actual); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("unexpected response body:\nexpected: %v\nactual: %v", expected, actual)
	}
}
