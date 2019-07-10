package cluster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func serverEnvFrom(cnts []*Container) (string, error) {
	m := map[string]string{}
	for _, cnt := range cnts {
		m[cnt.Name] = cnt.NetworkSettings.IPAddress
	}

	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// StartAppServer ...
func StartAppServer(
	root, name string,
	mc []*Container) (*Container, error) {

	env, err := serverEnvFrom(mc)
	if err != nil {
		return nil, err
	}

	c := exec.Command("docker",
		"run", "--rm", "-d",
		"--label=app=memcache",
		"--label=role=app",
		fmt.Sprintf("--name=%s", name),
		"-e", fmt.Sprintf("MEMCACHE_SERVERS=%s", env),
		"-v", fmt.Sprintf("%s:/app", filepath.Join(root, "pub")),
		"-v", fmt.Sprintf("%s:/etc/php/7.2/fpm/pool.d/application.conf", filepath.Join(root, "etc/application.conf")),
		"-p", "8080:80",
		"webdevops/php-nginx")

	var buf bytes.Buffer
	c.Stderr = os.Stderr
	c.Stdout = &buf

	if err := c.Run(); err != nil {
		return nil, err
	}

	return LoadContainer(strings.TrimSpace(buf.String()))
}
