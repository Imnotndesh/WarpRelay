package ConfigParser

/*
what am thinking is the following format for endpoints
<config.YAML>

proxyPort: ...
endpoints:
	- backendUrl: ..
		proxyEndpoint: ..
*/
import (
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
)

type EndpointConfig struct {
	BackendUrl    string `yaml:"backendUrl"`
	ProxyEndpoint string `yaml:"proxyEndpoint"`
	ParsedUrl     *url.URL
}
type Content struct {
	ProxyPort string           `yaml:"proxyPort"`
	Endpoints []EndpointConfig `yaml:"endpoints"`
}
type ConfigParser struct {
	ConfigLocation string
	YamlContent    Content
}

func (p *ConfigParser) ParseConfig() error {
	yamlFile, err := os.ReadFile(p.ConfigLocation)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &p.YamlContent)
	if err != nil {
		return err
	}
	return nil
}
func (p *ConfigParser) GetProxyPort() string {
	return p.YamlContent.ProxyPort
}
func (p *ConfigParser) GetEndpoints() []EndpointConfig {
	return p.YamlContent.Endpoints
}
