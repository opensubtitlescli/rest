package rest

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalsAndMarshalsUser(t *testing.T) {
	a0 := &User{}
	b0 := "{}"
	equalJSON(t, a0, b0)

	a1 := &UserResponse{}
	b1 := "{}"
	equalJSON(t, a1, b1)

	a2 := &UserResponse{
		Data: &User{
			AllowedDownloads: AllocateInt(100),
			AllowedTranslations: AllocateInt(20),
			DownloadsCount: AllocateInt(1),
			ExtInstalled: AllocateBool(false),
			Level: AllocateString("Sub leecher"),
			RemainingDownloads: AllocateInt(99),
			UserID: AllocateID(66),
			Username: AllocateString("user"),
			VIP: AllocateBool(false),
		},
	}
	b2 := `{
		"data": {
			"allowed_downloads": 100,
			"allowed_translations": 20,
			"downloads_count": 1,
			"ext_installed": false,
			"level": "Sub leecher",
			"remaining_downloads": 99,
			"user_id": 66,
			"username": "user",
			"vip": false
		}
	}`
	equalJSON(t, a2, b2)
}

func TestGetsUser(t *testing.T) {
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

	e := &UserResponse{
		Data: &User{
			UserID: AllocateID(66),
		},
	}
	ctx := context.Background()
	a, _, err := client.Users.Get(ctx)
	require.NoError(t, err)
	assert.Equal(t, e, a)
}
