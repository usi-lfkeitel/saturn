package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/naoina/toml"
)

// Config is a project-wide struct that holds configuration information
type Config struct {
	sourceFile string
	Core       struct {
		Debug         bool
		SpecialDebug  bool
		TempDir       string
		ModuleDir     string
		RemoteBaseDir string
		KeepTempFiles bool
		HaltOnError   bool
	}
	SSH struct {
		Username   string
		Password   string
		PrivateKey string
		Timeout    string
	}
	Hosts    []*ConfigHost
	HostsMap map[string]*ConfigHost
}

type ConfigHost struct {
	Name          string            `json:"name"`
	Address       string            `json:"address"`
	Username      string            `json:"username"`
	Password      string            `json:"-"`
	PrivateKey    string            `json:"-"`
	Disable       bool              `json:"disabled"`
	SSHConnection *ssh.Client       `json:"-"`
	SSHConfig     *ssh.ClientConfig `json:"-"`
}

func (c *ConfigHost) ConnectSSH(clientConfig *ssh.ClientConfig) error {
	if c.SSHConfig == nil {
		// Make a copy for this host
		hostSSHConfig := &ssh.ClientConfig{}
		*hostSSHConfig = *clientConfig

		// Replace the global connection options with host specific ones if given
		if c.Password != "" || c.PrivateKey != "" {
			hostSSHConfig.Auth = make([]ssh.AuthMethod, 0, 1)

			if c.Password != "" {
				hostSSHConfig.Auth = append(hostSSHConfig.Auth, ssh.Password(c.Password))
			}

			if c.PrivateKey != "" {
				sshPrivateKey, err := ioutil.ReadFile(c.PrivateKey)
				if err != nil {
					return err
				}

				signer, err := ssh.ParsePrivateKey(sshPrivateKey)
				if err != nil {
					return err
				}
				hostSSHConfig.Auth = append(hostSSHConfig.Auth, ssh.PublicKeys(signer))
			}
		}

		if c.Username != "" {
			hostSSHConfig.User = c.Username
		}

		c.SSHConfig = hostSSHConfig
	}

	if c.SSHConfig.User == "" || len(c.SSHConfig.Auth) == 0 {
		return fmt.Errorf("Username or password not given for host %s", c.Name)
	}

	var err error
	c.SSHConnection, err = ssh.Dial("tcp", c.Address+":22", c.SSHConfig)
	if err != nil {
		return fmt.Errorf("Login failed on %s. Check username or password.", c.Name)
	}
	return nil
}

// FindConfigFile will locate the a configuration file by looking at the following places
// and choosing the first: SATURN_CONFIG env variable, $PWD/config.toml, $PWD/config/config.toml,
// $HOME/.saturn/config.toml, and /etc/saturn/config.toml.
func FindConfigFile() string {
	if os.Getenv("SATURN_CONFIG") != "" && FileExists(os.Getenv("SATURN_CONFIG")) {
		return os.Getenv("SATURN_CONFIG")
	} else if FileExists("./config.toml") {
		return "./config.toml"
	} else if FileExists("./config/config.toml") {
		return "./config/config.toml"
	} else if FileExists(os.ExpandEnv("$HOME/.saturn/config.toml")) {
		return os.ExpandEnv("$HOME/.saturn/config.toml")
	} else if FileExists("/etc/saturn/config.toml") {
		return "/etc/saturn/config.toml"
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
	c.SSH.Timeout = setStringOrDefault(c.SSH.Timeout, "20s")
	if _, err := time.ParseDuration(c.SSH.Timeout); err == nil {
		c.SSH.Timeout = "20s"
	}
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
