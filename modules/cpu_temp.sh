#!/bin/bash
if [ `which sensors` ]; then
  returnString=`sensors`
  #amd
  if [[ "${returnString/"k10"}" != "${returnString}" ]] ; then
    echo -n ${returnString##*k10} | cut -d ' ' -f 6 | cut -c 2- | cut -c 1-4
  #intel
  elif [[ "${returnString/"core"}" != "${returnString}" ]] ; then
    fromcore=${returnString##*"coretemp"}
    echo -n ${fromcore##*Physical}  | cut -d ' ' -f 3 |  cut -c 2-5
  fi
else
  echo -n "0.0"
fi
