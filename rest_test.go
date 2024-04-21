package rest

import (
	"bytes"
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
	"time"

	"github.com/google/go-querystring/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersion_IsSetToTheLatest(t *testing.T) {
	assert.Equal(t, "0.0.1", Version)
}

func TestConstants_AreSetToTheDefaults(t *testing.T) {
	assert.Equal(t, "", defaultAPIKey)
	assert.Equal(t, "me.vanyauhalin.opensubtitlescli.rest v0.0.1", defaultUserAgent)

	assert.Equal(t, "https://api.opensubtitles.com/api/v1/", defaultBaseURL)
	assert.Equal(t, "https://vip-api.opensubtitles.com/api/v1/", vipBaseURL)

	assert.Equal(t, "Api-Key", apiKeyHeader)
	assert.Equal(t, "X-OpenSubtitles-Message", messageHeader)

	assert.Equal(t, "X-RateLimit-Limit", headerRateLimit)
	assert.Equal(t, "X-RateLimit-Remaining", headerRateRemaining)
	assert.Equal(t, "X-RateLimit-Reset", headerRateReset)
}

func TestNewClient_InitializesTheClientWithDefaults(t *testing.T) {
	c := NewClient(nil)

	assert.Equal(t, defaultAPIKey, c.APIKey)
	assert.Equal(t, defaultUserAgent, c.UserAgent)
	assert.Equal(t, defaultBaseURL, c.BaseURL.String())

	assert.Equal(t, http.DefaultClient, c.Client())

	assert.Equal(t, &AuthService{c}, c.Auth)
	assert.Equal(t, &FeaturesService{c}, c.Features)
	assert.Equal(t, &FormatsService{c}, c.Formats)
	assert.Equal(t, &LanguagesService{c}, c.Languages)
	assert.Equal(t, &SubtitlesService{c}, c.Subtitles)
	assert.Equal(t, &UsersService{c}, c.Users)
}

func TestNewClient_InitializesTheClientWithACustomClient(t *testing.T) {
	e := &http.Client{}
	c := NewClient(e)
	a0 := c.Client()
	assert.Equal(t, e.Timeout, a0.Timeout)

	e = &http.Client{
		Timeout: time.Hour,
	}
	c = NewClient(e)
	a1 := c.Client()
	assert.Equal(t, e.Timeout, a1.Timeout)

	assert.NotEqual(t, a0.Timeout, a1.Timeout)
}

func TestWithAuthToken_AppliesTheToken(t *testing.T) {
	m := http.NewServeMux()
	s := httptest.NewServer(m)
	defer s.Close()

	m.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		e := "Bearer xxx"
		a := r.Header.Get("Authorization")
		assert.Equal(t, e, a)
	})

	c := NewClient(nil)
	c = c.WithAuthToken("xxx")
	_, err := c.client.Get(s.URL)
	require.NoError(t, err)
}

func TestNewURL_InitializesTheURL(t *testing.T) {
	c := NewClient(nil)
	e := c.BaseURL.String()
	a, err := c.NewURL("", nil)
	require.NoError(t, err)
	assert.Equal(t, e, a)
}

func TestNewURL_InitializesTheURLWithThePath(t *testing.T) {
	c := NewClient(nil)
	e := c.BaseURL.String() + "p"
	a, err := c.NewURL("p", nil)
	require.NoError(t, err)
	assert.Equal(t, e, a)
}

func TestNewURL_InitializesTheURLWithParameters(t *testing.T) {
	c := NewClient(nil)
	e := c.BaseURL.String() + testQuery()
	a, err := c.NewURL("", testInstance())
	require.NoError(t, err)
	assert.Equal(t, e, a)
}

func TestNewURL_InitializesTheURLWithOmittedParameters(t *testing.T) {
	c := NewClient(nil)
	e := c.BaseURL.String()
	a, err := c.NewURL("", &testStruct{})
	require.NoError(t, err)
	assert.Equal(t, e, a)
}

func TestNewURL_ReturnsAnErrorIfTheBaseURLDoesNotHaveATrailingSlash(t *testing.T) {
	c := NewClient(nil)
	e := useBadBaseURL(c)
	_, a := c.NewURL("", nil)
	assert.EqualError(t, a, e)
}

func TestNewURL_ReturnsAnErrorIfThePathHaveALeadingSlash(t *testing.T) {
	c := NewClient(nil)
	e := `rest: url path must not have a leading slash, but "/" does`
	_, a := c.NewURL("/", nil)
	assert.EqualError(t, a, e)
}

func TestNewURL_ReturnsAnErrorIfTheBaseURLFailsToParseThePath(t *testing.T) {
	c := NewClient(nil)
	e := `parse ":": missing protocol scheme`
	_, a := c.NewURL(":", nil)
	assert.EqualError(t, a, e)
}

func TestNewURL_ReturnsAnErrorIfItFailsToParseParameters(t *testing.T) {
	c := NewClient(nil)
	e := "query: Values() expects struct input. Got int"
	_, a := c.NewURL("", 0)
	assert.EqualError(t, a, e)
}

func TestNewRequest_InitializesTheRequestWithDefaultHeaders(t *testing.T) {
	c := NewClient(nil)

	u, err := c.NewURL("", nil)
	require.NoError(t, err)

	req, err := c.NewRequest("", u, nil)
	require.NoError(t, err)

	assert.Equal(t, "application/json", req.Header.Get("Accept"))
	assert.Equal(t, defaultAPIKey, req.Header.Get("Api-Key"))
	assert.Equal(t, defaultUserAgent, req.Header.Get("User-Agent"))
}

func TestNewRequest_InitializesTheRequestWithCustomHeaders(t *testing.T) {
	c := NewClient(nil)
	c.APIKey = "xxx"
	c.UserAgent = "app v0.0.0"

	u, err := c.NewURL("", nil)
	require.NoError(t, err)

	req, err := c.NewRequest("", u, nil)
	require.NoError(t, err)

	assert.Equal(t, "xxx", req.Header.Get("Api-Key"))
	assert.Equal(t, "app v0.0.0", req.Header.Get("User-Agent"))
}

func TestNewRequest_InitializesTheRequestWithBodyHeaders(t *testing.T) {
	c := NewClient(nil)

	u, err := c.NewURL("", nil)
	require.NoError(t, err)

	req, err := c.NewRequest("", u, &testStruct{})
	require.NoError(t, err)

	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
}

func TestNewRequest_InitializesTheRequestWithoutAMethod(t *testing.T) {
	c := NewClient(nil)

	u, err := c.NewURL("", nil)
	require.NoError(t, err)

	req, err := c.NewRequest("", u, nil)
	require.NoError(t, err)

	assert.Equal(t, "GET", req.Method)
}

func TestNewRequest_InitializesTheRequestWithAMethod(t *testing.T) {
	c := NewClient(nil)

	u, err := c.NewURL("", nil)
	require.NoError(t, err)

	req, err := c.NewRequest("POST", u, nil)
	require.NoError(t, err)

	assert.Equal(t, "POST", req.Method)
}

func TestNewRequest_InitializesTheRequestWithoutAURL(t *testing.T) {
	c := NewClient(nil)

	u, err := c.NewURL("", nil)
	require.NoError(t, err)

	req, err := c.NewRequest("", u, nil)
	require.NoError(t, err)

	assert.Equal(t, u, req.URL.String())
}

func TestNewRequest_InitializesTheRequestWithAURL(t *testing.T) {
	c := NewClient(nil)

	u, err := c.NewURL("p", nil)
	require.NoError(t, err)

	req, err := c.NewRequest("", u, nil)
	require.NoError(t, err)

	assert.Equal(t, u, req.URL.String())
}

func TestNewRequest_InitializesTheRequestWithoutTheBody(t *testing.T) {
	c := NewClient(nil)

	u, err := c.NewURL("", nil)
	require.NoError(t, err)

	req, err := c.NewRequest("", u, nil)
	require.NoError(t, err)

	assert.Nil(t, req.Body)
}

func TestNewRequest_InitializesTheRequestWithTheBody(t *testing.T) {
	c := NewClient(nil)

	u, err := c.NewURL("", nil)
	require.NoError(t, err)

	req, err := c.NewRequest("", u, testInstance())
	require.NoError(t, err)

	e := testInstance()

	d := json.NewDecoder(req.Body)
	var a testStruct
	err = d.Decode(&a)
	require.NoError(t, err)

	assert.Equal(t, e, &a)
}

func TestNewRequest_ReturnsAnErrorIfTheMethodIsInvalid(t *testing.T) {
	c := NewClient(nil)
	e := `net/http: invalid method "\xbd"`
	_, a := c.NewRequest("\xbd", "", nil)
	assert.EqualError(t, a, e)
}

func TestNewRequest_ReturnsAnErrorIfItFailsToMarshalTheBody(t *testing.T) {
	c := NewClient(nil)

	type body struct {
		A map[interface {}]interface {}
	}

	_, err := c.NewRequest("", "", &body{})
	_, ok := err.(*json.UnsupportedTypeError)
	assert.True(t, ok)
}

func TestUserAgentErrorError_ReturnsError(t *testing.T) {
	e, r := testResponseError()
	a := UserAgentError(*r)
	assert.Equal(t, e, a.Error())
}

func TestAPIKeyErrorError_ReturnsError(t *testing.T) {
	e, r := testResponseError()
	a := APIKeyError(*r)
	assert.Equal(t, e, a.Error())
}

func TestAuthTokenErrorError_ReturnsError(t *testing.T) {
	e, r := testResponseError()
	a := AuthTokenError(*r)
	assert.Equal(t, e, a.Error())
}

func TestCredentialsErrorError_ReturnsError(t *testing.T) {
	e, r := testResponseError()
	a := CredentialsError(*r)
	assert.Equal(t, e, a.Error())
}

func TestFileErrorError_ReturnsError(t *testing.T) {
	e, r := testResponseError()
	a := FileError(*r)
	assert.Equal(t, e, a.Error())
}

func TestQuotaErrorError_ReturnsError(t *testing.T) {
	e, r := testResponseError()
	a := QuotaError{ResponseError: *r}
	assert.Equal(t, e, a.Error())
}

func TestLinkErrorError_ReturnsError(t *testing.T) {
	e, r := testResponseError()
	a := LinkError(*r)
	assert.Equal(t, e, a.Error())
}

func TestRateLimitErrorError_ReturnsError(t *testing.T) {
	e, r := testResponseError()
	a := RateLimitError(*r)
	assert.Equal(t, e, a.Error())
}

func TestErrorResponseError_ReturnsError(t *testing.T) {
	e, r := testResponseError()
	a := ErrorResponse{ResponseError: *r}
	assert.Equal(t, e, a.Error())
}

func TestResponseErrorError_ReturnsError(t *testing.T) {
	e, a := testResponseError()
	assert.Equal(t, e, a.Error())
}

func TestResponseErrorError_ReturnsErrorWithoutResponse(t *testing.T) {
	e := "error"
	a := ResponseError{Message: e}
	assert.Equal(t, e, a.Error())
}

func TestResponseErrorError_ReturnsErrorWithoutRequest(t *testing.T) {
	e := "400 error"
	a := ResponseError{
		Response: &http.Response{
			StatusCode: http.StatusBadRequest,
		},
		Message: "error",
	}
	assert.Equal(t, e, a.Error())
}

func TestCheckResponse_ReturnsNilIfTheResponseIsSuccessful(t *testing.T) {
	for s := 200; s <= 299; s += 1 {
		r := &http.Response{
			StatusCode: s,
		}
		err := CheckResponse(r)
		assert.Nil(t, err)
	}
}

func TestCheckResponse_ReturnsAnErrorIfTheResponseIsUnsuccessful(t *testing.T) {
	toMessage := func (s string) io.ReadCloser {
		return toBody(`{
			"message": "` + s + `"
		}`)
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
	newErr := func (r *http.Response, m string) *ErrorResponse {
		return &ErrorResponse{
			ResponseError: ResponseError{
				Response: r,
				Message: m,
			},
		}
	}
	try := func (r *http.Response) *ErrorResponse {
		err := CheckResponse(r)
		require.NotNil(t, err)
		var er *ErrorResponse
		if !errors.As(err, &er) {
			require.Fail(t, "error is not ErrorResponse")
		}
		return er
	}
	equal := func (e *ErrorResponse, r *http.Response) {
		a := try(r)
		assert.Equal(t, e, a)
	}

	r := newRes()
	e := newErr(r, "User-Agent header is empty; set it to App name with version eg: MyApp v1.2.3")
	r.Header.Set(messageHeader, e.Message)
	e.Errors = append(e.Errors, (*UserAgentError)(&e.ResponseError))
	equal(e, r)

	r = newRes()
	e = newErr(r, "")
	re0 := &AuthTokenError{
		Response: r,
		Message: "Invalid token a1",
	}
	e.Errors = append(e.Errors, re0)
	re1 := &AuthTokenError{
		Response: r,
		Message: "No token in request",
	}
	e.Errors = append(e.Errors, re1)
	re2 := &ResponseError{
		Response: r,
		Message: "Query is too short",
	}
	e.Errors = append(e.Errors, re2)
	reM := []string{re0.Message, re1.Message, re2.Message}
	e.Message = strings.Join(reM, "; ")
	reD, reDErr := json.Marshal(reM)
	require.NoError(t, reDErr)
	r.Body = toBody(`{"errors": ` + string(reD) + `}`)
	equal(e, r)

	r = newRes()
	e = newErr(r, "Internal Server Error")
	r.Body = toBody(`{"error": "` + e.Message + `"}`)
	e.Errors = append(e.Errors, (*ResponseError)(&e.ResponseError))
	equal(e, r)

	r = newRes()
	e = newErr(r, "Not enough parameters")
	r.Body = toMessage(e.Message)
	e.Errors = append(e.Errors, (*ResponseError)(&e.ResponseError))
	equal(e, r)

	r = newRes()
	e = newErr(r, "Error, invalid username/password")
	r.Body = toMessage(e.Message)
	e.Errors = append(e.Errors, (*CredentialsError)(&e.ResponseError))
	equal(e, r)

	r = newRes()
	e = newErr(r, "You cannot consume this service")
	r.Body = toMessage(e.Message)
	e.Errors = append(e.Errors, (*APIKeyError)(&e.ResponseError))
	equal(e, r)

	r = newRes()
	e = newErr(r, "Invalid file_id")
	r.Body = toMessage(e.Message)
	e.Errors = append(e.Errors, (*FileError)(&e.ResponseError))
	equal(e, r)

	r = newRes()
	e = newErr(r, "You have downloaded your allowed 20 subtitles for 24h.Your quota will be renewed in 23 hours and 57 minutes (2022-01-30 06:00:53 UTC) ")
	q := &QuotaError{
		ResponseError: e.ResponseError,
		Quota: Quota{
			Remaining: -1,
			Requests: 21,
			ResetTime: "23 hours and 57 minutes",
			ResetTimeUTC: *AllocateTime("2022-01-30T06:00:53.000Z"),
		},
	}
	r.Body = toBody(
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
	equal(e, r)

	r = newRes()
	e = newErr(r, "invalid token")
	r.Body = toMessage(e.Message)
	e.Errors = append(e.Errors, (*AuthTokenError)(&e.ResponseError))
	equal(e, r)

	r = newRes()
	e = newErr(r, "Throttle limit reached. Retry later.")
	r.Body = toMessage(e.Message)
	e.Errors = append(e.Errors, (*RateLimitError)(&e.ResponseError))
	equal(e, r)

	r = newRes()
	e = newErr(r, "Error 403 - User-Agent header is empty; set it to App name with version eg: MyApp v1.2.3")
	r.Header.Set("Content-Type", "text/html")
	r.Body = toBody(`<!DOCTYPE html><html><head><title>` + e.Message + `</title></head></html>`)
	e.Errors = append(e.Errors, (*UserAgentError)(&e.ResponseError))
	equal(e, r)

	r = newRes()
	e = newErr(r, "Error 410 - Invalid or expired link")
	r.Header.Set("Content-Type", "text/txt")
	r.Body = toBody(e.Message)
	e.Errors = append(e.Errors, (*LinkError)(&e.ResponseError))
	equal(e, r)

	r = newRes()
	e = newErr(r, "rest: api-key header is empty")
	r.Request.Header.Del(apiKeyHeader)
	e.Errors = append(e.Errors, (*APIKeyError)(&e.ResponseError))
	equal(e, r)

	r = newRes()
	e = newErr(r, "rest: user-agent header is empty")
	r.Request.Header.Set("User-Agent", "")
	e.Errors = append(e.Errors, (*UserAgentError)(&e.ResponseError))
	equal(e, r)

	r = newRes()
	e = newErr(r, "rest: authorization header is empty")
	r.Request.Header.Del("Authorization")
	e.Errors = append(e.Errors, (*AuthTokenError)(&e.ResponseError))
	equal(e, r)

	r = newRes()
	e = newErr(r, "rest: authorization token is empty")
	r.Request.Header.Set("Authorization", "Bearer ")
	e.Errors = append(e.Errors, (*AuthTokenError)(&e.ResponseError))
	equal(e, r)

	r = newRes()
	e = newErr(r, "")
	e0 := &UserAgentError{
		Response: r,
		Message: "User-Agent header is empty; set it to App name with version eg: MyApp v1.2.3",
	}
	r.Header.Set(messageHeader, e0.Message)
	e.Errors = append(e.Errors, e0)
	e1 := &ResponseError{
		Response: r,
		Message: "Not enough parameters",
	}
	r.Body = toMessage(e1.Message)
	e.Errors = append(e.Errors, e1)
	e2 := &APIKeyError{
		Response: r,
		Message: "rest: api-key header is empty",
	}
	r.Request.Header.Del(apiKeyHeader)
	e.Errors = append(e.Errors, e2)
	e.Message = strings.Join([]string{e0.Message, e1.Message, e2.Message}, "; ")
	equal(e, r)
}

func TestCheckResponse_ReturnsAnOpenBodyIfTheResponseIsSuccessful(t *testing.T) {
	e := testJSON()

	r := &http.Response{
		StatusCode: http.StatusOK,
		Body: toBody(e),
	}
	err := CheckResponse(r)
	require.NoError(t, err)

	d, err := io.ReadAll(r.Body)
	require.NoError(t, err)

	err = r.Body.Close()
	require.NoError(t, err)

	a := string(d)
	assert.Equal(t, e, a)
}

func TestCheckResponse_ReturnsAnOpenBodyIfTheResponseIsUnsuccessful(t *testing.T) {
	e := testJSON()

	r := &http.Response{
		Request: &http.Request{
			Header: http.Header{},
		},
		StatusCode: http.StatusBadRequest,
		Header: http.Header{
			"Content-Type": {"application/json"},
		},
		Body: toBody(e),
	}
	err := CheckResponse(r)
	require.NotNil(t, err)

	d, err := io.ReadAll(r.Body)
	require.NoError(t, err)

	err = r.Body.Close()
	require.NoError(t, err)

	a := string(d)
	assert.Equal(t, e, a)
}

func TestPagination_UnmarshalsAndMarshals(t *testing.T) {
	a := &Pagination{
		Page: 0,
		PerPage: 0,
		TotalCount: 0,
		TotalPages: 0,
	}
	b := "{}"
	equalUnmarshal(t, a, b)

	a = &Pagination{
		Page: 1,
		PerPage: 10,
		TotalCount: 100,
		TotalPages: 10,
	}
	b = `{
		"page": 1,
		"per_page": 10,
		"total_count": 100,
		"total_pages": 10
	}`
	equalJSON(t, a, b)
}

func TestQuota_UnmarshalsAndMarshals(t *testing.T) {
	a := &Quota{
		Remaining: 0,
		Requests: 0,
		ResetTime: "",
		ResetTimeUTC: *AllocateTime("0001-01-01T00:00:00Z"),
	}
	b := "{}"
	equalUnmarshal(t, a, b)

	a = &Quota{
		Remaining: 5,
		Requests: 4,
		ResetTime: "1 hour",
		ResetTimeUTC: *AllocateTime("2022-01-30T06:00:53Z"),
	}
	b = `{
		"remaining": 5,
		"requests": 4,
		"reset_time": "1 hour",
		"reset_time_utc": "2022-01-30T06:00:53Z"
	}`
	equalJSON(t, a, b)
}

func TestBareDo_DoesARequest(t *testing.T) {
	c, m, teardown := setup()
	defer teardown()

	m.HandleFunc("/hi", func (w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/hi" + testQuery(), r.RequestURI)
		w.WriteHeader(http.StatusOK)
	})

	ctx := context.Background()

	u, err := c.NewURL("hi", testInstance())
	require.NoError(t, err)

	req, err := c.NewRequest("POST", u, nil)
	require.NoError(t, err)

	res, err := c.BareDo(ctx, req)
	require.NoError(t, err)

	err = res.Body.Close()
	require.NoError(t, err)
}

func TestBareDo_ReturnsASuccessfulResponse(t *testing.T) {
	c, m, teardown := setup()
	defer teardown()

	m.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {})

	ctx := context.Background()

	e := Response{
		Response: &http.Response{
			StatusCode: http.StatusOK,
		},
	}

	u, err := c.NewURL("", nil)
	require.NoError(t, err)

	r, err := c.NewRequest("", u, nil)
	require.NoError(t, err)

	a, err := c.BareDo(ctx, r)
	require.NoError(t, err)

	err = a.Body.Close()
	require.NoError(t, err)

	assert.Equal(t, e.StatusCode, a.StatusCode)
}

func TestBareDo_ReturnsAUnsuccessfulResponse(t *testing.T) {
	c, m, teardown := setup()
	defer teardown()

	m.HandleFunc("/100", func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{
			"error": "Internal Server Error"
		}`)
	})

	ctx := context.Background()

	e := Response{
		Response: &http.Response{
			StatusCode: http.StatusBadRequest,
		},
	}

	u, err := c.NewURL("100", nil)
	require.NoError(t, err)

	r, err := c.NewRequest("", u, nil)
	require.NoError(t, err)

	a, ae := c.BareDo(ctx, r)

	err = a.Body.Close()
	require.NoError(t, err)

	var er *ErrorResponse
	require.ErrorAs(t, ae, &er)

	var re *ResponseError
	assert.ErrorAs(t, er.Errors[0], &re)

	var ke *APIKeyError
	assert.ErrorAs(t, er.Errors[1], &ke)

	var te *AuthTokenError
	assert.ErrorAs(t, er.Errors[2], &te)

	assert.Equal(t, e.StatusCode, a.StatusCode)
}

func TestBareDo_ReturnsAnOpenBodyIfTheResponseIsSuccessful(t *testing.T) {
	c, m, teardown := setup()
	defer teardown()

	m.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testJSON())
	})

	ctx := context.Background()

	e := testJSON()

	u, err := c.NewURL("", nil)
	require.NoError(t, err)

	r, err := c.NewRequest("", u, nil)
	require.NoError(t, err)

	res, err := c.BareDo(ctx, r)
	require.NoError(t, err)

	d, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	err = res.Body.Close()
	require.NoError(t, err)

	a := string(d)
	assert.Equal(t, e, a)
}

func TestBareDo_ReturnsAnOpenBodyIfTheResponseIsUnsuccessful(t *testing.T) {
	c, m, teardown := setup()
	defer teardown()

	m.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, testJSON())
	})

	ctx := context.Background()

	e := testJSON()

	u, err := c.NewURL("", nil)
	require.NoError(t, err)

	req, err := c.NewRequest("", u, nil)
	require.NoError(t, err)

	res, err := c.BareDo(ctx, req)
	require.NotNil(t, res)
	require.NotNil(t, err)

	d, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	err = res.Body.Close()
	require.NoError(t, err)

	a := string(d)
	assert.Equal(t, e, a)
}

func TestBareDo_ReturnsAResponseWithRate(t *testing.T) {
	c, m, teardown := setup()
	defer teardown()

	m.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set(headerRateLimit, "5")
		w.Header().Set(headerRateRemaining, "4")
		w.Header().Set(headerRateReset, "1")
	})

	ctx := context.Background()

	e := Response{
		Rate: Rate{
			Limit: 5,
			Remaining: 4,
			Reset: 1,
		},
	}

	u, err := c.NewURL("", nil)
	require.NoError(t, err)

	r, err := c.NewRequest("", u, nil)
	require.NoError(t, err)

	a, err := c.BareDo(ctx, r)
	require.Nil(t, err)

	err = a.Body.Close()
	require.NoError(t, err)

	assert.Equal(t, e.Rate, a.Rate)
}

func TestBareDo_ReturnsAResponseWithPagination(t *testing.T) {
	c, m, teardown := setup()
	defer teardown()

	m.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, `{
			"page": 1
		}`)
	})

	ctx := context.Background()

	e := Response{
		Pagination: Pagination{
			Page: 1,
		},
	}

	u, err := c.NewURL("", nil)
	require.NoError(t, err)

	r, err := c.NewRequest("", u, nil)
	require.NoError(t, err)

	a, err := c.BareDo(ctx, r)
	require.Nil(t, err)

	err = a.Body.Close()
	require.NoError(t, err)

	assert.Equal(t, e.Rate, a.Rate)
}

func TestBareDo_ReturnsAResponseWithQuota(t *testing.T) {
	c, m, teardown := setup()
	defer teardown()

	m.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, `{
			"remaining": 5
		}`)
	})

	ctx := context.Background()

	e := Response{
		Quota: Quota{
			Remaining: 5,
		},
	}

	u, err := c.NewURL("", nil)
	require.NoError(t, err)

	r, err := c.NewRequest("", u, nil)
	require.NoError(t, err)

	a, err := c.BareDo(ctx, r)
	require.Nil(t, err)

	err = a.Body.Close()
	require.NoError(t, err)

	assert.Equal(t, e.Quota, a.Quota)
}

func TestBareDo_ThrowsAPanicIfTheContextIsNil(t *testing.T) {
	c := NewClient(nil)
	assert.Panics(t, func () {
		var r http.Request
		_, _ = c.BareDo(nil, &r) // nolint:staticcheck
	})
}

func TestDo_DoesARequest(t *testing.T) {
	c, m, teardown := setup()
	defer teardown()

	m.HandleFunc("/hi", func (w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/hi" + testQuery(), r.RequestURI)
		w.WriteHeader(http.StatusOK)
	})

	ctx := context.Background()

	u, err := c.NewURL("hi", testInstance())
	require.NoError(t, err)

	r, err := c.NewRequest("POST", u, nil)
	require.NoError(t, err)

	_, err = c.Do(ctx, r, nil)
	require.NoError(t, err)
}

func TestDo_ReturnsASuccessfulResponse(t *testing.T) {
	c, m, teardown := setup()
	defer teardown()

	m.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {})

	ctx := context.Background()

	e := Response{
		Response: &http.Response{
			StatusCode: http.StatusOK,
		},
	}

	u, err := c.NewURL("", nil)
	require.NoError(t, err)

	r, err := c.NewRequest("", u, nil)
	require.NoError(t, err)

	a, err := c.Do(ctx, r, nil)
	require.NoError(t, err)

	assert.Equal(t, e.StatusCode, a.StatusCode)
}

func TestDo_ReturnsAUnsuccessfulResponse(t *testing.T) {
	c, m, teardown := setup()
	defer teardown()

	m.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	ctx := context.Background()

	e := Response{
		Response: &http.Response{
			StatusCode: http.StatusBadRequest,
		},
	}

	u, err := c.NewURL("", nil)
	require.NoError(t, err)

	r, err := c.NewRequest("", u, nil)
	require.NoError(t, err)

	a, err := c.Do(ctx, r, nil)

	var er *ErrorResponse
	require.ErrorAs(t, err, &er)

	assert.Equal(t, e.StatusCode, a.StatusCode)
}

func TestDo_WritesToTheBuffer(t *testing.T) {
	c, m, teardown := setup()
	defer teardown()

	m.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testJSON())
	})

	ctx := context.Background()

	e := testJSON()

	u, err := c.NewURL("", nil)
	require.NoError(t, err)

	r, err := c.NewRequest("", u, nil)
	require.NoError(t, err)

	var buf bytes.Buffer

	_, err = c.Do(ctx, r, &buf)
	require.NoError(t, err)

	a := buf.String()
	assert.Equal(t, e, a)
}

func TestDo_DecodesAnInterface(t *testing.T) {
	c, m, teardown := setup()
	defer teardown()

	m.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testJSON())
	})

	ctx := context.Background()

	e := testInstance()

	u, err := c.NewURL("", nil)
	require.NoError(t, err)

	r, err := c.NewRequest("", u, nil)
	require.NoError(t, err)

	var a testStruct
	_, err = c.Do(ctx, r, &a)
	require.NoError(t, err)
	assert.Equal(t, e, &a)
}

func TestDo_ReturnsAClosedBodyIfTheResponseIsSuccessful(t *testing.T) {
	c, m, teardown := setup()
	defer teardown()

	m.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testJSON())
	})

	ctx := context.Background()

	u, err := c.NewURL("", nil)
	require.NoError(t, err)

	req, err := c.NewRequest("", u, nil)
	require.NoError(t, err)

	res, err := c.Do(ctx, req, nil)
	require.NoError(t, err)

	_, err = io.ReadAll(res.Body)
	assert.Error(t, err)
}

func TestDo_ReturnsAnOpenBodyIfTheResponseIsUnsuccessful(t *testing.T) {
	c, m, teardown := setup()
	defer teardown()

	m.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, testJSON())
	})

	ctx := context.Background()

	e := testJSON()

	u, err := c.NewURL("", nil)
	require.NoError(t, err)

	req, err := c.NewRequest("", u, nil)
	require.NoError(t, err)

	res, err := c.Do(ctx, req, nil)
	require.NotNil(t, res)
	require.NotNil(t, err)

	d, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	err = res.Body.Close()
	assert.NoError(t, err)

	a := string(d)
	assert.Equal(t, e, a)
}

func TestDo_ThrowsAPanicIfTheContextIsNil(t *testing.T) {
	c := NewClient(nil)
	assert.Panics(t, func () {
		var r http.Request
		_, _ = c.Do(nil, &r, nil) // nolint:staticcheck
	})
}

func TestIDString_ReturnsString(t *testing.T) {
	e := "0"
	a := ID(0)
	assert.Equal(t, e, a.String())
}

func TestID_EncodesValues(t *testing.T) {
	type s struct {
		ID ID `url:"id"`
	}

	a := &s{0}
	b := `id=0`
	equalQuery(t, a, b)
}

func TestID_UnmarshalsAndMarshals(t *testing.T) {
	type s struct {
		ID ID `json:"id"`
	}

	a := &s{0}
	b := `{
		"id": 0
	}`
	equalJSON(t, a, b)

	a = &s{1}
	b = `{
		"id": "1"
	}`
	equalUnmarshal(t, a, b)

	a = &s{1}
	b = `{
		"id": 1
	}`
	equalMarshal(t, a, b)
}

func TestAllocateID_Allocates(t *testing.T) {
	e := ID(0)
	a := AllocateID(0)
	assert.Equal(t, e, *a)
}

func TestAllocateBool_Allocates(t *testing.T) {
	e := false
	a := AllocateBool(false)
	assert.Equal(t, e, *a)
}

func TestAllocateInt_Allocates(t *testing.T) {
	e := 0
	a := AllocateInt(0)
	assert.Equal(t, e, *a)
}

func TestAllocateString_Allocates(t *testing.T) {
	e := ""
	a := AllocateString("")
	assert.Equal(t, e, *a)
}

func TestAllocateTime_AllocatesWithTime(t *testing.T) {
	e := time.Time{}
	a := AllocateTime("0001-01-01T00:00:00Z")
	assert.Equal(t, e, *a)
}

func setup() (*Client, *http.ServeMux, func ()) {
	m := http.NewServeMux()
	s := httptest.NewServer(m)

	c := NewClient(nil)
	c.BaseURL, _ = url.Parse(s.URL + "/")

	return c, m, s.Close
}

type testStruct struct {
	A *int    `url:"a,omitempty" json:"a,omitempty"`
	B *int    `url:"b,omitempty" json:"b,omitempty"`
	C *string `url:"c,omitempty" json:"c,omitempty"`
}

func testInstance() *testStruct {
	return &testStruct{
		A: AllocateInt(0),
		B: AllocateInt(1),
		C: AllocateString("a b c"),
	}
}

func testJSON() string {
	return `{
		"a": 0,
		"b": 1,
		"c": "a b c"
	}`
}

func testQuery() string {
	return "?&a=0&b=1&c=a+b+c"
}

func testResponseError() (string, *ResponseError) {
	e := "GET http://localhost/: 400 error"
	r := &ResponseError{
		Response: &http.Response{
			Request: &http.Request{
				Method: "GET",
			},
			StatusCode: http.StatusBadRequest,
		},
		Message: "error",
	}
	r.Response.Request.URL, _ = url.Parse("http://localhost/")
	return e, r
}

func useBadBaseURL(c *Client) string {
	c.BaseURL, _ = url.Parse("http://localhost/v2")
	return `rest: base url must have a trailing slash, but "http://localhost/v2" does not`
}

func equalQuery(t *testing.T, a interface {}, b string) {
	q, err := query.Values(a)
	require.NoError(t, err)
	assert.Equal(t, q.Encode(), b)
}

func equalJSON(t *testing.T, a interface {}, b string) {
	equalUnmarshal(t, a, b)
	equalMarshal(t, a, b)
}

func equalUnmarshal(t *testing.T, a interface {}, b string) {
	es := a
	as := reflect.New(reflect.TypeOf(a).Elem()).Interface()
	err := json.Unmarshal([]byte(b), &as)
	require.NoError(t, err)
	assert.Equal(t, es, as)
}

func equalMarshal(t *testing.T, a interface {}, b string) {
	ej, err := json.MarshalIndent(a, "", "  ")
	require.NoError(t, err)
	var aj bytes.Buffer
	err = json.Indent(&aj, []byte(b), "", "  ")
	require.NoError(t, err)
	assert.JSONEq(t, string(ej), aj.String())
}

func equalBody(t *testing.T, b io.ReadCloser, e string) {
	d, err := io.ReadAll(b)
	require.NoError(t, err)

	err = b.Close()
	require.NoError(t, err)

	a := string(d)
	assert.JSONEq(t, e, a)
}

func toBody(s string) io.ReadCloser {
	r := strings.NewReader(s)
	return io.NopCloser(r)
}
