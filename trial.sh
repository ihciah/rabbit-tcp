#!/bin/bash
# trial.sh BIN -> fresh server+client, storm, report
pkill -f "rabbit-(orig|fix) -mode"; sleep 0.5
./$1 -mode s -password reprotest -rabbit-addr 127.0.0.1:19891 -verbose 0 > /dev/null 2>&1 &
./$1 -mode c -password reprotest -rabbit-addr 127.0.0.1:19891 -listen 127.0.0.1:19892 -dest 127.0.0.1:18000 -tunnelN 4 -verbose 0 > /dev/null 2>&1 &
sleep 1.5; ./storm.sh 19892 5 | awk '{print $3}' | tr '\n' ' '
