package memcache

import (
	"log"
	"runtests/cluster"
	"runtests/shutdown"
)

// StartAll ...
func StartAll(
	sm *shutdown.Manager,
	dumpLogs bool,
	basePort int,
	names ...string) ([]*Client, error) {
	clients := make([]*Client, 0, len(names))
	for i, name := range names {
		cnt, err := cluster.StartMemcached(name, basePort+i)
		if err != nil {
			return nil, err
		}

		sm.AtExit(func() {
			if dumpLogs {
				cnt.DumpLogs()
			}
			if err := cnt.Stop(); err != nil {
				log.Panic(err)
			}
		})

		clients = append(clients, &Client{
			Container: cnt,
			Name:      name,
			Port:      basePort + i,
		})
	}

	return clients, nil
}

// GetContainersFrom ...
func GetContainersFrom(clients []*Client) []*cluster.Container {
	cnts := make([]*cluster.Container, 0, len(clients))
	for _, client := range clients {
		cnts = append(cnts, client.Container)
	}
	return cnts
}
