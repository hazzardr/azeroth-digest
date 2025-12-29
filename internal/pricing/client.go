package pricing

import (
	"context"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
)

type Client struct {
	httpClient *http.Client
}

type TSMTokenSource struct {
	tokenURL *url.URL
	apiKey   string
	clientID string
}

func NewTSMTokenSource(apiKey string, clientID string, tokenURL *url.URL) TSMTokenSource {
	return TSMTokenSource{
		tokenURL: tokenURL,
		apiKey:   apiKey,
		clientID: clientID,
	}
}

func (ts TSMTokenSource) Token() (*oauth2.Token, error) {

}

func NewClient(t string, url string, clientID string, aid int) *Client {
	oc := oauth2.NewClient(
		context.Background(),
	)
}
