#ifndef LAT_VEC_H
#define LAT_VEC_H

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

using namespace std;




//! Purpose: Storing Trace information
/*!
  For containing trace (object) information In the triple of <object,attribute,relation> trace files are objects
*/


class Trace{
	int id; /*!< Unique id of each trace */
	string name; /*!< Value of each trace*/
	static int traceCount; /*!< Total number of unique traces*/
	static map<string,int> ttable; /*!< Table to store trace values with their ID*/
public:
	//! Copy Constructor
	Trace();

	//! Destructor
	~Trace();

	//! Default Constructor
	/*!
	  Takes a string and add it to the table to receive unique ID
	*/
	Trace(string trace);

	//! Initialization of the Attribute Table
	static map<string,int> initTable(){
		map<string,int> tt ;
		tt = {} ;
		return tt;
	}

	//! Add trace to table to receive a unique ID for that table
	int addToTable(string trace);

	//! Return the ID of current trace
	int getID();

	//! Return the name of current trace
	string getName();

	//! Returns a string-style representation of current Trace Object
	string toString();

	//! Returns a Component-string-style representation of current Trace Object
	string toCompString();

	//! Returns a string-style representation of Trace Table
	string tableString();

	//! Returns a Componen-string-style representation of current Trace Table
	string tableCompString();
};


//! Purpose: Storing Concept information
/*!
  Each Concept has a unique ID, the object(trace) ID associated with it and a set of attribute IDs
*/
class Concept{
	int id; /*!< Unique id of each concept  */
	int objectId; /*!< Object(trace) ID of the concept*/
	set<int> attributeIds; /*!< All Attribute IDs of the concept*/
	static string relation; /*!< Relationship between Object(trace) and Attributes of a concept*/
	static int conceptCount; /*!< Total number of concepts */
public:
	//! Copy Constructor
	Concept();

	//! Default Constructor
	Concept(int objId , set<int> attrId);

	//! Returns a Component-string-style representation of current Concept
	string toCompString();

	//! Returns an string-style representation of current Concept
	string toString();

	//! Return the ID of current concept
	int getID();

	//! Return the object ID of current concept
	int getObjectID();

	//! Return a set of attribute IDs of current concept
	set<int> getAttributeIDs();



};



//! Purpose: Storing Vertex information
/*!
  After creating concepts, we feed each concept to the lattice one by one
  and form a bunch of verteices and edges to construct a graph (lattice).
  Each Vertex has a unique ID, a set of object(trace) IDs and a set of attribute IDs. Also it stores the ids of current vertex children and parents
*/
class Vertex{
	int id; /*!< Unique id of each vertex  */
	set<int> objIds; /*!< All Object(trace) IDs of the vertex */
	set<int> atrIds; /*!< All Attributes IDs of the vertex */
	map<int,double> childIds; /*!< To store each vertex children (int id) and its cost(double) */
	map<int,double> parentIds; /*!< To store each vertex parent (int id) and its cost(double) */
	//set<int> reachableUp; /*!< */
	static int vertexCount; /*!< Total number of vertices */
public:
	//! Copy Constructor
	Vertex();

	//! Default Constructor
	Vertex(set<int> _objIds,set<int> _atrIds);

	//! Returns an string-style representation of current Vertex (LabelID: Summarize lists? 0: No , 1: Yes , 2: Size of List)
	string toString(int labelID);

	//! Returns an string-style representation of printing current Vertex
	string toCompPrinting();

	//! Return the ID of current Vertex
	int getID();

	//! Return a set of object IDs of current vertex
	set<int> getObjectIDs() ;

	//! Return a set of attribute IDs of current vertex
	set<int> getAttributeIDs();

	//! Add a vertex(id) to current vertex's children
	void addToChilds(int childID, double cost);

	//! Add a vertex(id) to current vertex's parents
	void addToParents(int parID, double cost);

	//! Remove a vertex(id) to current vertex's parents
	void removeFromParents(int id);

	//! Remove a vertex(id) to current vertex's children
	void removeFromChilds(int id);

	//! Returns list of all children of current vertex
	map<int,double> getChilds();

	//! Returns list of all parents of current vertex
	map<int,double> getParents();

	//! Insert some attribute IDs to current vertex attribute IDs
	void ins2aids(const set<int>& _attrIDs){
		this->atrIds.insert(_attrIDs.begin(),_attrIDs.end());
	}

	//! Insert some object IDs to current vertex object IDs
	void ins2oids(const set<int>& _objIDs){
		this->objIds.insert(_objIDs.begin(),_objIDs.end());
	}


};


//! Purpose: Storing Edge information
/*!
  After creating concepts, we feed each concept to the lattice one by one
  and form a bunch of verteices and edges to construct a graph (lattice).
  Each Edge has a unique ID, the vertex ID of source, the vertex ID of destination and their corresponding cost
*/
class Edge{
	int id; /*!< Unique id of each Edge*/
	int srcId; /*!< Vertex ID of source (head of the edge) */
	int destId; /*!< Vertex ID of destination (end of the edge) */
	double cost; /*!< Cost to go from src to dest*/
	static int edgeCount; /*!< Total number of unique edges*/
public:
	//! Copy Constructor
	Edge();

	//! Default Constructor
	Edge(int src,int dest,double cost);

	//! Returns an string-style representation of current Edge
	string toString();

	//! Return the ID of current Edge
	int getID();

	//! Return the Vertex ID of source
	int getSrcID();

	//! Return the Vertex ID of destination
	int getDestID();

	//! Return the cost of edge
	double getCost();
};


#endif
