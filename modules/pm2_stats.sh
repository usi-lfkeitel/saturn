#!/bin/bash
pm_command=$(which pm2 2>/dev/null)

if [ -z $pm_command ]; then
    echo -n "[]"
    exit
fi

#get data
data="$($pm_command list)"

if [ -z "$data" ]; then
    echo -n "[]"
    exit
fi

#start processing data on line 4
#don't process last 2 lines
json=$( echo "$data" | tail -n +4 | head -n +2 \
| awk 	'{print "{"}\
    {print "\"appName\":\"" $2 "\","} \
    {print "\"id\":\"" $4 "\","} \
    {print "\"mode\":\"" $6 "\","} \
    {print "\"pid\":\"" $8 "\","}\
    {print "\"status\":\"" $10 "\","}\
    {print "\"restart\":\"" $12 "\","}\
    {print "\"uptime\":\"" $14 "\","}\
    {print "\"memory\":\"" $16 $17 "\","}\
    {print "\"watching\":\"" $19 "\""}\
    {print "},"}')
#make sure to remove last comma and print in array
echo -n "[" ${json%?} "]"
