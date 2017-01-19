#!/bin/bash
#gen:module a interface:string,tx:int,rx:int

/bin/cat /proc/net/dev \
| awk 'BEGIN {print "["} NR>2 {print "{ \"interface\": \"" $1 "\"," \
          " \"tx\": " $2 "," \
          " \"rx\": " $10 " }," } END {print "]"}' \
| /bin/sed 'N;$s/,\n/\n/;P;D' \
| tr -d '\n'
