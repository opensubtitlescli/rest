package rest

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploader_UnmarshalsAndMarshals(t *testing.T) {
	a := &Uploader{}
	b := "{}"
	equalJSON(t, a, b)

	a = &Uploader{
		Name: AllocateString("scooby007"),
		Rank: AllocateString("translator"),
		UploaderID: AllocateID(47823),
	}
	b = `{
		"name": "scooby007",
		"rank": "translator",
		"uploader_id": 47823
	}`
	equalJSON(t, a, b)
}

func TestRelatedLink_UnmarshalsAndMarshals(t *testing.T) {
	a := &RelatedLink{}
	b := "{}"
	equalJSON(t, a, b)

	a = &RelatedLink{
		IMGURL: AllocateString("https://www.opensubtitles.com/"),
		Label: AllocateString("opensubtitles"),
		URL: AllocateString("https://www.opensubtitles.com/"),
	}
	b = `{
		"img_url": "https://www.opensubtitles.com/",
		"label": "opensubtitles",
		"url": "https://www.opensubtitles.com/"
	}`
	equalJSON(t, a, b)
}

func TestFile_UnmarshalsAndMarshals(t *testing.T) {
	a := &File{}
	b := "{}"
	equalJSON(t, a, b)

	a = &File{
		CDNumber: AllocateInt(1),
		FileID: AllocateID(1),
		FileName: AllocateString("aliens"),
	}
	b = `{
		"cd_number": 1,
		"file_id": 1,
		"file_name": "aliens"
	}`
	equalJSON(t, a, b)
}

func TestFeatureDetails_UnmarshalsAndMarshals(t *testing.T) {
	a := &FeatureDetails{}
	b := "{}"
	equalJSON(t, a, b)

	a = &FeatureDetails{
		EpisodeNumber: AllocateInt(1),
		FeatureID: AllocateID(38367),
		FeatureType: AllocateString("Episode"),
		IMDBID: AllocateID(583459),
		MovieName: AllocateString("Friends - S01E01  The Pilot"),
		ParentFeatureID: AllocateID(7251),
		ParentIMDBID: AllocateID(108778),
		ParentTitle: AllocateString("Friends"),
		ParentTMDBID: AllocateID(1668),
		SeasonNumber: AllocateInt(1),
		Title: AllocateString("The Pilot"),
		TMDBID: AllocateID(85987),
		Year: AllocateInt(1994),
	}
	b = `{
		"episode_number": 1,
		"feature_id": 38367,
		"feature_type": "Episode",
		"imdb_id": 583459,
		"movie_name": "Friends - S01E01  The Pilot",
		"parent_feature_id": 7251,
		"parent_imdb_id": 108778,
		"parent_title": "Friends",
		"parent_tmdb_id": 1668,
		"season_number": 1,
		"title": "The Pilot",
		"tmdb_id": 85987,
		"year": 1994
	}`
	equalJSON(t, a, b)
}

func TestSubtitle_UnmarshalsAndMarshals(t *testing.T) {
	a := &Subtitle{}
	b := "{}"
	equalJSON(t, a, b)

	a = &Subtitle{
		AITranslated: AllocateBool(false),
		DownloadCount: AllocateInt(697844),
		FeatureDetails: &FeatureDetails{
			EpisodeNumber: AllocateInt(1),
		},
		Files: []*File{
			{
				FileID: AllocateID(1923552),
			},
		},
		ForeignPartsOnly: AllocateBool(false),
		FPS: AllocateFloat32(23.976),
		FromTrusted: AllocateBool(true),
		HD: AllocateBool(false),
		HearingImpaired: AllocateBool(false),
		Language: AllocateString("en"),
		MachineTranslated: AllocateBool(false),
		Ratings: AllocateFloat32(6),
		Release: AllocateString("Season 1 (Whole) DVDrip.XviD-SAiNTS"),
		UploadDate: AllocateTime("2009-09-04T19:36:00Z"),
		Uploader: &Uploader{
			Name: AllocateString("scooby007"),
		},
		Votes: AllocateInt(4),
	}
	b = `{
		"ai_translated": false,
		"download_count": 697844,
		"feature_details": {
			"episode_number": 1
		},
		"files": [
			{
				"file_id": 1923552
			}
		],
		"foreign_parts_only": false,
		"fps": 23.976,
		"from_trusted": true,
		"hd": false,
		"hearing_impaired": false,
		"language": "en",
		"machine_translated": false,
		"ratings": 6,
		"release": "Season 1 (Whole) DVDrip.XviD-SAiNTS",
		"upload_date": "2009-09-04T19:36:00Z",
		"uploader": {
			"name": "scooby007"
		},
		"votes": 4
	}`
	equalJSON(t, a, b)
}

func TestSubtitleEntity_UnmarshalsAndMarshals(t *testing.T) {
	a := &SubtitleEntity{}
	b := "{}"
	equalJSON(t, a, b)

	a = &SubtitleEntity{
		Attributes: &Subtitle{
			AITranslated: AllocateBool(false),
		},
		ID: AllocateID(9000),
	}
	b = `{
		"attributes": {
			"ai_translated": false
		},
		"id": 9000
	}`
	equalJSON(t, a, b)
}

func TestSubtitlesDownloadParameters_UnmarshalsAndMarshals(t *testing.T) {
	a := &SubtitlesDownloadParameters{}
	b := "{}"
	equalJSON(t, a, b)

	a = &SubtitlesDownloadParameters{
		FileID: 1,
		FileName: "custom",
		ForceDownload: true,
		InFPS: 1,
		OutFPS: 1,
		SubFormat: "srt",
		Timeshift: 1,
	}
	b = `{
		"file_id": 1,
		"file_name": "custom",
		"force_download": true,
		"in_fps": 1,
		"out_fps": 1,
		"sub_format": "srt",
		"timeshift": 1
	}`
	equalJSON(t, a, b)
}

func TestSubtitlesDownloadResponse_UnmarshalsAndMarshals(t *testing.T) {
	a := &SubtitlesDownloadResponse{}
	b := "{}"
	equalJSON(t, a, b)

	a = &SubtitlesDownloadResponse{
		FileName: AllocateString("aliens"),
		Link: AllocateString("https://www.opensubtitles.com/"),
	}
	b = `{
		"file_name": "aliens",
		"link": "https://www.opensubtitles.com/"
	}`
	equalJSON(t, a, b)
}

func TestSubtitlesServiceDownload_DownloadsSubtitles(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/download", func (w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/download", r.RequestURI)
		equalBody(t, r.Body, `{
			"file_id": 1
		}`)
		fmt.Fprint(w, `{
			"link": "https://www.opensubtitles.com/"
		}`)
	})

	ctx := context.Background()

	e := &SubtitlesDownloadResponse{
		Link: AllocateString("https://www.opensubtitles.com/"),
	}
	p := &SubtitlesDownloadParameters{
		FileID: 1,
	}
	a, _, err := client.Subtitles.Download(ctx, p)
	require.NoError(t, err)
	assert.Equal(t, e, a)
}

func TestSubtitlesServiceDownload_ReturnsAnErrorIfCannotCreateAURL(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	ctx := context.Background()

	e := useBadBaseURL(client)
	_, _, a := client.Subtitles.Download(ctx, nil)
	assert.EqualError(t, a, e)
}

func TestSubtitlesServiceDownload_ReturnsAUnsuccessfulResponse(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/download", func (w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	ctx := context.Background()

	_, _, err := client.Subtitles.Download(ctx, nil)
	var a *ErrorResponse
	require.ErrorAs(t, err, &a)
	assert.Equal(t, a.Response.StatusCode, http.StatusBadRequest)
}

func TestSubtitlesLatestParameters_EncodesValues(t *testing.T) {
	a := &SubtitlesLatestParameters{}
	b := ""
	equalQuery(t, a, b)

	a = &SubtitlesLatestParameters{
		Languages: []string{"en", "ru"},
		Type: "all",
	}
	b =
		"languages=en%2Cru&" +
		"type=all"
	equalQuery(t, a, b)
}

func TestSubtitlesServiceLatest_DiscoversLatestSubtitles(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/discover/latest", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/discover/latest?&type=all", r.RequestURI)
		fmt.Fprint(w, `{
			"data": [
				{
					"id": "9000"
				}
			]
		}`)
	})

	ctx := context.Background()

	e := []*SubtitleEntity{
		{
			ID: AllocateID(9000),
		},
	}
	p := &SubtitlesLatestParameters{
		Type: "all",
	}
	a, _, err := client.Subtitles.Latest(ctx, p)
	require.NoError(t, err)
	assert.Equal(t, e, a)
}

func TestSubtitlesServiceLatest_ReturnsAnErrorIfCannotCreateAURL(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	ctx := context.Background()

	e := useBadBaseURL(client)
	_, _, a := client.Subtitles.Latest(ctx, nil)
	assert.EqualError(t, a, e)
}

func TestSubtitlesServiceLatest_ReturnsAUnsuccessfulResponse(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/discover/latest", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	ctx := context.Background()

	_, _, err := client.Subtitles.Latest(ctx, nil)
	var a *ErrorResponse
	require.ErrorAs(t, err, &a)
	assert.Equal(t, a.Response.StatusCode, http.StatusBadRequest)
}

func TestSubtitlesPopularParameters_EncodesValues(t *testing.T) {
	a := &SubtitlesPopularParameters{}
	b := ""
	equalQuery(t, a, b)

	a = &SubtitlesPopularParameters{
		Languages: []string{"en", "ru"},
		Type: "all",
	}
	b =
		"languages=en%2Cru&" +
		"type=all"
	equalQuery(t, a, b)
}

func TestSubtitlesServicePopular_DiscoversPopularSubtitles(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/discover/most_downloaded", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/discover/most_downloaded?&type=all", r.RequestURI)
		fmt.Fprint(w, `{
			"data": [
				{
					"id": "9000"
				}
			]
		}`)
	})

	ctx := context.Background()

	e := []*SubtitleEntity{
		{
			ID: AllocateID(9000),
		},
	}
	p := &SubtitlesPopularParameters{
		Type: "all",
	}
	a, _, err := client.Subtitles.Popular(ctx, p)
	require.NoError(t, err)
	assert.Equal(t, e, a)
}

func TestSubtitlesServicePopular_ReturnsAnErrorIfCannotCreateAURL(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	ctx := context.Background()

	e := useBadBaseURL(client)
	_, _, a := client.Subtitles.Popular(ctx, nil)
	assert.EqualError(t, a, e)
}

func TestSubtitlesServicePopular_ReturnsAUnsuccessfulResponse(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/discover/most_downloaded", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	ctx := context.Background()

	_, _, err := client.Subtitles.Popular(ctx, nil)
	var a *ErrorResponse
	require.ErrorAs(t, err, &a)
	assert.Equal(t, a.Response.StatusCode, http.StatusBadRequest)
}

func TestSubtitlesSearchParameters_EncodesValues(t *testing.T) {
	a := &SubtitlesSearchParameters{}
	b := ""
	equalQuery(t, a, b)

	a = &SubtitlesSearchParameters{
		AITranslated: "include",
		EpisodeNumber: 1,
		ForeignPartsOnly: "include",
		HearingImpaired: "include",
		ID: 1,
		IMDBID: 1,
		Languages: []string{"en", "ru"},
		MachineTranslated: "exclude",
		Moviehash: "b4d8",
		MoviehashMatch: "include",
		OrderBy: "download_count",
		OrderDirection: "asc",
		Page: 1,
		ParentFeatureID: 1,
		ParentIMDBID: 1,
		ParentTMDBID: 1,
		Query: "friends",
		SeasonNumber: 1,
		TMDBID: 1,
		TrustedSources: "include",
		Type: "all",
		UserID: 1,
		Year: 1994,
	}
	b =
		"ai_translated=include&" +
		"episode_number=1&" +
		"foreign_parts_only=include&" +
		"hearing_impaired=include&" +
		"id=1&" +
		"imdb_id=1&" +
		"languages=en%2Cru&" +
		"machine_translated=exclude&" +
		"moviehash=b4d8&" +
		"moviehash_match=include&" +
		"order_by=download_count&" +
		"order_direction=asc&" +
		"page=1&" +
		"parent_feature_id=1&" +
		"parent_imdb_id=1&" +
		"parent_tmdb_id=1&" +
		"query=friends&" +
		"season_number=1&" +
		"tmdb_id=1&" +
		"trusted_sources=include&" +
		"type=all&" +
		"user_id=1&" +
		"year=1994"
	equalQuery(t, a, b)
}

func TestSubtitlesServiceSearch_SearchesSubtitles(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/subtitles", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/subtitles?&ai_translated=include", r.RequestURI)
		fmt.Fprint(w, `{
			"data": [
				{
					"id": "9000"
				}
			]
		}`)
	})

	ctx := context.Background()

	e := []*SubtitleEntity{
		{
			ID: AllocateID(9000),
		},
	}
	p := &SubtitlesSearchParameters{
		AITranslated: "include",
	}
	a, _, err := client.Subtitles.Search(ctx, p)
	require.NoError(t, err)
	assert.Equal(t, e, a)
}

func TestSubtitlesServiceSearch_ReturnsAnErrorIfCannotCreateAURL(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	ctx := context.Background()

	e := useBadBaseURL(client)
	_, _, a := client.Subtitles.Search(ctx, nil)
	assert.EqualError(t, a, e)
}

func TestSubtitlesServiceSearch_ReturnsAUnsuccessfulResponse(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/subtitles", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	ctx := context.Background()

	_, _, err := client.Subtitles.Search(ctx, nil)
	var a *ErrorResponse
	require.ErrorAs(t, err, &a)
	assert.Equal(t, a.Response.StatusCode, http.StatusBadRequest)
}
