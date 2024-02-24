package rest

import (
	"context"
	"encoding/json"
	"net/http"
)

type AuthService service

type Credentials struct {
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
}

type Login struct {
	// Returns the base URL with the protocol and version, which can be utilized
	// with the SetBaseURL method of the Client. It is distinct from the BaseURL
	// because the latter includes the service host without the protocol and API
	// version.
	ClientBaseURL string  `json:"-"`

	BaseURL       *string `json:"base_url,omitempty"`
	Token         *string `json:"token,omitempty"`
	User          *User   `json:"user,omitempty"`

	// Ignored because it can be accessed from the response object.
	// Status *int `json:"status,omitempty"`
}

func (l *Login) UnmarshalJSON(data []byte) error {
	type a Login
	b := &struct { *a }{ (*a)(l) }

	err := json.Unmarshal(data, b.a)
	if err != nil {
		return err
	}

	if b.User != nil && b.User.VIP != nil && *b.User.VIP {
		l.ClientBaseURL = vipBaseURL
	} else {
		l.ClientBaseURL = defaultBaseURL
	}

	return nil
}

// Creates a token to authenticate a user.
//
// [OpenSubtitles Reference]
//
// [OpenSubtitles Reference]: https://opensubtitles.stoplight.io/docs/opensubtitles-api/73acf79accc0a-login
func (s *AuthService) Login(ctx context.Context, c *Credentials) (*Login, *http.Response, error) {
	u, err := s.client.NewURL("login", nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("POST", u, c)
	if err != nil {
		return nil, nil, err
	}

	var l *Login
	res, err := s.client.Do(ctx, req, &l)
	if err != nil {
		return nil, res, err
	}

	return l, res, nil
}

// Destroys a token to end a session.
//
// [OpenSubtitles Reference]
//
// [OpenSubtitles Reference]: https://opensubtitles.stoplight.io/docs/opensubtitles-api/9fe4d6d078e50-logout
func (s *AuthService) Logout(ctx context.Context) (*http.Response, error) {
	u, err := s.client.NewURL("logout", nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
