package rest

import (
	"context"
)

type LanguagesService service

type languagesResponse struct {
	Data []*Language `json:"data,omitempty"`
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
func (s *LanguagesService) List(ctx context.Context) ([]*Language, *Response, error) {
	u, err := s.client.NewURL("infos/languages", nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var r *languagesResponse
	res, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, res, err
	}

	return r.Data, res, nil
}
