package utils

import (
	"fmt"
	"log"
)

func CheckHosts(config *Config, hostList []string) (map[string]*ConfigHost, error) {
	hosts := make(map[string]*ConfigHost)

	if len(hostList) != 0 {
		for _, host := range hostList {
			configHost, ok := config.HostsMap[host]
			if !ok {
				return nil, fmt.Errorf("Host %s not found", host)
			}
			if configHost.Disable {
				log.Printf("%s disabled in config, skipping", host)
				continue
			}
			if configHost.Name == "" {
				configHost.Name = configHost.Address
			}
			hosts[host] = configHost
		}
	} else {
		for hostname, host := range config.HostsMap {
			if host.Disable {
				log.Printf("%s disabled in config, skipping", hostname)
				continue
			}
			if host.Name == "" {
				host.Name = host.Address
			}
			hosts[hostname] = host
		}
	}

	return hosts, nil
}
