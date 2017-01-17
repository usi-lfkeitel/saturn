#!/bin/bash
grepCmd=`which grep`
ifconfigCmd=`which ifconfig`

json='['

for item in $($ifconfigCmd | $grepCmd -oP "^[a-zA-Z0-9:\-]*"); do
    address=`$ifconfigCmd $item | $grepCmd -Po 't addr:\K[\d.]+'`
    json="$json{\"interface\":\""$item"\",\"ip\":\"$address\"},"
done

echo -n ${json%?}"]"
