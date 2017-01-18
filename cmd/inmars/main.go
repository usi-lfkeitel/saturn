package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lfkeitel/inmars/src/remote"
	"github.com/lfkeitel/inmars/src/utils"
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
	// Flags
	hostStatList stringFlagList
	moduleList   stringFlagList
	configFile   string
	outputMode   string
)

const (
	remoteBasePath string = "$HOME"
	keepTempFile   bool   = false
)

func init() {
	flag.Var(&hostStatList, "limit", "Hosts from which to get stats")

	flag.Var(&moduleList, "module", "Modules to run on the hosts")
	flag.Var(&moduleList, "m", "Modules to run on the hosts")

	flag.StringVar(&configFile, "c", "", "Configuration file location")

	flag.StringVar(&outputMode, "o", "plain", "Output mode: plain, json")
}

func main() {
	flag.Parse()

	if configFile == "" {
		configFile = utils.FindConfigFile()
	}

	config, err := utils.NewConfig(configFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := remote.LoadPrivateKey(config.SSH.PrivateKey); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	checkedHosts, err := utils.CheckHosts(config, hostStatList)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	tempFileName, err := remote.GenerateScript(config, moduleList)

	if err := remote.UploadScript(config, checkedHosts, tempFileName); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	resps, err := remote.ExecuteScript(config, checkedHosts, tempFileName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if !keepTempFile {
		os.Remove(tempFileName)
	}

	printResults(resps)
}

func printResults(resps []*utils.HostResponse) {
	switch outputMode {
	case "json":
		out, err := json.Marshal(resps)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(out))
	case "plain":
		fallthrough
	default:
		for _, resp := range resps {
			fmt.Printf("Stats for %s:\n\n", resp.Host.Address)
			resp.Print()
		}
	}
}
