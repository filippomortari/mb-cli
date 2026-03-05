package client

import (
	"fmt"
	"io"
)

// RunNativeQuery executes a native SQL query against the specified database.
func (c *Client) RunNativeQuery(databaseID int, sql string) (*QueryResult, error) {
	query := DatasetQuery{
		Database: databaseID,
		Type:     "native",
		Native:   NativeQuery{Query: sql},
	}

	resp, err := c.Post("/api/dataset/", query)
	if err != nil {
		return nil, err
	}

	var result QueryResult
	if err := c.DecodeJSON(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ExportNativeQuery executes a native SQL query and returns the result in the specified export format.
func (c *Client) ExportNativeQuery(databaseID int, sql string, format string) ([]byte, error) {
	query := DatasetQuery{
		Database: databaseID,
		Type:     "native",
		Native:   NativeQuery{Query: sql},
	}

	// The export endpoint wraps the dataset query in a "query" key.
	body := map[string]any{
		"query": query,
	}

	resp, err := c.Post(fmt.Sprintf("/api/dataset/%s", format), body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read export response: %w", err)
	}

	return data, nil
}
