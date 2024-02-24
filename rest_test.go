package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-querystring/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitializesClient(t *testing.T) {
	client := NewClient("")
	assert.Equal(t, defaultBaseURL, client.BaseURL())
	assert.Equal(t, defaultUserAgent, client.UserAgent)
	assert.NotNil(t, client.Auth)
	assert.NotNil(t, client.Features)
	assert.NotNil(t, client.Formats)
	assert.NotNil(t, client.Languages)
	assert.NotNil(t, client.Subtitles)
	assert.NotNil(t, client.Users)
}

func TestSetsBaseURL(t *testing.T) {
	client := NewClient("")
	e := "http://localhost/v2/"

	err := client.SetBaseURL(e)
	require.NoError(t, err)

	a := client.BaseURL()
	assert.Equal(t, e, a)
}

func TestReturnsAnErrorIfItFailsToParseBaseURL(t *testing.T) {
	client := NewClient("")
	u := ":"
	e := `parse "` + u + `": missing protocol scheme`
	a := client.SetBaseURL(u)
	assert.EqualError(t, a, e)
}

func TestReturnsAnErrorIfBaseURLEndsNotHaveATrailingSlash(t *testing.T) {
	client := NewClient("")
	u := "http://localhost/v2"
	e := `rest: base url must have a trailing slash, but "` + u + `" does not`
	a := client.SetBaseURL(u)
	assert.EqualError(t, a, e)
}

func TestAppliesAuthToken(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	client := NewClient("")
	token := "xxx"
	e := "Bearer " + token

	mux.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		a := r.Header.Get("Authorization")
		assert.Equal(t, e, a)
	})

	client = client.WithAuthToken(token)
	_, err := client.client.Get(server.URL)
	require.NoError(t, err)
}

func TestInitializesURL(t *testing.T) {
	client := NewClient("")
	e := client.BaseURL()

	u, err := client.NewURL("", nil)
	require.NoError(t, err)

	a := u.String()
	assert.Equal(t, e, a)
}

func TestInitializesURLWithPath(t *testing.T) {
	client := NewClient("")
	p := "s"
	e := client.BaseURL() + p

	u, err := client.NewURL(p, nil)
	require.NoError(t, err)

	a := u.String()
	assert.Equal(t, e, a)
}

func TestInitializesURLWithParameters(t *testing.T) {
	client := NewClient("")
	p := mockStruct()
	e := client.BaseURL() + mockQuery()

	u, err := client.NewURL("", p)
	require.NoError(t, err)

	a := u.String()
	assert.Equal(t, e, a)
}

func TestInitializesURLWithOmittedParameters(t *testing.T) {
	client := NewClient("")
	p := &mockParameters{}
	e := client.BaseURL()

	u, err := client.NewURL("", p)
	require.NoError(t, err)

	a := u.String()
	assert.Equal(t, e, a)
}

func TestInitializesURLWithPathAndParameters(t *testing.T) {
	client := NewClient("")
	path := "s"
	params := mockStruct()
	e := client.BaseURL() + path + mockQuery()

	u, err := client.NewURL(path, params)
	require.NoError(t, err)

	a := u.String()
	assert.Equal(t, e, a)
}

func TestInitializesURLWithPathAndOmittedParameters(t *testing.T) {
	client := NewClient("")
	path := "s"
	params := &mockParameters{}
	e := client.BaseURL() + path

	u, err := client.NewURL(path, params)
	require.NoError(t, err)

	a := u.String()
	assert.Equal(t, e, a)
}

func TestReturnsAnErrorIfItFailsToParseParameters(t *testing.T) {
	client := NewClient("")
	e := "query: Values() expects struct input. Got int"
	_, a := client.NewURL("", 0)
	assert.EqualError(t, a, e)
}

func TestReturnsAnErrorIfURLPathHaveALeadingSlash(t *testing.T) {
	client := NewClient("")
	p := "/"
	e := `rest: url path must not have a leading slash, but "` + p + `" does`
	_, a := client.NewURL(p, nil)
	assert.EqualError(t, a, e)
}

func TestInitializesRequest(t *testing.T) {
	apiKey := "xxx"
	client := NewClient(apiKey)

	u, err := client.NewURL("", nil)
	require.NoError(t, err)

	req, err := client.NewRequest("", u, nil)
	require.NoError(t, err)

	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, u.String(), req.URL.String())
	assert.Equal(t, "application/json", req.Header.Get("Accept"))
	assert.Equal(t, apiKey, req.Header.Get("Api-Key"))
	assert.Equal(t, "me.vanyauhalin.opensubtitlescli.rest v0.1.0", req.Header.Get("User-Agent"))
}

func TestInitializesRequestWithMethod(t *testing.T) {
	client := NewClient("")

	u, err := client.NewURL("", nil)
	require.NoError(t, err)

	req, err := client.NewRequest("POST", u, nil)
	require.NoError(t, err)

	assert.Equal(t, "POST", req.Method)
}

func TestInitializesRequestWithBody(t *testing.T) {
	client := NewClient("")

	e := mockStruct()

	u, err := client.NewURL("", nil)
	require.NoError(t, err)

	req, err := client.NewRequest("", u, e)
	require.NoError(t, err)

	d := json.NewDecoder(req.Body)
	a := &mockParameters{}
	err = d.Decode(a)
	require.NoError(t, err)

	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
	assert.Equal(t, e, a)
}

func TestInitializesRequestWithAll(t *testing.T) {
	client := NewClient("")

	e := mockStruct()

	u, err := client.NewURL("", nil)
	require.NoError(t, err)

	req, err := client.NewRequest("POST", u, e)
	require.NoError(t, err)

	d := json.NewDecoder(req.Body)
	a := &mockParameters{}
	err = d.Decode(a)
	require.NoError(t, err)

	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
	assert.Equal(t, e, a)
}

func TestReturnsAnErrorIfItFailsMarshalBody(t *testing.T) {
	client := NewClient("")

	type body struct {
		A map[interface{}]interface{}
	}

	_, err := client.NewRequest("", &url.URL{}, &body{})
	_, ok := err.(*json.UnsupportedTypeError)
	assert.True(t, ok)
}

func TestReturnsAnErrorIfURLIsNil(t *testing.T) {
	client := NewClient("")
	e := "rest: url must not be nil"
	_, a := client.NewRequest("", nil, nil)
	assert.EqualError(t, a, e)
}

func TestReturnsAnErrorIfMethodIsInvalid(t *testing.T) {
	client := NewClient("")

	m := "\xbd"
	e := "net/http: invalid method \"\\xbd\""

	_, a := client.NewRequest(m, &url.URL{}, nil)
	assert.EqualError(t, a, e)
}

func TestReturnsNilIfResponseSuccess(t *testing.T) {
	for s := 200; s <= 299; s += 1 {
		res := &http.Response{
			StatusCode: s,
		}
		err := CheckResponse(res)
		assert.Nil(t, err)
	}
}

func TestReturnsErrorIfResponseFailure(t *testing.T) {
	toBody := func (s string) io.ReadCloser {
		r := strings.NewReader(s)
		return io.NopCloser(r)
	}
	toMessage := func (s string) io.ReadCloser {
		return toBody(`{ "message": "` + s + `" }`)
	}
	newRes := func () *http.Response {
		return &http.Response{
			Request: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Scheme: "http",
					Host: "localhost",
					Path: "/",
				},
				Header: http.Header{
					apiKeyHeader: {"xxx"},
					"User-Agent": {"app v0.0.0"},
					"Authorization": {"Bearer yyy"},
				},
			},
			StatusCode: http.StatusBadRequest,
			Header: http.Header{
				"Content-Type": {"application/json"},
			},
			Body: toBody(""),
		}
	}
	newErr := func (res *http.Response, m string) *ErrorResponse {
		return &ErrorResponse{
			ResponseError: ResponseError{
				Response: res,
				Message: m,
			},
		}
	}
	try := func (res *http.Response) *ErrorResponse {
		err := CheckResponse(res)
		require.NotNil(t, err)
		var errRes *ErrorResponse
		if !errors.As(err, &errRes) {
			require.Fail(t, "error is not ErrorResponse")
		}
		return errRes
	}
	equal := func (e *ErrorResponse, res *http.Response) {
		a := try(res)
		assert.Equal(t, e, a)
	}

	res := newRes()
	e := newErr(res, "User-Agent header is empty; set it to App name with version eg: MyApp v1.2.3")
	res.Header.Set(messageHeader, e.Message)
	e.Errors = append(e.Errors, (*UserAgentError)(&e.ResponseError))
	equal(e, res)

	res = newRes()
	e = newErr(res, "")
	re0 := &AuthTokenError{
		Response: res,
		Message: "Invalid token a1",
	}
	e.Errors = append(e.Errors, re0)
	re1 := &AuthTokenError{
		Response: res,
		Message: "No token in request",
	}
	e.Errors = append(e.Errors, re1)
	re2 := &ResponseError{
		Response: res,
		Message: "Query is too short",
	}
	e.Errors = append(e.Errors, re2)
	reM := []string{re0.Message, re1.Message, re2.Message}
	e.Message = strings.Join(reM, "; ")
	reD, reDErr := json.Marshal(reM)
	require.NoError(t, reDErr)
	res.Body = toBody(`{ "errors": ` + string(reD) + `}`)
	equal(e, res)

	res = newRes()
	e = newErr(res, "Internal Server Error")
	res.Body = toBody(`{ "error": "` + e.Message + `" }`)
	e.Errors = append(e.Errors, (*ResponseError)(&e.ResponseError))
	equal(e, res)

	res = newRes()
	e = newErr(res, "Not enough parameters")
	res.Body = toMessage(e.Message)
	e.Errors = append(e.Errors, (*ResponseError)(&e.ResponseError))
	equal(e, res)

	res = newRes()
	e = newErr(res, "Error, invalid username/password")
	res.Body = toMessage(e.Message)
	e.Errors = append(e.Errors, (*CredentialsError)(&e.ResponseError))
	equal(e, res)

	res = newRes()
	e = newErr(res, "You cannot consume this service")
	res.Body = toMessage(e.Message)
	e.Errors = append(e.Errors, (*APIKeyError)(&e.ResponseError))
	equal(e, res)

	res = newRes()
	e = newErr(res, "Invalid file_id")
	res.Body = toMessage(e.Message)
	e.Errors = append(e.Errors, (*FileError)(&e.ResponseError))
	equal(e, res)

	res = newRes()
	e = newErr(res, "You have downloaded your allowed 20 subtitles for 24h.Your quota will be renewed in 23 hours and 57 minutes (2022-01-30 06:00:53 UTC) ")
	q := &QuotaError{
		ResponseError: e.ResponseError,
		Quota: Quota{
			Remaining: -1,
			Requests: 21,
			ResetTime: "23 hours and 57 minutes",
			ResetTimeUTC: *AllocateTime("2022-01-30T06:00:53.000Z"),
		},
	}
	res.Body = toBody(
		fmt.Sprintf(
			`{
				"message": "%s",
				"remaining": %d,
				"requests": %d,
				"reset_time": "%s",
				"reset_time_utc": "%s"
			}`,
			q.Message,
			q.Remaining,
			q.Requests,
			q.ResetTime,
			q.ResetTimeUTC.Format("2006-01-02T15:04:05.000Z"),
		),
	)
	e.Errors = append(e.Errors, q)
	equal(e, res)

	res = newRes()
	e = newErr(res, "invalid token")
	res.Body = toMessage(e.Message)
	e.Errors = append(e.Errors, (*AuthTokenError)(&e.ResponseError))
	equal(e, res)

	res = newRes()
	e = newErr(res, "Throttle limit reached. Retry later.")
	res.Body = toMessage(e.Message)
	e.Errors = append(e.Errors, (*RateLimitError)(&e.ResponseError))
	equal(e, res)

	res = newRes()
	e = newErr(res, "Error 403 - User-Agent header is empty; set it to App name with version eg: MyApp v1.2.3")
	res.Header.Set("Content-Type", "text/html")
	res.Body = toBody(`<!DOCTYPE html><html><head><title>` + e.Message + `</title></head></html>`)
	e.Errors = append(e.Errors, (*UserAgentError)(&e.ResponseError))
	equal(e, res)

	res = newRes()
	e = newErr(res, "Error 410 - Invalid or expired link")
	res.Header.Set("Content-Type", "text/txt")
	res.Body = toBody(e.Message)
	e.Errors = append(e.Errors, (*LinkError)(&e.ResponseError))
	equal(e, res)

	res = newRes()
	e = newErr(res, "rest: api-key header is empty")
	res.Request.Header.Del(apiKeyHeader)
	e.Errors = append(e.Errors, (*APIKeyError)(&e.ResponseError))
	equal(e, res)

	res = newRes()
	e = newErr(res, "rest: user-agent header is empty")
	res.Request.Header.Set("User-Agent", "")
	e.Errors = append(e.Errors, (*UserAgentError)(&e.ResponseError))
	equal(e, res)

	res = newRes()
	e = newErr(res, "rest: user-agent is wrong")
	res.Request.Header.Set("User-Agent", "xxx")
	e.Errors = append(e.Errors, (*UserAgentError)(&e.ResponseError))
	equal(e, res)

	res = newRes()
	e = newErr(res, "rest: authorization header is empty")
	res.Request.Header.Del("Authorization")
	e.Errors = append(e.Errors, (*AuthTokenError)(&e.ResponseError))
	equal(e, res)

	res = newRes()
	e = newErr(res, "rest: authorization token is empty")
	res.Request.Header.Set("Authorization", "Bearer ")
	e.Errors = append(e.Errors, (*AuthTokenError)(&e.ResponseError))
	equal(e, res)

	res = newRes()
	e = newErr(res, "")
	e0 := &UserAgentError{
		Response: res,
		Message: "User-Agent header is empty; set it to App name with version eg: MyApp v1.2.3",
	}
	res.Header.Set(messageHeader, e0.Message)
	e.Errors = append(e.Errors, e0)
	e1 := &ResponseError{
		Response: res,
		Message: "Not enough parameters",
	}
	res.Body = toMessage(e1.Message)
	e.Errors = append(e.Errors, e1)
	e2 := &APIKeyError{
		Response: res,
		Message: "rest: api-key header is empty",
	}
	res.Request.Header.Del(apiKeyHeader)
	e.Errors = append(e.Errors, e2)
	e.Message = strings.Join([]string{e0.Message, e1.Message, e2.Message}, "; ")
	equal(e, res)
}

func TestDoesBare(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	ctx := context.Background()
	e := http.StatusOK

	mux.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(e)
	})

	u, err := client.NewURL("", nil)
	require.NoError(t, err)

	r, err := client.NewRequest("", u, nil)
	require.NoError(t, err)

	res, err := client.BareDo(ctx, r)
	require.NoError(t, err)

	a := res.StatusCode
	assert.Equal(t, e, a)
}

func TestReturnsAnErrorIfContextIsNil(t *testing.T) {
	client := NewClient("")
	e := "rest: context must not be nil"
	_, a := client.BareDo(nil, &http.Request{}) // nolint:staticcheck
	assert.EqualError(t, a, e)
}

func TestDoes(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	ctx := context.Background()
	e := mockStruct()

	mux.HandleFunc("/e", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockJSON())
	})

	u, err := client.NewURL("e", nil)
	require.NoError(t, err)

	r, err := client.NewRequest("GET", u, nil)
	require.NoError(t, err)

	a := &mockParameters{}
	_, err = client.Do(ctx, r, &a)
	require.NoError(t, err)

	assert.Equal(t, e, a)
}

func TestEncodesIDValues(t *testing.T) {
	type s struct {
		ID ID `url:"id"`
	}

	a0 := s{0}
	b0 := `id=0`
	equalQuery(t, a0, b0)
}

func TestUnmarshalsAndMarshalsID(t *testing.T) {
	type s struct {
		ID ID `json:"id"`
	}

	a0 := &s{0}
	b0 := `{
		"id": 0
	}`
	equalJSON(t, a0, b0)

	a1 := &s{1}
	b1 := `{
		"id": "1"
	}`
	equalJSON(t, a1, b1)
}

func TestAllocatesID(t *testing.T) {
	e := ID(0)
	a := AllocateID(0)
	assert.Equal(t, e, *a)
}

func setup() (*Client, *http.ServeMux, func ()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	client := NewClient("")
	_ = client.SetBaseURL(server.URL + "/")

	return client, mux, server.Close
}

type mockParameters struct {
	A *int    `url:"a,omitempty" json:"a,omitempty"`
	B *int    `url:"b,omitempty" json:"b,omitempty"`
	C *string `url:"c,omitempty" json:"c,omitempty"`
}

func mockStruct() *mockParameters {
	return &mockParameters{
		A: AllocateInt(0),
		B: AllocateInt(1),
		C: AllocateString("a b c"),
	}
}

func mockJSON() string {
	return `{
		"a": 0,
		"b": 1,
		"c": "a b c"
	}`
}

func mockQuery() string {
	return "?&a=0&b=1&c=a+b+c"
}

func equalQuery(t *testing.T, a interface {}, b string) {
	q, err := query.Values(a)
	require.NoError(t, err)
	assert.Equal(t, q.Encode(), b)
}

func equalJSON(t *testing.T, a interface {}, b string) {
	s := reflect.New(reflect.TypeOf(a).Elem()).Interface()
	err := json.Unmarshal([]byte(b), &s)
	require.NoError(t, err)

	assert.Equal(t, a, s)

	x, err := json.MarshalIndent(a, "", "  ")
	require.NoError(t, err)

	y, err := json.MarshalIndent(s, "", "  ")
	require.NoError(t, err)

	assert.Equal(t, string(x), string(y))

	// i := reflect.New(reflect.TypeOf(a)).Interface()
	// err := json.Unmarshal([]byte(b), &i)
	// require.NoError(t, err)

	// x, err := json.MarshalIndent(a, "", "  ")
	// require.NoError(t, err)

	// y, err := json.MarshalIndent(i, "", "  ")
	// require.NoError(t, err)

	// assert.Equal(t, string(x), string(y))
}
