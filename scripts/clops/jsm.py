#!/usr/bin/env python
# Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2018, All rights reserved
# Code: genMatrix.py
# Description: Generate Jaccard similarity matrix from CL dot file (includes CL redundant removal and LCA) (input: dot)

import matplotlib
matplotlib.use('Agg')
#matplotlib.use('tkagg')

import glob
import sys,subprocess
import os
import numpy as np
import seaborn as sns ; sns.set(font_scale=0.9)
#import seaborn as sns
import matplotlib.pyplot as plt



from tabulate import tabulate
import math
from sets import Set
from collections import defaultdict




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



# corr = np.corrcoef(mat)
# mask = np.zeros_like(corr)
# mask[np.triu_indices_from(mask)] = True
# with sns.axes_style("white"):
# 	ax = sns.heatmap(corr, mask=mask,vmin=0, vmax=1,annot=True, fmt="f")

#ax = sns.heatmap(mat1, annot=True,vmin=0.64, vmax=1)
ax = sns.heatmap(mat1, annot=True)
fig = ax.get_figure()
#fig.tight_layout()
fig.savefig(out+"JSM.pdf")
