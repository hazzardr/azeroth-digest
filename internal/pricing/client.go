package pricing

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/oauth2"
)

const (
	TSMTokenURL = "https://auth.tradeskillmaster.com/oauth2/token" //nolint:gosec // token url
	TSMClientID = "c260f00d-1071-409a-992f-dda2e5498536"
)

type TSMTokenSource struct {
	tokenURL *url.URL
	apiKey   string
	clientID string
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func (ts TSMTokenSource) Token() (*oauth2.Token, error) {
	client := http.DefaultClient
	type tokenRequest struct {
		ClientID  string `json:"client_id"`
		GrantType string `json:"grant_type"`
		Scope     string `json:"scope"`
		Token     string `json:"token"`
	}
	tr := tokenRequest{
		ClientID:  ts.clientID,
		GrantType: "api_token",
		Scope:     "app:realm-api app:pricing-api",
		Token:     ts.apiKey,
	}

	body, err := json.Marshal(tr)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ts.tokenURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve access token: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode <= 200 && resp.StatusCode >= 300 {
		msg := fmt.Errorf("error retrieving TSM OAuth token: %s", resp.Status)
		if resp.Body != nil {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				msg = fmt.Errorf("error reading response body: %w", err)
			}
			return nil, errors.Join(msg, errors.New(string(body)))
		}
	}
	var otoken oauth2.Token
	err = json.UnmarshalRead(resp.Body, &otoken)
	if err != nil {
		return nil, fmt.Errorf("failed to parse access token: %w", err)
	}

	return &otoken, nil
}

func NewClient(apiKey string, baseURL string) (*Client, error) {
	tokenURL, err := url.Parse(TSMTokenURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token url: %w", err)
	}
	ts := oauth2.ReuseTokenSource(nil, &TSMTokenSource{
		tokenURL: tokenURL,
		clientID: TSMClientID,
		apiKey:   apiKey,
	})

	// Retrieving all pricing data for an AH can take a long time.
	// The underlying http client needs to have a significant timeout
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{
		Timeout: 5 * time.Minute,
	})
	oc := oauth2.NewClient(
		ctx,
		ts,
	)

	return &Client{
		baseURL:    baseURL,
		httpClient: oc,
	}, nil
}
