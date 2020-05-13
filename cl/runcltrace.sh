#!/bin/bash

#print "USAGE:\n\t " +sys.argv[0]+" config-file pathToJobSub tool experiment_name #ofRuns psize node flag runName"

for b in noBug allRed1wrgOp-1-all-x allRed1wrgSize-1-all-x allRed1wrgSize-all-all-x allRed2wrgOp-1-all-x allRed2wrgSize-1-all-x allRed2wrgSize-all-all-x bcastWrgSize-1-all-x bcastWrgSize-all-all-x misCrit-1-1-x misCrit-1-all-x misCrit-all-1-x misCrit-all-all-x misCrit2-1-1-x misCrit2-1-all-x misCrit2-all-1-x misCrit2-all-all-x misCrit3-all-all-x misCrit3-1-all-x infLoop-1-1-1
do 
	for i in m a
	do
		for p in 8 16
		do
			python genTCscript.py config.toml $SCRATCH/jobSub ilcsTSP $b $i $p auto
		done
    	done
done




#python newGenJobSub.py nascon2.toml $PROJECT/workspace pinMain pinMain 3 C 1 o1 run-00
