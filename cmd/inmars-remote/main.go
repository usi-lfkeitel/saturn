package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/lfkeitel/inmars/src/utils"

	"golang.org/x/crypto/ssh"
)

type stringFlagList []string

func (i *stringFlagList) String() string {
	return "List of hosts"
}

func (i *stringFlagList) Set(value string) error {
	set := strings.Split(value, ",")
	*i = append(*i, set...)
	return nil
}

var (
	hostStatList    stringFlagList
	moduleList      stringFlagList
	configFile      string
	sshClientConfig *ssh.ClientConfig
	checkedHosts    map[string]*utils.ConfigHost

	config *utils.Config
)

const (
	remoteBasePath string = "$HOME"
)

func init() {
	flag.Var(&hostStatList, "limit", "Hosts from which to get stats")

	flag.Var(&moduleList, "module", "Modules to run on the hosts")
	flag.Var(&moduleList, "m", "Modules to run on the hosts")

	flag.StringVar(&configFile, "c", "", "Configuration file location")
}

func main() {
	flag.Parse()

	if configFile == "" {
		configFile = utils.FindConfigFile()
	}

	var err error
	config, err = utils.NewConfig(configFile)
	if err != nil {
		panic(err)
	}

	if err := loadPrivateKey(); err != nil {
		panic(err)
	}

	if err := checkHosts(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	tempFile, err := ioutil.TempFile(config.Core.TempDir, "")
	if err != nil {
		panic(err)
	}
	tempFileName := tempFile.Name()

	if err := generateRemoteScript(tempFile, config.Core.ModuleDir, moduleList); err != nil {
		tempFile.Close()
		panic(err)
	}
	tempFile.Close()

	if err := os.Chmod(tempFileName, 0755); err != nil {
		fmt.Println("Failed make script executable")
		os.Exit(1)
	}

	if err := uploadScript(tempFileName); err != nil {
		panic(err)
	}

	if err := executeScript(tempFileName); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	printResults()
}

func loadPrivateKey() error {
	sshPrivateKey, err := ioutil.ReadFile(config.SSH.PrivateKey)
	if err != nil {
		return err
	}

	signer, err := ssh.ParsePrivateKey(sshPrivateKey)
	if err != nil {
		return err
	}

	sshClientConfig = &ssh.ClientConfig{
		User: "lfkeitel",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}
	return nil
}

func checkHosts() error {
	hosts := make(map[string]*utils.ConfigHost)

	if len(hostStatList) != 0 {
		realHosts := config.GetRealHosts()
		for _, host := range hostStatList {
			configHost, ok := realHosts[host]
			if !ok {
				return fmt.Errorf("Host %s not found", host)
			}
			if configHost.Disable {
				continue
			}
			hosts[host] = configHost
		}
	} else {
		hosts = config.GetRealHosts()
	}

	checkedHosts = hosts
	return nil
}

func generateRemoteScript(file *os.File, modulesDir string, modules []string) error {
	file.WriteString("#!/bin/bash\n\n")

	file.WriteString("MODULES=")
	file.WriteString(`(` + strings.Join(modules, " ") + `)`)

	file.WriteString(`
main() {
	echo -n '{'

	i=1
	for var in "${MODULES[@]}"; do
		echo -n "\"$var\": "
		echo -n $($var)
		if [ $i -lt ${#MODULES[@]} ]; then
				i=$[i + 1]
				echo -n ', '
		fi
	done

	echo -n '}'
}

`)

	goodModules := make(map[string]bool)

	for _, module := range modules {
		moduleFile := filepath.Join(modulesDir, module+".sh")

		if !utils.FileExists(moduleFile) {
			fmt.Printf("Module %s not found\n", module)
			continue
		}

		m, err := ioutil.ReadFile(moduleFile)
		if err != nil {
			fmt.Println(err)
			continue
		}

		goodModules[module] = true

		file.WriteString(module + "() {\n")
		file.Write(m)
		file.WriteString("\n}\n\n")
	}

	file.WriteString(`main

if [ "$1" = "-d" ]; then
	rm "$0"
fi
`)

	return nil
}

func uploadScript(genFilename string) error {
	f, err := os.Open(genFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	s, err := f.Stat()
	if err != nil {
		return err
	}

	for _, host := range checkedHosts {
		if host.Disable {
			continue
		}

		_, err = f.Seek(0, 0) // rewind file reader
		if err != nil {
			return err
		}

		if host.Address == "localhost" || host.Address == "127.0.0.1" {
			if err := uploadLocalScript(f, s); err != nil {
				fmt.Println(err.Error())
			}
			continue
		}

		if err := uploadRemoteScript(host, f, s); err != nil {
			fmt.Println(err.Error())
		}
	}

	return nil
}

func uploadRemoteScript(host *utils.ConfigHost, f *os.File, s os.FileInfo) error {
	if host.SSHConnection == nil {
		if err := host.ConnectSSH(sshClientConfig); err != nil {
			return err
		}
	}

	session, err := host.SSHConnection.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()
		fmt.Fprintln(w, "D0755", 0, ".inmars") // mkdir
		fmt.Fprintf(w, "C%#o %d %s\n", s.Mode().Perm(), s.Size(), path.Base(f.Name()))
		io.Copy(w, f)
		fmt.Fprint(w, "\x00")
	}()

	cmd := fmt.Sprintf("scp -rt %s", remoteBasePath)
	return session.Run(cmd)
}

func uploadLocalScript(f *os.File, s os.FileInfo) error {
	expandedPath := os.ExpandEnv(remoteBasePath)

	if err := os.MkdirAll(path.Join(expandedPath, ".inmars"), 0755); err != nil {
		return err
	}

	out, err := os.OpenFile(
		path.Join(expandedPath, ".inmars", path.Base(f.Name())),
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		0755,
	)
	if err != nil {
		return err
	}
	defer out.Close()

	io.Copy(out, f)
	return nil
}

func executeScript(filename string) error {
	filename = path.Base(filename)
	for _, host := range checkedHosts {
		if host.Disable {
			continue
		}

		if host.SSHConnection == nil {
			if err := host.ConnectSSH(sshClientConfig); err != nil {
				return err
			}
		}

		session, err := host.SSHConnection.NewSession()
		if err != nil {
			return err
		}

		var stdoutBuf bytes.Buffer
		var stderrBuf bytes.Buffer
		session.Stdout = &stdoutBuf
		session.Stderr = &stderrBuf

		if err := session.Run("/bin/bash /home/lfkeitel/.inmars/" + filename + " -d"); err != nil {
			fmt.Println(err.Error())
			fmt.Println(stderrBuf.String())
			session.Close()
			continue
		}
		session.Close()

		var response utils.HostResponse
		if err := json.Unmarshal(stdoutBuf.Bytes(), &response); err != nil {
			fmt.Println(err.Error())
			continue
		}

		host.SetResponse(&response)
	}
	return nil
}

func printResults() {
	for hostname, host := range checkedHosts {
		if host.Disable {
			continue
		}

		fmt.Printf("Stats for %s:\n\n", hostname)
		host.Response.Print()
	}
}
