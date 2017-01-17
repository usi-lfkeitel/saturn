#!/bin/bash
result=$(/usr/bin/lscpu \
    | /usr/bin/awk -F: '{print "\""$1"\":\""$2"\","}' \
    | sed 's/\"\s\+/\"/g')

echo -n "{" ${result%?} "}"