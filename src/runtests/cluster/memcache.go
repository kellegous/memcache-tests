package cluster

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// StartMemcached ...
func StartMemcached(name string, port int) (*Container, error) {
	c := exec.Command("docker",
		"run", "-t", "--rm", "-d",
		"--label=app=memcache",
		"--label=role=memcached",
		fmt.Sprintf("--name=%s", name),
		"-p", fmt.Sprintf("%d:8080", port),
		"memcached-proxy")
	var buf bytes.Buffer
	c.Stderr = os.Stderr
	c.Stdout = &buf

	if err := c.Run(); err != nil {
		return nil, err
	}

	return LoadContainer(strings.TrimSpace(buf.String()))
}
