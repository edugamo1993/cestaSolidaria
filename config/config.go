package config

import (
	"fmt"
	"go-solidary/mongo"
	"io/ioutil"
)

/*Config Type for config data*/
type Config struct {
	Server ServerConfig `json:"server"`
	mongo.Mongo
}
type ServerConfig struct {
	TLS      bool   `json:"tls"`
	CertPath string `json:"certPath"`
	Cert     string `json:"cert"`
	KeyPath  string `json:"keyPath"`
	Key      string `json:"key"`
}

//IsHTTPS returns true if is https and false if is http
func (c *Config) IsHTTPS() bool {
	return c.Server.TLS
}

//GetCertString returns string of the cert
func (c *Config) GetCertString() ([]byte, error) {
	if c.Server.Cert != "" {
		return []byte(c.Server.Cert), nil
	}
	if c.Server.CertPath == "" {
		return []byte(""), fmt.Errorf("No certs found in config file")
	}
	f, err := ioutil.ReadFile(c.Server.CertPath)
	if err != nil {
		return []byte(""), err
	}
	return f, nil
}

//GetKeyString returns string of the key
func (c *Config) GetKeyString() ([]byte, error) {
	if c.Server.Key != "" {
		return []byte(c.Server.Key), nil
	}
	if c.Server.KeyPath == "" {
		return []byte(""), fmt.Errorf("No certs found in config file")
	}
	f, err := ioutil.ReadFile(c.Server.KeyPath)
	if err != nil {
		return []byte(""), err
	}
	return f, nil
}
