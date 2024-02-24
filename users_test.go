package rest

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_UnmarshalsAndMarshals(t *testing.T) {
	a := &User{}
	b := "{}"
	equalJSON(t, a, b)

	a = &User{
		AllowedDownloads: AllocateInt(100),
		AllowedTranslations: AllocateInt(20),
		DownloadsCount: AllocateInt(1),
		ExtInstalled: AllocateBool(false),
		Level: AllocateString("Sub leecher"),
		RemainingDownloads: AllocateInt(99),
		UserID: AllocateID(66),
		Username: AllocateString("user"),
		VIP: AllocateBool(false),
	}
	b = `{
		"allowed_downloads": 100,
		"allowed_translations": 20,
		"downloads_count": 1,
		"ext_installed": false,
		"level": "Sub leecher",
		"remaining_downloads": 99,
		"user_id": 66,
		"username": "user",
		"vip": false
	}`
	equalJSON(t, a, b)
}

func TestUsersService_GetsUser(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/infos/user", func (w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		fmt.Fprint(w, `{
			"data": {
				"user_id": 66
			}
		}`)
	})

	ctx := context.Background()
	e := &User{
		UserID: AllocateID(66),
	}
	a, _, err := client.Users.Get(ctx)
	require.NoError(t, err)
	assert.Equal(t, e, a)
}

func TestUsersService_ReturnsAnErrorIfCannotCreateAURL(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	ctx := context.Background()

	e := useBadBaseURL(client)
	_, _, a := client.Users.Get(ctx)
	assert.EqualError(t, a, e)
}

func TestUsersService__ReturnsAUnsuccessfulResponse(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/infos/user", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	ctx := context.Background()

	_, _, err := client.Users.Get(ctx)
	var a *ErrorResponse
	require.ErrorAs(t, err, &a)
	assert.Equal(t, a.Response.StatusCode, http.StatusBadRequest)
}
