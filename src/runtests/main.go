package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtests/api"
	"time"

	"runtests/cluster"
	"runtests/memcache"
	"runtests/shutdown"
)

func buildMemcacheImage(name, root string) error {
	c := exec.Command("make", "etc/memcache/bin/proxy")
	// c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Dir = root
	if err := c.Run(); err != nil {
		return err
	}

	c = exec.Command("docker",
		"build",
		"-t", name,
		"etc/memcache")
	c.Stderr = os.Stderr
	// c.Stdout = os.Stdout
	c.Dir = root

	return c.Run()
}

func startApp(
	sm *shutdown.Manager,
	root, name string,
	dumpLogs bool,
	mc []*cluster.Container) (*cluster.Container, error) {
	c, err := cluster.StartAppServer(root, name, mc)
	if err != nil {
		return nil, err
	}

	sm.AtExit(func() {
		if dumpLogs {
			c.DumpLogs()
		}
		if err := c.Stop(); err != nil {
			log.Panic(err)
		}
	})

	return c, nil
}

func getRoot() (string, error) {
	return filepath.Abs(
		filepath.Join(filepath.Dir(os.Args[0]), ".."))
}

func main() {
	flagDumpLogs := flag.Bool("dump-logs",
		false, "Dump container logs on exit")
	flag.Parse()

	root, err := getRoot()
	if err != nil {
		log.Panic(err)
	}

	if err := buildMemcacheImage(
		"memcached-proxy",
		root); err != nil {
		log.Panic(err)
	}

	s := shutdown.Create()
	defer s.Exit()

	mcs, err := memcache.StartAll(
		s,
		*flagDumpLogs,
		9000,
		"memcache-00",
		"memcache-01",
		"memcache-02")
	if err != nil {
		log.Panic(err)
	}

	_, err = startApp(
		s,
		root,
		"app",
		*flagDumpLogs,
		memcache.GetContainersFrom(mcs))
	if err != nil {
		log.Panic(err)
	}

	time.Sleep(10 * time.Second)
	runAll(&api.Client{
		BaseURL: "http://localhost:8080",
	}, mcs)
}
