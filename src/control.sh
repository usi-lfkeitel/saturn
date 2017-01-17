#!/bin/bash

SCRIPT="$1"
shift

echo -n '{'

for var in "$@"; do
    echo -n "\"$var\": "
    echo $($SCRIPT $var)
    echo ', '
done

echo -n '}'