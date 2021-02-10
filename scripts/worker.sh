#!/bin/bash

for d in 0 1 2 3 4 5
do
  for app in mutex_circularWait_abba chanMutex_select chanMutex_circularWait chanMutex_circularWait2 chanMutexCond_circularWait.go
  do
    goat -app=$app.go -src=schedTest -cmd=test -iter=100 -depth=$d;
  done
done
