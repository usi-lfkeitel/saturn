#!/bin/bash
#gen:module o cores:int

echo -n '{"cores":' $(/bin/grep -c 'model name' /proc/cpuinfo)'}'
