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
		Username: AllocateString("username"),
		Password: AllocateString("password"),
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
		User: &User{
			AllowedDownloads: AllocateInt(100),
			AllowedTranslations: AllocateInt(5),
			Level: AllocateString("Sub leecher"),
			UserID: AllocateID(66),
			ExtInstalled: AllocateBool(false),
			VIP: AllocateBool(false),
		},
		BaseURL: AllocateString("api.opensubtitles.com"),
		Token: AllocateString("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9"),
	}
	b = `{
		"base_url": "api.opensubtitles.com",
		"token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9",
		"user": {
			"allowed_downloads": 100,
			"allowed_translations": 5,
			"ext_installed": false,
			"level": "Sub leecher",
			"user_id": 66,
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
		Username: AllocateString("username"),
		Password: AllocateString("password"),
	}
	a, _, err := client.Auth.Login(ctx, c)
	require.NoError(t, err)
	assert.Equal(t, e, a)
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
