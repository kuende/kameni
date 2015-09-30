package main

// Config keeps config file settings
type Config struct {
	Marathon      string   `toml:"marathon"`
	EtcdServers   []string `toml:"etcd_servers"`
	addr          string   `toml:"addr"`
	kameniPrefix  string   `toml:"kameni_prefix"`
	vulcandPrefix string   `toml:"vulcand_prefix"`
}

// ListenAddr returns address to listen
func (c Config) ListenAddr() string {
	if c.addr == "" {
		return ":7373"
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
