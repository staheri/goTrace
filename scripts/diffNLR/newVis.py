
class diffMrrGraph:
	def __init__(self,name):
		self.name = name
		self.nodes = {}
		self.nodes["s"] = "[label = \"{Start}\" , group=g0]"
		self.nodes["f"] = "[label = \"{End}\" , group=g0]"
		self.edges = []
		self.eseq = []
		self.invisNodes = {}
	def addNodes(self,eseq,showC):
		self.eseq = eseq
		if len(eseq) <= 0:
			print "ERROR"
			sys.exit(-1)
		else:
			for i in range(0,len(eseq)):
				# Add C
				if len(eseq[i].c) != 0:
					nodeContent = "[label = \"{"
					if showC == 1:
						seq = specialCharFilter(eseq[i].c)
						for item in seq:
							nodeContent = nodeContent + item+"\\l"
					else:
						nodeContent = nodeContent + "hidden"
					nodeContent = nodeContent + "}\" , group=g0, color=green4]"
					self.nodes["c"+`i`] = nodeContent
				else:
					if i != 0:
						print "BIG ERROR"
						break

				# Add A
				if len(eseq[i].a) == 0:
					self.invisNodes["a"+`i`] = "[label = \"{garbage}\",group=g1 , color=blue, style=invis]"
				else:
					nodeContent = "[label = \"{"
					seq = specialCharFilter(eseq[i].a)
					for item in seq:
						nodeContent = nodeContent + item+"\\l"
					nodeContent = nodeContent + "}\" , group=g1, color=blue, style = bold]"
					self.nodes["a"+`i`] = nodeContent

				# Add B
				if len(eseq[i].b) == 0:
					self.invisNodes["b"+`i`] = "[label = \"{garbage}\",group=g2 , color=red, style=invis]"
				else:
					nodeContent = "[label = \"{"
					seq = specialCharFilter(eseq[i].b)
					for item in seq:
						nodeContent = nodeContent + item+"\\l"
					nodeContent = nodeContent + "}\" , group=g2, color=red, style = dashed]"
					self.nodes["b"+`i`] = nodeContent
	def addEdges(self):
		alist = []
		blist = []
		clist = []
		for item in self.nodes.keys():
			if item.startswith('a'):
				alist.append(item)
			elif item.startswith('b'):
				blist.append(item)
			elif item.startswith('c'):
				clist.append(item)
			elif item.startswith('s') or item.startswith('f'):
				continue
			else:
				print "Error. Node started with something other than a,b,c,s or f"
		# from Start to others
		flg = [0,0,0]
		if "c0" in clist:
			self.edges.append("s -> c0")
			flg[0] = 1
		else:
			if "a0" in alist:
				self.edges.append("s -> a0")
				flg[1] = 1
			if "b0" in blist:
				self.edges.append("s -> b0")
				flg[2] = 1
		if (flg == [0,1,0] or flg == [0,0,1]) and "c1" in clist:
			self.edges.append("s -> c1")
		# from others to End

		#rest
		# Add invisible edges of As and Bs
		for i in range(0,len(self.eseq)-1):
			for case in ["a","b"]:
				self.edges.append(case+`i` + " -> " + case+`i+1` + " [style = invis]")

		# Add edges from Cs to their As and Bs or the next C
		clistSorted = sorted(clist,key = lambda k: int(k[1:]))
		for ll in range(0,len(clistSorted)):
			c = clistSorted[ll]
			i = int(c[1:])
			if "a"+`i` in alist:
				self.edges.append(c+" -> "+"a"+`i`)
			if "b"+`i` in blist:
				self.edges.append(c+" -> "+"b"+`i`)
			if "b"+`i` not in blist or "a"+`i` not in alist:
				if "c"+`i+1` in clistSorted:
					self.edges.append(c+" -> "+"c"+`i+1`)
		# Add edges of As and Bs to the next Cs
		alistSorted = sorted(alist,key = lambda k: int(k[1:]))
		for ll in range(0,len(alistSorted)):
			a = alistSorted[ll]
			i = int(a[1:])
			if "c"+`i+1` in clist:
				self.edges.append("a"+`i` +"-> c"+`i+1`)
		blistSorted = sorted(blist,key = lambda k: int(k[1:]))
		for ll in range(0,len(blistSorted)):
			b = blistSorted[ll]
			i = int(b[1:])
			if "c"+`i+1` in clist:
				self.edges.append("b"+`i` +"-> c"+`i+1`)

		# END part
		flg =[0,0,0]
		if "a"+`len(self.eseq)-1` in alist:
			self.edges.append("a"+`len(self.eseq)-1` + " -> f")
			flg[1] = 1
		if "b"+`len(self.eseq)-1` in blist:
			self.edges.append("b"+`len(self.eseq)-1` + " -> f")
			flg[2] = 1
		if flg == [0,1,0] or flg == [0,0,1] or flg == [0,0,0]:
			self.edges.append("c"+`len(self.eseq)-1` + " -> f")
	def toDot(self):
		s = "{\n\tnode[shape=record]\n\n\t"
		#s = s + genLegend() + "\n\n\t"
		for node,cont in sorted(self.nodes.items(),key=lambda k: k[1:]):
			print node
			if (node[0] == "a" and "b"+node[1:] in self.nodes.keys()) or (node[0] == "b" and "a"+node[1:] in self.nodes.keys()):
				s = s + "{rank = same ; a"+node[1:]+self.nodes["a"+node[1:]]+" ; b" + node[1:]+self.nodes["b"+node[1:]]+"}\n\t"
			else:
				s = s + node + " " + cont + "\n\t"
		for node,cont in sorted(self.invisNodes.items(),key=lambda k: k[1:]):
			print node
			s = s + node + " " + cont + "\n\t"

		s = s + "\n\t"
		for edge in self.edges:
			s = s + edge + "\n\t"
		s = s + "\n}"
		print s
		return s




class editSeq:
	def __init__(self,a,b,c):
		self.a = a
		self.b = b
		self.c = c
	def toString(self):
		s = "C: %s\n"%self.c
		s = s + "\tA: %s\n"%self.a
		s = s + "\tB: %s\n"%self.b
		return s

def genLegend():
	s = "subgraph cluster_legend{\n\t\trankdir= TP\n\t\t"
	s = s+"	label = \"Legend\" ;\n\t\t"
	s = s+"shape=rectangle  ;\n\t\t"
	s = s+"color = black  ;\n\t\t"
	s = s + "\"Block of Native Thread\" [shape=record , color=blue] ; \n\t\t"
	s = s + "\"Block of Buggy Thread\" [shape=record , color=red] ; \n\t\t"
	s = s + "\"Common Block in both\" [shape=record , color=green4] ; \n\t\t"
	s = s + "\"Common Block in both\" -> \"Block of Native Thread\" [style = invis]; \n\t\t"
	s = s + "\"Block of Native Thread\" -> \"Block of Buggy Thread\"[style = invis]; \n\t\t"
	s = s+"}"
	return s
def processBuf(buf):
	ret = {}
	inserts = []
	deletes = []
	if len(buf) != 0:
		for item in buf:
			if item[0] == 1:
				#add to inserts
				for tt in [x.strip().strip("'") for x in item[1].rpartition("[")[2].partition("]")[0].split(",")] :
					inserts.append(tt)
			else:
				for tt in [x.strip().strip("'") for x in item[1].rpartition("[")[2].partition("]")[0].split(",")] :
					deletes.append(tt)
				#add to deletes
	ret["A"]=deletes
	ret["B"]=inserts
	return ret

def mergeCs(li):
	i = 0
	line = []
	while i < len(li):
		if li[i].startswith("C:"):
			#Check next ones, look for merging options
			j = i + 1
			befc = [x.strip().strip("'") for x in li[i].rpartition("[")[2].partition("]")[0].split(",") if len(x) > 0]
			#print befc
			s = "C: ["
			while j < len(li)-1 and li[j].startswith("C:"):
				#print j
				#print len(li)
				print "\t INSIDE WHILE\n\tLine[%d] to process: %s"%(j,li[j])
				toAdd = [x.strip().strip("'") for x in li[j].rpartition("[")[2].partition("]")[0].split(",") if len(x) > 0]
				for item in toAdd:
					befc.append(item)
				j = j + 1
			i = j
			i = j
			for k in range(0,len(befc)):
				#print befc[k]
				if k != len(befc) - 1:
					s = s + "'"+befc[k]+"',"
				else:
					s = s + "'"+befc[k]+"'"
			s = s + "]"
			#print "\tS: "+s
			line.append(s)
		else:
			line.append(li[i])
			i = i + 1
	# Removing empty Cs, As and Bs
	ret = []
	for item in line:
		if len(item.rpartition("[")[2].rpartition("]")[0]) != 0:
			ret.append(item)
	return ret

def edit2eseq(li):
	line = mergeCs(li)
	i = 0
	prevC = -1
	buf = []
	eseqObjs = []
	while i < len(line):
		#print "Line to process: %s"%(line[i])
		if not line[i].startswith("C:"):
			if line[i].startswith("B:"):
				buf.append((1,line[i]))
			elif line[i].startswith("A:"):
				buf.append((0,line[i]))
		elif line[i].startswith("C:"):
			if prevC == -1 and len(buf) != 0:
				# Empty C should be inserted
				# Create editSeq object using contents of buf (for A and B)
				a = processBuf(buf)["A"]
				b = processBuf(buf)["B"]
				c = []
				obj = editSeq(a,b,c)
				#print obj.toString()
				eseqObjs.append(obj)
				prevC = i
				buf = []
			elif prevC == -1 and len(buf) == 0:
				# edit starts with C, only update prevC, do nothing
				prevC = i
			elif prevC != -1 and len(buf) != 0:
				# Create editSeq object using prevC and contents of buf (for A and B)
				a = processBuf(buf)["A"]
				b = processBuf(buf)["B"]
				c = [x.strip().strip("'") for x in line[prevC].rpartition("[")[2].partition("]")[0].split(",") if len(x) > 0]
				prevC = i
				obj = editSeq(a,b,c)
				#print obj.toString()
				eseqObjs.append(obj)
				buf = []
			else:
				print "Error, two consequitive Cs in the edit"
				sys.exit(-1)
		else:
			print "Error, Line starts with something other than A, B or C"
			sys.exit(-1)
		i = i + 1

	if len(buf) == 0:
		a = []
		b = []
	else:
		a = processBuf(buf)["A"]
		b = processBuf(buf)["B"]
	if prevC != -1:
		c = [x.strip().strip("'") for x in line[prevC].rpartition("[")[2].partition("]")[0].split(",") if len(x) > 0]
	else:
		c = []
	obj = editSeq(a,b,c)
	#print obj.toString()
	eseqObjs.append(obj)
	return eseqObjs

def edit2dot(lcs,name,showC):
	li = [x.strip() for x in lcs.split("\n") if len(x) != 0]
	es = edit2eseq(li)
	dmg = diffMrrGraph(name)
	dmg.addNodes(es,showC)
	dmg.addEdges()
	return dmg.toDot()


def specialCharFilter(seq):
	special = ["}","{",">","<"]
	retSeq = []
	for item in seq:
		retItem = ""
		for i in range(0,len(item)):
			if item[i] in special:
				retItem = retItem + "-"
			else:
				retItem = retItem + item[i]
		retSeq.append(retItem)
	return retSeq
