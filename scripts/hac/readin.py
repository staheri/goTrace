#!/usr/bin/env python
# Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2020, All rights reserved
# Code: readin.py
# Description: Read in cl files


def readTable(list):
	ret = {}
	for line in list:
		ll = line.split("|")
		if len(ll) == 3:
			value = ll[1].strip()
			ret[ll[0].strip()] = value
	return ret

def readin(exp):
	global objTable
	global attrTable
	global cmat
	global latmatc
	objTable = readTable(open(exp+".objTable.txt","r").read().split("\n")[3:])
	attrTable = readTable(open(exp+".attrTable.txt","r").read().split("\n")[3:])
	cmat = [x for x in open(exp+".context.txt","r").read().split("\n") if len(x)>0]
	latmatc =  [x for x in open(exp+".latmat.txt","r").read().split("\n") if len(x)>0]
