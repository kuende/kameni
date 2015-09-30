package main

// Config keeps config file settings
type Config struct {
	Marathon      string
	EtcdServers   []string
	addr          string
	kameniPrefix  string
	vulcandPrefix string
}

// ListenAddr returns address to listen
func (c Config) ListenAddr() string {
	if c.addr == "" {
		return ":7334"
	}

	return c.addr
}

// KameniPrefix returns kameni etcd namespace or "kameni"
func (c Config) KameniPrefix() string {
	if c.kameniPrefix == "" {
		return "kameni"
	}

	return c.kameniPrefix
}

// VulcandPrefix returns vulcand etcd namespace or "vulcand"
func (c Config) VulcandPrefix() string {
	if c.vulcandPrefix == "" {
		return "vulcand"
	}

	return c.kameniPrefix
}
