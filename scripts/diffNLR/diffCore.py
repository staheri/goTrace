#!/usr/bin/env python
# Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2018, All rights reserved
# Code: diffCore.py
# Description: Takes 2 MRR text files return their diff

import sys,subprocess
import os
import math
import decimal


# Summarize sequence of integers
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


# Generate an edit graph with respect to length of seq1 and seq2 for diff operations
def genEditGraph(seq1,seq2):
	eg = []
	for i in range(0,len(seq2)):
		eg_row = [0]
		for j in range(0,len(seq1)):
			if seq1[j] == seq2[i]:
				eg_row.append(1)
			else:
				eg_row.append(0)
		eg.append(eg_row)
	eg.insert(0,[0 for i in range(0,len(seq1)+1)])
	return eg

# Get element i and j of EG
def getEGelement(eg,i,j):
	return eg[j][i]

# Print EG
def printEG(eg):
	for row in eg:
		print row
		
# Find middle snake of A and B. n=len(a), m=len(b)
def midSnake(a,n,b,m,eg):
	ret = {}
	vf = {}
	vr = {}
	vfd = {}
	vrd = {}
	vf[1] = 0
	delta = n - m
	vr[n-m-1] = n
	
	for d in range(0,int(math.ceil((m+n)*1.0/2))+1):
		vfd[d] = {}
		vrd[d] = {}
		for k in range(-d,d+1,2):
			if k == -d or k != d and vf[k-1] < vf[k+1]: 
				x = vf[k+1]
			else: 
				x = vf[k-1]+1
			y = x - k
			xp = x
			while x < n and y < m and getEGelement(eg,x+1,y+1) == 1:
				x = x+1
				y = y+1
			vf[k] = x
			vfd[d][k] = (x,y)
			if delta%2==1 and delta - d + 1 <= k <= delta + d - 1  :
				if overlaps(x,y,vrd[d-1][k][0],vrd[d-1][k][1],"f"):
					ret = [2*d-1,xp,xp-k,x,x-k]
					return ret
		for k in range(d,-d-1,-2):
			kd = k + delta
			if k == d:
				x = vr[kd-1]
				# add insertion to path
			elif k == -d:
				x = vr[kd+1]-1
				# add deletion to path
			else:
				if k!=d and vr[kd-1] < vr[kd+1] :
					x = vr[kd-1]
				else:
					x = vr[kd+1] - 1
					# add insertion to path
			y = x - kd
			xp = x
			while x >= 1 and y >= 1 and  getEGelement(eg,x,y) == 1:
				x = x-1
				y = y-1
			vr[kd] = x
			vrd[d][kd] = (x,y)
			if delta % 2 == 0 and (-1) * d <= kd and kd <= d :
				if overlaps(x,y,vfd[d][kd][0],vfd[d][kd][1],"r"):
					ret = [2*d,x,x-kd,xp,xp-kd]
					return ret
					
					
	print "NO LS"
	return [-1,-1,-1,-1,-1]
	#return ret

# Check if two middle snakes overlaps (t = f(orward) ? r(reverse))
def overlaps(x,y,u,v,t):
	if t == "f":
		return x-y == u-v and x >= u 
	else:
		return x-y == u-v and x <= u


# Returns appropriate edit command
def func1(a,b):
	s = ""
	assert(len(a)-len(b)==1)
	dif = -1
	for i in range(0,len(a)-1):
		if a[i] != b[i]:
			dif = i
			break
		dif = i + 1
		#print i
	#print dif
	if dif == len(a)-1:
		s = s + "C: %s\n"%(b)
		s = s + "A: %s\n"%(a[dif:])
	else:
		s = s + "C: %s\n"%(b[:dif])
		s = s + "A: %s\n"%(a[dif:dif+1])
		s = s + "C: %s\n"%(b[dif:])
	return s

# Returns appropriate edit command
def func2(a,b):
	s = ""
	assert(len(b)-len(a)==1)
	dif = -1
	for i in range(0,len(b)-1):
		if a[i] != b[i]:
			dif = i
			break
		dif = i + 1
	if dif == len(b)-1:
		s = s + "C: %s\n"%(a)
		s = s + "B: %s\n"%(b[dif:])
	else:
		s = s + "C: %s\n"%(a[:dif])
		s = s + "B: %s\n"%(b[dif:dif+1])
		s = s + "C: %s\n"%(a[dif:])
	return s
	
	
	if len(dif) == 1:
		for i in range(0,len(b)):
			if b[i] == dif[0]:
				#print "INS %s"%([b[i]])
				s = s + "B: %s\n"%([b[i]])
			else:
				#print "*** %s"%([b[i]])
				s = s + "C: %s\n"%([b[i]])
	else:
		assert(len(dif)==0)
		assert(len(b)>len(set(b)))
		#print "*** %s"%([b[0]])
		#print "INS %s"%([b[1]])
		s = s + "C: %s\n"%([b[0]])
		s = s + "B: %s\n"%([b[1]])
		#print "xxx"
	return s

# Recursive function to find least common edit script of A and B

def lcs(a,b):
	n = len(a)
	m = len(b)
	#print "LCS:\nA:"
	eg = genEditGraph(a,b)
	if n > 0 and m > 0:
		r = midSnake(a,n,b,m,eg)
		d = r[0]
		x = r[1]
		y = r[2]
		u = r[3]
		v = r[4]
		#print "D:%d MID-Snake from (%d,%d) to (%d,%d)"%(d,x,y,u,v)
		if d > 1:
			#print "MID-Snake from (%d,%d) to (%d,%d)"%(x,y,u,v)
			#print "\n\tLL LCS %s , %s"%(a[:x],b[:y])
			s = ""
			s = s + lcs(a[:x],b[:y])
			if len(a[x:u]) != 0:
				s = s + "C: %s\n"%(a[x:u])
				#print "*** %s"%(a[x:u])
			#print "\n\tRR LCS %s , %s"%(a[u:],b[v:])
			s = s + lcs(a[u:],b[v:])
			return s
		elif m > n:
		#	print "INSERT %s"%(b)
			#print "*** Insert something to %s"%(a)
			return func2(a,b)
			#print a
		elif m < n:
			#print "DELETE %s"%(a)
			#print "*** Delete something from %s"%(a)
			return func1(a,b)
			#print b
		else:
			#print a
			return "C: %s\n"%(a)

	elif m > n:
		#print "INS %s"%(b)
		return "B: %s\n"%(b)
		#print a
	else:
		#print "DEL %s\n"%(a)
		return "A: %s\n"%(a)
		#print b

