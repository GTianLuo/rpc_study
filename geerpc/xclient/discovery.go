package xclient

import (
	"errors"
	"math"
	"math/rand"
	"sync"
	"time"
)

type SelectMode int

const (
	RandomSelect SelectMode = iota
	RoundRobinSelect
)

type Discovery interface {
	Refresh() error
	Update(servers []string) error
	Get(mode SelectMode) (string, error)
	GetAll() ([]string, error)
}

type MultiServerDiscovery struct {
	r       *rand.Rand
	mu      sync.Mutex
	servers []string
	index   int
}

func NewMultiServerDiscovery(servers []string) *MultiServerDiscovery {
	discovery := &MultiServerDiscovery{
		servers: servers,
		r:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	discovery.index = discovery.r.Intn(math.MaxInt32)
	return discovery
}

func (m *MultiServerDiscovery) Refresh() error {
	return nil
}

func (m *MultiServerDiscovery) Update(servers []string) error {
	m.mu.Lock()
	m.servers = servers
	m.mu.Unlock()
	return nil
}

func (m *MultiServerDiscovery) Get(mode SelectMode) (string, error) {

	n := len(m.servers)
	if n == 0 {
		return "", errors.New("rpc discovery: no available server")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	switch mode {
	case RandomSelect:
		return m.servers[m.r.Intn(n)], nil
	case RoundRobinSelect:
		s := m.servers[m.index%n]
		m.index = (m.index + 1) % n
		return s, nil
	default:
		return "", errors.New("rpc discovery: not support this select mode")
	}
}

func (m *MultiServerDiscovery) GetAll() ([]string, error) {
	return m.servers, nil
}
