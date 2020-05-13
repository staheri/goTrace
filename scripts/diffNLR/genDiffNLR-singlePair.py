#!/usr/bin/env python
# Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2018, All rights reserved
# Code: cmrx.py
# Description: Read texTrace and convert them to MRR (minimal reptetive representation)

import glob
import sys,subprocess
from tabulate import tabulate
import math
import re
from sets import Set
import mrrCore as mrr
import diffCore as diff
import newVis as nvis
import toml

conf={}

def loopExpansion(l,ltab):
	ret = []
	for item in l:
		ls = item.partition("^")
		if len(ls[1]) > 0 and len(ls[2]) > 0 :
			# it is a loop, replace it
			print ls
			for lc in range(0,int(ls[2])):
				for lb in ltab[ls[0][1:]]:
					ret.append(lb)
		else:
			ret.append(item)
	return ret

def readLtab(lf):
	ret = {}
	lines = [x for x in open(lf,"r").readlines()]
	for l in lines:
		lid = l.partition(":")[0]
		lbs = l.partition(":")[2].strip().split(" - ")
		ret[lid] = lbs
	return ret



def ltab2html(ltab,pre):
	lines = [x for x in open(ltab,"r").readlines()]
	s = "{\n\n\t"
	s = s + pre+"HtmlTable [\n\t\t"
	s = s + "shape=plaintext\n\t\t"
	if pre == "b":
		s = s + "color=Red\n\t\t"
	elif pre == "nb":
		s = s + "color=Blue\n\t\t"
	s = s + "label=<\n\t\t\t"
	s = s + "<table border=\'0\' cellborder=\'1\'>\n\t\t\t\t "
	s = s + "<tr><td> Loop </td> <td> Body"
	if pre == "b":
		s = s + " (buggy)"
	elif pre == "nb":
		s = s + " (noBug)"
	s = s + " </td></tr>\n\t\t\t\t "
	for line in lines:
		s = s + "<tr><td> L"+line.partition(":")[0]+" </td> <td> "+line.partition(":")[2].strip()+" </td></tr>\n\t\t\t\t "
	s = s + "</table>\n\t>];\n}"
	return s



if len(sys.argv) != 4:
	print "USAGE:\n\t " +sys.argv[0]+" A(noBug) B(buggy) out"
	sys.exit(-1)


pathA = sys.argv[1] # noBug
pathB = sys.argv[2] # buggy
pathO = sys.argv[3] # Out

#pathAltab = readLtab(pathA.rpartition("/")[0]+"/ltab.txt")
#pathAexpanded = loopExpansion([x.strip() for x in open(pathA,"r").readlines() if len(x) > 0],pathAltab)

#pathBltab = readLtab(pathB.rpartition("/")[0]+"/ltab.txt")
#pathBexpanded = loopExpansion([x.strip() for x in open(pathB,"r").readlines() if len(x) > 0],pathBltab)
pathAlines = [x.strip() for x in open(pathA,"r").readlines() if len(x) > 0]
pathBlines = [x.strip() for x in open(pathB,"r").readlines() if len(x) > 0]

nm = pathA.rpartition("/")[2].rpartition(".")[0]+"_"+pathB.rpartition("/")[2].rpartition(".")[0]
print nm
fe = diff.lcs(pathAlines,pathBlines)
print fe
#ltab = ltab2html(pathA.rpartition("/")[0]+"/ltab.txt","b")
#ltab_nb = ltab2html(pathB.rpartition("/")[0]+"/ltab.txt","nb")
diffNLR = nvis.edit2dot(fe,nm,1)

finalDot= "digraph \"diffNLR\""+diffNLR+"\n\t"
#finalDot= finalDot +"subgraph diffNLR"+diffNLR+"\n\t"
#finalDot= finalDot +"subgraph ltab_nb"+ltab_nb+"\n\t"
#finalDot= finalDot +"}"

f = open(nm+".dot","w")
f.write(finalDot)
f.close()
process = subprocess.Popen("dot -Tpdf "+nm+".dot -o "+pathO+".pdf", stdout=subprocess.PIPE,shell=True)
si, err = process.communicate()
