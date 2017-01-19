package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"

	"github.com/naoina/toml"
)

// Config is a project-wide struct that holds configuration information
type Config struct {
	sourceFile string
	Core       struct {
		TempDir       string
		ModuleDir     string
		RemoteBaseDir string
		KeepTempFiles bool
	}
	SSH struct {
		Username   string
		Password   string
		PrivateKey string
	}
	Hosts    []*ConfigHost
	HostsMap map[string]*ConfigHost
}

type ConfigHost struct {
	Name          string      `json:"name"`
	Address       string      `json:"address"`
	Disable       bool        `json:"disabled"`
	SSHConnection *ssh.Client `json:"-"`
}

func (c *ConfigHost) ConnectSSH(clientConfig *ssh.ClientConfig) error {
	var err error
	c.SSHConnection, err = ssh.Dial("tcp", c.Address+":22", clientConfig)
	return err
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

	if len(con.Hosts) == 0 {
		return nil, errors.New("No hosts defined")
	}

	con.sourceFile = configFile

	con.HostsMap = make(map[string]*ConfigHost)
	for _, host := range con.Hosts {
		if _, exists := con.HostsMap[host.Name]; exists {
			return nil, fmt.Errorf("Host %s duplicated in configuration", host.Name)
		}
		con.HostsMap[host.Name] = host
	}

	return setConfigDefaults(&con)
}

func setConfigDefaults(c *Config) (*Config, error) {
	// Anything not set here implies its zero value is the default

	c.Core.TempDir = setStringOrDefault(c.Core.TempDir, "./tmp")
	c.Core.ModuleDir = setStringOrDefault(c.Core.ModuleDir, "./modules")
	c.Core.RemoteBaseDir = setStringOrDefault(c.Core.RemoteBaseDir, "$HOME")

	c.SSH.Username = setStringOrDefault(c.SSH.Username, "root")
	return c, nil
}

// Given string s, if it is empty, return v else return s.
func setStringOrDefault(s, v string) string {
	if s == "" {
		return v
	}
	return s
}

// Given integer s, if it is 0, return v else return s.
func setIntOrDefault(s, v int) int {
	if s == 0 {
		return v
	}
	return s
}

func FileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}
