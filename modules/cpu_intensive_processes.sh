#!/bin/bash
#gen:module a pid:int,user:string,cpu_percent:float64,rss:int,vsz:int,cmd:string

result=$(/bin/ps axo pid,user,pcpu,rss,vsz,comm --sort -pcpu,-rss,-vsz \
      | head -n 15 \
      | /usr/bin/awk 'BEGIN{OFS=":"} NR>1 {print "{ \"pid\": " $1 \
              ", \"user\": \"" $2 "\"" \
              ", \"cpu_percent\": " $3 \
              ", \"rss\": " $4 \
              ", \"vsz\": " $5 \
              ", \"cmd\": \"" $6 "\"" "},"\
            }')

echo -n "["${result%?}"]"
