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
					"language_code": "en",
					"language_name": "English"
				}
			],
			"status": 200
		}`)
	})

	ctx := context.Background()
	e := []*Language{
		{
			LanguageCode: AllocateString("en"),
			LanguageName: AllocateString("English"),
		},
	}
	a, _, err := client.Languages.List(ctx)
	require.NoError(t, err)
	assert.Equal(t, e, a)
}
