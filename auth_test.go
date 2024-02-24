package rest

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalsAndMarshalsCredentials(t *testing.T) {
	a0 := &Credentials{}
	b0 := "{}"
	equalJSON(t, a0, b0)

	a1 := &Credentials{
		Username: AllocateString("username"),
		Password: AllocateString("password"),
	}
	b1 := `{
		"username": "username",
		"password": "password"
	}`
	equalJSON(t, a1, b1)
}

func TestUnmarshalsAndMarshalsLogin(t *testing.T) {
	a0 := &Login{
		ClientBaseURL: defaultBaseURL,
	}
	b0 := "{}"
	equalJSON(t, a0, b0)

	a1 := &Login{
		ClientBaseURL: defaultBaseURL,
		User: &User{
			VIP: AllocateBool(false),
		},
	}
	b1 := `{
		"user": {
			"vip": false
		}
	}`
	equalJSON(t, a1, b1)

	a2 := &Login{
		ClientBaseURL: vipBaseURL,
		User: &User{
			VIP: AllocateBool(true),
		},
	}
	b2 := `{
		"user": {
			"vip": true
		}
	}`
	equalJSON(t, a2, b2)

	a3 := &Login{
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
	b3 := `{
		"user": {
			"allowed_downloads": 100,
			"allowed_translations": 5,
			"level": "Sub leecher",
			"user_id": 66,
			"ext_installed": false,
			"vip": false
		},
		"base_url": "api.opensubtitles.com",
		"token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9"
	}`
	equalJSON(t, a3, b3)
}

func TestLogins(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/login", func (w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		fmt.Fprint(w, `{
			"token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9"
		}`)
	})

	e := &Login{
		ClientBaseURL: defaultBaseURL,
		Token: AllocateString("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9"),
	}
	ctx := context.Background()
	c := &Credentials{
		Username: AllocateString("username"),
		Password: AllocateString("password"),
	}
	a, _, err := client.Auth.Login(ctx, c)
	require.NoError(t, err)
	assert.Equal(t, e, a)
}

func TestLogouts(t *testing.T) {
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
