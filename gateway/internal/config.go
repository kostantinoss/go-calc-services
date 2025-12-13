package internal

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

// This is the type holding the configurable values
// for rhe gateway service
type Config struct {
	ApiGateway ApiGatewayConfig `yaml:"api_gateway"`
}

type ApiGatewayConfig struct {
	GatewayServer GatewayServerConfig `yaml:"gateway_server"`
	TargetServers []TargetServer      `yaml:"target_servers"`
	Routes        []Route             `yaml:"routing"`
}

type GatewayServerConfig struct {
	Port string `yaml:"port"`
}

type TargetServer struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

type Route struct {
	Path    string   `yaml:"path"`
	Methods []string `yaml:"methods"`
	Server  string   `yaml:"server"`
}

// Load configuration from config.yaml
func (c *Config) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, c)
}

func (c *Config) GetTargetServerUrl(name string) (string, error) {
	for _, server := range c.ApiGateway.TargetServers {
		if server.Name == name {
			return server.Url, nil
		}
	}

	return "", fmt.Errorf("target server %s not found", name)
}
