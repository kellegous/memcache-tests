package shutdown

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Manager ...
type Manager struct {
	fns []func()
	lck sync.Mutex
}

// Exit ...
func (m *Manager) Exit() {
	m.lck.Lock()
	defer m.lck.Unlock()
	for _, fn := range m.fns {
		fn()
	}
	m.fns = nil
}

// AtExit ...
func (m *Manager) AtExit(fn func()) {
	m.lck.Lock()
	defer m.lck.Unlock()
	m.fns = append(m.fns, fn)
}

// Create ...
func Create() *Manager {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	var mgr Manager
	go func() {
		<-ch
		mgr.Exit()
		os.Exit(0)
	}()
	return &mgr
}
