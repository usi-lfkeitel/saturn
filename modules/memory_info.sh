#!/bin/bash
#gen:module o MemTotal:string,MemFree:string,MemAvailable:string,Buffers:string,Cached:string,SwapCached:string,Active:string,Inactive:string,Active(anon):string,Inactive(anon):string,Active(file):string,Inactive(file):string,Unevictable:string,Mlocked:string,SwapTotal:string,SwapFree:string,Dirty:string,Writeback:string,AnonPages:string,Mapped:string,Shmem:string,Slab:string,SReclaimable:string,SUnreclaim:string,KernelStack:string,PageTables:string,NFS_Unstable:string,Bounce:string,WritebackTmp:string,CommitLimit:string,Committed_AS:string,VmallocTotal:string,VmallocUsed:string,VmallocChunk:string,HardwareCorrupted:string,AnonHugePages:string,CmaTotal:string,CmaFree:string,HugePages_Total:string,HugePages_Free:string,HugePages_Rsvd:string,HugePages_Surp:string,Hugepagesize:string,DirectMap4k:string,DirectMap2M:string,DirectMap1G:string

/bin/cat /proc/meminfo \
  | /usr/bin/awk -F: 'BEGIN {print "{"} {print "\"" $1 "\": \"" $2 "\"," } END {print "}"}' \
  | /bin/sed 'N;$s/,\n/\n/;P;D' \
  | /bin/sed 's/\"\s\+/\"/g'
