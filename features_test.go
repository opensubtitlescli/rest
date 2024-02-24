package rest

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodesFeaturesSearchParametersValues(t *testing.T) {
	a0 := &FeaturesSearchParameters{}
	b0 := ""
	equalQuery(t, a0, b0)

	a1 := &FeaturesSearchParameters{
		FeatureID: AllocateID(1),
		IMDBID: AllocateID(1),
		Query: AllocateString("hi"),
		TMDBID: AllocateID(1),
		Type: AllocateString("all"),
		Year: AllocateInt(2009),
	}
	b1 := "feature_id=1&imdb_id=1&query=hi&tmdb_id=1&type=all&year=2009"
	equalQuery(t, a1, b1)
}

func TestUnmarshalsAndMarshalsFeature(t *testing.T) {
	a0 := &Feature{}
	b0 := "{}"
	equalJSON(t, a0, b0)

	a1 := &FeatureEntity{}
	b1 := "{}"
	equalJSON(t, a1, b1)

	a2 := &FeaturesSearchResponse{}
	b2 := "{}"
	equalJSON(t, a2, b2)

	a3 := &FeaturesSearchResponse{
		Data: []*FeatureEntity{
			{
				ID: AllocateID(1),
				Attributes: &Feature{
					EpisodeNumber: AllocateInt(1),
					FeatureType: AllocateString("movie"),
					IMDBID: AllocateID(1),
					ParentIMDBID: AllocateID(1),
					ParentTitle: AllocateString("hi"),
					SeasonNumber: AllocateInt(1),
					Title: AllocateString("hola"),
					TMDBID: AllocateID(1),
					Year: AllocateString("2009"),
				},
			},
		},
	}
	b3 := `{
		"data": [
			{
				"id": "1",
				"attributes": {
					"episode_number": 1,
					"feature_type": "movie",
					"imdb_id": 1,
					"parent_imdb_id": 1,
					"parent_title": "hi",
					"season_number": 1,
					"title": "hola",
					"tmdb_id": 1,
					"year": "2009"
				}
			}
		]
	}`
	equalJSON(t, a3, b3)
}

func TestSearchesFeatures(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/features", func (w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/features?&feature_id=0", r.RequestURI)
		fmt.Fprint(w, `{
			"data": [
				{
					"id": "126826"
				}
			]
		}`)
	})

	e := &FeaturesSearchResponse{
		Data: []*FeatureEntity{
			{
				ID: AllocateID(126826),
			},
		},
	}
	ctx := context.Background()
	p := &FeaturesSearchParameters{
		FeatureID: AllocateID(0),
	}
	a, _, err := client.Features.Search(ctx, p)
	require.NoError(t, err)
	assert.Equal(t, e, a)
}
