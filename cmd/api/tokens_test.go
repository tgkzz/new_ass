package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"greenlight.bcc/internal/assert"
)

func TestCreateAuthenticationToken(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	const (
		validEmail    = "mock@test.com"
		validPassword = "examplepassword"
	)

	tests := []struct {
		name     string
		Email    string
		Password string
		wantCode int
	}{
		{
			name:     "Valid submission",
			Email:    validEmail,
			Password: validPassword,
			wantCode: http.StatusCreated,
		},
		{
			name:     "test for wrong input",
			Email:    validEmail,
			Password: validPassword,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Invalid credentials",
			Email:    validEmail,
			Password: "helloWorld",
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "Wrong validation",
			Email:    "wrong@email",
			Password: "wrong",
			wantCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputData := struct {
				Email    string
				Password string
			}{
				Email:    tt.Email,
				Password: tt.Password,
			}

			b, err := json.Marshal(&inputData)
			if err != nil {
				t.Fatal("wrong input data")
			}
			if tt.name == "test for wrong input" {
				b = append(b, 'a')
			}

			code, _, _ := ts.postForm(t, "/v1/tokens/authentication", b)

			assert.Equal(t, code, tt.wantCode)
		})
	}
}
