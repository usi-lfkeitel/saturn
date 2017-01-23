#!/bin/bash
#gen:module o os:string,current_kernel:string,latest_kernel:string,hostname:string,uptime:string,server_time:string,path:string

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

distrofamily=$(grep -oP '^ID="?.*"?' /etc/os-release | cut -d"=" -f2 | tr -d '" ')

debian() {
	latestkernel=$(dpkg -l | grep linux-headers | grep ii | awk '{print $3}' | tail -n 1)
}
redhat() {
	latestkernel=$(rpm -q kernel | tail -1 | sed 's/kernel-//')
}

case "$distrofamily" in
debian|ubuntu) debian ;;
fedora|rhel|centos|ol) redhat ;;
esac

hostname=$(/bin/hostname)
uptime_seconds=$(/bin/cat /proc/uptime | awk '{print $1}')
server_time=$(date)

echo -n "{\"os\":\"$distro\",\"current_kernel\":\"$uname\",\"latest_kernel\":\"$latestkernel\",\"hostname\":\"$hostname\",\"uptime\":\"$(displaytime ${uptime_seconds%.*})\",\"server_time\":\"$server_time\",\"path\":\"$PATH\"}"
