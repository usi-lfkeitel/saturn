#!/bin/bash
#gen:module a process:string,pid:string,type:string,proto:string,address:string,port:string

result=$(lsof -i -P | grep LISTEN | /usr/bin/awk '
    {
        split($9, addrport, ":");
        print "{\"process\":\""$1"\",\"pid\":\""$2"\",\"type\":\""$5"\",\"proto\":\""$8"\",\"address\":\""addrport[1]"\",\"port\":\""addrport[2]"\"},"
    }
')

echo -n [${result%?}]
