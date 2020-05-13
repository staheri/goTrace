import math
import sys

# Recursive function for constructing MRR
def recList(seq,tp,i,j):

	if i == j:
		return seq[i-1]
	if tp[(i,j)][0] == 1:
		z = tp[(i,j)][1]
		c = recList(seq,tp,i,i+z-1)
		ret = ""
		for x in range(0,(j-i+1)/z):
			if x != ((j-i+1)/z) - 1:
				ret = ret + c + " -> "
			else:
				ret = ret + c
		#return c + " * "+ `(j-i+1)/z`
		return ret
	if tp[(i,j)][0] == 2:
		d = tp[(i,j)][1]
		a = recList(seq,tp,i,i+d)
		b = recList(seq,tp,i+d+1,j)
		return a + " -> "+ b

# Construction of cmr-like pattern through a recursive function
def constMRList(seq,tp):
	n = len(seq)
	return recList(seq,tp,1,n)

# Used in CMRX, returns the cmr-like pattern
def cmrList(seq):
	wr = 1
	we = 10
	l = {}
	tp = {}
	n = len(seq)
	for i in range(1,n+1):
		l[(i,i)] = we
	for i in range(1,n+1):
		for j in range(2,n+1):
			if i < j :
				l[(i,j)] = sys.maxint
	for z in range(1,n+1):
		for i in range(1,n+1-z+1):
			for d in range(0,z-2+1):
				if l[(i,i+d)] + l [(i+d+1,i+z-1)] < l[(i,i+z-1)]:
					l[(i,i+z-1)] = l[(i,i+d)] + l[(i+d+1,i+z-1)]
					tp[(i,i+z-1)] = (2,d)
			for h in range(1,(n-i+1)/z+1):
				ii = i + (h-1)*z
				seq1 = []
				seq2 = []
				if  ii-1 == ii+z-1-1:
					seq1 = seq[ii-1]
				else:
					seq1 = seq[ii-1:ii+z-1]
				if i-1 == i+z-1-1:
					seq2 = seq[i-1]
				else:
					seq2 = seq[i-1:i+z-1]
				if seq1 != seq2:
					break
				if l[(i,i+z-1)] + wr < l[(i,ii+z-1)]:
					l[(i,ii+z-1)] = l[i,i+z-1] + wr
					tp[(i,ii+z-1)] = (1,z)
	
	return constMRList(seq,tp)

# Used in topWins for finding the frequency of each repetitive pattern
def vectorFreq(vec):
	ret = {}
	for i in range(0,len(vec)):
		if vec[i] in ret.keys():
			ret[vec[i]].append(i)
		else:
			tmp = []
			tmp.append(i)
			ret[vec[i]] = tmp
	return ret

# Returns the winner of DBF (start and end point of most frequent repetitive pattern)
def topWins(inp,thr,maxRep):
	data = {}
	N = len(inp["name"])
	#print inp

	# Analyzing elements of DBF...
	#print "Analyzing elements of DBF..."
	for repLen in range(2,maxRep+1):
		seq = inp["name"][repLen]
		data[repLen]={}
		freqtab = {}
		loctab = {}
		extab = {}
		globExist = 0
		for element,locations in sorted(vectorFreq(seq).items()):
			extab[element] = 0
			if len(locations) >= thr:
				globExist = 1
				extab[element] = 1
			freqtab[element] = len(locations)
			loctab[element] = locations
		data[repLen]["exist"] = globExist
		data[repLen]["extab"] = extab	
		data[repLen]["freqtab"]=freqtab
		data[repLen]["loctab"]=loctab

	hdrs = ["Rep Len","First Loc","Actual Max Freq Elem","Max Freq","repLen * MaxFreq"]
	globalMaxFreq = 0
	globalMaxFreqRep = 0
	ttmp = []
	anal = {}
	for repLen,tabs in sorted(data.items()):
		if repLen != 1:
			tmp = ""
			tmp2 = []
			totFreq = 0
			maxFreq = 0
			firstLoc = 0
			maxFreqElem = ""
			maxmax = 0
			for elem,freq in sorted(tabs["freqtab"].items()):
				if freq >= thr:
					#tmp = tmp + "("+ `tabs["loctab"][elem][0]` +":"+ `freq` + ") , "
					totFreq = totFreq + freq
					if freq > maxFreq:
						maxFreq = freq
						firstLoc = tabs["loctab"][elem][0]
						maxFreqElem = elem
			tmp2.append(repLen)			
			tmp2.append(firstLoc)
			if maxFreqElem == '':
				tmp2.append("NULL" + " -> " + "NULL")
			else:
				tmp2.append("o|o" + " -> " + "-|-")
			tmp2.append(maxFreq)
			tmp2.append(maxFreq*repLen)
			ttmp.append(tmp2)
			anal[repLen] = tmp2

	#print tabulate(ttmp, hdrs, tablefmt="grid")
	ret = []
	maxmax = 0
	maxRepLen = 0
	for key,val in anal.items():
		#print key
		#print val
		if val[4] > maxmax:
			maxmax = val[4]
			maxRepLen = key
	#input("KIRRRRR")
	while maxRepLen > 1:
		if data[maxRepLen]["exist"] == 1:
			return anal[maxRepLen]
		else:
			maxRepLen = maxRepLen - 1
	return [0,0,0]

# Rename procedure of KMR algorithm
def rename(text):
	dict = {}
	y = []
	ret = {}
	temp = {}
	#print "Rename: Calculating Y"
	for i in range(0,len(text)):
		y.append((text[i],i+1))
	newKey = 1
	ltmp = dict.keys()
	for i in range(0,len(y)):
		if y[i][0] in ltmp:
			continue
		else:
			dict[y[i][0]] = newKey
			newKey = newKey + 1

	result = []
	for i in range(0,len(y)):
		result.append(0)
	for i in range(0,len(y)):
		result[y[i][1]-1] = dict[y[i][0]]
	ret["result"] = result
	ret["dictionary"] = temp
	return ret

# Implementation of KMR algorithm. Input: list S, int MaxRep, output: a data structure with "name" and "dict"
# as DBF (dictionary of basic factors) 	
def kmr(s,maxRep):
	data = {}
	#print "Creating Name table....\nlen: %d , maxRep: %d"%(len(s),maxRep)
	data["name"] = {}
	data["dict"] = {}
	#print "s:" + s 
	kk1 = rename(s)

	data["dict"][1] = kk1["dictionary"]
	data["name"][1] = kk1["result"]
	k = 2
	n = len(data["name"][1])
	for k in range(2,maxRep+1):
		#print "K: " + `k`
		
		if math.log(k,2).is_integer():
			y = []
			for j in range (0,n-k+1):
				y.append(0)
			for j in range(0,n-k+1):
				y[j] = (data["name"][k/2][j],data["name"][k/2][j+k/2])
			kk2 = rename(y)
			data["dict"][k] = kk2["dictionary"]
			data["name"][k] = kk2["result"]

		else:
			t = int(math.pow(2,int(math.log(k,2))))
			offset = k - t
			y = []
			for j in range (0,n-k+1):
				y.append(0)
			for j in range(0,n-k+1):
				#print "k/2:%s , j:%s , n:%s" %(k/2,j,n)
				y[j] = (data["name"][t][j],data["name"][t][j+offset])
			kk2 = rename(y)
			data["dict"][k] = kk2["dictionary"]
			data["name"][k] = kk2["result"]

	return data

# Returns the MRR representation of seq
def cmrx(seq):
	#print "Processing(cmrx) sequence of length %d ..."%(len(seq))
	n = len(seq)
	maxRep = 2
	dbf = kmr(seq,maxRep)
	thresh = 2
	ref = topWins(dbf,thresh,maxRep)
	if ref[0] != 0:
		#print "INSIDE CMRX - TOPWINS 1"
		pstart = ref[1]
		pend = pstart + ref[0]
		pattern = seq[pstart:pend]
		cmrPattern = cmrList(pattern)
		plen = ref[0]
		i = 0
		s = []
		cnt = 0
		cntb = 0
		maxcntb = 0 
		toBeAdded = ""
		while (i<len(seq)):
			#print "i: %d, len(seq): %d"%(i,len(seq))
			if seq[i:i+plen] == pattern:
				#print "Pattern[%d,%d]"%(i,i+plen)
				
				cnt = cnt + 1
				i = i + plen
			else:
				#print "No Pattern[%d]: %s"%(i,seq[i])
				toBeAdded = seq[i]
				cntb = cnt
				cnt = 0
				i = i + 1
			if cnt == 0 :
				if cntb == 0:
					s.append(toBeAdded)
				else:
					if cntb == 1:
						s.append(cmrPattern)
					else:
						s.append("(" + cmrPattern + ")^" + `cntb`)
					s.append(toBeAdded)
				if maxcntb < cntb:
					maxcntb = cntb
		if cnt !=0:
			cntb = cnt
			if cntb == 1:
				s.append(cmrPattern)
			else:
				s.append("(" + cmrPattern + ")^" + `cntb`)
			if maxcntb < cntb:
				maxcntb = cntb
		#print s
		return (s,maxcntb)
	else:
		#print "INSIDE CMRX - TOPWINS 0"
		return([],0)

def toMRR(sst):
	dataIter = {}
	dataIter[0] = {}
	dataIter[0]["seq"] = sst
	cout = cmrx(sst)
	dataIter[0]["cmrSeq"] = cout[0]
	i = 1
	while len(cout[0]) != 0:
		dataIter[i] = {}
		dataIter[i]["seq"] = dataIter[i-1]["cmrSeq"]
		cout = cmrx(dataIter[i-1]["cmrSeq"])
		dataIter[i]["cmrSeq"] = cout[0]
		if len(cout[0]) == 0:
			break
		else:
			i = i + 1
	if len (dataIter[i-1]["cmrSeq"]) != 0:
		return dataIter[i-1]["cmrSeq"]
	else:
		return dataIter[i-1]["seq"]