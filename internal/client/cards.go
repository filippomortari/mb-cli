package client

import (
	"fmt"
	"net/http"
	"net/url"
)

// ListCards retrieves all saved questions (cards).
func (c *Client) ListCards() ([]Card, error) {
	resp, err := c.Get("/api/card/", nil)
	if err != nil {
		return nil, err
	}

	var cards []Card
	if err := c.DecodeJSON(resp, &cards); err != nil {
		return nil, err
	}

	return cards, nil
}

// GetCard retrieves a single card by ID.
func (c *Client) GetCard(id int) (*Card, error) {
	params := url.Values{}
	params.Set("legacy-mbql", "true")

	resp, err := c.Get(fmt.Sprintf("/api/card/%d", id), params)
	if err != nil {
		return nil, err
	}

	var card Card
	if err := c.DecodeJSON(resp, &card); err != nil {
		return nil, err
	}

	return &card, nil
}

// RunCard executes a saved question and returns the query result.
func (c *Client) RunCard(id int) (*QueryResult, error) {
	resp, err := c.Post(fmt.Sprintf("/api/card/%d/query", id), nil)
	if err != nil {
		return nil, err
	}

	return c.decodeCardQueryResult(resp)
}

// RunCardWithParams executes a saved question with parameter values.
func (c *Client) RunCardWithParams(id int, params map[string]string) (*QueryResult, error) {
	if len(params) == 0 {
		return c.RunCard(id)
	}

	card, err := c.GetCard(id)
	if err != nil {
		return nil, err
	}

	body := map[string]any{
		"parameters": buildCardQueryParameters(card, params),
	}

	resp, err := c.Post(fmt.Sprintf("/api/card/%d/query", id), body)
	if err != nil {
		return nil, err
	}

	return c.decodeCardQueryResult(resp)
}

func (c *Client) decodeCardQueryResult(resp *http.Response) (*QueryResult, error) {
	var result QueryResult
	if err := c.DecodeJSON(resp, &result); err != nil {
		return nil, err
	}

	if c.RedactPII {
		RedactQueryResult(&result)
	}

	return &result, nil
}
