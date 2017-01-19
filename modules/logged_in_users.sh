#!/bin/bash
#gen:module a user:string,from:string,when:string

result=$(COLUMNS=300 /usr/bin/w -h | /usr/bin/awk '{print "{\"user\":\""$1"\",\"from\":\""$3"\",\"when\":\""$4"\"},"}')

echo -n [${result%?}]
