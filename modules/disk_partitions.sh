#!/bin/bash
#gen:module a file_system:string,size:string,used:string,avail:string,used_percent:int:%,mounted:string

result=$(/bin/df -Ph | awk 'NR>1 {print "{\"file_system\":\""$1"\",\"size\":\""$2"\",\"used\":\""$3"\",\"avail\":\""$4"\",\"used_percent\":"$5",\"mounted\":\""$6"\"},"}' | tr -d '%')

echo -n [${result%?}]
