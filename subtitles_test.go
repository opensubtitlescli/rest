package rest

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubtitlesDownloadParameters_UnmarshalsAndMarshals(t *testing.T) {
	a := &SubtitlesDownloadParameters{}
	b := "{}"
	equalJSON(t, a, b)

	a = &SubtitlesDownloadParameters{
		FileID: AllocateID(1),
		FileName: AllocateString("custom"),
		ForceDownload: AllocateBool(false),
		InFPS: AllocateInt(1),
		OutFPS: AllocateInt(1),
		SubFormat: AllocateString("srt"),
		Timeshift: AllocateInt(1),
	}
	b = `{
		"file_id": 1,
		"file_name": "custom",
		"force_download": false,
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
		fmt.Fprint(w, `{
			"link": "https://www.opensubtitles.com/"
		}`)
	})

	ctx := context.Background()
	e := &SubtitlesDownloadResponse{
		Link: AllocateString("https://www.opensubtitles.com/"),
	}
	p := &SubtitlesDownloadParameters{
		FileID: AllocateID(1),
		FileName: AllocateString("custom"),
		ForceDownload: AllocateBool(false),
		InFPS: AllocateInt(1),
		OutFPS: AllocateInt(1),
		SubFormat: AllocateString("srt"),
		Timeshift: AllocateInt(1),
	}
	a, _, err := client.Subtitles.Download(ctx, p)
	require.NoError(t, err)
	assert.Equal(t, e, a)
}

func TestSubtitlesSearchParameters_EncodesValues(t *testing.T) {
	a := &SubtitlesSearchParameters{}
	b := ""
	equalQuery(t, a, b)

	a = &SubtitlesSearchParameters{
		AITranslated: AllocateString("include"),
		EpisodeNumber: AllocateInt(0),
		ForeignPartsOnly: AllocateString("include"),
		HearingImpaired: AllocateString("include"),
		ID: AllocateID(0),
		IMDBID: AllocateID(0),
		Languages: AllocateString("en,ru"),
		MachineTranslated: AllocateString("exclude"),
		Moviehash: AllocateString("b4d8"),
		MoviehashMatch: AllocateString("include"),
		OrderBy: AllocateString("download_count"),
		OrderDirection: AllocateString("asc"),
		Page: AllocateInt(0),
		ParentFeatureID: AllocateID(0),
		ParentIMDBID: AllocateID(0),
		ParentTMDBID: AllocateID(0),
		Query: AllocateString("friends"),
		SeasonNumber: AllocateInt(0),
		TMDBID: AllocateID(0),
		TrustedSources: AllocateString("include"),
		Type: AllocateString("all"),
		UserID: AllocateID(0),
		Year: AllocateInt(1994),
	}
	b =
		"ai_translated=include&" +
		"episode_number=0&" +
		"foreign_parts_only=include&" +
		"hearing_impaired=include&" +
		"id=0&" +
		"imdb_id=0&" +
		"languages=en%2Cru&" +
		"machine_translated=exclude&" +
		"moviehash=b4d8&" +
		"moviehash_match=include&" +
		"order_by=download_count&" +
		"order_direction=asc&" +
		"page=0&" +
		"parent_feature_id=0&" +
		"parent_imdb_id=0&" +
		"parent_tmdb_id=0&" +
		"query=friends&" +
		"season_number=0&" +
		"tmdb_id=0&" +
		"trusted_sources=include&" +
		"type=all&" +
		"user_id=0&" +
		"year=1994"
	equalQuery(t, a, b)
}

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

func TestSearchesSubtitles(t *testing.T) {
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
		AITranslated: AllocateString("include"),
	}
	a, _, err := client.Subtitles.Search(ctx, p)
	require.NoError(t, err)
	assert.Equal(t, e, a)
}
