package rest

import (
	"context"
)

type FeaturesService service

type featuresResponse struct {
	Data []*FeatureEntity `json:"data,omitempty"`
}

type FeatureEntity struct {
	Attributes *Feature `json:"attributes,omitempty"`
	ID         *ID      `json:"id,omitempty"`
	Type       *string  `json:"type,omitempty"`
}

type Feature struct {
	EpisodeNumber   *int            `json:"episode_number,omitempty"`
	FeatureID       *ID             `json:"feature_id,omitempty"`
	FeatureType     *string         `json:"feature_type,omitempty"`
	IMDBID          *ID             `json:"imdb_id,omitempty"`
	IMGURL          *string         `json:"img_url,omitempty"`
	OriginalTitle   *string         `json:"original_title,omitempty"`
	ParentIMDBID    *ID             `json:"parent_imdb_id,omitempty"`
	ParentTitle     *string         `json:"parent_title,omitempty"`
	SeasonNumber    *int            `json:"season_number,omitempty"`
	Seasons         *Season         `json:"seasons,omitempty"`
	SeasonsCount    *int            `json:"seasons_count,omitempty"`
	SubtitlesCount  *int            `json:"subtitles_count,omitempty"`
	SubtitlesCounts *map[string]int `json:"subtitles_counts,omitempty"`
	Title           *string         `json:"title,omitempty"`
	TitleAka        []*string       `json:"title_aka,omitempty"`
	TMDBID          *ID             `json:"tmdb_id,omitempty"`
	URL             *string         `json:"url,omitempty"`
	Year            *string         `json:"year,omitempty"`
}

type Season struct {
	Episodes     []*Episode `json:"episodes,omitempty"`
	SeasonNumber *int       `json:"season_number,omitempty"`
}

type Episode struct {
	EpisodeNumber *int    `json:"episode_number,omitempty"`
	FeatureID     *ID     `json:"feature_id,omitempty"`
	FeatureIMDBID *ID     `json:"feature_imdb_id,omitempty"`
	Title         *string `json:"title,omitempty"`
}

type FeaturesPopularParameters struct {
	Languages []string `url:"languages,omitempty" del:","`
	Type      string   `url:"type,omitempty"`
}

// Discovers popular features, according to last 30 days downloads.
//
// [OpenSubtitles Reference]
//
// [OpenSubtitles Reference]: https://opensubtitles.stoplight.io/docs/opensubtitles-api/6d285998026d0-popular-features
func (s *FeaturesService) Popular(ctx context.Context, p *FeaturesPopularParameters) ([]*FeatureEntity, *Response, error) {
	u, err := s.client.NewURL("discover/popular", &p)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var r *featuresResponse
	res, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, res, err
	}

	return r.Data, res, nil
}

type FeaturesSearchParameters struct {
	FeatureID ID     `url:"feature_id,omitempty"`
	IMDBID    ID     `url:"imdb_id,omitempty"`
	Query     string `url:"query,omitempty"`
	TMDBID    ID     `url:"tmdb_id,omitempty"`
	Type      string `url:"type,omitempty"`
	Year      int    `url:"year,omitempty"`
}

// Searches for features.
//
// [OpenSubtitles Reference]
//
// [OpenSubtitles Reference]: https://opensubtitles.stoplight.io/docs/opensubtitles-api/f5eb2608c8fc7-search-for-features
func (s *FeaturesService) Search(ctx context.Context, p *FeaturesSearchParameters) ([]*FeatureEntity, *Response, error) {
	u, err := s.client.NewURL("features", &p)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var r *featuresResponse
	res, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, res, err
	}

	return r.Data, res, nil
}
