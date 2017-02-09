#!/bin/bash
#gen:module2 a
#type address
#  address string
#  broadcast string
#  mask string
#endtype
#key interface string
#key mac_address string
#key ipv4 []address
#key ipv6 []address
#!gen:module2

PATH="$PATH:/usr/sbin:/sbin"
grepCmd=`which grep 2>/dev/null`
ipCmd=`which ip 2>/dev/null`
ifconfigCmd=`which ifconfig 2>/dev/null`

# From: https://forum.openwrt.org/viewtopic.php?pid=220781#p220781
# Which is a consolidated version from: https://forums.gentoo.org/viewtopic-t-888736-start-0.html
cdr2mask () { # Takes one argument which mask number (8, 16, 24, 32)
    if [ -z "$1" ]; then return; fi
    # Number of args to shift, 255..255, first non-255 byte, zeroes
    set -- $(( 5 - ($1 / 8) )) 255 255 255 255 $(( (255 << (8 - ($1 % 8))) & 255 )) 0 0 0
    [ $1 -gt 1 ] && shift $1 || shift
    echo ${1-0}.${2-0}.${3-0}.${4-0}
}

useIPCmd() {
    interfaces=$(ip address show | $grepCmd -oP '^\d+\: .*?\:' | cut -d' ' -f2 | tr -d ':')

    for item in $interfaces; do
        item=$(echo "$item" | cut -d'@' -f1) # Get the main name of a subinterface
        interface=$(ip address show dev $item)
        if [[ ! "$interface" =~ inet ]]; then
            continue
        fi

        # Get MAC address
        macaddr=''
        if [[ $interface =~ link/ether[[:space:]](([a-f0-9]{2}\:?){6}) ]]; then
            macaddr=${BASH_REMATCH[1]}
        fi

        # Get lists of IPv4 and v6 addresses
        v4addresses=$(echo "$interface" | $grepCmd -E 'inet ([[:digit:]]{1,3}\.?){4}/[[:digit:]]{1,2}')
        v6addresses=$(echo "$interface" | $grepCmd -oE 'inet6 ([a-f0-9\:]+)')

        json="$json{\"interface\":\"$item\",\"mac_address\":\"$macaddr\",\"ipv4\":["

        # Process the IPv4 addresses
        while read -r address; do
            a=""
            ad=""
            sub=""
            mask=""
            b=""

            if [[ $address =~ inet[[:space:]](([[:digit:]]{1,3}\.?){4}/[[:digit:]]{1,2}) ]]; then
                a=${BASH_REMATCH[1]}
            fi

            ad=$(echo $a | cut -d'/' -f1)
            sub=$(echo $a | cut -d'/' -f2)
            mask=$(cdr2mask $sub)

            if [[ $address =~ brd[[:space:]](([[:digit:]]{1,3}\.?){4}) ]]; then
                b=${BASH_REMATCH[1]}
            fi

            json="$json{\"address\":\"$ad\",\"mask\":\"$mask\",\"broadcast\":\"$b\"},"
        done <<< "$v4addresses"
        json="${json%?}],\"ipv6\":["

        # Process the IPv6 addresses
        while read -r address6; do
            add=$(echo $address6 | cut -d' ' -f2)
            json="$json{\"address\":\"$add\",\"mask\":\"\",\"broadcast\":\"\"},"
        done <<< "$v6addresses"
        json="${json%?}]},"
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
