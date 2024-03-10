package rest

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatsListResponse_UnmarshalsAndMarshals(t *testing.T) {
	a := &FormatsListResponse{}
	b := "{}"
	equalJSON(t, a, b)

	a = &FormatsListResponse{
		OutputFormats: []*string{
			AllocateString("srt"),
		},
	}
	b = `{
		"output_formats": ["srt"]
	}`
	equalJSON(t, a, b)
}

func TestLFormatsServiceList_ListsFormats(t *testing.T) {
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

	ctx := context.Background()
	e := &FormatsListResponse{
		OutputFormats: []*string{
			AllocateString("srt"),
		},
	}
	a, _, err := client.Formats.List(ctx)
	require.NoError(t, err)
	assert.Equal(t, e, a)
}
