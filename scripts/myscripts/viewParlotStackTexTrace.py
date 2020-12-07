import sys

def printStack(s):
    for item in s:
        print item,
    print


if len(sys.argv) != 2:
    print "USAGE:\n\t " +sys.argv[0]+" traceFile"
    sys.exit(-1)

f= open(sys.argv[1],"r")

lines = [l[:-1] for l in f.readlines()]
stack=[]

cnt = 0

for x in lines:
    cnt = cnt + 1
    #if cnt%100==0:
    #    input()
    if x == "[ret]":
        stack= stack[:-1]
    else:
        stack.append(x)
        print cnt,
        printStack(stack)


