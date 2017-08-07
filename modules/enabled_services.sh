#!/bin/bash
#gen:module a name:string

sysctl="$(which systemctl 2>/dev/null)"
if [ -z "$sysctl" ]; then
    echo "systemctl not available"
    exit 1
fi

systemctl list-unit-files --no-page |\
    grep '.service' |\
    grep 'enabled' |\
    cut -d' ' -f1 |\
    /usr/bin/awk -F: 'BEGIN {print "["} {print "{\"name\": \"" $1 "\"}," } END {print "]"}' |\
    /bin/sed 'N;$s/,\n/\n/;P;D'