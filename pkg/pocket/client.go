package pocket

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	basePath            = "https://getpocket.com/v3"
	authorizeUrlPattern = "https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s"

	endpointAdd             = "/add"
	endpointGetRequestToken = "/oauth/request"
	endpointAuthorize       = "/oauth/authorize"

	//xErrorHeader used for parse errors from Headers on non-2XX response
	xErrorHeader = "X-Error"

	defaultTimeout = 5 * time.Second
)

type (
	requestTokenRequest struct {
		ConsumerKey string `json:"consumer_key"`
		RedirectURI string `json:"redirect_uri"`
	}

	authorizationRequest struct {
		ConsumerKey string `json:"consumer_key"`
		Code        string `json:"code"`
	}

	AuthorizationResponse struct {
		AccessToken string `json:"access_token"`
		Username    string `json:"username"`
	}

	addRequest struct {
		URL         string `json:"url"`
		Title       string `json:"title,omitempty"`
		Tags        string `json:"tags,omitempty"`
		AccessToken string `json:"access_token"`
		ConsumerKey string `json:"consumer_key"`
	}

	AddInput struct {
		URL         string
		Title       string
		Tags        []string
		AccessToken string
	}

	AuthorizationUrl string
)

func (ai AddInput) validate() error {
	if ai.URL == "" {
		return errors.New("`URL` should not be empty")
	}

	if ai.AccessToken == "" {
		return errors.New("`AccessToken` should not be empty")
	}

	return nil
}

func (ai AddInput) makeRequest(consumerKey string) addRequest {
	return addRequest{
		URL:         ai.URL,
		Title:       ai.Title,
		Tags:        strings.Join(ai.Tags, ","),
		AccessToken: ai.AccessToken,
		ConsumerKey: consumerKey,
	}
}

type Client struct {
	client      *http.Client
	consumerKey string
}

func NewClient(consumerKey string) (*Client, error) {
	if consumerKey == "" {
		return nil, errors.New("consumer key should not be empty")
	}

	return &Client{
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		consumerKey: consumerKey,
	}, nil
}

func (c *Client) GetRequestToken(ctx context.Context, redirectUrl string) (string, error) {
	inp := &requestTokenRequest{
		ConsumerKey: c.consumerKey,
		RedirectURI: redirectUrl,
	}

	values, err := c.sendRequest(ctx, endpointGetRequestToken, inp)
	if err != nil {
		return "", err
	}

	requestToken := values.Get("code")
	if requestToken == "" {
		return "", errors.New("empty request token in api response")
	}

	return requestToken, nil
}

func (c *Client) MakeAuthorizationUrl(requestToken, redirectUrl string) (AuthorizationUrl, error) {
	if requestToken == "" || redirectUrl == "" {
		return "", errors.New("empty param(-s)")
	}

	return AuthorizationUrl(fmt.Sprintf(authorizeUrlPattern, requestToken, redirectUrl)), nil
}

func (c *Client) Authorize(ctx context.Context, requestToken string) (*AuthorizationResponse, error) {
	if requestToken == "" {
		return nil, errors.New("request token should not be empty")
	}

	inp := &authorizationRequest{
		ConsumerKey: c.consumerKey,
		Code:        requestToken,
	}

	values, err := c.sendRequest(ctx, endpointAuthorize, inp)
	if err != nil {
		return nil, err
	}

	accessToken, username := values.Get("access_token"), values.Get("username")
	if accessToken == "" {
		return nil, errors.New("empty access token in api response")
	}

	return &AuthorizationResponse{
		AccessToken: accessToken,
		Username:    username,
	}, nil
}

func (c *Client) Add(ctx context.Context, input AddInput) error {
	if err := input.validate(); err != nil {
		return err
	}

	_, err := c.sendRequest(ctx, endpointAdd, input.makeRequest(c.consumerKey))

	return err
}

func (c *Client) sendRequest(ctx context.Context, endpointUri string, body interface{}) (url.Values, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return url.Values{}, errors.New(fmt.Sprintf("failed to marshal input body. error: %s", err.Error()))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, basePath+endpointUri, bytes.NewBuffer(b))
	if err != nil {
		return url.Values{}, errors.New(fmt.Sprintf("failed to create new request. error: %s", err.Error()))
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF8")

	res, err := c.client.Do(req)
	if err != nil {
		return url.Values{}, errors.New(fmt.Sprintf("failed to send http request. error: %s", err.Error()))
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return url.Values{}, errors.New(fmt.Sprintf("API Error: %s", res.Header.Get(xErrorHeader)))
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return url.Values{}, errors.New(fmt.Sprintf("failed to read response body. error: %s", err.Error()))
	}

	values, err := url.ParseQuery(string(resBody))
	if err != nil {
		return url.Values{}, errors.New(fmt.Sprintf("failed to parse response body. error: %s", err.Error()))
	}

	return values, nil
}
