package pricing

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type Price struct {
	AuctionHouseID int  `json:"auctionHouseId"`
	ItemID         int  `json:"itemId"`
	PetSpeciesID   *int `json:"petSpeciesId"`
	MinBuyout      int  `json:"minBuyout"`
	Quantity       int  `json:"quantity"`
	MarketValue    int  `json:"marketValue"`
	Historical     int  `json:"historical"`
	NumAuctions    int  `json:"numAuctions"`
}

// WriteAHPricesToFile queries the TSM API for pricing data for a given AH and dumps the response to a file.
func (c *Client) WriteAHPricesToFile(id int) (string, error) {
	fullURL, err := url.Parse(fmt.Sprintf("%s/ah/%d", c.baseURL, id))
	if err != nil {
		return "", err
	}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    fullURL,
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error retrieving pricing from AH %d: %w", id, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		msg := fmt.Errorf("error retrieving pricing from AH %d: %s", id, resp.Status)
		if resp.Body != nil {
			body, _ := io.ReadAll(resp.Body)
			return "", errors.Join(msg, errors.New(string(body)))
		}
	}

	ts := time.Now().Format("2006/01/02/15:04")
	filePath := fmt.Sprintf("./data/pricing/%s.json", ts)
	err = os.MkdirAll(filepath.Dir(filePath), 0777)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	f, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to stream pricing data to file: %w", err)
	}
	return filePath, nil
}
