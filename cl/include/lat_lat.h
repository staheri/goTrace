#ifndef LAT_LAT_H
#define LAT_LAT_H

#include <iostream>
#include <vector>
#include <map>
#include <vector>
#include <set>
#include <string>
#include <string.h>
#include <stdio.h>
#include <algorithm>
#include <fstream>
#include <sstream>
#include <iterator>
#include <queue>
#include <typeinfo>
#include <iomanip>
#include "util.h"
#include "lat_vec.h"

using namespace std;



class Concept;
class Trace;
template <class T> class Attribute;
class Vertex;
class Edge;


//! Purpose: Storing Lattice information
/*!
  It contains all information about different concepts and the generated vertices and edges.
  Also all of lattice operations such as adding new concept, new vertex, new edge and any other
  operatio are taking place in this class.
*/
class Lattice{
	string name; /*!< Name of the lattice*/
	string constDot; /*!< Subgraphs of the lattice during construction (addConcept) in dot format*/
	static map<int,Concept> CV; /*!< Stores all concepts*/
	static map<int,Vertex> VV; /*!< Stores all vertices*/
	static map<int,Edge> EV; /*!< Stores all edges*/

	static int supID; /*!< Vertex ID of lattice SUP */
	static int infID; /*!< Vertex ID of lattice INF */
	static int maxAttribute; /*!< Maximum number of attributes among all vertices*/

public:
	//! Default constructor
	Lattice();

	//! Copy Constructor
	Lattice(string _name);

	// IO methods
	//! Returns an string-style representation of current Lattice
	string toString();

	//! Returns a dot graph representation of current Lattice (LabelID: Summarize lists? 0: No , 1: Yes , 2: Size of List)
	string toDot(string label, int labelID);

	//! Returns a dot graph representation of current Lattice (only edges, no header/wrapper) (LabelID: Summarize lists? 0: No , 1: Yes , 2: Size of List)
	string toDotEdges(int cid, int labelID);

	//! Returns the lattice construction graph (subgraphs)
	string toDotConst();

	//! Returns a string-style representation of all Concepts within the lattice
	string toConceptString();

	//! Generate context bit-matrix
	string toContextBitmax();

	//! Generate Lattice Matrix
	string toLatMat();

	//! Add the subgraph for ConstSubgraph by adding each concpet
	void addSubgraph(string subg,string label);

	//! Add subgraph to constDot for each concept
	void addSubgraph(string subg);

	//! Prints lattice components
	void printLatticeComponents();


	//! Add new concept to lattice
	void addConcept(Concept& c);

	//! Add new vertex to lattice
	void addVertex(Vertex v);

	//! Add new edge to lattice
	void addEdge(Edge e);

	//! Delete an edge from lattice
	void deleteEdge(int src,int dest);


	//! Set suprimom of the lattice
	void setSup(int sid);

	//! Set infimum of the lattice
	void setInf(int iid);

	//! Set maximum number of attributes of the lattice (maxAttribute)
	void setMaxAttribute(int val);


	//! Get suprimum id
	int getSup();

	//! Get infimum id
	int getInf();

	//! Get maxAttribute
	int getMaxAttribute();

	// Returns true if Vertex ID a is parent of Vertex ID b in this lattice
	bool isParent(int a,int b);
};


#endif
