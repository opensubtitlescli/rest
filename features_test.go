package rest

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeaturesSearchParameters_EncodesValues(t *testing.T) {
	a := &FeaturesSearchParameters{}
	b := ""
	equalQuery(t, a, b)

	a = &FeaturesSearchParameters{
		FeatureID: AllocateID(1),
		IMDBID: AllocateID(1),
		Query: AllocateString("hi"),
		TMDBID: AllocateID(1),
		Type: AllocateString("all"),
		Year: AllocateInt(2009),
	}
	b = "feature_id=1&imdb_id=1&query=hi&tmdb_id=1&type=all&year=2009"
	equalQuery(t, a, b)
}

func TestFeature_UnmarshalsAndMarshals(t *testing.T) {
	a := &Feature{}
	b := "{}"
	equalJSON(t, a, b)

	a = &Feature{
		EpisodeNumber: AllocateInt(1),
		FeatureType: AllocateString("movie"),
		IMDBID: AllocateID(1),
		ParentIMDBID: AllocateID(1),
		ParentTitle: AllocateString("hi"),
		SeasonNumber: AllocateInt(1),
		Title: AllocateString("hola"),
		TMDBID: AllocateID(1),
		Year: AllocateString("2009"),
	}
	b = `{
		"episode_number": 1,
		"feature_type": "movie",
		"imdb_id": 1,
		"parent_imdb_id": 1,
		"parent_title": "hi",
		"season_number": 1,
		"title": "hola",
		"tmdb_id": 1,
		"year": "2009"
	}`
	equalJSON(t, a, b)
}

func TestFeatureEntity_UnmarshalsAndMarshals(t *testing.T) {
	a := &FeatureEntity{}
	b := "{}"
	equalJSON(t, a, b)

	a = &FeatureEntity{
		Attributes: &Feature{
			EpisodeNumber: AllocateInt(1),
		},
		ID: AllocateID(1),
	}
	b = `{
		"attributes": {
			"episode_number": 1
		},
		"id": 1
	}`
	equalJSON(t, a, b)
}

func TestFeaturesServiceSearch_SearchesFeatures(t *testing.T) {
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

	ctx := context.Background()
	e := []*FeatureEntity{
		{
			ID: AllocateID(126826),
		},
	}
	p := &FeaturesSearchParameters{
		FeatureID: AllocateID(0),
	}
	a, _, err := client.Features.Search(ctx, p)
	require.NoError(t, err)
	assert.Equal(t, e, a)
}
