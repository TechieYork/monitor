#!/bin/bash

app_name=dm_monitor_agent

echo "Starting $app_name ... "

echo "Current directory:"$(dirname $(readlink -f $0))

cd `dirname $0`

../bin/dm_monitor_agent &

echo "Already started $app_name"
