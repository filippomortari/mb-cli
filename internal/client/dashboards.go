package client

import "fmt"

// ListDashboards retrieves all dashboards.
func (c *Client) ListDashboards() ([]Dashboard, error) {
	resp, err := c.Get("/api/dashboard/", nil)
	if err != nil {
		return nil, err
	}

	var dashboards []Dashboard
	if err := c.DecodeJSON(resp, &dashboards); err != nil {
		return nil, err
	}

	return dashboards, nil
}

// GetDashboard retrieves a single dashboard by ID.
func (c *Client) GetDashboard(id int) (*Dashboard, error) {
	resp, err := c.Get(fmt.Sprintf("/api/dashboard/%d", id), nil)
	if err != nil {
		return nil, err
	}

	var dashboard Dashboard
	if err := c.DecodeJSON(resp, &dashboard); err != nil {
		return nil, err
	}

	return &dashboard, nil
}
