#!/bin/bash
#gen:module o Architecture:string,CPU op-mode(s):string,Byte Order:string,CPU(s):string,On-line CPU(s) list:string,Thread(s) per core:string,Core(s) per socket:string,Socket(s):string,NUMA node(s):string,Vendor ID:string,CPU family:string,Model:string,Model name:string,Stepping:string,CPU MHz:string,CPU max MHz:string,CPU min MHz:string,BogoMIPS:string,Virtualization:string,L1d cache:string,L1i cache:string,L2 cache:string,L3 cache:string,NUMA node0 CPU(s):string,Flags:string

result=$(/usr/bin/lscpu \
    | /usr/bin/awk -F: '{print "\""$1"\":\""$2"\","}' \
    | sed 's/\"\s\+/\"/g')

echo -n "{" ${result%?} "}"