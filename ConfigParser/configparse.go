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

type endpointConfig struct {
	backendUrl    string `yaml:"backendUrl"`
	proxyEndpoint string `yaml:"proxyEndpoint"`
	parsedUrl     *url.URL
}
type content struct {
	ProxyPort string           `yaml:"proxyPort"`
	Endpoints []endpointConfig `yaml:"endpoints"`
}
type ConfigParser struct {
	ConfigLocation string
	YamlContent    content
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
func (p *ConfigParser) GetEndpoints() []endpointConfig {
	return p.YamlContent.Endpoints
}
