package rest

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestFeaturesPopularParameters_EncodesValues(t *testing.T) {
	a := &FeaturesPopularParameters{}
	b := ""
	equalQuery(t, a, b)

	a = &FeaturesPopularParameters{
		Languages: []string{"en", "ru"},
		Type: "all",
	}
	b =
		"languages=en%2Cru&" +
		"type=all"
	equalQuery(t, a, b)
}

func TestFeaturesServicePopular_DiscoversPopularFeatures(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/discover/popular", func (w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/discover/popular?&type=all", r.RequestURI)
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
	p := &FeaturesPopularParameters{
		Type: "all",
	}
	a, _, err := client.Features.Popular(ctx, p)
	require.NoError(t, err)
	assert.Equal(t, e, a)
}

func TestFeaturesServicePopular_ReturnsAnErrorIfCannotCreateAURL(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	ctx := context.Background()

	e := useBadBaseURL(client)
	_, _, a := client.Features.Popular(ctx, nil)
	assert.EqualError(t, a, e)
}

func TestFeaturesServicePopular_ReturnsAUnsuccessfulResponse(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/discover/popular", func (w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	ctx := context.Background()

	_, _, err := client.Features.Popular(ctx, nil)
	var a *ErrorResponse
	require.ErrorAs(t, err, &a)
	assert.Equal(t, a.Response.StatusCode, http.StatusBadRequest)
}

func TestFeaturesSearchParameters_EncodesValues(t *testing.T) {
	a := &FeaturesSearchParameters{}
	b := ""
	equalQuery(t, a, b)

	a = &FeaturesSearchParameters{
		FeatureID: 1,
		IMDBID: 1,
		Query: "hi",
		TMDBID: 1,
		Type: "all",
		Year: 2009,
	}
	b = "feature_id=1&imdb_id=1&query=hi&tmdb_id=1&type=all&year=2009"
	equalQuery(t, a, b)
}

func TestFeaturesServiceSearch_SearchesFeatures(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/features", func (w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/features?&feature_id=1", r.RequestURI)
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
		FeatureID: 1,
	}
	a, _, err := client.Features.Search(ctx, p)
	require.NoError(t, err)
	assert.Equal(t, e, a)
}

func TestFeaturesServiceSearch_ReturnsAnErrorIfCannotCreateAURL(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	ctx := context.Background()

	e := useBadBaseURL(client)
	_, _, a := client.Features.Search(ctx, nil)
	assert.EqualError(t, a, e)
}

func TestFeaturesServiceSearch_ReturnsAUnsuccessfulResponse(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/features", func (w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	ctx := context.Background()

	_, _, err := client.Features.Search(ctx, nil)
	var a *ErrorResponse
	require.ErrorAs(t, err, &a)
	assert.Equal(t, a.Response.StatusCode, http.StatusBadRequest)
}
