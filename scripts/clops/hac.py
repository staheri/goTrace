#!/usr/bin/env python
# Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2018, All rights reserved
# Code: genMatrix.py
# Description: Generate Jaccard similarity matrix from CL dot file (includes CL redundant removal and LCA) (input: dot)


import glob
import sys,subprocess
import os
import numpy as np

import matplotlib.pyplot as plt
from scipy.cluster.hierarchy import *


from tabulate import tabulate
import math
from sets import Set




def jac2pdist(m):
	ret = []
	for i in range(0,len(m)-1):
		for j in range(i+1,len(m)):
			#print "(%d,%d)=%.2f"%(i,j,m[i][j])
			ret.append(1-m[i][j])
	print ret
	return ret



if len(sys.argv) != 3:
	print "USAGE:\n\t " +sys.argv[0]+" matfile1 out"
	sys.exit(-1)

matFile1 = sys.argv[1]
out =  sys.argv[2]


try:
	data1 = open(matFile1,"r").read().split("\n")[:-1]
	dataSize1 = len(data1)
	mat1 = np.ones((len(data1),len(data1)))
	for i in range(0,len(data1)):
		item = [x for x in data1[i].split(",") if len(x) > 0]
		#print item
		for j in range(0,len(item)):
			if float(item[j]) != 1:
				mat1[i][j] = float(item[j])
except Exception, e:
	print "ERROR %s"%e
	mat1 = np.zeros((dataSize1,dataSize1))
print mat1


print ward(jac2pdist(mat1))
#dendrogram(ward(jac2pdist(mat1)), truncate_mode='level', p=4)
plt.title("HAC")
plt.xlabel('Goroutine IDs')
plt.ylabel('Distance')
#dendrogram(ward(jac2pdist(mat1)))
#dendrogram(ward(jac2pdist(mat1)), truncate_mode='level', p=4)
#plt.show()

print fcluster(ward(jac2pdist(mat1)),4,criterion='maxclust')
