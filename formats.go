package rest

import (
	"context"
	"net/http"
)

type FormatsService service

type FormatsResponse struct {
	Data *FormatsData `json:"data,omitempty"`

	// Ignored because it can be accessed from the response object.
	// Status *int `json:"status,omitempty"`
}

type FormatsData struct {
	OutputFormats []*string `json:"output_formats,omitempty"`
}

// Lists subtitle formats.
//
// [OpenSubtitles Reference]
//
// [OpenSubtitles Reference]: https://opensubtitles.stoplight.io/docs/opensubtitles-api/69b286fc7506e-subtitle-formats
func (s *FormatsService) List(ctx context.Context) (*FormatsResponse, *http.Response, error) {
	u, err := s.client.NewURL("infos/formats", nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var r *FormatsResponse
	res, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, res, err
	}

	return r, res, nil
}
