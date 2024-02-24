package rest

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalsAndMarshalsFormat(t *testing.T) {
	a0 := &FormatsData{}
	b0 := "{}"
	equalJSON(t, a0, b0)

	a1 := &FormatsResponse{}
	b1 := "{}"
	equalJSON(t, a1, b1)

	a2 := &FormatsResponse{
		Data: &FormatsData{
			OutputFormats: []*string{
				AllocateString("srt"),
			},
		},
	}
	b2 := `{
		"data": {
			"output_formats": ["srt"]
		}
	}`
	equalJSON(t, a2, b2)
}

func TestListsFormats(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/infos/formats", func (w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		fmt.Fprint(w, `{
			"data": {
				"output_formats": ["srt"]
			},
			"status": 200
		}`)
	})

	e := &FormatsResponse{
		Data: &FormatsData{
			OutputFormats: []*string{
				AllocateString("srt"),
			},
		},
	}
	ctx := context.Background()
	a, _, err := client.Formats.List(ctx)
	require.NoError(t, err)
	assert.Equal(t, e, a)
}
