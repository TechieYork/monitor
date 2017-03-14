#!/bin/bash

app_name=dm_monitor_server

echo "Starting $app_name ... "

echo "Current directory:"$(dirname $(readlink -f $0))

cd `dirname $0`

../bin/dm_monitor_server &

echo "Already started $app_name"
