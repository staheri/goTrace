#!/bin/bash

for d in 0 1 2 3 4 5
do
  for app in chanMutexCond_circularWait
  do
    goat -app=$app.go -src=schedTest -cmd=test -iter=100 -depth=$d;
  done
done
