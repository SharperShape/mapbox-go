package mapbox

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Client struct {
	username string
	token    string
	baseURL  url.URL

	HttpClient *http.Client
}

func NewClient(username, token string) *Client {

	baseURL, _ := url.Parse("https://api.mapbox.com")

	httpClient := http.Client{
		Timeout: 15 * time.Second,
	}

	return &Client{
		username:   username,
		token:      token,
		baseURL:    *baseURL,
		HttpClient: &httpClient,
	}
}

func (c *Client) do(method string, url url.URL, body io.Reader, value interface{}) (*http.Response, error) {

	url = c.addAuthentication(url)
	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "requesting")
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(value)
	if err != nil {
		return nil, errors.Wrap(err, "decoding json")
	}

	return resp, nil
}

// addAuthentication adds an authentication token to the URL.
func (c *Client) addAuthentication(url url.URL) url.URL {

	q := url.Query()
	q.Add("access_token", c.token)
	url.RawQuery = q.Encode()
	return url
}

// nextPageURL finds the link pointing to the next page of data in the header,
// and adds authentication.
func (c *Client) nextPageURL(header http.Header) *url.URL {

	nextRegex := regexp.MustCompile("<(.*)>")

	link := header.Get("Link")
	if strings.Contains(link, "next") {
		requestURL, err := url.Parse(nextRegex.FindStringSubmatch(link)[1])
		if err != nil {
			return nil
		}
		return requestURL
	} else {
		return nil
	}
}
