package client

import "fmt"

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
	resp, err := c.Get(fmt.Sprintf("/api/card/%d", id), nil)
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

	var result QueryResult
	if err := c.DecodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
