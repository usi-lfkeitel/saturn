#!/bin/bash
#gen:module a interface:string,ip:string

# Some newer systems have ifconfig here and for some reason which sometimes doesn't find it
PATH="$PATH:/usr/sbin:/sbin"
grepCmd=`which grep`
ifconfigCmd=`which ifconfig`

json='['

for item in $($ifconfigCmd | $grepCmd -oP "^[a-zA-Z0-9:\-]*"); do
    address=""
    if [ "${item: -1}" = ":" ]; then
        address=`$ifconfigCmd ${item: -1} | $grepCmd -oP 'inet \K[\d\.]+'`
    else
        address=`$ifconfigCmd $item | $grepCmd -oP 't addr:\K[\d.]+'`
    fi
    json="$json{\"interface\":\"$item\",\"ip\":\"$address\"},"
done

echo -n ${json%?}"]"
