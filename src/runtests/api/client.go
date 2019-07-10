package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Client ...
type Client struct {
	BaseURL string
}

// SetResult ...
type SetResult struct {
	Keys          []string          `json:"keys"`
	Val           string            `json:"val"`
	ActiveServers map[string]string `json:"active_servers"`
	ServersByKey  map[string]string `json:"servers_by_key"`
	Result        bool              `json:"result"`
	ResultCode    int               `json:"result_code"`
	ResultMessage string            `json:"result_message"`
}

// GetResult ...
type GetResult struct {
	Keys          []string          `json:"keys"`
	ActiveServers map[string]string `json:"active_servers"`
	ServersByKey  map[string]string `json:"servers_by_key"`
	Result        map[string]string `json:"result"`
	ResultCode    int               `json:"result_code"`
	ResultMessage string            `json:"result_message"`
}

// Set ...
func (c *Client) Set(keys []string, val string) (*SetResult, error) {
	res, err := http.PostForm(strings.TrimRight(c.BaseURL, "/")+"/set/",
		url.Values{
			"keys": {strings.Join(keys, ",")},
			"val":  {val},
		})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d", res.StatusCode)
	}

	var sr SetResult
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return nil, err
	}

	return &sr, nil
}

// Get ...
func (c *Client) Get(keys []string) (*GetResult, error) {
	res, err := http.Get(
		strings.TrimRight(c.BaseURL, "/") + "/get/?" + url.Values{
			"keys": {strings.Join(keys, ",")},
		}.Encode())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d", res.StatusCode)
	}

	var gr GetResult
	if err := json.NewDecoder(res.Body).Decode(&gr); err != nil {
		return nil, err
	}

	return &gr, nil
}
