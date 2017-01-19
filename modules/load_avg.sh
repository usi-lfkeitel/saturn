#!/bin/bash
#gen:module o avg_1_min:int,avg_5_min:int,avg_15_min:int

grepCmd=`which grep`
awkCmd=`which awk`
catCmd=`which cat`

numberOfCores=$($grepCmd -c 'processor' /proc/cpuinfo)

if [ $numberOfCores -eq 0 ]; then
  numberOfCores=1
fi

result=$($catCmd /proc/loadavg | $awkCmd '{print "{\"avg_1_min\":"($1*100)/'$numberOfCores'",\"avg_5_min\":"($2*100)/'$numberOfCores'",\"avg_15_min\":"($3*100)/'$numberOfCores' "},"}')

echo -n ${result%?}
