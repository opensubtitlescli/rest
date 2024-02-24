package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

	defaultVersion   = "v1"
	defaultUserAgent = "me.vanyauhalin.opensubtitlescli.rest v0.1.0"

	defaultBaseURL = "https://api.opensubtitles.com/api/" + defaultVersion + "/"
	vipBaseURL     = "https://vip-api.opensubtitles.com/api/" + defaultVersion + "/"

	apiKeyHeader  = "Api-Key"
	messageHeader = "X-OpenSubtitles-Message"
)

type Client struct {
	client *http.Client

	apiKey  string
	baseURL *url.URL

	UserAgent string

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

func NewClient(apiKey string) *Client {
	c := &Client{
		client: http.DefaultClient,
		apiKey: apiKey,
	}
	c.init()
	return c
}

func (c *Client) BaseURL() string {
	return c.baseURL.String()
}

func (c *Client) SetBaseURL(u string) error {
	b, err := url.Parse(u)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(b.Path, "/") {
		return fmt.Errorf("rest: base url must have a trailing slash, but %q does not", u)
	}

	c.baseURL = b
	return nil
}

func (c *Client) WithAuthToken(t string) *Client {
	cp := c.copy()
	defer cp.init()

	tr := cp.client.Transport
	if tr == nil {
		tr = http.DefaultTransport
	}

	cp.client.Transport = roundTripperFunc(
		func (req *http.Request) (*http.Response, error) {
			req = req.Clone(req.Context())
			req.Header.Set("Authorization", "Bearer " + t)
			return tr.RoundTrip(req)
		},
	)

	return cp
}

// Creates a RoundTripper (transport).
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}

func (c *Client) init() {
	if c.client == nil {
		c.client = http.DefaultClient
	}
	if c.baseURL == nil {
		c.baseURL, _ = url.Parse(defaultBaseURL)
	}
	if c.UserAgent == "" {
		c.UserAgent = defaultUserAgent
	}
	c.internal.client = c
	c.Auth = (*AuthService)(&c.internal)
	c.Features = (*FeaturesService)(&c.internal)
	c.Formats = (*FormatsService)(&c.internal)
	c.Languages = (*LanguagesService)(&c.internal)
	c.Subtitles = (*SubtitlesService)(&c.internal)
	c.Users = (*UsersService)(&c.internal)
}

func (c *Client) copy() *Client {
	return &Client{
		client: c.client,
		apiKey: c.apiKey,
		baseURL: c.baseURL,
		UserAgent: c.UserAgent,
	}
}

func (c *Client) NewURL(path string, params interface {}) (*url.URL, error) {
	if strings.HasPrefix(path, "/") {
		return nil, fmt.Errorf("rest: url path must not have a leading slash, but %q does", path)
	}

	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	if params == nil {
		return u, nil
	}

	q, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	s := q.Encode()
	if len(s) > 0 {
		// For some unknown reason, some parameters may not come leading.
		u.RawQuery = "&" + s
	}

	return u, nil
}

func (c *Client) NewRequest(method string, url *url.URL, body interface {}) (*http.Request, error) {
	if url == nil {
		return nil, errors.New("rest: url must not be nil")
	}

	var buf io.ReadWriter
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, url.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set(apiKeyHeader, c.apiKey)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface {}) (*http.Response, error) {
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

func (c *Client) BareDo(ctx context.Context, req *http.Request) (*http.Response, error) {
	if ctx == nil {
		return nil, errors.New("rest: context must not be nil")
	}

	req = req.WithContext(ctx)

	res, err := c.client.Do(req)
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

	err = CheckResponse(res)
	if err != nil {
		res.Body.Close()
	}

	return res, nil
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

	// Ignored because it can be accessed from the response object.
	// Status *int `json:"status,omitempty"`
}

func CheckResponse(res *http.Response) error {
	if 200 <= res.StatusCode && res.StatusCode <= 299 {
		return nil
	}

	errRes := &ErrorResponse{
		ResponseError: ResponseError{
			Response: res,
		},
	}

	messages := []string{}

	v := errRes.Response.Header.Get(messageHeader)
	if v != "" {
		messages = append(messages, v)
	}

	data, err := io.ReadAll(res.Body)
	if err == nil && data != nil {
		t := errRes.Response.Header.Get("Content-Type")
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
				Response: errRes.ResponseError.Response,
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
				q := &Quota{}
				qErr := json.Unmarshal(data, q)
				if qErr == nil && q.Remaining < 0 {
					err = &QuotaError{
						ResponseError: *resErr,
						Quota: *q,
					}
				}

			case strings.Contains(m, "invalid or expired link"):
				err = (*LinkError)(resErr)

			case strings.Contains(m, "throttle limit reached"):
				err = (*RateLimitError)(resErr)

			default:
				err = (*ResponseError)(resErr)
			}

			errRes.Errors = append(errRes.Errors, err)
		}
	}

	v = errRes.Response.Request.Header.Get(apiKeyHeader)
	if v == "" {
		err := &APIKeyError{
			Response: errRes.ResponseError.Response,
			Message: "rest: api-key header is empty",
		}
		messages = append(messages, err.Message)
		errRes.Errors = append(errRes.Errors, err)
	}

	v = errRes.Response.Request.Header.Get("User-Agent")
	if v == "" {
		err := &UserAgentError{
			Response: errRes.ResponseError.Response,
			Message: "rest: user-agent header is empty",
		}
		messages = append(messages, err.Message)
		errRes.Errors = append(errRes.Errors, err)
	} else {
		re, err := regexp.Compile(`\S+ v\d+\.\d+\.\d+`)
		if err == nil && !re.MatchString(v) {
			err := &UserAgentError{
				Response: errRes.ResponseError.Response,
				Message: "rest: user-agent is wrong",
			}
			messages = append(messages, err.Message)
			errRes.Errors = append(errRes.Errors, err)
		}
	}

	v = errRes.Response.Request.Header.Get("Authorization")
	if v == "" {
		err := &AuthTokenError{
			Response: errRes.ResponseError.Response,
			Message: "rest: authorization header is empty",
		}
		messages = append(messages, err.Message)
		errRes.Errors = append(errRes.Errors, err)
	} else if !(strings.HasPrefix(v, "Bearer ") && len(v) > len("Bearer ")) {
		err := &AuthTokenError{
			Response: errRes.ResponseError.Response,
			Message: "rest: authorization token is empty",
		}
		messages = append(messages, err.Message)
		errRes.Errors = append(errRes.Errors, err)
	}

	errRes.Message = strings.Join(messages, "; ")

	return errRes
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
	b := struct { *a }{ (*a)(id) }

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
