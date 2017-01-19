#!/bin/bash
#gen:module o bytes:int,bytes_read:int,bytes_written:int

echo "stats" \
  | /bin/nc -w 1 127.0.0.1 11211 \
  | /bin/grep 'bytes' \
  | /usr/bin/awk 'BEGIN {print "{"} {print "\"" $2 "\": " $3 } END {print "}"}' \
  | /usr/bin/tr '\r' ',' \
  | /bin/sed 'N;$s/,\n/\n/;P;D'
