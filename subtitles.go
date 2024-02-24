package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type SubtitlesService service

type SubtitlesDownloadParameters struct {
	FileID        *ID     `json:"file_id,omitempty"`
	FileName      *string `json:"file_name,omitempty"`
	ForceDownload *bool   `json:"force_download,omitempty"`
	InFPS         *int    `json:"in_fps,omitempty"`
	OutFPS        *int    `json:"out_fps,omitempty"`
	SubFormat     *string `json:"sub_format,omitempty"`
	Timeshift     *int    `json:"timeshift,omitempty"`
}

type SubtitlesDownloadResponse struct {
	Quota    *Quota  `json:"-"`
	FileName *string `json:"file_name,omitempty"`
	Link     *string `json:"link,omitempty"`

	// Quota replaces them without the message. Usually, you do not need to see
	// the message unless you encounter an error.
	// Message      *string    `json:"message,omitempty"`
	// Remaining    *int       `json:"remaining,omitempty"`
	// Requests     *int       `json:"requests,omitempty"`
	// ResetTime    *string    `json:"reset_time,omitempty"`
	// ResetTimeUTC *time.Time `json:"reset_time_utc,omitempty"`

	// Found in the response body, but I do not know what they are for.
	// TS  *int    `json:"ts,omitempty"`
	// UID *int    `json:"uid,omitempty"`
	// UK  *string `json:"uk,omitempty"`
}

func (r *SubtitlesDownloadResponse) UnmarshalJSON(data []byte) error {
	type a SubtitlesDownloadResponse
	b := struct { *a }{ (*a)(r) }

	err := json.Unmarshal(data, b.a)
	if err != nil {
		return err
	}

	var q *Quota
	err = json.Unmarshal(data, &q)
	if err != nil {
		return err
	}

	r.Quota = q
	return nil
}

type Quota struct {
	Remaining    int       `json:"remaining"`
	Requests     int       `json:"requests"`
	ResetTime    string    `json:"reset_time"`
	ResetTimeUTC time.Time `json:"reset_time_utc"`
}

// Requests a download URL for a subtitles.
//
// [OpenSubtitles Reference]
//
// [OpenSubtitles Reference]: https://opensubtitles.stoplight.io/docs/opensubtitles-api/6be7f6ae2d918-download
func (s *SubtitlesService) Download(ctx context.Context, p *SubtitlesDownloadParameters) (*SubtitlesDownloadResponse, *http.Response, error) {
	u, err := s.client.NewURL("download", nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("POST", u, &p)
	if err != nil {
		return nil, nil, err
	}

	buf := &bytes.Buffer{}
	res, err := s.client.Do(ctx, req, buf)
	if err != nil {
		return nil, res, err
	}

	var r *SubtitlesDownloadResponse
	data := buf.Bytes()
	err = json.Unmarshal(data, &r)
	if err != nil {
		return nil, res, err
	}

	return r, res, nil
}

type SubtitlesSearchParameters struct {
	AITranslated      *string `url:"ai_translated,omitempty"`
	EpisodeNumber     *int    `url:"episode_number,omitempty"`
	ForeignPartsOnly  *string `url:"foreign_parts_only,omitempty"`
	HearingImpaired   *string `url:"hearing_impaired,omitempty"`
	ID                *ID     `url:"id,omitempty"`
	IMDBID            *ID     `url:"imdb_id,omitempty"`
	Languages         *string `url:"languages,omitempty"`
	MachineTranslated *string `url:"machine_translated,omitempty"`
	Moviehash         *string `url:"moviehash,omitempty"`
	MoviehashMatch    *string `url:"moviehash_match,omitempty"`
	OrderBy           *string `url:"order_by,omitempty"`
	OrderDirection    *string `url:"order_direction,omitempty"`
	Page              *int    `url:"page,omitempty"`
	ParentFeatureID   *ID     `url:"parent_feature_id,omitempty"`
	ParentIMDBID      *ID     `url:"parent_imdb_id,omitempty"`
	ParentTMDBID      *ID     `url:"parent_tmdb_id,omitempty"`
	Query             *string `url:"query,omitempty"`
	SeasonNumber      *int    `url:"season_number,omitempty"`
	TMDBID            *ID     `url:"tmdb_id,omitempty"`
	TrustedSources    *string `url:"trusted_sources,omitempty"`
	Type              *string `url:"type,omitempty"`
	UserID            *ID     `url:"user_id,omitempty"`
	Year              *int    `url:"year,omitempty"`
}

type SubtitlesSearchResponse struct {
	Data       []*SubtitleEntity `json:"data,omitempty"`
	Page       *int              `json:"page,omitempty"`
	PerPage    *int              `json:"per_page,omitempty"`
	TotalCount *int              `json:"total_count,omitempty"`
	TotalPages *int              `json:"total_pages,omitempty"`
}

type SubtitleEntity struct {
	Attributes *Subtitle `json:"attributes,omitempty"`
	ID         *ID       `json:"id,omitempty"`
	Type       *string   `json:"type,omitempty"`
}

type Subtitle struct {
	AITranslated      *bool           `json:"ai_translated,omitempty"`
	Comments          *string         `json:"comments,omitempty"`
	DownloadCount     *int            `json:"download_count,omitempty"`
	FeatureDetails    *FeatureDetails `json:"feature_details,omitempty"`
	Files             []*File         `json:"files,omitempty"`
	ForeignPartsOnly  *bool           `json:"foreign_parts_only,omitempty"`
	FPS               *float32        `json:"fps,omitempty"`
	FromTrusted       *bool           `json:"from_trusted,omitempty"`
	HD                *bool           `json:"hd,omitempty"`
	HearingImpaired   *bool           `json:"hearing_impaired,omitempty"`
	Language          *string         `json:"language,omitempty"`
	MachineTranslated *bool           `json:"machine_translated,omitempty"`
	NewDownloadCount  *int            `json:"new_download_count,omitempty"`
	Ratings           *float32        `json:"ratings,omitempty"`
	RelatedLinks      []*RelatedLink  `json:"related_links,omitempty"`
	Release           *string         `json:"release,omitempty"`
	SubtitleID        *ID             `json:"subtitle_id,omitempty"`
	UploadDate        *time.Time      `json:"upload_date,omitempty"`
	Uploader          *Uploader       `json:"uploader,omitempty"`
	URL               *string         `json:"url,omitempty"`
	Votes             *int            `json:"votes,omitempty"`
}

type FeatureDetails struct {
	EpisodeNumber   *int    `json:"episode_number,omitempty"`
	FeatureID       *ID     `json:"feature_id,omitempty"`
	FeatureType     *string `json:"feature_type,omitempty"`
	IMDBID          *ID     `json:"imdb_id,omitempty"`
	MovieName       *string `json:"movie_name,omitempty"`
	ParentFeatureID *ID     `json:"parent_feature_id,omitempty"`
	ParentIMDBID    *ID     `json:"parent_imdb_id,omitempty"`
	ParentTitle     *string `json:"parent_title,omitempty"`
	ParentTMDBID    *ID     `json:"parent_tmdb_id,omitempty"`
	SeasonNumber    *int    `json:"season_number,omitempty"`
	Title           *string `json:"title,omitempty"`
	TMDBID          *ID     `json:"tmdb_id,omitempty"`
	Year            *int    `json:"year,omitempty"`
}

type File struct {
	CDNumber *int    `json:"cd_number,omitempty"`
	FileID   *ID     `json:"file_id,omitempty"`
	FileName *string `json:"file_name,omitempty"`
}

type RelatedLink struct {
	IMGURL *string `json:"img_url,omitempty"`
	Label  *string `json:"label,omitempty"`
	URL    *string `json:"url,omitempty"`
}

type Uploader struct {
	Name       *string `json:"name,omitempty"`
	Rank       *string `json:"rank,omitempty"`
	UploaderID *ID     `json:"uploader_id,omitempty"`
}

// Searches for subtitles.
//
// [OpenSubtitles Reference]
//
// [OpenSubtitles Reference]: https://opensubtitles.stoplight.io/docs/opensubtitles-api/a172317bd5ccc-search-for-subtitles
func (s *SubtitlesService) Search(ctx context.Context, p *SubtitlesSearchParameters) (*SubtitlesSearchResponse, *http.Response, error) {
	u, err := s.client.NewURL("subtitles", &p)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var r *SubtitlesSearchResponse
	res, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, res, err
	}

	return r, res, nil
}
