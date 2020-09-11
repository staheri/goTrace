#!/bin/bash

APP=../CodeBenchmark/misc/powser.go
CMD=hac
OUT=../results
SRC=latest
CAT=CHNL\ GRTN\ GRTN,CHNL\ WGRP\ CHNL,WGRP\ GRTN,CHNL,WGRP\ GRTN,WGRP\ MUTX\ GCMM

# for S in 4 8 16 32 64 128 256 512 1024 2048
# do
#   echo $SIZE_FFT;
#   export SIZE_FFT=$S
#   ./src -app=$APP -cmd=$CMD -outdir=$OUT -src=$SRC CHNL;
# done

for C in 2 3 4 5 6 7 8
do
 for A in 2 3 4
 do
   ./src -app=$APP -cmd=$CMD -outdir=$OUT -src=$SRC -cons=$C -atrmode=$A $CAT ;
 done
done
