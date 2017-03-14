#!/bin/bash

app_name=dm_monitor_server

echo "Stopping $app_name ... "

keys=`(ps -ef |grep "$app_name" |grep -v "grep") | awk '{print $2}'`

for key in ${keys[*]}
do
    echo "Killing pid -> "$key
    kill -9 $key
done

echo "Already stopped $app_name"
