package rest

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalsAndMarshalsSubtitlesDownloadParameters(t *testing.T) {
	a0 := &SubtitlesDownloadParameters{}
	b0 := "{}"
	equalJSON(t, a0, b0)

	a1 := &SubtitlesDownloadParameters{
		FileID: AllocateID(1),
		FileName: AllocateString("custom"),
		ForceDownload: AllocateBool(false),
		InFPS: AllocateInt(1),
		OutFPS: AllocateInt(1),
		SubFormat: AllocateString("srt"),
		Timeshift: AllocateInt(1),
	}
	b1 := `{
		"file_id": 1,
		"file_name": "custom",
		"force_download": false,
		"in_fps": 1,
		"out_fps": 1,
		"sub_format": "srt",
		"timeshift": 1
	}`
	equalJSON(t, a1, b1)
}

func TestUnmarshalsAndMarshalsSubtitlesDownload(t *testing.T) {
	a0 := &Quota{
		Remaining: 0,
		Requests: 0,
		ResetTime: "",
	}
	b0 := "{}"
	equalJSON(t, a0, b0)

	a1 := &SubtitlesDownloadResponse{
		Quota: a0,
	}
	b1 := "{}"
	equalJSON(t, a1, b1)

	a2 := &SubtitlesDownloadResponse{
		Quota: &Quota{
			Remaining: -1,
			Requests: 21,
			ResetTime: "23 hours and 57 minutes",
			ResetTimeUTC: *AllocateTime("2022-01-30T06:00:53.000Z"),
		},
		FileName: AllocateString("aliens"),
		Link: AllocateString("https://www.opensubtitles.com/"),
	}
	b2 := `{
		"remaining": -1,
		"requests": 21,
		"reset_time": "23 hours and 57 minutes",
		"reset_time_utc": "2022-01-30T06:00:53.000Z",
		"file_name": "aliens",
		"link": "https://www.opensubtitles.com/"
	}`
	equalJSON(t, a2, b2)
}

func TestDownloadsSubtitles(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/download", func (w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/download", r.RequestURI)
		fmt.Fprint(w, `{
			"link": "https://www.opensubtitles.com/"
		}`)
	})

	e := &SubtitlesDownloadResponse{
		Quota: &Quota{
			Remaining: 0,
			Requests: 0,
			ResetTime: "",
		},
		Link: AllocateString("https://www.opensubtitles.com/"),
	}
	ctx := context.Background()
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

func TestEncodesSubtitlesSearchParametersValues(t *testing.T) {
	a0 := &SubtitlesSearchParameters{}
	b0 := ""
	equalQuery(t, a0, b0)

	a1 := &SubtitlesSearchParameters{
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
	b1 := "ai_translated=include&episode_number=0&foreign_parts_only=include&hearing_impaired=include&id=0&imdb_id=0&languages=en%2Cru&machine_translated=exclude&moviehash=b4d8&moviehash_match=include&order_by=download_count&order_direction=asc&page=0&parent_feature_id=0&parent_imdb_id=0&parent_tmdb_id=0&query=friends&season_number=0&tmdb_id=0&trusted_sources=include&type=all&user_id=0&year=1994"
	equalQuery(t, a1, b1)
}

func TestUnmarshalsAndMarshalsSubtitles(t *testing.T) {
	a0 := &Uploader{}
	b0 := "{}"
	equalJSON(t, a0, b0)

	a1 := &RelatedLink{}
	b1 := "{}"
	equalJSON(t, a1, b1)

	a2 := &File{}
	b2 := "{}"
	equalJSON(t, a2, b2)

	a3 := &FeatureDetails{}
	b3 := "{}"
	equalJSON(t, a3, b3)

	a4 := &Subtitle{}
	b4 := "{}"
	equalJSON(t, a4, b4)

	a5 := &SubtitleEntity{}
	b5 := "{}"
	equalJSON(t, a5, b5)

	a6 := &SubtitlesSearchResponse{}
	b6 := "{}"
	equalJSON(t, a6, b6)

	a7 := &SubtitlesSearchResponse{
		Data: []*SubtitleEntity{
			{
				ID: AllocateID(9000),
				Attributes: &Subtitle{
					Language: AllocateString("en"),
					DownloadCount: AllocateInt(697844),
					HearingImpaired: AllocateBool(false),
					HD: AllocateBool(false),
					FPS: AllocateFloat32(23.976),
					Votes: AllocateInt(4),
					Ratings: AllocateFloat32(6),
					FromTrusted: AllocateBool(true),
					ForeignPartsOnly: AllocateBool(false),
					UploadDate: AllocateTime("2009-09-04T19:36:00Z"),
					AITranslated: AllocateBool(false),
					MachineTranslated: AllocateBool(false),
					Release: AllocateString("Season 1 (Whole) DVDrip.XviD-SAiNTS"),
					Uploader: &Uploader{
						UploaderID: AllocateID(47823),
						Name: AllocateString("scooby007"),
						Rank: AllocateString("translator"),
					},
					FeatureDetails: &FeatureDetails{
						FeatureID: AllocateID(38367),
						FeatureType: AllocateString("Episode"),
						Year: AllocateInt(1994),
						Title: AllocateString("The Pilot"),
						MovieName: AllocateString("Friends - S01E01  The Pilot"),
						IMDBID: AllocateID(583459),
						TMDBID: AllocateID(85987),
						SeasonNumber: AllocateInt(1),
						EpisodeNumber: AllocateInt(1),
						ParentIMDBID: AllocateID(108778),
						ParentTitle: AllocateString("Friends"),
						ParentTMDBID: AllocateID(1668),
						ParentFeatureID: AllocateID(7251),
					},
					Files: []*File{
						{
							FileID: AllocateID(1923552),
							FileName: AllocateString("Friends.S01E01.DVDrip.XviD-SAiNTS_(ENGLISH)_DJJ.HOME.SAPO.PT"),
						},
					},
				},
			},
		},
		TotalCount: AllocateInt(1),
	}
	b7 := `{
		"data": [
			{
				"id": "9000",
				"attributes": {
					"language": "en",
					"download_count": 697844,
					"hearing_impaired": false,
					"hd": false,
					"fps": 23.976,
					"votes": 4,
					"ratings": 6,
					"from_trusted": true,
					"foreign_parts_only": false,
					"upload_date": "2009-09-04T19:36:00Z",
					"ai_translated": false,
					"machine_translated": false,
					"release": "Season 1 (Whole) DVDrip.XviD-SAiNTS",
					"uploader": {
						"uploader_id": 47823,
						"name": "scooby007",
						"rank": "translator"
					},
					"feature_details": {
						"feature_id": 38367,
						"feature_type": "Episode",
						"year": 1994,
						"title": "The Pilot",
						"movie_name": "Friends - S01E01  The Pilot",
						"imdb_id": 583459,
						"tmdb_id": 85987,
						"season_number": 1,
						"episode_number": 1,
						"parent_imdb_id": 108778,
						"parent_title": "Friends",
						"parent_tmdb_id": 1668,
						"parent_feature_id": 7251
					},
					"files": [
						{
							"file_id": 1923552,
							"file_name": "Friends.S01E01.DVDrip.XviD-SAiNTS_(ENGLISH)_DJJ.HOME.SAPO.PT"
						}
					]
				}
			}
		],
		"total_count": 1
	}`
	equalJSON(t, a7, b7)
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

	e := &SubtitlesSearchResponse{
		Data: []*SubtitleEntity{
			{
				ID: AllocateID(9000),
			},
		},
	}
	ctx := context.Background()
	p := &SubtitlesSearchParameters{
		AITranslated: AllocateString("include"),
	}
	a, _, err := client.Subtitles.Search(ctx, p)
	require.NoError(t, err)
	assert.Equal(t, e, a)
}
