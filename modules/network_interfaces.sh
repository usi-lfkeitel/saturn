#!/bin/bash
#gen:module a interface:string,ipv4_address:string,ipv6_address:string,mac_address:string,broadcast:string,subnet_mask:string

grepCmd=`which grep 2>/dev/null`
ipCmd=`which ip 2>/dev/null`
ifconfigCmd=`which ifconfig 2>/dev/null`

# From: https://forum.openwrt.org/viewtopic.php?pid=220781#p220781
# Which is a consolidated version from: https://forums.gentoo.org/viewtopic-t-888736-start-0.html
cdr2mask () { # Takes one argument which mask number (8, 16, 24, 32)
   # Number of args to shift, 255..255, first non-255 byte, zeroes
   set -- $(( 5 - ($1 / 8) )) 255 255 255 255 $(( (255 << (8 - ($1 % 8))) & 255 )) 0 0 0
   [ $1 -gt 1 ] && shift $1 || shift
   echo ${1-0}.${2-0}.${3-0}.${4-0}
}

useIPCmd() {
    interfaces=$($ipCmd address show | $grepCmd -oP '^\d+\: .*?\:' | cut -d' ' -f2 | tr -d ':')

    for item in $interfaces; do
        interface=$($ipCmd address show dev $item)
        if [[ ! "$interface" =~ inet ]]; then
            continue
        fi

        address4=''
        cidrmask=0
        if [[ $interface =~ inet[[:space:]](([[:digit:]]{1,3}\.?){4}/[[:digit:]]{1,2}) ]]; then
            address4=$(echo ${BASH_REMATCH[1]} | cut -d'/' -f1)
            cidrmask=$(echo ${BASH_REMATCH[1]} | cut -d'/' -f2)
        fi

        address6=''
        if [[ $interface =~ inet6[[:space:]]([a-f0-9\:]+) ]]; then
            address6=${BASH_REMATCH[1]}
        fi

        macaddr=''
        if [[ $interface =~ link/ether[[:space:]](([a-f0-9]{2}\:?){6}) ]]; then
            macaddr=${BASH_REMATCH[1]}
        fi

        broadcast=''
        if [[ $interface =~ brd[[:space:]](([[:digit:]]{1,3}\.?){4}) ]]; then
            broadcast=${BASH_REMATCH[1]}
        fi

        subnetmask="$(cdr2mask $cidrmask)"

        json="$json{\"interface\":\"$item\",\"ipv4_address\":\"$address4\",\"ipv6_address\":\"$address6\",\"mac_address\":\"$macaddr\",\"broadcast\":\"$broadcast\",\"subnet_mask\":\"$subnetmask\"},"
    done
}

useIFConfigCmd() {
    for item in $($ifconfigCmd | $grepCmd -oP "^[a-zA-Z0-9:\-]*"); do
        address=""
        if [ "${item: -1}" = ":" ]; then
            address=`$ifconfigCmd ${item: -1} | $grepCmd -oP 'inet \K[\d\.]+'`
        else
            address=`$ifconfigCmd $item | $grepCmd -oP 't addr:\K[\d.]+'`
        fi
        json="$json{\"interface\":\"$item\",\"ipv4_address\":\"$address\"},"
    done
}

if [ -n "$ipCmd" ]; then
    useIPCmd
elif [ -n "$ifconfigCmd" ]; then
    useIFConfigCmd
fi

echo -n '['${json%?}']'
