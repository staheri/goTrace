#!/usr/bin/env python
# Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2020, All rights reserved
# Code: main.py
# Description:
#        - Generate Jaccard similarity matrix from CL dot file (includes CL redundant removal and LCA) (input: dot)
#        - Generate reduced_label concept lattice graph
#        - Hierarchicaly cluster data based on the jacmat


if len(sys.argv) != 3:
	print "USAGE:\n\t " +sys.argv[0]+" path out"
	sys.exit(-1)

outName = sys.argv[2]
outdir = ""



for ccl in glob.glob(sys.argv[1]+"/cl/*.dot"):
	print "Generating Matrices...\npath:"+ccl
	#READ DOT FILE
	exp = ccl.rpartition(".")[0]
	fi = open(ccl,"r").read().split("\n")
	f = [x for x in fi if "->" in x]


	# READ TABLES, CONTEXT AND LATMAT
	objTable = readTable(open(exp+".objTable.txt","r").read().split("\n")[3:])
	attrTable = readTable(open(exp+".attrTable.txt","r").read().split("\n")[3:])

	if len(attrTable) == 0 :
		simmax2(len(objTable.keys()),ccl)
		continue
	cmat = [x for x in open(exp+".context.txt","r").read().split("\n") if len(x)>0]
	latmatc =  [x for x in open(exp+".latmat.txt","r").read().split("\n") if len(x)>0]
	fullMatrix = latmatToFullMat(latmatc)
	lat = Lattice(exp)

	# Initialize and Generate Lite Lattice Object
	ll = LiteLat(len(latmatc))
	for i in range(0,len(latmatc)):
		if latmatc[i] != "-":
			#print latmatc[i].split(" ")
			for u in latmatc[i].split(" "):
				if len(u) > 0:
					ll.addEdge(i,int(u))
		else:
			ll.addEdge(i,-1)

	# Initialize and Generate Original Lattice Object from DOT files
	for item in f:
		edge = item.split("->")
		snode = edge[0]
		dnode = edge[1]
		sid = int(snode.partition(":")[0].partition("\"")[2])-1
		did = int(dnode.partition(":")[0].partition("\"")[2])-1
		lat.addNode(sid,snode)
		lat.addNode(did,dnode)
		lat.addEdge(sid,did)


	# Operations on original Lattice (lat) for reducing labels

	# Prepare Lattice : detect sup, remove redundants
	lat.supinfDetection()
	lat.assignLabel()

	#print lat.toString()
	#print lat.toReducedString()
	out = sys.argv[2]
	outname = out+"-"+ccl.rpartition(".")[0].rpartition("/")[2]
	f = open(outname+".dot","w")
	f.write(lat.toReucedFancyDot())
	f.close()

	#cmd = "python /home/saeed/diffTrace/scripts/genReport/genSingleJSM.py " + f + " " + ex + ";"

	cmd = "dot -Tpdf "+outname+".dot -o "+outname+".pdf"
	print cmd
	process = subprocess.Popen([cmd], stdout=subprocess.PIPE,shell=True)
	si, err = process.communicate()
	cmd = "open "+outname+".pdf"
	print cmd
	process = subprocess.Popen([cmd], stdout=subprocess.PIPE,shell=True)
	si, err = process.communicate()


	# Prepration for LCA
	ll.DFS(lat.supID)
	ll.LCA_1partition()
	ll.LCA_2createLists()


	# Compute Jaccard Similarity Matrix
	simmax(lat,ll,ccl)
