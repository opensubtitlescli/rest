package rest

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCredentials_UnmarshalsAndMarshals(t *testing.T) {
	a := &Credentials{}
	b := "{}"
	equalJSON(t, a, b)

	a = &Credentials{
		Username: "username",
		Password: "password",
	}
	b = `{
		"username": "username",
		"password": "password"
	}`
	equalJSON(t, a, b)
}

func TestLogin_UnmarshalsAndMarshals(t *testing.T) {
	a := &Login{
		ClientBaseURL: defaultBaseURL,
	}
	b := "{}"
	equalJSON(t, a, b)

	a = &Login{
		ClientBaseURL: defaultBaseURL,
		User: &User{
			VIP: AllocateBool(false),
		},
	}
	b = `{
		"user": {
			"vip": false
		}
	}`
	equalJSON(t, a, b)

	a = &Login{
		ClientBaseURL: vipBaseURL,
		User: &User{
			VIP: AllocateBool(true),
		},
	}
	b = `{
		"user": {
			"vip": true
		}
	}`
	equalJSON(t, a, b)

	a = &Login{
		ClientBaseURL: defaultBaseURL,
		BaseURL: AllocateString("api.opensubtitles.com"),
		Token: AllocateString("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9"),
		User: &User{
			VIP: AllocateBool(false),
		},
	}
	b = `{
		"base_url": "api.opensubtitles.com",
		"token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9",
		"user": {
			"vip": false
		}
	}`
	equalJSON(t, a, b)
}

func TestAuthServiceLogin_Logins(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/login", func (w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		equalBody(t, r.Body, `{
			"username": "username"
		}`)
		fmt.Fprint(w, `{
			"token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9"
		}`)
	})

	ctx := context.Background()

	e := &Login{
		ClientBaseURL: defaultBaseURL,
		Token: AllocateString("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9"),
	}
	c := &Credentials{
		Username: "username",
	}
	a, _, err := client.Auth.Login(ctx, c)
	require.NoError(t, err)
	assert.Equal(t, e, a)
}

func TestAuthServiceLogin_ReturnsAnErrorIfCannotCreateAURL(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	ctx := context.Background()

	e := useBadBaseURL(client)
	_, _, a := client.Auth.Login(ctx, nil)
	assert.EqualError(t, a, e)
}

func TestAuthServiceLogin_ReturnsAUnsuccessfulResponse(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/login", func (w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	ctx := context.Background()

	_, _, err := client.Auth.Login(ctx, nil)
	var a *ErrorResponse
	require.ErrorAs(t, err, &a)
	assert.Equal(t, a.Response.StatusCode, http.StatusBadRequest)
}

func TestAuthServiceLogout_Logouts(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/logout", func (w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		fmt.Fprint(w, `{
			"message": "token successfully destroyed",
			"status": 200
		}`)
	})

	ctx := context.Background()

	_, err := client.Auth.Logout(ctx)
	require.NoError(t, err)
}

func TestAuthServiceLogout_ReturnsAnErrorIfCannotCreateAURL(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	ctx := context.Background()

	e := useBadBaseURL(client)
	_, a := client.Auth.Logout(ctx)
	assert.EqualError(t, a, e)
}

func TestAuthServiceLogout_ReturnsAUnsuccessfulResponse(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/logout", func (w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	ctx := context.Background()

	_, err := client.Auth.Logout(ctx)
	var a *ErrorResponse
	require.ErrorAs(t, err, &a)
	assert.Equal(t, a.Response.StatusCode, http.StatusBadRequest)
}
