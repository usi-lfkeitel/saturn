#!/bin/bash
#gen:module a type:string,user:string,home:string

result=$(/usr/bin/awk -F: '{ \
        if ($3<=499) {userType="system";} \
        else {userType="user";} \
        print "{ \"type\": \"" userType "\"" ", \"user\": \"" $1 "\", \"home\": \"" $6 "\" }," }' < /etc/passwd
    )

if [ ${#result} -eq 0 ]; then
  result=$(getent passwd | /usr/bin/awk -F: '{ \
      if ($3<=499) {userType="system";} \
      else {userType="user";} \
      print "{ \"type\": \"" userType "\"" ", \"user\": \"" $1 "\", \"home\": \"" $6 "\" }," }'
    )
fi

echo -n [ ${result%?} ]
