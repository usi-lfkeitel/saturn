#!/bin/bash
#gen:module a device:string,reads:int,writes:int,in_prog:int,time:int

result=$(/bin/cat /proc/diskstats | /usr/bin/awk \
        '{ if($4==0 && $8==0 && $12==0 && $13==0) next } \
        {print "{\"device\": \"" $3 "\", \"reads\": "$4", \"writes\": "$8 ", \"in_prog\": "$12 ", \"time\": "$13"},"}'
    )

echo -n [${result%?}]
