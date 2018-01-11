#!/bin/bash
#gen:module o os:string,os_version:string,current_kernel:string,latest_kernel:string,hostname:string,uptime:string,server_time:string,path:string,installed:string

function displaytime {
  local T=$1
  local D=$((T/60/60/24))
  local H=$((T/60/60%24))
  local M=$((T/60%60))
  local S=$((T%60))

  (( $D > 0 )) && if (( $D < 10 )); then printf '0%d:' $D; else printf '%d:' $D; fi
  if (( $H < 10 )); then echo -n "0$H:"; else echo -n "$H:"; fi
  if (( $M < 10 )); then echo -n "0$M:"; else echo -n "$M:"; fi
  if (( $S < 10 )); then echo -n "0$S"; else echo -n "$S"; fi
}

distro=$(grep -oP '^NAME="?.*"?' /etc/os-release | cut -d"=" -f2 | tr -d '"')
os_version=$(grep -oP '^VERSION_ID="?.*"?' /etc/os-release | cut -d"=" -f2 | tr -d '"')
uname=$(/bin/uname -r | sed -e 's/^"//'  -e 's/"$//')

distrofamily=$(grep -oP '^ID="?.*"?' /etc/os-release | cut -d"=" -f2 | tr -d '" ')

debian() {
  latestkernel=$(dpkg -l | grep 'linux-image-[[:digit:]]' | grep ii | awk '{print $2}' | tail -n 1 | sed 's/linux-image-//')
  installed_date="$(ls -lt /var/log/installer/ | head -n2 | tail -n1 | awk -e '{print $6 " " $7 " " $8}')"
}
redhat() {
  latestkernel=$(rpm -q kernel | tail -1 | sed 's/kernel-//')
  installed_date="$(rpm -qi setup | grep Install | awk -e '{print $5 " " $4 " " $6}')"
}

case "$distrofamily" in
debian|ubuntu) debian ;;
fedora|rhel|centos|ol) redhat ;;
esac

hostname=$(/bin/hostname)
uptime_seconds=$(/bin/cat /proc/uptime | awk '{print $1}')
server_time=$(date)

echo -n "{\"os\":\"$distro\",\"os_version\":\"$os_version\",\"installed\":\"$installed_date\",\"current_kernel\":\"$uname\",\"latest_kernel\":\"$latestkernel\",\"hostname\":\"$hostname\",\"uptime\":\"$(displaytime ${uptime_seconds%.*})\",\"server_time\":\"$server_time\",\"path\":\"$PATH\"}"
