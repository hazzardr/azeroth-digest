package pricing

import (
	"bytes"
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

const (
	TSMTokenURL = "https://auth.tradeskillmaster.com/oauth2/token"
	TSMClientID = "c260f00d-1071-409a-992f-dda2e5498536"
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func NewClient(apiKey string, baseURL string, auctionHouseID int) (*Client, error) {
	tokenURL, err := url.Parse(TSMTokenURL)
	if err != nil {
		return nil, errors.Join(errors.New("failed to parse token url"), err)
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
	bu, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return &Client{
		baseURL:    bu,
		httpClient: oc,
	}, nil
}

type TSMTokenSource struct {
	tokenURL *url.URL
	apiKey   string
	clientID string
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
	resp, err = client.Post(
		ts.tokenURL.String(),
		"application/json",
		bytes.NewBuffer(body),
	)
}
