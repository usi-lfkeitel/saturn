# Saturn

Saturn is a system that allows system administrators to quickly gather information and statistics about their systems.

## Install

All you need to get started is the release package found under Releases. It contains the binary for the cli client, and the license and readme files.

## Build

Clone the repo and run `make generate && make`. A binary will be created at `./bin/saturn`.

## Configuration

Saturn requires a TOML configuration file to run. This specifies SSH connection parameters, application settings, and hosts. Here's a sample configuration:

```
[core]

[SSH]
Username = "jsmith"
#Password = ""
PrivateKey = "/home/jsmith/.ssh/id_rsa"

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

Hosts have the following settings:

- name: string - The display name of the host
- address: string - The hostname or ip address of the host
- disable: bool - Disables the host

Using the disable option can be nice so you don't have to delete a host entry. This can be used as a maintenance mode.

SSH supports both password and key authentication. If using a key, the setting must be the file path to the private key.

## Running Saturn

So, you have your configuration ready with all your hosts defined. Now let's actually get some data. Here's all the flags that Saturn can use:

- -limit  Comma separated list of hosts to run on. Use this option for quick data grabs without having to disable all clients in a configuration. Listed clients will be used except those disabled by the config.
- -m, -module  Comma separated list of modules to run.
- -c  Configuration file to use.
- -o  Output mode. One of: json, plain (default)
- -s  Print shorter output. Only affects plain output mode.
- -v  Print version information.

Here's a quick demo:

`./saturn -c config.toml -limit Server1 -m io_stats,cpu_info`

This will run the io_stats and cpu_info modules only on Server1 (unless it's disabled in the config). It will then output information from the modules in a structured, readable format.

## Modules

Saturn is based on a set of modules which are nothing more than shell scripts that output JSON. The modules are embedded into the binary at build time. At this time there's no way to run or use custom modules. If you would like to use a different module, you'll need to recompile the binary.

A module file must output only JSON to stdout. A module may output either a JSON array or object. A module file must have a `gen:module` comment on the second line so it can be added to the binary. The syntax of this line is as follows: `#gen:module [a|o] key1:type1,key2:type2,...`. The first parameter is either "a" or "o" depending on if the module outputs a JSON Array or Object. The letter must be lower-case. The key type pairs specify the schema of the JSON. Keys cannot start with a number or a comma/semicolon. Supported types are `string`, `int`, `bool`, and `float64`. The key names will be transformed as needed to conform with Go's variable naming syntax.

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
- **general_info**:
- **io_stats**:
- **ip_addresses**:
- **load_avg**:
- **logged_in_users**:
- **memcached**:
- **memory_info**:
- **needs_upgrades**:
- **network_connections**:
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
