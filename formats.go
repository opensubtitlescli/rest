package rest

import (
	"context"
)

type FormatsService service

type formatsResponse struct {
	Data *FormatsListResponse `json:"data,omitempty"`
}

type FormatsListResponse struct {
	OutputFormats []*string `json:"output_formats,omitempty"`
}

// Lists subtitle formats.
//
// [OpenSubtitles Reference]
//
// [OpenSubtitles Reference]: https://opensubtitles.stoplight.io/docs/opensubtitles-api/69b286fc7506e-subtitle-formats
func (s *FormatsService) List(ctx context.Context) (*FormatsListResponse, *Response, error) {
	u, err := s.client.NewURL("infos/formats", nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var r *formatsResponse
	res, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, res, err
	}

	return r.Data, res, nil
}
