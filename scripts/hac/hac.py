#!/usr/bin/env python
# Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2020, All rights reserved
# Code: hac.py
# Description:
#        - Generate Jaccard similarity matrix from CL dot file (includes CL redundant removal and LCA) (input: dot)
#        - Generate reduced_label concept lattice graph
#        - Hierarchicaly cluster data based on the jacmat

import numpy as np
import matplotlib.pyplot as plt
from scipy.cluster.hierarchy import *
from tabulate import tabulate






def jac2pdist(m):
	ret = []
	for i in range(0,len(m)-1):
		for j in range(i+1,len(m)):
			#print "(%d,%d)=%.2f"%(i,j,m[i][j])
			ret.append(1-m[i][j])
	#print ret
	return ret


def cluster(m,maxc,out):
	#print ward(jac2pdist(m))
	#dendrogram(ward(jac2pdist(mat1)), truncate_mode='level', p=4)
	plt.title("HAC")
	plt.xlabel('Goroutine IDs')
	plt.ylabel('Distance')
	dendrogram(ward(jac2pdist(m)))
	#dendrogram(ward(jac2pdist(mat1)), truncate_mode='level', p=4)
	#plt.show()
	plt.savefig(out+"-dend.pdf")
	ret = fcluster(ward(jac2pdist(m)),maxc,criterion='maxclust')
	f=open(out+"-C"+`maxc`+"-rep.txt","w")
	f.write(clusterTable(ret))
	f.close()
	return ret

def clusterTable(c):
	data = {}
	for i in range(0,len(c)):
		if c[i] in data.keys():
			data[c[i]].append(i)
		else:
			data[c[i]]=[i]
	hdrs=["Cluster","Goroutines"]
	tab = []
	for k,v in data.items():
		ttab = []
		ttab.append(k)
		s = ""
		for item in v:
			s = s + "g"+`item`+", "
		ttab.append(s)
		tab.append(ttab)
	#print tabulate(tab,headers=hdrs,tablefmt="fancy_grid")
	return tabulate(tab,headers=hdrs)
