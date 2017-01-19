#!/bin/bash
#gen:module a cname:string,pid:int,user:string,cpu_percent:float64,mem_percent:float64,cmd:string

docker=/usr/bin/docker

if [ ! -x $docker ]; then
  echo "[]"
  exit
fi

result=""
containers="$($docker ps | awk '{if(NR>1) print $NF}')"
for i in $containers; do
  result="$result $($docker top $i axo pid,user,pcpu,pmem,comm --sort -pcpu,-pmem \
      | head -n 15 \
      | /usr/bin/awk -v cnt="$i" 'BEGIN{OFS=":"} NR>1 {print "{\"cname\": \""cnt \
              "\",\"pid\":"$1 \
              ",\"user\":\""$2"\"" \
              ",\"cpu_percent\":"$3 \
              ",\"mem_percent\":"$4 \
              ",\"cmd\":\""$5"\"""},"\
            }')"
done

echo -n "[${result%?}]"
