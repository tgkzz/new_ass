package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"greenlight.bcc/internal/assert"
)

func TestRegisterUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	const (
		validName     = "Giorno"
		validEmail    = "giovanna@gmail.com"
		validPassword = "GoldenExperience123"
	)

	tests := []struct {
		name     string
		Name     string
		Email    string
		Password string
		wantCode int
	}{
		{
			name:     "valid submission",
			Name:     validName,
			Email:    validEmail,
			Password: validPassword,
			wantCode: http.StatusCreated,
		},
		{
			name:     "empty name",
			Name:     "",
			Email:    validEmail,
			Password: validPassword,
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "email address not valid",
			Name: validName,
			// Email:    "not@valid",
			Email:    "not@valid.",
			Password: validPassword,
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "test for wrong input",
			Name:     validName,
			Email:    validEmail,
			Password: validPassword,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Duplicate email",
			Name:     "Mock Test",
			Email:    "duplicate@test.com",
			Password: "12345678",
			wantCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputData := struct {
				Name     string
				Email    string
				Password string
			}{
				Name:     tt.Name,
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

			code, _, _ := ts.postForm(t, "/v1/users", b)

			assert.Equal(t, code, tt.wantCode)
		})
	}
}

func TestActivateUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	const (
		validTokenPlaintext = "GZDAMCLCS6QHZ6FBTVB4WKH3EI"
	)

	tests := []struct {
		name           string
		TokenPlaintext string
		wantCode       int
	}{
		{
			name:           "Valid submission",
			TokenPlaintext: validTokenPlaintext,
			wantCode:       http.StatusOK,
		},
		{
			name:           "Empty Token Plaintext",
			TokenPlaintext: "",
			wantCode:       http.StatusUnprocessableEntity,
		},
		{
			name:           "Invalid(expired) Token",
			TokenPlaintext: "FGTUP3VTWFI7HZDYS3IPWN25E",
			wantCode:       http.StatusUnprocessableEntity,
		},
		{
			name:           "test for wrong input",
			TokenPlaintext: validTokenPlaintext,
			wantCode:       http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputData := struct {
				TokenPlaintext string `json:"token"`
			}{
				TokenPlaintext: tt.TokenPlaintext,
			}

			b, err := json.Marshal(&inputData)

			if err != nil {
				t.Fatal("wrong input data")
			}
			if tt.name == "test for wrong input" {
				b = append(b, 'a')
			}

			code, _, _ := ts.putForm(t, "/v1/users/activated", b)

			assert.Equal(t, code, tt.wantCode)
		})
	}
}
