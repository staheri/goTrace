#include "lat_lat.h"
#include "lat_atr.h"
//#include "myUtil.h"
#include <string>

using namespace std;

//////////////////////////////////////////////////////////////////////////////
// Class LATTICE definitions
//////////////////////////////////////////////////////////////////////////////

map<int,Concept> Lattice::CV; /*!< Holds all the concepts of current lattice */
map<int,Vertex> Lattice::VV; /*!< Holds all the verices of current lattice */
map<int,Edge> Lattice::EV; /*!< Holds all the edges of current lattice */

int Lattice::supID; /*!< The Concept ID of sup*/
int Lattice::infID; /*!< The Concept ID of inf*/
int Lattice::maxAttribute; /*!< Maximum number of attributes */

Lattice::Lattice(){}

Lattice::Lattice(string _name){
	this->name = _name;
	this->constDot = "";
	this->infID = 0;
	this->supID = 0 ;
	this->maxAttribute = 0;
}


string Lattice::toString(){
	string s;

	return s;
}
string Lattice::toDot(string label, int labelID){
	map<int,Edge>::iterator edgeIter;

	string s = "digraph { \n\tlabel = \""+label+"\"\n\t";
	if (this->EV.size() != 0 ){
		for (edgeIter = this->EV.begin() ; edgeIter != this->EV.end() ; edgeIter++){
			s = s + "\""+this->VV[edgeIter->second.getSrcID()].toString(labelID);
			s = s + "\" -> \"" ;
			s = s + this->VV[edgeIter->second.getDestID()].toString(labelID);
			s = s + "\" ;\n\t";
		}
	} else if (this->VV.size() == 1 ){
		s = s + "\""+this->VV[1].toString(labelID);
		s = s + "\" -> \"" ;
		s = s + this->VV[1].toString(labelID);
		s = s + "\" ;\n\t";
	} else {
		s = s + "\"NULL\"";
		s = s + " -> " ;
		s = s + "\"NULL\"";
		s = s + " ;\n\t";
	}
	s = s + "\n}\n";

	return s;
}

string Lattice::toLatMat(){ // return s with |V| line, each line show index of childs, "-" means no child (inf)
	map<int,Edge>::iterator eit;
	//iterate over edges
	//latmat[s][d] = 1
	vector<vector<int>> latmat(this->VV.size(), vector<int>(this->VV.size()));
	for (eit = EV.begin(); eit != EV.end(); eit++){
		latmat[eit->second.getSrcID() - 1][eit->second.getDestID() - 1] = 1;
	}
	string ss = "";
	bool flg ;
	for (unsigned int i=0;i<latmat.size();i++){
		flg = false;
		for (unsigned int j=0;j<latmat[i].size();j++){
			if (latmat[i][j] == 1){
				ss = ss + intToString(j) + " ";
				flg = true;
			}
		}
		if (flg){
			ss = ss + "\n";
		}
		else{
			ss = ss + "-\n";
		}

	}
	//printf("\n%s",ss.c_str());
	return ss;
}
string Lattice::toDotConst(){
	string s = "digraph G{\n\tcompound=true;\n";
	s = s + this->constDot + "\n}\n";
	return s;
}

/*string Lattice::toConceptString(){
	string s = ""
	map<int,Concept>::iterator conceptIter;
	for (conceptIter = this->CV.begin() ; conceptIter != this->CV.end() ; conceptIter++){
		s = s + conceptIter->
	}
} */

string Lattice::toDotEdges(int cid, int labelID){
	map<int,Edge>::iterator edgeIter;

	string s = "\t\t";
	if (this->EV.size() != 0 ){
		for (edgeIter = this->EV.begin() ; edgeIter != this->EV.end() ; edgeIter++){
			s = s + "\"("+intToString(cid)+")"+this->VV[edgeIter->second.getSrcID()].toString(labelID) + "\"";
			s = s + " -> " ;
			s = s + "\"("+intToString(cid)+")"+this->VV[edgeIter->second.getDestID()].toString(labelID) + "\"";
			s = s + " ;\n\t\t";
		}
		s = s + "\"("+intToString(cid)+")"+this->VV[this->supID].toString(labelID) + "\" -> \"("+this->CV[cid].toString()+")\" [style=dotted];\n\t\t";
	} else if (this->VV.size() == 1 ){
		s = s + "\"("+intToString(cid)+")"+this->VV[1].toString(labelID)+ "\"";
		s = s + " -> " ;
		s = s + "\"("+intToString(cid)+")"+this->VV[1].toString(labelID)+ "\"";
		s = s + " ;\n\t\t";
		s = s + "\"("+intToString(cid)+")"+this->VV[this->supID].toString(labelID) + "\" -> \"("+this->CV[cid].toString()+")\" [style=dotted];\n\t\t";
	} else {
		s = s + "\"NULL\"";
		s = s + " -> " ;
		s = s + "\"NULL\"";
		s = s + " ;\n\t\t";
		s = s + "\"("+intToString(cid)+")"+this->VV[this->supID].toString(labelID) + "\" -> \"("+this->CV[cid].toString()+")\" [style=dotted];\n\t\t";
	}

	s = s + "\n";

	return s;
}


string Lattice::toContextBitmax(){
	map<int,Concept>::iterator cit;
	set<int>::iterator sit;
	set<int> tmpAtr;
	unsigned int i,j;
	vector<vector<int>> bitmax(this->CV.size(), vector<int>(this->maxAttribute));
	for (cit = CV.begin(); cit != CV.end(); cit++){
		//printf("%s\n",cit->second.toString().c_str());
		tmpAtr = cit->second.getAttributeIDs();
		for (sit = tmpAtr.begin() ; sit != tmpAtr.end() ; sit++){
			i = cit->first;
			j = *sit;
			bitmax[i-1][j-1] = 1;
		}
	}
	string s ="";
	for (i=0;i<bitmax.size();i++){
		for (j=0;j<bitmax[i].size();j++){
			s = s + intToString(bitmax[i][j]);
		}
		s = s + "\n";
	}
	//printf("%s",s.c_str());
	return s;

}

void Lattice::addSubgraph(string subg,string label){
	string s = "";
	s = s + "\tsubgraph \""+label[0]+"\" {\n";
	s = s + "\t\tlabel = \""+label+"\";\n";
	s = s + subg+"\n\t}\n";
	this->constDot = this->constDot + s;
}


void Lattice::printLatticeComponents(){
	map<int,Concept>::iterator cit;
	map<int,Vertex>::iterator vit;
	map<int,Edge>::iterator eit;

	// print set of concepts
	printf("All Concepts of the lattice...\n");
	for (cit = CV.begin(); cit != CV.end(); cit++){
		printf("%s\n",cit->second.toString().c_str());
	}
	// print set of vertices
	printf("All Vertices of the lattice...\n");
	for (vit = VV.begin(); vit != VV.end(); vit++){
		printf("%s\n",vit->second.toCompPrinting().c_str());
	}
	// print set of edges
	printf("All Edges of the lattice...\n");
	for (eit = EV.begin(); eit != EV.end(); eit++){
		printf("%s\n",eit->second.toString().c_str());
	}
}




void Lattice::addEdge(Edge e){
	//printf("Adding Edge...\n");
	//printf("Edge to Add:\n%s\n",e.toString().c_str());
	this->EV[e.getID()] = e;

	// Set SRC childs
	this->VV[e.getSrcID()].addToChilds(e.getDestID(),e.getCost());
	// Set Dest Parents
	this->VV[e.getDestID()].addToParents(e.getSrcID(),e.getCost());
}



void Lattice::deleteEdge(int src,int dest){
	//printf("Deleting Edge...\n");
	int id2del;

	// find which edge to delete
	map<int,Edge>::iterator iter;
	for (iter = this->EV.begin() ; iter != this->EV.end() ; iter++){
		if (iter->second.getSrcID() == src && iter->second.getDestID() == dest ){
			// Edge ID to delete:
			id2del = iter->second.getID();
			break;
		}
	}

	// Remove from Edge Vector
	assert(iter != this->EV.end());
	this->EV.erase(id2del);

	// Fix SRC and DEST childs and parents
	this->VV[src].removeFromChilds(dest);
	this->VV[dest].removeFromParents(src);
}




void Lattice::addVertex(Vertex v){
//	printf("Adding Vertex...\n");
	//printf("Vertex to Add:\n%s\n",v.toString(0).c_str());
	this->VV[v.getID()] = v;
}





void Lattice::addConcept(Concept& c){

//	printf("\t\tAdding Concept...\n");
	//printf("\t\tConcept to Add:\n\t\t%s\n",c.toString().c_str());


	// Needed vars
	typedef vector<Vertex*> Vervec;
	vector<Vertex*>::iterator vervecIter;

	map<int,Vervec> buk;
	map<int,Vervec> bbuk;
	map<int,Vervec>::iterator bukIter;


	int sizeAttr;
	bool isGenerator;
	bool isParent;
	vector<int> intersection;
	vector<int> tmp1,tmp2;

	set<int> tempSet1,tempSet2;

	set<int> intersect;

	bool test;

	//Adding Concept c to CV
	//printf("\t\t**************************************************************************** \n");
//	printf("\t\tAdding Concept c to Lattic(this)->CV \n");
	this->CV[c.getID()]=c;
	//printf("\t\tSet Attributes of new Concept \n");
	set<int> ca = c.getAttributeIDs() ; // Attributes of Concept c
	int co = c.getObjectID();

	//if (this->supID == 0){ // Adjusting SUP. Check if anything assigned to it before
	if (this->getSup() == 0){ // Adjusting SUP. Check if anything assigned to it before
		//printf("\t\t > Replace SUP with New Vertex\n");
		set<int> tempV;
		tempV.insert(c.getObjectID());
		Vertex newVertex = Vertex(tempV,c.getAttributeIDs());
		//printf("\t\t > Vertex Added:\n%s\n",newVertex.toString(0).c_str());
		addVertex(newVertex);
		//printf("\t\t > SUP ID NOW IS : %d \n",newVertex.getID());
		//this->supID = newVertex.getID();
		this->setSup(newVertex.getID());
	}else{
		// To check if NOT concept attributes (c.attributes()) are subset of of SUP(vertex) attributes
		//set<int> sa =  this->VV[this->supID].getAttributeIDs();        // Attributes of Sup (vertex)
		set<int> sa =  this->VV[this->getSup()].getAttributeIDs();        // Attributes of Sup (vertex)
		//printf("\t\t >> To check if NOT concept attributes (c.attributes()) are subset of of SUP(vertex) attributes \n");
		if (!std::includes(sa.begin(),sa.end(),ca.begin(),ca.end())) {

			//to check if value of the map changes by my function
			/*
			printf("\tBefore ins2aids");
			printf("\tSize of attrs: %d\n",this->VV[this->supID].getAttributeIDs().size());
			set<int> tmp = {1000};
			this->VV[this->supID].ins2aids(tmp);
			printf("\tAfter ins2aids");
			printf("\tSize of attrs: %d\n",this->VV[this->supID].getAttributeIDs().size());
			assert(this->supID != 0);
			*/

			// Check if SUP.objects = NULL (empty)
			//printf("\t\t >>> Check if SUP.objects = NULL (empty)\n");
			if (this->VV[this->getSup()].getObjectIDs().size() == 0 ){
				// Union SUP.attributeIDS and concept
				//printf("\t\t >>>> SUP.objects = NULL (empty)\n");
			//	printf("\t\t >>>> Union SUP.attributeIDS and concept\n");
				this->VV[this->getSup()].ins2aids(ca);
			} else{
				// add new pair H {becomes sup(G*)}: (empty,X'(sup(G)) UNION f(x*):
				// New Vertex (H): ObjIDS: EMPTY, AttrIDS: Union SUP.attributeIDS and concept
				// New Edge: from current  SUP to H(new vertex)
				// Update SUP: to new Vertex

				//printf("\t\t >>>> NOT NULL SUP.objects\n");
				//printf("\t\t >>>> New Vertex, New Edge, Update SUP\n");

				// New Vertex
				Vertex Hs;

				// First Argument of New Vertex
				set<int> emptySet;
				emptySet.clear();

				//Second Argument of New Vertex
				set<int> newAtrSet;
				newAtrSet = this->VV[this->getSup()].getAttributeIDs();
				newAtrSet.insert(ca.begin(),ca.end());

				Hs = Vertex(emptySet,newAtrSet);
				addVertex(Hs);

				// New Edge
				Edge E;
				int srcID = this->VV[this->getSup()].getID();
				int destID = this->VV[Hs.getID()].getID();;
				double cost = 0;

				E = Edge(srcID,destID,cost);
				addEdge(E);

				this->setSup(Hs.getID());
			}
			//printf("\t\t >>> NOT/After concept attributes subset of of SUP(vertex) attributes \n");
		}
		// Classifying Vertices into buckets based on their size of attributes
		buk.clear();
		bbuk.clear();
		//printf("\t\t >> Classifying ..\n");
		map<int,Vertex>::iterator itm;
		//printf("\t\t >>> Touching Lat->Vertices ..\n");
		for (itm = this->VV.begin(); itm != this->VV.end(); ++itm) {
			// itm->first // key(id)
			// itm->second // val(Vertex)
			//printf("\t\t >>> Checking Vertex ...%s\n",itm->second.toString(0).c_str() );
			sizeAttr = (itm->second).getAttributeIDs().size();
			buk[sizeAttr].push_back(&(itm->second));
		}

		for (bukIter = buk.begin(); bukIter != buk.end() ; bukIter++){
			for (vervecIter = bukIter->second.begin(); vervecIter != bukIter->second.end() ; vervecIter++ ){
				//printf("\t\t >>>> BukID: %d \n\t Object Within it:\n%s\n",bukIter->first,(*vervecIter)->toString(0).c_str() );
			}
		}
		// end of classification

		//printf("\t\t >>> End of classification\n");
		//printf("\t\t >>> Iterate Over Buk\n");
		for (bukIter = buk.begin(); bukIter != buk.end() ; bukIter++){
		  //printf("\t\t >>>> Now Buk[%d]\n",bukIter->first);
			for (vervecIter = bukIter->second.begin(); vervecIter != bukIter->second.end() ; vervecIter++ ){

				Vertex* H = *vervecIter;
				set<int>::iterator k3;
				// if H.attributes subset of ca (c.attributes) -->> modified pair
				//printf("\t\t >>>>> For Each Pair H in Buk[%d] \n %s \n",bukIter->first,H->toString().c_str());

				/* FOR TESTING INCLUDES() MECHANISM
				printf("\n\n\n*****************\n" );
				printf("new Concept Attributes: \n");
				for(k3 = ca.begin();k3 != ca.end(); k3++){
					printf("ca : %d\n",*k3);
				}
				printf("H Attributes: \n");
				for(k3 = H->getAttributeIDs().begin();k3 != H->getAttributeIDs().end(); k3++){
					printf("H : %d\n",*k3);
				}
				printf("\t\t *****************\n\n\n\n" );
				printf("\t\t >>>>> Check if H->attributes subset of C where C is:\n %s \n",c.toString().c_str());
				bool test = false;
				printf("\t\t >>>>> init bool = %d\n", test);
				*/
				tempSet1.clear();
				tempSet1 = ca;
				tempSet2.clear();
				tempSet2 = H->getAttributeIDs();
				test = std::includes(tempSet1.begin(),tempSet1.end(),tempSet2.begin(),tempSet2.end());

				//printf("\t\t >>>>> Modified pair??? = %d\n", test);

				if (test) {
					//printf("\t\t >>>>>>> H.attributes subset of ca (c.attributes) -->> modified pair\n");
					// WE NEED TO MODIFY VERTEX
					// BE CAREFUL ABOUT POINTERS

					//add c.obj to H.obj

					set<int> ttt;
					ttt.insert(co);
					H->ins2oids(ttt);

					//add H to bbuk[i]
					bbuk[bukIter->first].push_back(H);
				}
				// if H.attributes = c.attributes -->> Exit Algorithm
				if (H->getAttributeIDs() == ca){
					//printf("\t\t >>>>>>> H.attributes == c.attributes \n\t\t >>>>>>>>>>>>>> BREAK2\n");
					return;
				}else{ // OLD PAIR
					//Intersection
					//printf("\t\t >>>>>>> OLD PAIR  \n");

					tmp1.clear();

					set<int>::iterator fiter;
					for (fiter = ca.begin();fiter != ca.end(); fiter++ ){
						tmp1.push_back(*fiter);
					}

					tmp2.clear();
					set<int> kireKhar = H->getAttributeIDs();
					for (fiter = kireKhar.begin();fiter != kireKhar.end(); fiter++ ){
						tmp2.push_back(*fiter);
					}
					intersection.clear();
					//printf("\t\t >>>>>>> INTERSECTION:  \n");
					set_intersection(tmp1.begin(), tmp1.end(), tmp2.begin(), tmp2.end(), std::back_inserter(intersection));

					// Print Intersection
					//vector<int>::iterator pit;
					//for (pit = intersection.begin(); pit != intersection.end() ; pit++ ){
					//	printf ("\t\t >>>>>>>> %d \n",*pit);
					//}

				}

				// Convert Vector-Intersection to Set-intersect
				//printf("\t\t >>>>> Convert Intersection...\n");

				intersect.clear();
				vector<int>::iterator vecIter;
				for (vecIter = intersection.begin() ; vecIter != intersection.end() ; vecIter++){
					intersect.insert(*vecIter);
				}
				//printf("\t\t >>>>> Block Check (if not exist bbuk[]) IS GENERATOR?\n");

				// BLOCK CHECK
				//if not exist a member (H') of bbuk[(size of intersect)] such that
				// H'.attributes = intersection -->> H is generator

				map<int,Vervec>::iterator tmpBukIter;

				vector<Vertex*>::iterator tmpVervecIter;

				tmpBukIter = bbuk.find(intersection.size());

				isGenerator = true;
				if (tmpBukIter != bbuk.end()) {
					//printf("\t\t >>>>>> bbuk[size(intersec)] found \n");
					for (tmpVervecIter = tmpBukIter->second.begin(); tmpVervecIter != tmpBukIter->second.end(); tmpVervecIter++) {
						if ((*tmpVervecIter)->getAttributeIDs() == intersect) {
							//printf("\t\t >>>>>>> bbuk[size(intersec)] found NOT GENERATOR\n");
							isGenerator = false;
						}
					}
				}
				//printf("\t\t >>>>> END BLOCK CHECK\n");
				// END BLOCK CHECK

				if (isGenerator){
					// New Vertex (Hn): objids: H.objids union c.objid | attrids: intersection
					// Add Hn to bbuk[||int||]
					// New Edge: (Hn) -> H
					//printf("\t\t >>>>>> IS Generator\n");
					Vertex Hn;

					//first Argument
					set<int> newObjSet = H->getObjectIDs();
					newObjSet.insert(co);

					//Second Argument of New Vertex : INTERSECT
					Hn = Vertex(newObjSet,intersect);
					addVertex(Hn);


					// Adding HN to BBUK[||int||]
					bbuk[intersect.size()].push_back(&(this->VV[Hn.getID()]));

					Edge E;
					int srcID = Hn.getID();
					int destID = this->VV[H->getID()].getID();;
					double cost = 0;

					E = Edge(srcID,destID,cost);
					addEdge(E);


					//Modify Edges
					//printf("\t\t >>>>>> Modifying Edges\n");
					// for j: 0 - ||int||-1
						//for each Ha belongs to bbuk[j]
							//if Ha.attributes subset of intersect
								//Ha is potential parent
								// isParent = true
								// for each Hd child of Ha
									//if Hd.attributes subset of intersect
										//isParent = false;
										// exit for
								// if isParent
									//if Ha is parent of H
										// delete edge between Ha and H
									//add edge ha > hn
					//printf("\t\t >>>>>> for j: 0 - ||int||-1\n");
					for (unsigned int j = 0 ; j < intersect.size() ; j++ ){
						map<int,Vervec>::iterator bviter;
						bviter = bbuk.find(j);
						//printf("\t\t >>>>>>> Checking BBUK\n");
						if (bviter != bbuk.end()){
							//printf("\t\t >>>>>>>> J : %d\n",j);
							vector<Vertex*>::iterator viter;
							for (viter = bviter->second.begin() ; viter != bviter->second.end(); viter++){
								Vertex* Ha = *viter;
								//printf("\t\t >>>>>>>>> for each Ha belongs to bbuk[j]\n");

								//Set Include Block
								tempSet1.clear();
								tempSet1 = intersect;
								tempSet2.clear();
								tempSet2 = Ha->getAttributeIDs();
								test = std::includes(tempSet1.begin(),tempSet1.end(),tempSet2.begin(),tempSet2.end());


								if (test) {
									//printf("\t\t >>>>>>>>>> Ha Potential Parent\n");
									isParent = true;
									//for each child of HA
									map<int,double> childMap = Ha->getChilds();
									map<int,double>::iterator miter;
									//printf("\t\t >>>>>>>>>> For Each child of Ha (ifParent)\n");
									for (miter = childMap.begin(); miter != childMap.end() ; miter++){
										int childID = miter->first;
										Vertex Hd = this->VV[childID];

										//Set Include Block
										tempSet1.clear();
										tempSet1 = intersect;
										tempSet2.clear();
										tempSet2 = Hd.getAttributeIDs();
										test = std::includes(tempSet1.begin(),tempSet1.end(),tempSet2.begin(),tempSet2.end());

										if (test) {
											//printf("\t\t >>>>>>>>>>> Not Parent\n");
											isParent = false;
											break;
										}
									}
									if (isParent){
										//printf("\t\t >>>>>>>>>>> Parent\n");
										// if Ha is a parent of H
										map<int,double> childMap2 = Ha->getChilds();
										map<int,double>::iterator miter2;
										miter2 = childMap2.find(H->getID());
									//	printf("\t\t >>>>>>>>>>> if Ha is a parent of H\n");
										if ( miter2 != childMap2.end()){
											//printf("\t\t >>>>>>>>>>>> delete edge between HA -> H\n");
											//delete edge between HA -> H
											deleteEdge(Ha->getID(),H->getID());
										}
										//printf("\t\t >>>>>>>>>>> ADD EDGE HA->HN \n");
										Edge E2 = Edge(Ha->getID(),Hn.getID(),0);
										//add edge ha - > Hn
										addEdge(E2);
									}
								}
							}
						}
						//printf("\t\t >>>>>>> Done with BBUK\n");
					}
					if (intersect == ca ){
						//printf("\t\t >>>>> BREAK\n");
						return;
					}
				}
			}
		}
	}
}

void Lattice::setSup(int sid){
	this->supID = sid;
}

void Lattice::setInf(int iid){
	this->infID = iid;
}

void Lattice::setMaxAttribute(int val){
	this->maxAttribute = val;
}

int Lattice::getSup(){
	return this->supID;
}
int Lattice::getInf(){
	return this->infID;
}

int Lattice::getMaxAttribute(){
	return this->maxAttribute;
}
