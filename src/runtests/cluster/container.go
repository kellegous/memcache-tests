package cluster

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
)

// Container ...
type Container struct {
	ID     string `json:"Id"`
	Name   string
	Config struct {
		Labels map[string]string
	}
	NetworkSettings struct {
		IPAddress string
	}
}

// LoadContainer ...
func LoadContainer(id string) (*Container, error) {
	c := exec.Command(
		"docker",
		"inspect",
		id,
		"--format={{json .}}")
	var buf bytes.Buffer
	c.Stderr = os.Stderr
	c.Stdout = &buf

	if err := c.Run(); err != nil {
		return nil, err
	}

	var cnt Container
	if err := json.NewDecoder(&buf).Decode(&cnt); err != nil {
		return nil, err
	}

	return &cnt, nil
}

// DumpLogs ...
func (c *Container) DumpLogs() error {
	cmd := exec.Command(
		"docker", "logs", c.ID)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// Stop ...
func (c *Container) Stop() error {
	return exec.Command(
		"docker",
		"stop",
		c.ID).Run()
}

// Restart ...
func (c *Container) Restart() error {
	return exec.Command(
		"docker",
		"restart",
		c.ID).Run()
}
