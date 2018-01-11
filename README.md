# Saturn

Saturn is a system that allows system administrators to quickly gather information and statistics about their systems.

## Install

All you need to get started is the release package found under Releases. It contains the binary for the cli client, and the license and readme files.

## Build

Run `go get github.com/usi-lfkeitel/saturn/cmd/saturn` or clone the repo to the appropriate place in your $GOPATH and run `make generate && make`. A binary will be created at `./bin/saturn`.

## Configuration

Saturn requires a TOML configuration file to run. This specifies SSH connection parameters, application settings, and hosts. Here's a sample configuration:

```
[core]

[SSH]
Username = "jsmith"
#Password = ""
PrivateKey = "/home/jsmith/.ssh/id_rsa"
UseAgent = true # Use an SSH agent with the environment variable SSH_AUTH_SOCK

[[hosts]]
name = "Localhost"
address = "localhost"

[[hosts]]
name = "Server1"
address = "192.168.1.1"

[[hosts]]
name = "Server2"
address = "192.168.1.2"
```

Please refer to the example configuration file for details on all available settings.

## Running Saturn

So, you have your configuration ready with all your hosts defined. Now let's actually get some data. Here's all the flags that Saturn can use:

- `-c`:  Configuration file to use.
- `-d,-dd`: Enable debug mode
- `-limit host1,host2`:  Comma separated list of hosts to run on. Use this option for quick data grabs without having to disable
all clients in a configuration. Listed clients will be used except those disabled by the config.
- `-list`: List available modules.
- `-m,-module module1, module2`: Comma separated list of modules to run.
- `-o MODE`:  Output mode. One of: json, plain (default)
- `-s`:  Print shorter output. Only affects plain output mode.
- `-v`:  Print version information.

Here's a quick demo:

`./saturn -c config.toml -limit Server1 -m io_stats,cpu_info`

This will run the io\_stats and cpu\_info modules only on Server1 (unless it's disabled in the config). It will then output
information from the modules in a structured, readable format.

## Modules

Saturn is based on a set of modules which are nothing more than shell scripts that output JSON. The modules are embedded into the
binary at build time. Custom modules can be used by placing them in the modules directory set in the configuration file. Module
named are based on their file name so it's recommended to keep them sort and use only alphanumeric characters, underscores, or
hyphens. Other characters are not guaranteed to work. Custom modules cannot override builtin modules.

A module file must output only JSON to stdout. A module may output either a JSON array or object. A module file must have a
`gen:module` comment on the second line so it can be added to the binary. The syntax of this line is as follows:
`#gen:module [a|o] key1:type1,key2:type2,...`. The first parameter is either "a" or "o" depending on if the module outputs a JSON
[a]rray or [o]bject. The letter must be lower-case. The key type pairs specify the schema of the JSON. Keys cannot start with a number
or a comma/semicolon. Supported types are `string`, `int`, `bool`, and `float64`. The key names will be transformed as needed to
conform with Go's variable naming syntax.

Here's a list of first-party modules:

- **arp_cache**:
- **bandwidth**:
- **common_applications**:
- **cpu_info**:
- **cpu_intensive_processes**:
- **cpu_temp**:
- **cpu_utilization**:
- **cron_history**:
- **current_ram**:
- **disk_partitions**:
- **docker_processes**:
- **download_transfer_rate**:
- **enabled_services**:
- **general_info**:
- **io_stats**:
- **listening_ports**:
- **load_avg**:
- **logged_in_users**:
- **memcached**:
- **memory_info**:
- **needs_upgrades**:
- **network_connections**:
- **network_interfaces**:
- **number_of_cpu_cores**:
- **ping**:
- **pm2_stats**:
- **ram_intensive_processes**:
- **recent_account_logins**:
- **redis**:
- **scheduled_crons**:
- **swap**:
- **upload_transfer_rate**:
- **user_accounts**:
