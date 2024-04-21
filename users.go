package rest

import (
	"context"
)

type UsersService service

type userResponse struct {
	Data *User `json:"data,omitempty"`
}

type User struct {
	AllowedDownloads    *int    `json:"allowed_downloads,omitempty"`
	AllowedTranslations *int    `json:"allowed_translations,omitempty"`
	DownloadsCount      *int    `json:"downloads_count,omitempty"`
	ExtInstalled        *bool   `json:"ext_installed,omitempty"`
	Level               *string `json:"level,omitempty"`
	RemainingDownloads  *int    `json:"remaining_downloads,omitempty"`
	UserID              *ID     `json:"user_id,omitempty"`
	Username            *string `json:"username,omitempty"`
	VIP                 *bool   `json:"vip,omitempty"`
}

// Gets a user's information.
//
// [OpenSubtitles Reference]
//
// [OpenSubtitles Reference]: https://opensubtitles.stoplight.io/docs/opensubtitles-api/ea912bb244ef0-user-informations
func (s *UsersService) Get(ctx context.Context) (*User, *Response, error) {
	u, err := s.client.NewURL("infos/user", nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var r *userResponse
	res, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, res, err
	}

	return r.Data, res, nil
}
