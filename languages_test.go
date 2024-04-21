package rest

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLanguage_UnmarshalsAndMarshals(t *testing.T) {
	a := &Language{}
	b := "{}"
	equalJSON(t, a, b)

	a = &Language{
		LanguageCode: AllocateString("en"),
		LanguageName: AllocateString("English"),
	}
	b = `{
		"language_code": "en",
		"language_name": "English"
	}`
	equalJSON(t, a, b)
}

func TestLanguagesServiceList_ListsLanguages(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/infos/languages", func (w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		fmt.Fprint(w, `{
			"data": [
				{
					"language_code": "en"
				}
			],
			"status": 200
		}`)
	})

	ctx := context.Background()

	e := []*Language{
		{
			LanguageCode: AllocateString("en"),
		},
	}
	a, _, err := client.Languages.List(ctx)
	require.NoError(t, err)
	assert.Equal(t, e, a)
}

func TestLanguagesServiceList_ReturnsAnErrorIfCannotCreateAURL(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	ctx := context.Background()

	e := useBadBaseURL(client)
	_, _, a := client.Languages.List(ctx)
	assert.EqualError(t, a, e)
}

func TestLanguagesServiceList_ReturnsAUnsuccessfulResponse(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/infos/languages", func (w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	ctx := context.Background()

	_, _, err := client.Languages.List(ctx)
	var a *ErrorResponse
	require.ErrorAs(t, err, &a)
	assert.Equal(t, a.Response.StatusCode, http.StatusBadRequest)
}
