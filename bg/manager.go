package bg

import (
	"context"
	"sync"
)

type Manager struct {
	BgProcesses []*Task
	stop        func()
	active      bool
	mu          sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		BgProcesses: make([]*Task, 0),
	}
}

func (m *Manager) Start(ctx context.Context, wg *sync.WaitGroup) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.active == true {
		return
	}
	m.active = true

	c, stop := context.WithCancel(ctx)
	m.stop = stop
	for _, p := range m.BgProcesses {
		p.Start(c, wg)
	}
}

func (m *Manager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stop()
	m.active = false
}

func (m *Manager) Status() []string {
	var ret []string
	for _, p := range m.BgProcesses {
		ret = append(ret, p.String())
	}
	return ret
}

func (m *Manager) Add(p *Task) {
	m.BgProcesses = append(m.BgProcesses, p)
}
