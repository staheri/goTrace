import networkx as nx
import matplotlib.pyplot as plt


G = nx.Graph()

G.add_node(0)
G.add_node(1)
G.add_node(2)
G.add_node(3)
G.add_node(4)
G.add_node(5)
G.add_node(6)
G.add_node(7)
G.add_node(8)
G.add_node(9)
G.add_node(10)
G.add_node(11)


seq = [(0,1),(0,2),(2,3),(3,4),(3,5),(3,5),(4,7),(5,8),(7,9),(8,10),(9,10),(10,11),(1,6),(6,11)]
msg = [(5,6),(7,8)]


# G.add_edge(0,1)
# G.add_edge(0,2)
# G.add_edge(2,3)
# G.add_edge(3,4)
# G.add_edge(3,5)
# G.add_edge(3,5)
# G.add_edge(4,7)
# G.add_edge(5,8)
# G.add_edge(7,9)
# G.add_edge(8,10)
# G.add_edge(9,10)
# G.add_edge(10,11)
# G.add_edge(1,6)
# G.add_edge(6,11)
# G.add_edge(5,6,style='dashed')

pos = {0 : [0,0],
       1 : [1,0],
       2 : [2,1],
       3 : [3,1],
       4 : [4,2],
       5 : [5,1],
       6 : [6,0],
       7 : [7,2],
       8 : [8,1],
       9 : [9,2],
       10: [10,1],
       11: [11,0]}

nx.draw_networkx_edges(G,pos,edgelist=seq)
nx.draw_networkx_edges(G,pos,edgelist=msg,style='dashed')
nx.draw_networkx(G,pos)
#plt.ylim(-2,15)
plt.show()
