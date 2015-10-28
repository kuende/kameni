package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kuende/kameni/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"github.com/kuende/kameni/Godeps/_workspace/src/github.com/mailgun/log"
)

// VulcandServer keeps marathon app data
type VulcandServer struct {
	ID  string `json:"-"`
	URL string `json:"URL"`
}

// BackendConfig is struct fetched from etcd
type BackendConfig struct {
	BackendID   string `json:"backend_id"`
	BackendType string `json:"type"`
}

var (
	etcdClient *etcd.Client
	// ErrEtcdValueNotPresent is returned when the value is not present in etcd
	ErrEtcdValueNotPresent = errors.New("Etcd value was not present")
)

// Type returns Backend type, vulcand default
func (b *BackendConfig) Type() string {
	if b.BackendType == "" {
		return "vulcand"
	}

	return b.BackendType
}

func addVulcandServer(appID string, server VulcandServer) error {
	backend, err := fetchBackend(appID)

	if err != nil {
		return err
	}

	path := serverPath(backend, server.ID)

	value, _ := json.Marshal(server)

	err = etcdSet(path, value)

	if err != nil {
		log.Errorf("Error setting key in etcd: %v", err)
		return err
	}

	log.Infof("Added server: %s to backend: %s, url: %s", server.ID, backend.BackendID, server.URL)
	return nil
}

func removeVulcandServer(appID string, server VulcandServer) error {
	backend, err := fetchBackend(appID)

	if err != nil {
		log.Errorf("Error fetching backend from etcd: %v", err)
		return err
	}

	path := serverPath(backend, server.ID)

	err = etcdDelete(path)

	if err != nil {
		log.Errorf("Error delete key from etcd: %v", err)
		return err
	}

	log.Infof("Removed server: %s from backend: %s, url: %s", server.ID, backend.BackendID, server.URL)
	return nil
}

func fetchBackend(appID string) (*BackendConfig, error) {
	value, err := etcdGet(appPath(appID))

	if err != nil {
		log.Errorf("Error fetching backend from etcd: %v", err)
		return nil, err
	}

	backend := BackendConfig{}
	err = json.Unmarshal(value, &backend)

	return &backend, err
}

func appPath(appID string) string {
	return fmt.Sprintf("/%s/apps/%s", config.KameniPrefix(), appID)
}

func backendPath(backendID string) string {
	return fmt.Sprintf("/%s/backends/%s/backend", config.VulcandPrefix(), backendID)
}

func serverPath(backendID *BackendConfig, serverID string) string {
	if backendID.Type() == "vulcand" {
		return fmt.Sprintf("/%s/backends/%s/servers/%s", config.VulcandPrefix(), backendID, serverID)
	}

	return fmt.Sprintf("%s/%s", backendID.BackendID, serverID)
}

func setupEtcd() {
	etcdClient = etcd.NewClient(config.EtcdServers)
}

func etcdGet(key string) ([]byte, error) {
	value, err := etcdClient.Get(key, false, false)

	if err != nil {
		return nil, err
	}

	if value.Node == nil {
		return nil, ErrEtcdValueNotPresent
	}

	return []byte(value.Node.Value), nil
}

func etcdSet(key string, value []byte) error {
	_, err := etcdClient.Set(key, string(value), 0)

	return err
}

func etcdDelete(key string) error {
	_, err := etcdClient.Delete(key, false)

	return err
}
