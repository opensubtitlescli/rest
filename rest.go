package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	Version = "0.0.1"

	defaultAPIKey    = ""
	defaultUserAgent = "me.vanyauhalin.opensubtitlescli.rest v" + Version

	defaultVersion = "v1"
	defaultBaseURL = "https://api.opensubtitles.com/api/" + defaultVersion + "/"
	vipBaseURL     = "https://vip-api.opensubtitles.com/api/" + defaultVersion + "/"

	apiKeyHeader  = "Api-Key"
	messageHeader = "X-OpenSubtitles-Message"

	headerRateLimit     = "X-RateLimit-Limit"
	headerRateRemaining = "X-RateLimit-Remaining"
	headerRateReset     = "X-RateLimit-Reset"

	// Found in several endpoint in response headers, but I do not use them.
	// vipHeader         = "X-Vip"
	// vipConsumerHeader = "X-Vip-Consumer"
	// vipUserHeader     = "X-Vip-User"
)

type Client struct {
	client *http.Client

	APIKey    string
	UserAgent string
	BaseURL   *url.URL

	internal service

	Auth      *AuthService
	Features  *FeaturesService
	Formats   *FormatsService
	Languages *LanguagesService
	Subtitles *SubtitlesService
	Users     *UsersService
}

type service struct {
	client *Client
}

func NewClient(client *http.Client) *Client {
	c := &Client{}

	if (client == nil) {
		client := *http.DefaultClient
		c.client = &client
	} else {
		client := *client
		c.client = &client
		if c.client.Transport == nil {
			c.client.Transport = http.DefaultTransport
		}
	}

	c.APIKey = defaultAPIKey
	c.UserAgent = defaultUserAgent
	c.BaseURL, _ = url.Parse(defaultBaseURL)

	c.internal.client = c

	c.Auth = (*AuthService)(&c.internal)
	c.Features = (*FeaturesService)(&c.internal)
	c.Formats = (*FormatsService)(&c.internal)
	c.Languages = (*LanguagesService)(&c.internal)
	c.Subtitles = (*SubtitlesService)(&c.internal)
	c.Users = (*UsersService)(&c.internal)

	return c
}

func (c *Client) Client() *http.Client {
	client := *c.client
	return &client
}

func (c *Client) WithAuthToken(t string) *Client {
	cp := c.copy()
	tr := cp.client.Transport
	cp.client.Transport = roundTripperFunc(
		func (req *http.Request) (*http.Response, error) {
			req = req.Clone(req.Context())
			req.Header.Set("Authorization", "Bearer " + t)
			return tr.RoundTrip(req)
		},
	)
	return cp
}

func (c *Client) copy() *Client {
	cp := NewClient(c.client)

	if c.APIKey != "" {
		cp.APIKey = c.APIKey
	}
	if c.UserAgent != "" {
		cp.UserAgent = c.UserAgent
	}
	if c.BaseURL != nil {
		cp.BaseURL = c.BaseURL
	}

	return cp
}

func (c *Client) NewURL(p string, v interface {}) (string, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return "", fmt.Errorf("rest: base url must have a trailing slash, but %q does not", c.BaseURL)
	}
	if strings.HasPrefix(p, "/") {
		return "", fmt.Errorf("rest: url path must not have a leading slash, but %q does", p)
	}

	u, err := c.BaseURL.Parse(p)
	if err != nil {
		return "", err
	}

	if v == nil {
		return u.String(), nil
	}

	q, err := query.Values(v)
	if err != nil {
		return "", err
	}

	s := q.Encode()
	if len(s) > 0 {
		// For some unknown reason, some parameters may not come leading.
		u.RawQuery = "&" + s
	}

	return u.String(), nil
}

func (c *Client) NewRequest(m string, u string, b interface {}) (*http.Request, error) {
	var body io.ReadWriter
	if b != nil {
		b, err := json.Marshal(b)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(m, u, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set(apiKeyHeader, c.APIKey)

	if b != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

type Response struct {
	*http.Response
	Pagination     Pagination
	Quota          Quota
	Rate           Rate
}

type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}

type Quota struct {
	Remaining    int       `json:"remaining"`
	Requests     int       `json:"requests"`
	ResetTime    string    `json:"reset_time"`
	ResetTimeUTC time.Time `json:"reset_time_utc"`
}

type Rate struct {
	Limit     int
	Remaining int
	Reset     int
}

func newResponse(r *http.Response) *Response {
	res := &Response{Response: r}

	var p Pagination
	var q Quota

	t := r.Header.Get("Content-Type")
	if t != "" && strings.Contains(t, "application/json") {
		data, err := io.ReadAll(r.Body)
		if err == nil && data != nil {
			err = json.Unmarshal(data, &p)
			if err != nil {
				p = Pagination{}
			}
			err = json.Unmarshal(data, &q)
			if err != nil {
				q = Quota{}
			}
		}
		r.Body = io.NopCloser(bytes.NewBuffer(data))
	}

	res.Pagination = p
	res.Quota = q
	res.Rate = parseRate(r)

	return res
}

func parseRate(res *http.Response) Rate {
	var r Rate

	h := res.Header.Get(headerRateLimit)
	if h != "" {
		r.Limit, _ = strconv.Atoi(h)
	}

	h = res.Header.Get(headerRateRemaining)
	if h != "" {
		r.Remaining, _ = strconv.Atoi(h)
	}

	h = res.Header.Get(headerRateReset)
	if h != "" {
		r.Reset, _ = strconv.Atoi(h)
	}

	return r
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface {}) (*Response, error) {
	res, err := c.BareDo(ctx, req)
	if err != nil {
		return res, err
	}
	defer res.Body.Close()

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, res.Body)
	default:
		d := json.NewDecoder(res.Body)
		dErr := d.Decode(v)
		if dErr == io.EOF {
			// Ignore EOF errors caused by empty response body.
			dErr = nil
		}
		if dErr != nil {
			err = dErr
		}
	}

	return res, err
}

func (c *Client) BareDo(ctx context.Context, req *http.Request) (*Response, error) {
	req = req.WithContext(ctx)

	r, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled, the context's
		// error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// continue
		}
		return nil, err
	}

	res := newResponse(r)

	err = CheckResponse(r)
	if err != nil {
		r.Body.Close()
	}

	return res, err
}

type UserAgentError ResponseError

func (e *UserAgentError) Error() string {
	return (*ResponseError)(e).Error()
}

type APIKeyError ResponseError

func (e *APIKeyError) Error() string {
	return (*ResponseError)(e).Error()
}

type AuthTokenError ResponseError

func (e *AuthTokenError) Error() string {
	return (*ResponseError)(e).Error()
}

type CredentialsError ResponseError

func (e *CredentialsError) Error() string {
	return (*ResponseError)(e).Error()
}

type FileError ResponseError

func (e *FileError) Error() string {
	return (*ResponseError)(e).Error()
}

type QuotaError struct {
	ResponseError
	Quota
}

func (e *QuotaError) Error() string {
	return e.ResponseError.Error()
}

type LinkError ResponseError

func (e *LinkError) Error() string {
	return (*ResponseError)(e).Error()
}

type RateLimitError ResponseError

func (e *RateLimitError) Error() string {
	return (*ResponseError)(e).Error()
}

type ErrorResponse struct {
	ResponseError
	Errors        []error
}

func (e *ErrorResponse) Error() string {
	return e.ResponseError.Error()
}

type ResponseError struct {
	Response *http.Response
	Message  string
}

func (e *ResponseError) Error() string {
	if e.Response != nil && e.Response.Request != nil {
		return fmt.Sprintf(
			"%v %v: %d %v",
			e.Response.Request.Method,
			e.Response.Request.URL,
			e.Response.StatusCode,
			e.Message,
		)
	}
	if e.Response != nil {
		return fmt.Sprintf("%d %v", e.Response.StatusCode, e.Message)
	}
	return e.Message
}

type internalError struct {
	Message *string   `json:"message,omitempty"`
	Error   *string   `json:"error,omitempty"`
	Errors  []*string `json:"errors,omitempty"`
}

func CheckResponse(res *http.Response) error {
	if 200 <= res.StatusCode && res.StatusCode <= 299 {
		return nil
	}

	er := &ErrorResponse{
		ResponseError: ResponseError{
			Response: res,
		},
	}

	var messages []string

	v := er.Response.Header.Get(messageHeader)
	if v != "" {
		messages = append(messages, v)
	}

	data, err := io.ReadAll(res.Body)
	if err == nil && data != nil {
		t := er.Response.Header.Get("Content-Type")
		if t != "" {
			switch {
			case strings.Contains(t, "application/json"):
				in := &internalError{}
				err := json.Unmarshal(data, in)
				if err == nil {
					if in.Message != nil {
						messages = append(messages, *in.Message)
					}
					if in.Error != nil {
						messages = append(messages, *in.Error)
					}
					for _, e := range in.Errors {
						if e == nil {
							continue
						}
						messages = append(messages, *e)
					}
				}

			case strings.Contains(t, "text/html"):
				re, err := regexp.Compile(`<title>(.*?)</title>`)
				if err == nil {
					m := string(data)
					m = strings.TrimSpace(m)
					ma := re.FindStringSubmatch(m)
					if len(ma) > 1 {
						m = ma[1]
					}
					messages = append(messages, m)
				}

			case strings.Contains(t, "text/txt"):
				m := string(data)
				m = strings.TrimSpace(m)
				messages = append(messages, m)

			default:
				// continue
			}
		}

		for _, m := range messages {
			resErr := &ResponseError{
				Response: er.ResponseError.Response,
				Message: m,
			}

			m = strings.ToLower(m)
			var err error

			switch {
			case strings.Contains(m, "user-agent header is wrong") ||
				strings.Contains(m, "user-agent header is empty"):
				err = (*UserAgentError)(resErr)

			case strings.Contains(m, "you cannot consume this service"):
				err = (*APIKeyError)(resErr)

			case strings.Contains(m, "invalid token") ||
				strings.Contains(m, "no token"):
				err = (*AuthTokenError)(resErr)

			case strings.Contains(m, "invalid username/password"):
				err = (*CredentialsError)(resErr)

			case strings.Contains(m, "invalid file_id"):
				err = (*FileError)(resErr)

			case strings.Contains(m, "you have downloaded your allowed"):
				var q Quota
				qErr := json.Unmarshal(data, &q)
				if qErr == nil && q.Remaining < 0 {
					err = &QuotaError{
						ResponseError: *resErr,
						Quota: q,
					}
				}

			case strings.Contains(m, "invalid or expired link"):
				err = (*LinkError)(resErr)

			case strings.Contains(m, "throttle limit reached"):
				err = (*RateLimitError)(resErr)

			default:
				err = (*ResponseError)(resErr)
			}

			er.Errors = append(er.Errors, err)
		}
	}
	res.Body = io.NopCloser(bytes.NewBuffer(data))

	v = er.Response.Request.Header.Get(apiKeyHeader)
	if v == "" {
		err := &APIKeyError{
			Response: er.ResponseError.Response,
			Message: "rest: api-key header is empty",
		}
		messages = append(messages, err.Message)
		er.Errors = append(er.Errors, err)
	}

	v = er.Response.Request.Header.Get("User-Agent")
	if v == "" {
		err := &UserAgentError{
			Response: er.ResponseError.Response,
			Message: "rest: user-agent header is empty",
		}
		messages = append(messages, err.Message)
		er.Errors = append(er.Errors, err)
	}

	v = er.Response.Request.Header.Get("Authorization")
	if v == "" {
		err := &AuthTokenError{
			Response: er.ResponseError.Response,
			Message: "rest: authorization header is empty",
		}
		messages = append(messages, err.Message)
		er.Errors = append(er.Errors, err)
	} else if !(strings.HasPrefix(v, "Bearer ") && len(v) > len("Bearer ")) {
		err := &AuthTokenError{
			Response: er.ResponseError.Response,
			Message: "rest: authorization token is empty",
		}
		messages = append(messages, err.Message)
		er.Errors = append(er.Errors, err)
	}

	er.Message = strings.Join(messages, "; ")

	return er
}

// The OpenSubtitles API is inconsistent as it may represent IDs in either
// string or int32. This type helps by casting any ID to int64 for consistency.
type ID int64

func (id *ID) String() string {
	return strconv.FormatInt(int64(*id), 10)
}

func (id *ID) EncodeValues(key string, v *url.Values) error {
	v.Add(key, id.String())
	return nil
}

func (id *ID) UnmarshalJSON(data []byte) error {
	type a ID
	b := struct {*a}{(*a)(id)}

	err := json.Unmarshal(data, b.a)
	if err == nil {
		return nil
	}

	var s string
	err = json.Unmarshal(data, &s)
	if err == nil {
		i, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			*id = ID(i)
			return nil
		}
	}

	return err
}

// Allocates a new ID value to store v and returns a pointer to it.
func AllocateID(v int64) *ID {
	id := ID(v)
	return &id
}

// Allocates a new bool value to store v and returns a pointer to it.
func AllocateBool(v bool) *bool {
	return &v
}

// Allocates a new float32 value to store v and returns a pointer to it.
func AllocateFloat32(v float32) *float32 {
	return &v
}

// Allocates a new int value to store v and returns a pointer to it.
func AllocateInt(v int) *int {
	return &v
}

// Allocates a new string value to store v and returns a pointer to it.
func AllocateString(v string) *string {
	return &v
}

// Allocates a new time.Time value to store v and returns a pointer to it.
func AllocateTime(v string) *time.Time {
	t, _ := time.Parse(time.RFC3339, v)
	return &t
}

// Creates a RoundTripper (transport).
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}
