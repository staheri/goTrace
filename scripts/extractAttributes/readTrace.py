#!/usr/bin/env python

# Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2020, All rights reserved
# Code: readTrace.py
# Description: Read gopherlyzer traces and generate attributes



import glob
import sys,subprocess
import os
from tabulate import tabulate
import math
import time
from enum import Enum

class Goroutine:
	goroutines = {} # key: Unique id, value : Goroutine object
	rgoroutines = {} # key: Goroutine Object.funcName, value: id
	signals = {}
	def __init__(self,fname):
		self.id = len(Goroutine.goroutines)
		Goroutine.goroutines[self.id] = self
		Goroutine.rgoroutines[fname] = self.id
		self.sigID = 1000
		self.funcName = fname
		self.events = []
		#Goroutine.goroutines[]
	def addEvent(self,e):
		if len(e.ops) == 1 and e.ops[0].op == "W":
			self.sigID = e.ops[0].op
		self.events.append(e)

class Resource:
	resources = {}
	rresources = {}
	def __init__(self,type,realID,cap):
		self.id = len(Resource.resources)+1
		Resource.resources[self.id] = self
		Resource.rresources[type+realID] = self.id
		self.realID = realID
		self.type= type
		self.cap= cap
		if type == "S":
			self.exID = realID
		elif type == "C":
			self.exID = type+realID
		else:
			self.exID = type+`self.id`

class Op:
	ops = {}
	def __init__(self,eid,rid,cap,op,loc):
		self.id = len(Op.ops)
		Op.ops[self.id] = self
		self.eid = eid
		self.rid = rid
		self.cap = cap
		self.op = op
		self.loc = loc
	def toString(self):
		s = "<" + `self.eid` + "-" + `self.id` + ">:"
		s = s + "[" + `self.rid` + "," + self.cap + "," + self.op + "," + self.loc + "]"
		return s

class Event:
	events = {}
	def __init__(self,eventLine):
		self.id = len(Event.events)
		Event.events[self.id] = self
		line = self.splitLine(eventLine)
		self.owner = newGoroutine(line[0])
		#self.ops = genOps(eventLine)
		self.predicate = line[1]
		self.partner = line[2].strip()
	def splitLine(self,line):
		#fun020,[(824634490936,1,+,newTests/dl/dl-lock.go:15)],P,-
		ret = []
		ret.append(line.partition(",")[0]) #owner
		#cont = line.rpartition(")")[0].partition("(")[2].split(",")
		self.ops = self.genOps(line.rpartition("]")[0].partition("[")[2],self.id)
		# ret.append(cont[0]) #resourceID
		# ret.append(cont[1]) #cap
		# ret.append(cont[2]) #op
		# ret.append(cont[3].rpartition("/")[2]) #loc
		ret.append(line.split(",")[-2]) #predicate
		ret.append(line.split(",")[-1]) #partner
		return ret
	def genOps(self,l,eid):
		#print l
		ret = []
		for it in l.split(")")[:-1]:
			t = it.partition("(")[2].split(",")
			#print t
			opObj = Op(eid,t[0],t[1],t[2],t[3].rpartition("/")[2])
			if opObj.op in ["+","*"]: #it is a lock
				print "before: M"+opObj.rid
				opObj.rid = newResource("M",opObj.rid,opObj.cap)
				print "after: M"+`opObj.rid`
			if opObj.op in ["?","!","C"]: #it is a channel
				print "before: C"+opObj.rid
				opObj.rid = newResource("C",opObj.rid,opObj.cap)
				print "after: C"+`opObj.rid`
			if opObj.op in ["W","S"]: #it is a channel
				print "before: S"+opObj.rid
				opObj.rid = newResource("S",opObj.rid,opObj.cap)
				print "after: S"+`opObj.rid`
			ret.append(opObj)
		return ret
	def toString(self):
		s = ""
		s = s + `self.id` + ": ("
		s = s + `self.owner`+":" + Goroutine.goroutines[self.owner].funcName + ")-> "
		for op in self.ops:
			s = s + op.toString() + ","
		s = s + self.predicate + "," + self.partner
		return s
	def toGrtnAtrFormat(self,withID):
		if len(self.ops) > 1: #select
			s = "SE-"
			for t in self.ops:
				if t.op == "!":
					s = s + "S"
				elif t.op == "?":
					s = s + "R"
				elif t.op == "+":
					s = s + "L"
				elif t.op == "*":
					s = s + "U"
				else:
					s = s + "X"
				s = s + ","
		else:
			s = self.predicate
			t = self.ops[0]
			if t.op == "!":
				s = s + "S"
			elif t.op == "?":
				s = s + "R"
			elif t.op == "+":
				s = s + "L"
			elif t.op == "*":
				s = s + "U"
			elif t.op == "S":
				s = s + "Signal"
			elif t.op == "W":
				s = s + "Wait"
			else:
				s = s + "X"
			if withID == 1:
				if s == "CR":
					if self.partner == "-":
						s = s + "-" + Resource.resources[t.rid].exID + "(def Select)"
					else:
						s = s + "-" + Resource.resources[t.rid].exID + "-g" + `Goroutine.goroutines[Goroutine.rgoroutines[self.partner]].id`
				else:
					s = s + "-" + Resource.resources[t.rid].exID
		if withID == 2:
			s = s + ":" + self.ops[0].loc
		if withID == 3:
			if self.ops[0].loc != "-":
				s = s + ":<source>." + self.ops[0].loc.partition(".")[2]
			else:
				s = s + ":" + self.ops[0].loc

		return s

def newResource(type,realID,cap):
	if type+realID in Resource.rresources.keys():
		return Resource.rresources[type+realID]
	else:
		print ">>>"+type+realID+" NOT IN RESOURCES - CREATE NEW"
		return Resource(type,realID,cap).id

def newGoroutine(fname):
	if fname in Goroutine.rgoroutines.keys(): #fname exists - return id
		return Goroutine.rgoroutines[fname]
	else:		#create new goroutine, return id
		return Goroutine(fname).id

def genAttributes(path,withID):
	if withID == 1:
		path = path+"id/"
	elif withID == 2:
		path = path+"loc/"
	elif withID == 3:
		path = path+"locx/"
	else:
		path = path + "/"
	process = subprocess.Popen(["mkdir -p "+path], stdout=subprocess.PIPE,shell=True)
	si, err = process.communicate()
	for id,g in Goroutine.goroutines.items():
		fname = path+"/g"+`id`+"-"+g.funcName+".txt"
		f = open(fname,"w")
		print "g"+`id`+"-"+g.funcName
		for e in g.events:
			f.write(e.toGrtnAtrFormat(withID)+"\n")
			print "\t"+e.toGrtnAtrFormat(withID)


if len(sys.argv) != 2:
	print "USAGE:\n\t " +sys.argv[0]+" gtrace<full-path>"
	sys.exit(-1)

gtrace = sys.argv[1]

lines = open(gtrace,"r").readlines()
for l in lines:
	e = Event(l)
	Goroutine.goroutines[e.owner].addEvent(e)

for k,v in Event.events.items():
	print k
	print v.toString()


hdrs = ["Obj","Attributes"]
tab = []

for k,v in Goroutine.goroutines.items():
	# k : gid
	# v : goroutine object
	ttab = []
	ttab.append(v.funcName)
	s = ""
	for e in v.events:
		s = s + e.toString() + "\n"
	ttab.append(s)
	tab.append(ttab)
#print tabulate(tab,headers=hdrs,tablefmt="fancy_grid")


#outpath = gtrace.rpartition("/")[0]+"/"+gtrace.rpartition("/")[2].rpartition(".")[0]+"/
outpath = gtrace.rpartition("/")[0]+"/gtrace"
genAttributes(outpath,1)
genAttributes(outpath,0)
genAttributes(outpath,2)
genAttributes(outpath,3)

print outpath
