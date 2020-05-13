#!/usr/bin/env python

# Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2017, All rights reserved
# Code: newGenSub.py
# Description: generates job submission scripts to run on Stampede for NAS applications with different flags {1/16,4/64,16/256,64/1024}

import toml
import argparse
import glob
import sys,subprocess



if len(sys.argv) != 2:
	print "USAGE:\n\t " +sys.argv[0]+" tracePath"
	sys.exit(-1)

path = sys.argv[1]
# Create directory to stor jobSubmission scripts

for f in glob.glob("/home/saeed/goTrace/traces/trace-patterns*/gtrace/cl/nlr10/sing.orig.w.jacmat.txt"):
	print f.partition("trace-patterns-")[2].partition("/gtrace")[0]
	ex = f.partition("trace-patterns-")[2].partition("/gtrace")[0]
	cmd = "python /home/saeed/diffTrace/scripts/genReport/genSingleJSM.py " + f + " " + ex + ";"
	print cmd
	process = subprocess.Popen([cmd], stdout=subprocess.PIPE,shell=True)
	si, err = process.communicate()

	#print "python /home/saeed/diffTrace/scripts/genReport/cl2jacmat.py "+f
	#process = subprocess.Popen(["python /home/saeed/diffTrace/scripts/genReport/cl2jacmat.py "+f], stdout=subprocess.PIPE,shell=True)
	#si, err = process.communicate()


	#print "/home/saeed/goTrace/cl/cltrace -m 6 -p "+f+" -a 1 -q 0 -n 0 -k 10"
	#process = subprocess.Popen(["/home/saeed/goTrace/cl/cltrace -m 6 -p "+f+" -a 1 -q 0 -n 0 -k 10"], stdout=subprocess.PIPE,shell=True)
	#si, err = process.communicate()
