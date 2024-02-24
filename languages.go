package rest

import (
	"context"
	"net/http"
)

type LanguagesService service

type LanguagesResponse struct {
	Data []*Language `json:"data,omitempty"`

	// Ignored because it can be accessed from the response object.
	// Status *int `json:"status,omitempty"`
}

type Language struct {
	LanguageCode *string `json:"language_code,omitempty"`
	LanguageName *string `json:"language_name,omitempty"`
}

// Lists subtitle languages.
//
// [OpenSubtitles Reference]
//
// [OpenSubtitles Reference]: https://opensubtitles.stoplight.io/docs/opensubtitles-api/1de776d20e873-languages
func (s *LanguagesService) List(ctx context.Context) (*LanguagesResponse, *http.Response, error) {
	u, err := s.client.NewURL("infos/languages", nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var r *LanguagesResponse
	res, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, res, err
	}

	return r, res, nil
}
