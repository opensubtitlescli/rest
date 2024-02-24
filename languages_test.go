package rest

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalsAndMarshalsLanguage(t *testing.T) {
	a0 := &Language{}
	b0 := "{}"
	equalJSON(t, a0, b0)

	a1 := &LanguagesResponse{}
	b1 := "{}"
	equalJSON(t, a1, b1)

	a2 := &LanguagesResponse{
		Data: []*Language{
			{
				LanguageCode: AllocateString("en"),
				LanguageName: AllocateString("English"),
			},
		},
	}
	b2 := `{
		"data": [{
			"language_code": "en",
			"language_name": "English"
		}]
	}`
	equalJSON(t, a2, b2)
}

func TestListsLanguages(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/infos/languages", func (w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		fmt.Fprint(w, `{
			"data": [{
				"language_code": "en",
				"language_name": "English"
			}],
			"status": 200
		}`)
	})

	e := &LanguagesResponse{
		Data: []*Language{
			{
				LanguageCode: AllocateString("en"),
				LanguageName: AllocateString("English"),
			},
		},
	}
	ctx := context.Background()
	a, _, err := client.Languages.List(ctx)
	require.NoError(t, err)
	assert.Equal(t, e, a)
}
