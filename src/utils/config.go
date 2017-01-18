package utils

import (
	"errors"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"

	"github.com/naoina/toml"
)

// Config is a project-wide struct that holds configuration information
type Config struct {
	sourceFile string
	Core       struct {
		TempDir   string
		ModuleDir string
	}
	SSH struct {
		PrivateKey string
	}
	Hosts map[string]ConfigHost
	hosts map[string]*ConfigHost
}

func (c *Config) GetRealHosts() map[string]*ConfigHost {
	if c.hosts == nil {
		c.hosts = make(map[string]*ConfigHost)
		for hostname, host := range c.Hosts {
			c.hosts[hostname] = &ConfigHost{
				Address: host.Address,
				Disable: host.Disable,
			}
		}
	}
	return c.hosts
}

type ConfigHost struct {
	Address       string
	Disable       bool
	SSHConnection *ssh.Client
	Response      *HostResponse
}

func (c *ConfigHost) ConnectSSH(clientConfig *ssh.ClientConfig) error {
	var err error
	c.SSHConnection, err = ssh.Dial("tcp", c.Address+":22", clientConfig)
	return err
}

func (c *ConfigHost) SetResponse(r *HostResponse) {
	c.Response = r
}

// FindConfigFile will locate the a configuration file by looking at the following places
// and choosing the first: INMARS_CONFIG env variable, $PWD/config.toml, $PWD/config/config.toml,
// $HOME/.inmars/config.toml, and /etc/inmars/config.toml.
func FindConfigFile() string {
	if os.Getenv("INMARS_CONFIG") != "" && FileExists(os.Getenv("INMARS_CONFIG")) {
		return os.Getenv("INMARS_CONFIG")
	} else if FileExists("./config.toml") {
		return "./config.toml"
	} else if FileExists("./config/config.toml") {
		return "./config/config.toml"
	} else if FileExists(os.ExpandEnv("$HOME/.inmars/config.toml")) {
		return os.ExpandEnv("$HOME/.inmars/config.toml")
	} else if FileExists("/etc/inmars/config.toml") {
		return "/etc/inmars/config.toml"
	}
	return ""
}

// NewConfig will load a Config object using the given TOML file.
func NewConfig(configFile string) (conf *Config, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()

	if configFile == "" {
		configFile = "config.toml"
	}

	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var con Config
	if err := toml.Unmarshal(buf, &con); err != nil {
		return nil, err
	}
	con.sourceFile = configFile
	return &con, nil
}

func FileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}
