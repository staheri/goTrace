#!/usr/bin/env python
# Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2020, All rights reserved
# Code: lattice.py
# Description:
#        - Generate Jaccard similarity matrix from CL dot file (includes CL redundant removal and LCA) (input: dot)
#        - Generate reduced_label concept lattice graph
#        - Hierarchicaly cluster data based on the jacmat
#        - Operations on lattice, jaccard matrix, fancy graphs, etc.

import sys,subprocess
import os
import numpy as np

from tabulate import tabulate
import math
from sets import Set
from collections import defaultdict



class LiteLat:

	# Constructor
	def __init__(self,vertices):

		# default dictionary to store graph
		self.graph = defaultdict(list)
		self.eulerTour = [] # From supremum
		self.level = [] # Levels of euler tour
		self.depths = defaultdict(list)
		self.firstOccur = []

		self.V = vertices #No. of vertices

		# data structures for LCA
		self.lca_spantree = defaultdict(list)
		self.lca_remaining = defaultdict(list)
		self.lca_revRemaining = defaultdict(list)
		self.lca_ancestors = defaultdict(list)

	# function to add an edge to graph
	def addEdge(self,u,v):
		if u not in self.graph.keys():
			self.firstOccur.append(-1)
		if v == -1:
			self.graph[u]=[]
		else:
			self.graph[u].append(v)

	# A function used by DFS
	def DFSUtil(self,v,visited):
		# Mark the current node as visited and print it
		visited[v]= True
		if self.firstOccur[v] == -1:
			self.firstOccur[v] = len(self.eulerTour)
		self.eulerTour.append(v)
		if len(self.level) == 1 and self.level[0] == -1:
			self.level[0] = 0
		else:
			self.level.append(self.level[-1]+1)

		# Recur for all the vertices adjacent to this vertex
		for i in self.graph[v]:
			if visited[i] == False:
				self.DFSUtil(i, visited)
				# This extra print is for Euler tour
				#print v+1,
				self.eulerTour.append(v)
				self.level.append(self.level[-1]-1)

	# The function to do DFS traversal. It uses recursive DFSUtil()
	def DFS(self,v):
		# Mark all the vertices as not visited
		visited = [False]*(len(self.graph))
		self.level.append(-1)
		# Call the recursive helper function to print
		# DFS traversal
		self.DFSUtil(v,visited)
		for i in range(0,len(self.eulerTour)):
			self.depths[self.eulerTour[i]] = self.level[i]

	# Using RMQ, it returns the LCA of nodes x and y in TREE
	def lca_t(self,x,y):
		return self.eulerTour[self.rmq(self.firstOccur[x],self.firstOccur[y])]

	# Using lca_t(), it returns the LCA of nodes x and y in GRAPH
	def lca_g(self,x,y):
		ll = []
		for xx in self.lca_ancestors[x]:
			for yy in self.lca_ancestors[y]:
				ll.append(self.lca_t(xx,yy))
		dd = -1
		lca = -1
		for item in ll:
			if self.depths[item] > dd:
				dd = self.depths[item]
				lca = item
		return lca

	def rmq(self,x,y):
		minx = sys.maxint
		idx = -1
		if x == y:
			#print "RMQ ERROR"
			return x
		elif x<y:
			for i in range(x,y+1):
				if self.level[i] < minx:
					minx = self.level[i]
					idx = i
			return idx
		else:
			for i in range(y,x+1):
				if self.level[i] < minx:
					minx = self.level[i]
					idx = i
			return idx

	# Partition the lattice to spanning tree and remaining for O(n3) LCA algorithm
	def LCA_1partition(self):
		stk = []
		seen = []
		for item in self.eulerTour:
			if len(stk) != 0:
				if item not in seen:
					if stk[-1] not in self.lca_spantree.keys():
						self.lca_spantree[stk[-1]] = [item]
					else:
						self.lca_spantree[stk[-1]].append(item)
					stk.append(item)
					seen.append(item)
				else:
					stk = stk[:-1]
			else:
				stk.append(item)
				seen.append(item)
		for key,val in self.graph.items():
			if key not in self.lca_spantree.keys():
				self.lca_remaining[key] = val
			else:
				self.lca_remaining[key] = []
				for item in val:
					if item not in self.lca_spantree[key]:
						self.lca_remaining[key].append(item)

		#print "G: %s"%(sorted(self.graph.items()))
		#print "T: %s"%(sorted(self.lca_spantree.items()))
		#print "D: %s"%(sorted(self.lca_remaining.items()))

	# Create list of ancestors of each node in remaining
	def LCA_2createLists(self):
		#reverse edges of self.remaining
		#recursively add edges to the list
		for parent,kids in self.lca_remaining.items():
			if len(kids) != 0:
				for kid in kids:
					if kid in self.lca_revRemaining.keys():
						self.lca_revRemaining[kid].append(parent)
					else:
						self.lca_revRemaining[kid] = [parent]
		for key in self.graph.keys():
			if key not in self.lca_revRemaining.keys():
				self.lca_revRemaining[key] = []
		for key,val in self.lca_revRemaining.items():
			self.lca_ancestors[key] = self.LCA_3listGen(key)

	# Generating lists of ancestors for LCA operations (copy of DFSUtil)
	def LCA_3listGenHelper(self,v,visited,alist):
		# Mark the current node as visited and print it
		visited[v]= True
		alist.append(v)
		# Recur for all the vertices adjacent to this vertex
		for i in self.lca_revRemaining[v]:
			if visited[i] == False:
				self.LCA_3listGenHelper(i, visited,alist)
		return alist

	# A copy of DFS for creating LCA lists
	def LCA_3listGen(self,v):
		# Mark all the vertices as not visited
		visited = [False]*(len(self.lca_revRemaining))
		alist = []
		# Call the recursive helper function to print
		# DFS traversal
		return self.LCA_3listGenHelper(v,visited,alist)

class Lattice:
	def __init__(self,name):
		self.name = name
		self.attributes = Set([])
		self.objects = Set([])
		self.nodes = {}
		self.rnodes = {}
		self.ncnt = 0
		self.edges = {}
		self.ecnt = 0
		self.atr2node = defaultdict(list)
		self.obj2node = defaultdict(list)
		self.supID = -1
		self.infID = -1

    def addNode(self,id,content):
		#print "Graph::AddNode(%s)..."%(id)
		if id not in self.nodes.keys():
			self.nodes[id] = {}
			self.rnodes[hash(content.partition(":")[2].partition("\"")[0].strip())] = id
			self.nodes[id]["objs"] = [x for x in content.partition(">")[0].partition("<")[2].split(",") if len(x) > 0]
			self.objects |= Set(self.nodes[id]["objs"])
			self.nodes[id]["atrs"] = [x for x in content.partition(")")[0].partition("(")[2].split(",") if len(x) > 0]
			self.attributes |= Set(self.nodes[id]["atrs"])
			self.nodes[id]["redObjs"] = []
			self.nodes[id]["redAtrs"] = []
			self.nodes[id]["label"] = "n/a"
			self.nodes[id]["childs"] = []
			self.nodes[id]["parents"] = []
			self.ncnt = self.ncnt + 1
		#else:
		#	print "Node(%s) already exist"%(id)

	def addEdge(self,sid,did):
		self.edges[self.ecnt] = (sid,did)
		self.ecnt = self.ecnt + 1
		self.nodes[sid]["childs"].append(did)
		self.nodes[did]["parents"].append(sid)

    def supinfDetection(self):
		for key,val in sorted(self.nodes.items()):
			if len(val["parents"]) == 0:
				self.supID = key
			if len(val["childs"]) == 0:
				self.infID = key
		if len(self.nodes) < 2:
			self.supID = 0
			self.infID = 0

    def assignLabel(self):
		#for key,val in self.rnodes.items():
		#	print key
		#	print val
		for atr in self.attributes:
			#find concept
			atr_p = firstClosure(int(atr),"a")
			atr_pp = secondClosure(atr_p,"o")
			cntt = "<"
			for item in atr_p:
				cntt = cntt + `item` + ","
			cntt = cntt + ">,("
			for item in atr_pp:
				cntt = cntt + `item` + ","
			cntt = cntt + ")"
			#print "Atr: %s -> %s"%(atr,cntt)
			idd = self.rnodes[hash(cntt)]
			self.nodes[idd]["redAtrs"].append(atr)
			self.atr2node[atr] = idd
		for obj in self.objects:
			#find concept
			obj_p = firstClosure(int(obj),"o")
			obj_pp = secondClosure(obj_p,"a")
			cntt = "<"
			for item in obj_pp:
				cntt = cntt + `item` + ","
			cntt = cntt + ">,("
			for item in obj_p:
				cntt = cntt + `item` + ","
			cntt = cntt + ")"
			#print "Obj: %s -> %s"%(obj,cntt)
			idd = self.rnodes[hash(cntt)]
			self.nodes[idd]["redObjs"].append(obj)
			self.obj2node[obj]=idd

	def toString(self):
		global attrTable
		print "Nodes:"
		hdrs = ["ID","Objects","Attributes","In Deg","Out Deg"]
		tab = []
		for key,val in sorted(self.nodes.items()):
			tmp = []
			tmp.append(key)
			tmp.append(val["objs"])
			ttt = [attrTable[x] for x in val["atrs"]]
			tmp.append(ttt)
			tmp.append(len(val["parents"]))
			tmp.append(len(val["childs"]))
			tab.append(tmp)
		print tabulate(tab,headers=hdrs,tablefmt="fancy_grid")
		print "Edges:"
		hdrs = ["ID","From","To"]
		tab = []
		for key,val in sorted(self.edges.items()):
			tmp = []
			tmp.append(key)
			tmp.append(val[0])
			tmp.append(val[1])
			tab.append(tmp)
		print tabulate(tab,headers=hdrs,tablefmt="fancy_grid")

	def toReducedString(self):
		global attrTable
		print "Reduced Nodes:"
		hdrs = ["ID","Objects","Attributes"]
		tab = []
		for key,val in sorted(self.nodes.items()):
			tmp = []
			tmp.append(key)
			tmp.append(val["redObjs"])
			ttt = [attrTable[x] for x in val["redAtrs"]]
			tmp.append(ttt)
			tab.append(tmp)
		print tabulate(tab,headers=hdrs,tablefmt="fancy_grid")
		print "Object and attribute nodes:"
		hdrs = ["Objects","ID"]
		tab = []
		for key,val in sorted(self.obj2node.items()):
			tmp = []
			tmp.append(key)
			tmp.append(val)
			tab.append(tmp)
		print tabulate(tab,headers=hdrs,tablefmt="fancy_grid")
		hdrs = ["Attributes","ID"]
		tab = []
		for key,val in sorted(self.atr2node.items()):
			tmp = []
			tmp.append(key)
			tmp.append(val)
			tab.append(tmp)
		print tabulate(tab,headers=hdrs,tablefmt="fancy_grid")

def readTable(list):
	ret = {}
	for line in list:
		ll = line.split("|")
		if len(ll) == 3:
			value = ll[1].strip()
			ret[ll[0].strip()] = value
	return ret

def firstClosure(id,type):
	global cmat
	if type == "a":
		m = id-1
		mp = []
		for i in range(0,len(cmat)):
			if cmat[i][m] == "1":
				mp.append(i+1)
		ret = mp
	else:
		g = id-1
		gp = []
		for i in range(0,len(cmat[g])):
			if cmat[g][i] == "1":
				gp.append(i+1)
		ret = gp
	return ret

def secondClosure(flist,ftype):
	global cmat
	if ftype == "a":
		mp = flist
		mpp = []
		for i in range (0,len(cmat)):
			flg = True
			for item in mp:
				if cmat[i][item-1] != "1":
					flg = False
			if flg:
				mpp.append(i+1)
		ret = mpp
	else:
		gp = flist
		gpp = []
		for j in range(0,len(cmat[0])):
			flg = True
			for i in gp:
				if cmat[i-1][j] != "1":
					flg = False
			if flg:
				gpp.append(j+1)
		ret = gpp
	return ret

def attrSummary(l,type):
	global attrTable
	s = ""
	if len(l) == 0:
		s = "[-]"
	else:
		if type == "StackMRR":
			for item in l:
				s = s + attrTable[item] + " \\l "
		else:
			for item in l:
				#print item
				#print attrTable[item]
				if ":" in attrTable[item] and attrTable[item].partition(":")[2] == "1":
					s = s + attrTable[item].partition(":")[0] + " \\l "
				else:
					s = s + attrTable[item] + " \\l "
	ss = ""
	for ch in s:
		if ch != ">":
			ss = ss + ch
	return ss

def setSummary(l):
	iprev= -1
	istart= 0
	iend = 0
	sistart = ""
	siend = ""
	tmps = ""
	if len(l) == 0:
		tmps = "[-]"
	else:
		for item in l:
			ito = int(item)
			if iprev == -1:
				istart = ito
				iend = ito
			else:
				if iprev == ito - 1:
					iend = ito
				else:
					sistart = `istart`
					siend = `iend`
					#wrap previouses
					if istart == iend:
						tmps = tmps + sistart + ","
					else:
						tmps = tmps + sistart + "-" + siend + ","
					istart = ito
					iend = ito
					#set new istart
			iprev = ito
		sistart = `istart`
		siend = `iend`
		if istart == iend:
			tmps = tmps + sistart
		else:
			tmps = tmps + sistart + "-" + siend
	return tmps

def fancyDot(lat,showSupAtr,showFull):
	s = "digraph { \n\tnode[shape=record,style=filled,fillcolor=gray95]\n\n"
	#for key,val in lat.edges.items():
	#	print key
	#	print val
	for key,val in lat.nodes.items():
		#print key
		if showFull:
			newObjs = sorted([int(x)-1 for x in val["objs"]])
		else:
			newObjs = sorted([int(x)-1 for x in val["redObjs"]])
		objects = "Rank(s) " + setSummary(newObjs)
		if showSupAtr == 0 and key == lat.supID:
			attributes = "X"
		else:
			if showFull:
				attributes = attrSummary(val["atrs"],"x")
			else:
				attributes = attrSummary(val["redAtrs"],"x")
		s = s + "\t" + `key` + " [label = \"{" + objects + " | " + attributes + "}\"]\n"
	s = s + "\n"
	# ADD edges
	l = [lat.supID]
	edges = []
	visited = []
	while len(l) > 0:
		n = l[0]
		visited.append(n)
		l = l[1:]
		for child in lat.nodes[n]["childs"]:
			if child not in visited:
				l.append(child)
		for child in lat.nodes[n]["childs"]:
			edge = `n` + " -> " + `child`
			if edge not in edges:
				s = s + "\t" + edge + "\n"
				edges.append(edge)
	s = s + "}"
	print s
	return s

def latmatToFullMat(lm):
	mat = []
	irow = [0]*len(lm)
	for i in range(0,len(lm)):
		irow = [0]*len(lm)
		irow[i] = 1
		if latmatc[i] != "-":
			#print latmatc[i].split(" ")
			for u in latmatc[i].split(" "):
				if len(u) > 0:
					irow[int(u)] = 1

		mat.append(irow)
	return mat

def transitiveClosure(mat):
	reach = [i[:] for i in mat]
	for k in range(len(mat)):
		for i in range(len(mat)):
			for j in range(len(mat)):
				reach[i][j] = reach[i][j] or (reach[i][k] and reach[k][j])
	return reach

def unSum(lat,i,j):
	# reverse parse in the lattice to sum the number of atributes
	toCheck = [i,j]
	seen = []
	sum = 0
	while len(toCheck) != 0:
		#print "To Check %s"%(toCheck)
		tt = toCheck[0]
		toCheck = toCheck[1:]
		if tt not in seen:
			seen.append(tt)
			sum = sum + len(lat.nodes[tt]["redAtrs"])
			for item in lat.nodes[tt]["parents"]:
				toCheck.append(item)
	#print sum
	return sum

def simmax(lat,ll,path):
	objects = sorted([int(x) for x in lat.objects])
	mtx = np.ones((len(objects),len(objects)))
	#print lat.obj2node
	#print objects
	for i in range(0,len(objects)):
		for j in range(i+1,len(objects)):
			iid = lat.obj2node[str(objects[i])]
			jid = lat.obj2node[str(objects[j])]
			if iid != jid: # i and j are in the same node in lattce, sim = 1
				lca = ll.lca_g(iid,jid)
				inter = len(lat.nodes[lca]["atrs"])
				union = unSum(lat,iid,jid)
				mtx[i][j] = inter * 1.0 / union
				mtx[j][i] = inter * 1.0 / union

	#print "%s"%mtx
	s = ""
	for i in range(0,len(mtx)):
		for j in range (0,len(mtx[i])):
			s = s + "%.3f,"%mtx[i][j]
		s = s + "\n"
	print "Writing to...\n"+path.rpartition(".")[0]+".jacmat.txt"
	fout = open(path.rpartition(".")[0]+".jacmat.txt","w")
	fout.write(s)
	fout.close()

def simmax2(lobj,path):
	mtx = np.zeros((lobj,lobj))
	#print "%s"%mtx
	s = ""
	for i in range(0,len(mtx)):
		for j in range (0,len(mtx[i])):
			s = s + "%.3f,"%mtx[i][j]
		s = s + "\n"
	print "Writing to...\n"+path.rpartition(".")[0]+".jacmat.txt"
	fout = open(path.rpartition(".")[0]+".jacmat.txt","w")
	fout.write(s)
	fout.close()
