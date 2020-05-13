#!/usr/bin/env python

# Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2020, All rights reserved
# Code: readTrace.py
# Description: read Go traces (schedTrace)

import sys,subprocess

class P:
    def __init__(self,line):
        #"P0: status=1 schedtick=1 syscalltick=0 m=0 runqsize=0 gfreecnt=0"
        self.id=line.partition(':')[0].strip()[1:]
        self.status=line.partition("status=")[2].partition(' ')[0]
        self.schedtick=line.partition("schedtick=")[2].partition(' ')[0]
        self.syscalltick=line.partition("syscalltick=")[2].partition(' ')[0]
        self.m=line.partition("m=")[2].partition(' ')[0]
        self.runqsize=line.partition("runqsize=")[2].partition(' ')[0]
        self.gfreecnt=line.rpartition("=")[2].strip()

class G:
    def __init__(self,line):
        #"status=4(sleep) m=-1 lockedm=-1"
        self.id=line.partition(':')[0].strip()[1:]
        self.statusId=line.partition("status=")[2].partition('(')[0]
        self.status =line.partition("(")[2].rpartition('(')[0]
        self.m=line.partition("m=")[2].partition(' ')[0]
        self.lockedm=line.rpartition("=")[2].strip()

class M:
    def __init__(self,line):
        #"p=-1 curg=-1 mallocing=0 throwing=0 preemptoff= locks=0 dying=0 spinning=false blocked=false lockedg=-1"
        self.id=line.partition(':')[0].strip()[1:]
        self.p=line.partition("p=")[2].partition(' ')[0]
        self.curg=line.partition("curg=")[2].partition(' ')[0]
        self.mallocing=line.partition("mallocing=")[2].partition(' ')[0]
        self.throwing=line.partition("throwing=")[2].partition(' ')[0]
        self.preemptoff=line.partition("preemptoff=")[2].partition(' ')[0]
        self.locks=line.partition("locks=")[2].partition(' ')[0]
        self.dying=line.partition("dying=")[2].partition(' ')[0]
        self.spinning=line.partition("spinning=")[2].partition(' ')[0]
        self.blocked=line.partition("blocked=")[2].partition(' ')[0]
        self.lockedg=line.rpartition("=")[2].strip()

class Sched: # for each time-frame, we create an object containing all information about the current snapshot of schedules
    #"SCHED 400ms: gomaxprocs=2 idleprocs=2 threads=5 spinningthreads=0 idlethreads=3 runqueue=0 gcwaiting=0 nmidlelocked=0 stopwait=0 sysmonwait=0"
    def __init__(self,line,ps,gs,ms):
        self.time=line.partition('ms:')[0].partition(' ')[2]
        self.gomaxprocs=line.partition("gomaxprocs=")[2].partition(' ')[0]
        self.idleprocs=line.partition("idleprocs=")[2].partition(' ')[0]
        self.threads=line.partition("threads=")[2].partition(' ')[0]
        self.spinningthreads=line.partition("spinningthreads=")[2].partition(' ')[0]
        self.idlethreads=line.partition("idleThreads=")[2].partition(' ')[0]
        self.runqueue=line.partition("runqueue=")[2].partition(' ')[0]
        self.gcwaiting=line.partition("gcwaiting=")[2].partition(' ')[0]
        self.nmidlelocked=line.partition("nmidlelocked=")[2].partition(' ')[0]
        self.stopwait=line.partition("stopwait=")[2].partition(' ')[0]
        self.sysmonwait=line.rpartition("=")[2].strip()
        self.ps = ps
        self.gs = gs
        self.ms = ms
    def toString(self):
        s = ""
        s = s + "Sched Time: %s\n"%self.time
        s = s + "#P: %s\n"%`len(self.ps)`
        s = s + "#G: %s\n"%`len(self.gs)`
        s = s + "#M: %s\n"%`len(self.ms)`
        return s
    def toDotSing(self):

        edges={}
        invedges={}
        nodes={}
        ret = "digraph{\n\t"
        ret = ret + "rankdir=LR;\n\n\t"
        ret = ret + "S0 [label = \"t: "+self.time+"ms\"]\n\n\t"
        # P nodes
        group = "p"
        nodes[group]=[]
        shape = "box"
        for n in self.ps:
            name=""
            label = ""
            name = "P"+n.id
            label = name+"\\lrq: "+n.runqsize+", gfcnt: "+n.gfreecnt+"\\l"
            ret = ret + name + " " + "[label=\""+label+"\", group="+group+", shape="+shape+"]\n\t"
            nodes[group].append(name)
        # G nodes
        group = "g"
        nodes[group]=[]
        shape = "circle"
        for n in self.gs:
            name=""
            label = ""
            status = ""
            style = ""
            color = ""

            # ifs
            if n.statusId == "0":
                color = "yellow"
                status = "idle"
            elif n.statusId == "1":
                color = "green4"
                status = "runnable"
            elif n.statusId == "2":
                color = "green"
                status = "running"
            elif n.statusId == "3":
                color = "blue"
                status = "syscall"
            elif n.statusId == "4":
                color = "red"
                status = "waiting"
            elif n.statusId == "5":
                color = "orange"
                status = "unused"
            elif n.statusId == "6":
                color = "black"
                status = "dead"
            elif n.statusId == "7":
                color = "pink"
                status = "enqueue"
            elif n.statusId == "8":
                color = "brown"
                status = "copystack"
            #endifs

            name = "G"+n.id
            #label = name+"\\lstatus: "+status+"\\ldet: "+n.status+"\\l"
            label = name
            ret = ret + name + " " + "[label=\""+label+"\", group="+group+", color="+color+", shape="+shape
            if n.lockedm == "1":
                ret = ret + ", style=dashed]\n\t"
            else:
                ret = ret +"]\n\t"
            if n.m != "-1":
                edges[name] = "M"+n.m
            nodes[group].append(name)
        # M nodes
        group = "m"
        nodes[group]=[]
        shape = "triangle"
        for n in self.ms:
            name=""
            label = ""
            style = ""
            filledcolor = ""
            name = "M"+n.id
            label = name
            ret = ret + name + " " + "[label=\""+label+"\", group="+group+", shape="+shape+", fillcolor="
            if n.blocked == "true":
                ret = ret + "yellow , style="
            else:
                ret = ret + "white, style="
            if n.spinning == "true":
                if n.locks != "0":
                    ret = ret + "\"rounded,dashed,filled\""
                else:
                    ret = ret + "\"rounded,filled\""
            else:
                if n.locks != "0":
                    ret = ret + "\"dashed,filled\""
                else:
                    ret = ret + "filled"
            ret = ret +"]\n\t"
            if n.p != "-1":
                edges[name] = "P"+n.p
            nodes[group].append(name)

        # ranksames

        ret = ret +"{rank=same; "
        for n in nodes["p"]:
            ret = ret + n+";"
        ret = ret + "}\n\t"

        ret = ret +"{rank=same; "
        for n in nodes["m"]:
            ret = ret + n+";"
        ret = ret + "}\n\t"

        ret = ret +"{rank=same; "
        for n in nodes["g"]:
            ret = ret + n+";"
        ret = ret + "}\n\t"

        # Invis edges
        # sorting
        nodes["g"]=sorted(nodes["g"], key= lambda x: int(x[1:]))
        nodes["m"]=sorted(nodes["m"], key= lambda x: int(x[1:]))
        nodes["p"]=sorted(nodes["p"], key= lambda x: int(x[1:]))
        # S edges
        ret = ret + "S0 -> " + nodes["g"][0] + " [style=invis]\n\t"
        ret = ret + "S0 -> " + nodes["m"][0] + " [style=invis]\n\t"
        ret = ret + "S0 -> " + nodes["p"][0] + " [style=invis]\n\t"

        ret = ret + nodes["g"][0] + " -> " + nodes["m"][0] + " [style=invis]\n\t"
        ret = ret + nodes["m"][0] + " -> " + nodes["p"][0] + " [style=invis]\n\t"

        for i in range(0,len(nodes["g"])-1):
            ret = ret + nodes["g"][i]+" -> " + nodes["g"][i+1] + " [style=invis]\n\t"
        for i in range(0,len(nodes["p"])-1):
            ret = ret + nodes["p"][i]+" -> " + nodes["p"][i+1] + " [style=invis]\n\t"
        for i in range(0,len(nodes["m"])-1):
            ret = ret + nodes["m"][i]+" -> " + nodes["m"][i+1] + " [style=invis]\n\t"

        # edges
        for k,v in edges.items():
            ret = ret + k +" -> " + v + "\n\t"







        ret = ret+"}"
        return ret
    def toPdf(self,path,name,t):
        st = self.toDotSing()
        #print path+"/"+name+"-"+`t`+".dot"
        f = open(path+"/"+name+"-"+`t`+".dot","w")
        f.write(st)
        f.close()
        #print "dot -Tpdf "+path+"/"+name+"-"+`t`+".dot -o "+name+"-"+`t`+".pdf"
        process = subprocess.Popen("dot -Tpdf "+path+"/"+name+"-"+`t`+".dot -o "+name+"-"+`t`+".pdf", stdout=subprocess.PIPE,shell=True)
        si, err = process.communicate()

def readTrace(fi):
    f = open(fi,"r")
    prevS = ""
    ps=[]
    ms=[]
    gs=[]
    ss=[]
    while True:
        line = f.readline()
        #print line
        if not line:
            #flush the buffers
            if prevS != "":
                newS = Sched(prevS,ps,gs,ms)
                ss.append(newS)
                prevS = line
                ms=[]
                gs=[]
                ps=[]
            break
        if line[0] == 'S': #new Schedule
            if prevS != "":

                newS = Sched(prevS,ps,gs,ms)
                ss.append(newS)
                prevS = line
                ms=[]
                gs=[]
                ps=[]
            else:
                prevS=line

        else: #MPG
            t = line.partition(':')[0].strip()
            if t[0] == "P":
                newP = P(line)
                ps.append(newP)
            elif t[0] == "M":
                newM = M(line)
                ms.append(newM)
            elif t[0] == "G":
                newG = G(line)
                gs.append(newG)
            else:
                print "!!!! line: %s"%line
    return ss


if len(sys.argv) != 2:
	print "USAGE:\n\t " +sys.argv[0]+" traceFile"
	sys.exit(-1)

s = readTrace(sys.argv[1])
for item in s:
    print item.toString()

path = "dots"
name = sys.argv[1].rpartition('/')[2].rpartition('.')[0]
print s[0].toDotSing()

for i in range(0,len(s)):
    s[i].toPdf(path,name,i)
