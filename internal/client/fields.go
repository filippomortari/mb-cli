package client

import (
	"fmt"
)

// GetField retrieves a single field by ID.
func (c *Client) GetField(id int) (*Field, error) {
	resp, err := c.Get(fmt.Sprintf("/api/field/%d", id), nil)
	if err != nil {
		return nil, err
	}

	var field Field
	if err := c.DecodeJSON(resp, &field); err != nil {
		return nil, err
	}

	return &field, nil
}

// GetFieldSummary retrieves summary statistics for a field.
// The API returns pairs like [["count",2],["distincts",2]].
func (c *Client) GetFieldSummary(id int) ([]FieldSummary, error) {
	resp, err := c.Get(fmt.Sprintf("/api/field/%d/summary", id), nil)
	if err != nil {
		return nil, err
	}

	var raw [][]any
	if err := c.DecodeJSON(resp, &raw); err != nil {
		return nil, err
	}

	summary := make([]FieldSummary, len(raw))
	for i, pair := range raw {
		if len(pair) >= 2 {
			if name, ok := pair[0].(string); ok {
				summary[i] = FieldSummary{Type: name, Value: pair[1]}
			}
		}
	}

	return summary, nil
}

// GetFieldValues retrieves distinct values for a field.
func (c *Client) GetFieldValues(id int) (*FieldValues, error) {
	resp, err := c.Get(fmt.Sprintf("/api/field/%d/values", id), nil)
	if err != nil {
		return nil, err
	}

	var values FieldValues
	if err := c.DecodeJSON(resp, &values); err != nil {
		return nil, err
	}

	return &values, nil
}
