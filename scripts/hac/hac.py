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
#import plotly.graph_objects as gobj

import readin





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
	dt = [0.1,0.2,0.4,0.5,0.8,1,2]
	mt = [2,3,4,5,6]
	#ret = fcluster(ward(jac2pdist(m)),maxc,criterion='maxclust')

	tab = []
	# d_criterion = []
	# d_thresh = []
	# d_single = []
	# d_ward = []
	# d_avg = []
	# d_comp= []
	hdrs = ["Criterion","T/Max","Ward","Single","Complete","Average"]

	table = "|Criterion|T/Max|Ward|Single|Complete|Average|\n"
	table = table + "|---:|---|---|---|---|---|\n"

	criterions = ["distance","inconsistent","maxclust"]
	for c in criterions:
		if c == "maxclust":
			t = mt[:]
		else:
			t = dt[:]
		for tt in t:
			ttab=[]
			# d_criterion.append(c)
			# d_thresh.append(tt)
			# d_ward.append(clusterText(fcluster(ward(jac2pdist(m)),tt,criterion=c)))
			# d_single.append(clusterText(fcluster(single(jac2pdist(m)),tt,criterion=c)))
			# d_comp.append(clusterText(fcluster(complete(jac2pdist(m)),tt,criterion=c)))
			# d_avg.append(clusterText(fcluster(average(jac2pdist(m)),tt,criterion=c)))
			table = table + "|" + c
			table = table + "|" + `tt`
			table = table + "|" + clusterText(fcluster(ward(jac2pdist(m)),tt,criterion=c))
			table = table + "|" + clusterText(fcluster(single(jac2pdist(m)),tt,criterion=c))
			table = table + "|" + clusterText(fcluster(complete(jac2pdist(m)),tt,criterion=c))
			table = table + "|" + clusterText(fcluster(average(jac2pdist(m)),tt,criterion=c))
			table = table + "|\n"

			ttab.append(c)
			ttab.append(tt)
			ttab.append(clusterText(fcluster(ward(jac2pdist(m)),tt,criterion=c)))
			ttab.append(clusterText(fcluster(single(jac2pdist(m)),tt,criterion=c)))
			ttab.append(clusterText(fcluster(complete(jac2pdist(m)),tt,criterion=c)))
			ttab.append(clusterText(fcluster(average(jac2pdist(m)),tt,criterion=c)))
			tab.append(ttab)
	f = open(out+"_hac.md","w")
	f.write(table)
	f.close()
	print tabulate(tab,headers=hdrs,tablefmt="plain")


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
			s = s + readin.objTable[`item+1`].rpartition(".")[0]+", "

			#s = s + "g"+`item`+", "
		ttab.append(s)
		tab.append(ttab)
	#print tabulate(tab,headers=hdrs,tablefmt="fancy_grid")
	return tabulate(tab,headers=hdrs)


def clusterText(c):
	data = {}
	for i in range(0,len(c)):
		if c[i] in data.keys():
			data[c[i]].append(i)
		else:
			data[c[i]]=[i]

	st = ""
	tab = []
	i = 0
	for k,v in data.items():
		st = st +  "["+`k` + "]: "
		s = ""
		for item in v:
			s = s + readin.objTable[`item+1`].rpartition(".")[0]+", "

		if i < len(data.keys()) - 1 :
			st = st + s + "<br>"
		else:
			st = st + s
		i = i + 1
	return st
