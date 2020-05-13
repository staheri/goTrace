#include "lat_vec.h"
#include "lat_atr.h"
#include <string>

using namespace std;


//////////////////////////////////////////////////////////////////////////////
// Class TRACE definitions
//////////////////////////////////////////////////////////////////////////////
int Trace::traceCount = 0;
map<string,int> Trace::ttable = initTable();

Trace::Trace(){}

Trace::~Trace(){}

Trace::Trace(string trace){
	//printf("\tNew Trace is constructing: %s\n\n", trace.c_str());
	this->id = addToTable(trace);
	this->name = trace;
}


int Trace::addToTable(string trace){
	map<string,int>::iterator it = this->ttable.find(trace);
	if (it == this->ttable.end() ){
		this->traceCount++;
		this->ttable[trace] = this->traceCount;
	}
	return this->ttable[trace];
}

int Trace::getID(){
	return this->id;
}

string Trace::getName(){
	return this->name;
}

string Trace::toString(){
	string s="\t---------------------";
	s = s + "\n\tTrace ID:" + intToString(this->id);
	s = s + "\n\t----------------------";
	s = s + "\n\tTrace Name:" + this->name;
	s = s + "\n\tTraceCount: " + intToString(this->traceCount);
	s = s + "\n\t-----------end---------\n";
	return s;
}

string Trace::tableCompString(){
	map<string,int>::iterator it;
	string s="\t--------------------- Trace Table ---------------------";
	s = s + "\n\tTrace ID:" + intToString(this->id);
	s = s + "\n\t-------------------------------------------------------";
	s = s + "\n\tTrace Name | Trace ID";
	s = s + "\n\t_____________________";
	for (it = this->ttable.begin() ; it != this->ttable.end() ; it++){
		s = s + "\n\t" + it->first + " | " + intToString(it->second);
	}
	s = s + "\n\t_____________________";
	s = s + "\n\t--------------------------end--------------------------\n";
	return s;
}

string Trace::tableString(){
	map<string,int>::iterator it;
	stringstream ss;
	string s;

	const char separator    = ' ';
    const int nameWidth     = 20;
    const int numWidth      = 5;

	ss << "--------------------------" << endl;
    ss << left << setw(numWidth) << setfill(separator) << "ID" << "|";
    ss << left << setw(nameWidth) << setfill(separator) << "Trace" << "|" << endl;
	ss << "--------------------------" << endl;
    s = s + ss.str();

	for (it = this->ttable.begin() ; it != this->ttable.end() ; it++){
		stringstream sss;
		sss << left << setw(numWidth) << setfill(separator) << it->second << "|";
		sss << left << setw(nameWidth) << setfill(separator) << it->first << "|" << endl;
		s = s + sss.str();
	}
	s = s + "--------------------------\n";
	return s;
}

//////////////////////////////////////////////////////////////////////////////
// Class CONCEPT definitions
//////////////////////////////////////////////////////////////////////////////

int Concept::conceptCount = 0;
string Concept::relation = "contains";

Concept::Concept(){}

Concept::Concept(int objid,set<int> attrid){
	//printf(">>> New Concept is constructing, objid : %d\n", objid);
	this->attributeIds = attrid;
	this->objectId = objid;
	this->conceptCount++;
	this->id = this->conceptCount;

}

string Concept::toCompString(){
	set<int>::iterator iter;
	string s="\t---------------------\n";
	s = s + "\tConcept ID:" + intToString(this->id);
	s = s + "\n\t----------------------\n";
	s = s + "\tObject ID:" + intToString(this->objectId);
	s = s + "\n\tAttribute IDs:\n\t";
	for (iter = this->attributeIds.begin() ; iter != this->attributeIds.end() ; iter++){
		//printf("%d\n", );
		s = s + intToString(*iter) + ", ";
	}
	s = s + "\n\t-----------end---------\n";
	return s;
}


string Concept::toString(){
	set<int>::iterator iter;
	string s=intToString(this->id) + " : ";
	s = s + setSummary(this->attributeIds,0);
	return s;
}

int Concept::getID(){
	return this->id;
}

int Concept::getObjectID(){
	return this->objectId;
}

set<int> Concept::getAttributeIDs(){
	return this->attributeIds;
}
//////////////////////////////////////////////////////////////////////////////
// Class VERTEX definitions
//////////////////////////////////////////////////////////////////////////////

int Vertex::vertexCount = 0;

Vertex::Vertex(){}

Vertex::Vertex(set<int> _objids,set<int> _attrids){
	//printf(">>> New Vertex is constructing\n");
	this->objIds = _objids;
	this->atrIds = _attrids;
	this->vertexCount++;
	this->id = this->vertexCount;
}

int Vertex::getID(){
	return this->id;
}

string Vertex::toString(int labelID){
	string tmps;
	set<int>::iterator ito;
	set<int> shrinkerSet;

	tmps = intToString(this->id);

	string s = tmps;
	s = s + ":<";
	tmps = setSummary(this->objIds,labelID);
	s = s + tmps + ">,(";
	tmps = setSummary(this->atrIds,labelID);
	//tmps = "TEMP";
	s = s + tmps + ")";

	return s;

}

string Vertex::toCompPrinting(){
	set<int>::iterator iter;
	string s="\t---------------------\n";
	s = s + "\tVertex ID:" + intToString(this->id);
	s = s + "\n\t----------------------\n";
	s = s + "\n\tObject IDs:\n\t";

	for (iter = this->objIds.begin() ; iter != this->objIds.end() ; iter++){
		s = s + intToString(*iter) + ", ";
	}

	s = s + "\n\tAttribute IDs:\n\t";
	for (iter = this->atrIds.begin() ; iter != this->atrIds.end() ; iter++){
		s = s + intToString(*iter) + ", ";
	}

	s = s + "\n\t-----------end---------\n";
	return s;
}

set<int> Vertex::getObjectIDs(){
	return this->objIds;
}

set<int> Vertex::getAttributeIDs(){
	set<int> t;
	t = this->atrIds;
	return t;
}

void Vertex::addToChilds(int childID, double cost){
	this->childIds[childID] = cost;
}
void Vertex::addToParents(int parID, double cost){
	this->parentIds[parID] = cost;
}

map<int,double> Vertex::getParents(){
	return this->parentIds;
}

map<int,double> Vertex::getChilds(){
	return this->childIds;
}

void Vertex::removeFromChilds(int id){
	this->childIds.erase(id);
}

void Vertex::removeFromParents(int id){
	this->parentIds.erase(id);
}

//////////////////////////////////////////////////////////////////////////////
// Class EDGE definitions
//////////////////////////////////////////////////////////////////////////////

int Edge::edgeCount = 0;

Edge::Edge(){};

Edge::Edge(int src,int dest,double cost){
	this->srcId = src;
	this->destId = dest;
	this->edgeCount++;
	this->id = this->edgeCount;
	this->cost = cost;
}

int Edge::getID(){
	return this->id;
}


string Edge::toString(){
	string s="\t---------------------\n";
	s = s + "\tEdge ID:" + intToString(this->id);
	s = s + "\n\t----------------------\n";
	s = s + "\n\tSource ID(vertex): " + intToString(this->srcId) ;
	s = s + "\n\tDestination ID(vertex): "+ intToString(this->destId) ;
	s = s + "\n\tCost: " + intToString(int(this->cost));
	s = s + "\n\t-----------end---------\n";
	return s;
}
int Edge::getSrcID(){
	return this->srcId;
}
int Edge::getDestID(){
	return this->destId;
}
double Edge::getCost(){
	return this->cost;
}
