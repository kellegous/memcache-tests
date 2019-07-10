package memcache

import (
	"fmt"
	"net/http"
	"net/url"
	"runtests/cluster"
)

// Client ...
type Client struct {
	Container *cluster.Container
	Name      string
	Port      int
}

func boolAsString(v bool) string {
	if v {
		return "on"
	}
	return "off"
}

// SetState ...
func (c *Client) SetState(enabled bool) error {
	res, err := http.PostForm(
		fmt.Sprintf("http://localhost:%d/", c.Port),
		url.Values{
			"state": {boolAsString(enabled)},
		})
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", res.StatusCode)
	}

	return nil
}
