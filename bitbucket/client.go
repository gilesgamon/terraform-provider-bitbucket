package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

// Error represents a error from the bitbucket api.
type Error struct {
	APIError struct {
		Message string `json:"message,omitempty"`
	} `json:"error,omitempty"`
	Type       string `json:"type,omitempty"`
	StatusCode int
	Endpoint   string
}

func (e Error) Error() string {
	return fmt.Sprintf("API Error: %d %s %s", e.StatusCode, e.Endpoint, e.APIError.Message)
}

const (
	// BitbucketEndpoint is the fqdn used to talk to bitbucket
	BitbucketEndpoint string = "https://api.bitbucket.org/"
)

// Client is the base internal Client to talk to bitbuckets API. This should be a username and password
// the password should be a app-password.
type Client struct {
	Username         *string
	Password         *string
	OAuthToken       *string
	OAuthTokenSource oauth2.TokenSource
	HTTPClient       *http.Client
}

// Do Will just call the bitbucket api but also add auth to it and some extra headers
func (c *Client) Do(method, endpoint string, payload *bytes.Buffer, contentType string) (*http.Response, error) {
	absoluteendpoint := BitbucketEndpoint + endpoint
	log.Printf("[DEBUG] Sending request to %s %s", method, absoluteendpoint)

	var bodyreader io.Reader

	if payload != nil {
		log.Printf("[DEBUG] With payload %s", payload.String())
		bodyreader = payload
	}

	req, err := http.NewRequest(method, absoluteendpoint, bodyreader)
	if err != nil {
		return nil, err
	}

	if c.Username != nil && c.Password != nil {
		log.Printf("[DEBUG] Setting Basic Auth")
		req.SetBasicAuth(*c.Username, *c.Password)
	}

	if c.OAuthToken != nil {
		log.Printf("[DEBUG] Setting Bearer Token")
		bearer := "Bearer " + *c.OAuthToken
		req.Header.Add("Authorization", bearer)
	}

	if c.OAuthTokenSource != nil {
		token, err := c.OAuthTokenSource.Token()
		if err != nil {
			return nil, err
		}

		token.SetAuthHeader(req)
	}

	if payload != nil && contentType != "" {
		// Can cause bad request when putting default reviews if set.
		req.Header.Add("Content-Type", contentType)
	}

	req.Close = true

	resp, err := c.HTTPClient.Do(req)
	log.Printf("[DEBUG] Resp: %v Err: %v", resp, err)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("no response received from %s %s", method, absoluteendpoint)
	}
	if resp.StatusCode >= 400 || resp.StatusCode < 200 {
		apiError := Error{
			StatusCode: resp.StatusCode,
			Endpoint:   endpoint,
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Resp Body: %s", string(body))

		err = json.Unmarshal(body, &apiError)
		if err != nil {
			apiError.APIError.Message = string(body)
		}

		return resp, error(apiError)

	}
	return resp, err
}

// Get is just a helper method to do but with a GET verb
func (c *Client) Get(endpoint string) (*http.Response, error) {
	return c.Do("GET", endpoint, nil, "application/json")
}

// GetPaginated retrieves every page of a paginated Bitbucket 2.0 collection
// endpoint by following the `next` links, returning the concatenated `values`
// entries as raw JSON messages. Bitbucket collection endpoints default to a
// small page size (10), so callers that need the full result set must paginate.
func (c *Client) GetPaginated(endpoint string) ([]json.RawMessage, error) {
	var values []json.RawMessage
	next := endpoint

	for next != "" {
		res, err := c.Get(next)
		if err != nil {
			return nil, err
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return nil, err
		}

		var page struct {
			Values []json.RawMessage `json:"values"`
			Next   string            `json:"next"`
		}
		if err := json.Unmarshal(body, &page); err != nil {
			return nil, err
		}

		values = append(values, page.Values...)
		next = toRelativeEndpoint(page.Next)
	}

	return values, nil
}

// toRelativeEndpoint converts an absolute Bitbucket API URL (as returned in the
// `next` field of paginated responses) into an endpoint relative to
// BitbucketEndpoint, which is what Client.Do expects.
func toRelativeEndpoint(raw string) string {
	if raw == "" {
		return ""
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
	}

	rel := parsed.Path
	for len(rel) > 0 && rel[0] == '/' {
		rel = rel[1:]
	}
	if parsed.RawQuery != "" {
		rel += "?" + parsed.RawQuery
	}
	return rel
}

// encodeQueryParams builds a deterministic, URL-encoded query string (including
// the leading "?") from the provided parameters. It returns an empty string when
// there are no parameters. Using url.Values ensures keys are sorted and values
// are escaped, avoiding malformed requests when values contain spaces or special
// characters (common in Bitbucket query language filters).
func encodeQueryParams(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return "?" + values.Encode()
}

// Post is just a helper method to do but with a POST verb
func (c *Client) Post(endpoint string, jsonpayload *bytes.Buffer) (*http.Response, error) {
	return c.Do("POST", endpoint, jsonpayload, "application/json")
}

// PostNonJson is just a helper method to do but with a POST verb without Json Header
func (c *Client) PostNonJson(endpoint string, payload *bytes.Buffer) (*http.Response, error) {
	return c.Do("POST", endpoint, payload, "")
}

// PostWithContentType is just a helper method to do but with a POST verb and a provided content type
func (c *Client) PostWithContentType(endpoint, contentType string, payload *bytes.Buffer) (*http.Response, error) {
	return c.Do("POST", endpoint, payload, contentType)
}

// Put is just a helper method to do but with a PUT verb
func (c *Client) Put(endpoint string, jsonpayload *bytes.Buffer) (*http.Response, error) {
	return c.Do("PUT", endpoint, jsonpayload, "application/json")
}

// PutOnly is just a helper method to do but with a PUT verb and a nil body
func (c *Client) PutOnly(endpoint string) (*http.Response, error) {
	return c.Do("PUT", endpoint, nil, "application/json")
}

// Delete is just a helper to Do but with a DELETE verb
func (c *Client) Delete(endpoint string) (*http.Response, error) {
	return c.Do("DELETE", endpoint, nil, "application/json")
}
