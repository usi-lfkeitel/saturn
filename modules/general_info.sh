#!/bin/bash
#gen:module o os:string,kernel:string,hostname:string,uptime:string,server_time:string

function displaytime {
  local T=$1
  local D=$((T/60/60/24))
  local H=$((T/60/60%24))
  local M=$((T/60%60))
  local S=$((T%60))
  [[ $D > 0 ]] && printf '%d days ' $D
  [[ $H > 0 ]] && printf '%d hours ' $H
  [[ $M > 0 ]] && printf '%d minutes ' $M
  [[ $D > 0 || $H > 0 || $M > 0 ]] && printf 'and '
  printf '%d seconds\n' $S
}

distro=$(grep -oP '^NAME="?.*"?' /etc/os-release | cut -d"=" -f2 | tr -d '"')
uname=$(/bin/uname -r | sed -e 's/^"//'  -e 's/"$//')
hostname=$(/bin/hostname)
uptime_seconds=$(/bin/cat /proc/uptime | awk '{print $1}')
server_time=$(date)

echo -n "{\"os\":\"$distro\",\"kernel\":\"$uname\",\"hostname\":\"$hostname\",\"uptime\":\"$(displaytime ${uptime_seconds%.*})\",\"server_time\":\"$server_time\"}"
