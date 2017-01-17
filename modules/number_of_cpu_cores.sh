#!/bin/bash
echo -n $(/bin/grep -c 'model name' /proc/cpuinfo)
