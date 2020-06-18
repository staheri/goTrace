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
#G.add_node("TEST")

# some math labels
labels = {}
labels[0] = r'$a$'
labels[1] = r'$b$'
labels[2] = r'$c$'
labels[3] = r'$d$'
labels[4] = r'$\alpha$'
labels[5] = r'$\beta$'
labels[6] = r'$\gamma$'
labels[7] = r'$\delta$'
labels[8] = "EvGoCreate"
labels[9] = "ChSend"
labels[10] = "ChRecv"
labels[11] = "GoEnd"



seq = [(0,1),(0,2),(2,3),(3,4),(3,5),(4,6),(5,8),(6,9),(8,10),(9,10),(10,11),(1,7),(7,11)]
msg = [(5,7),(6,8)]


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

pos = {0 : [0,30],
       1 : [0,29],
       2 : [1,28],
       3 : [1,27],
       4 : [2,26],
       5 : [1,25],
       6 : [2,24],
       7 : [0,23],
       8 : [1,22],
       9 : [2,21],
       10: [1,20],
       11: [0,19]}
      # "Test": [12,0]}

lab_pos = {0 : [.5,30],
       1 : [0.5,29],
       2 : [1.5,28],
       3 : [1.5,27],
       4 : [2.5,26],
       5 : [1.5,25],
       6 : [2.5,24],
       7 : [0.5,23],
       8 : [1.5,22],
       9 : [2.5,21],
       10: [1.5,20],
       11: [0.5,19]}

nx.draw_networkx_edges(G,pos,edgelist=seq)
nx.draw_networkx_edges(G,pos,edgelist=msg,style='dashed')
nx.draw_networkx_labels(G, lab_pos, labels, font_size=10)
nx.draw_networkx(G,pos, with_labels=False)
#plt.ylim(-2,15)
plt.show()
