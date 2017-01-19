#!/bin/bash
#gen:module a addr:string,hw_type:string,hw_addr:string,mask:string

# TODO: Fix parsed output. Mask doesn't show most of the time so the interface
# name is being used as the mask.
arpCommand=$(which arp 2>/dev/null)

result=$($arpCommand | awk 'BEGIN {print "["} NR>1 \
            {print "{ \"addr\": \"" $1 "\", " \
                  "\"hw_type\": \"" $2 "\", " \
                  "\"hw_addr.\": \"" $3 "\", " \
                  "\"mask\": \"" $5 "\" }, " \
                  } \
          END {print "]"}' \
      | /bin/sed 'N;$s/},/}/;P;D')

if [ -z "$result" ];  then
      echo -n {}
else
      echo -n $result
fi
