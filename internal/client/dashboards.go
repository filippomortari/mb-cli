package client

import (
	"fmt"
	"net/url"
)

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

// GetDashboardParamValues retrieves valid values for a dashboard parameter.
func (c *Client) GetDashboardParamValues(dashboardID int, paramKey string) (*ParameterValues, error) {
	resp, err := c.Get(fmt.Sprintf("/api/dashboard/%d/params/%s/values", dashboardID, url.PathEscape(paramKey)), nil)
	if err != nil {
		return nil, err
	}

	var values ParameterValues
	if err := c.DecodeJSON(resp, &values); err != nil {
		return nil, err
	}

	return &values, nil
}

// SearchDashboardParamValues searches dashboard parameter values.
func (c *Client) SearchDashboardParamValues(dashboardID int, paramKey string, query string) (*ParameterValues, error) {
	resp, err := c.Get(fmt.Sprintf("/api/dashboard/%d/params/%s/search/%s", dashboardID, url.PathEscape(paramKey), url.PathEscape(query)), nil)
	if err != nil {
		return nil, err
	}

	var values ParameterValues
	if err := c.DecodeJSON(resp, &values); err != nil {
		return nil, err
	}

	return &values, nil
}

// RunDashboardCard executes a dashboard card with parameter values.
func (c *Client) RunDashboardCard(dashboardID, dashcardID, cardID int, params map[string]string) (*QueryResult, error) {
	dashboard, err := c.GetDashboard(dashboardID)
	if err != nil {
		return nil, err
	}

	body := map[string]any{
		"parameters": buildDashboardQueryParameters(dashboard, params),
	}

	resp, err := c.Post(fmt.Sprintf("/api/dashboard/%d/dashcard/%d/card/%d/query", dashboardID, dashcardID, cardID), body)
	if err != nil {
		return nil, err
	}

	return c.decodeCardQueryResult(resp)
}
