/**
 * Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2019, All rights reserved
 * Code: entry.cpp
 * Description: decompression of ParLOT traces
 */
#include "nlr.h"

// how to retireve Loops from dtab
//    dtab[did] = info / L_lid
//    rdtab[L_lid] = did
//    rdtab[info] = did
//    ltab[lid] = LoopBody
//    rltab[LoopBody] = lid
int Entry::distinctElemenets = 0;
int Entry::ldataLen = 0;
int Entry::origLen = 0;
int Entry::fdataLen = 0;
int Entry::maxLoopBody = 0;
int Entry::maxLC = 0 ;

map<int,string> Entry::dtab ;
unordered_map<string,int> Entry::rdtab ;
map<int,string> Entry::ltab ;
unordered_map<string,int> Entry::rltab ;


Entry::Entry(){
  this->elements.clear();
  this->lc = 0;
}

Entry::~Entry(){}

void Entry::addElement(string el){
  int eli=addToDtab(el);
	this->elements.push_back(eli);
}

int Entry::addToDtab(string s){
  int ret;
  if (this->rdtab.count(s) > 0){
    // s exist
    ret = rdtab[s];
    return ret;
  }else{
    ret = this->dtab.size();
    dtab[ret] = s;
    rdtab[s] = ret;
    return ret;
  }
}

string Entry::retFromDtab(int i){
  return this->dtab[i];
}

int Entry::addToLtab(string s){
  int ret;
  if (this->rltab.count(s) > 0){
    // s exist
    ret = rltab[s];
    return ret;
  }else{
    ret = this->ltab.size();
    ltab[ret] = s;
    rltab[s] = ret;
    return ret;
  }
}

string Entry::retFromLtab(int i){
  return this->ltab[i];
}



void Entry::setLC(int lcc){
	this->lc = lcc;
}

int Entry::getLC(){
	return this->lc;
}

int Entry::getElementLen(){
  return this->elements.size();
}

bool Entry::operator==(const Entry& b){
  if (this->lc != b.lc){
    return false;
  }
  else if (this->elements.size() != b.elements.size()){
    return false;
  } else{
    for(unsigned int i = 0 ; i < this->elements.size() ; i++){
      if (this->elements[i] != b.elements[i]){
        return false;
      }
    }
    return true;
  }
}

void Entry::incLC(){
	this->lc++;
}

string Entry::toString(){
  printf("*****ALERT****\nLTAB:\n%s\nDTAB\n%s\n",ltabToString().c_str(),dtabToString().c_str());
	string s = "";
  string st = "";
	vector<int>::iterator vit;
	if (this->lc == 1){
    assert(this->getElementLen() == 1);
    printf("RERURN >> id[%d] : %s\n",this->elements[0],retFromDtab(this->elements[0]).c_str());
		return retFromDtab(this->elements[0]);
	} else{
		for (vit = (this->elements).begin();vit != (this->elements).end(); vit++){
      //tmp = retFromTable(*vit);
      st = st + retFromDtab(*vit);
      if (vit != (this->elements).end() - 1 ){
          st = st + " - ";
      }
		}
    printf("RERURN >> id[%d] : %s\n",this->elements[0],retFromDtab(this->elements[0]).c_str());
    return "L"+to_string(addToLtab(st))+"^"+to_string(this->lc);
    //return "L"+to_string(addToLtab(st))+"^"+to_string(this->lc);
	}
}

string Entry::toLString(){
  printf("*****ALERT****\nLTAB:\n%s\nDTAB\n%s\n",ltabToString().c_str(),dtabToString().c_str());
	string s = "";
  string st = "";
	vector<int>::iterator vit;
	if (this->lc == 1){
    assert(this->getElementLen() == 1);
    printf("RERURN >> id[%d] : %s\n",this->elements[0],retFromDtab(this->elements[0]).c_str());
		return retFromDtab(this->elements[0]);
	} else{
		for (vit = (this->elements).begin();vit != (this->elements).end(); vit++){
      //tmp = retFromTable(*vit);
      st = st + retFromDtab(*vit);
      if (vit != (this->elements).end() - 1 ){
          st = st + " - ";
      }
		}
    printf("RERURN >> id[%d] : %s\n",this->elements[0],retFromDtab(this->elements[0]).c_str());
    return "("+st+")^"+to_string(this->lc);
    //return "L"+to_string(addToLtab(st))+"^"+to_string(this->lc);
	}
}



vector<int> Entry::getElements(){
  return this->elements;
}


string Entry::ltabToString(){
  string s = "";
  map<int,string>::iterator mit;
  for(mit = (this->ltab).begin() ; mit != (this->ltab).end() ; mit++){
    s = s + to_string(mit->first) + ':' + mit->second + "\n";
  }
  return s;
}


string Entry::dtabToString(){
  string s = "";
  map<int,string>::iterator mit;
  for(mit = (this->dtab).begin() ; mit != (this->dtab).end() ; mit++){
    s = s + to_string(mit->first) + ':' + mit->second + "\n";
  }
  return s;
}

void Entry::loadDtab(string s){
  vector<string> sp = splitString(s,':');
  this->dtab[stoi(sp[0])] = sp[1];
  this->rdtab[sp[1]]= stoi(sp[0]);
}

void Entry::loadLtab(string s){
  vector<string> sp = splitString(s,':');
  this->ltab[stoi(sp[0])] = sp[1];
  this->rltab[sp[1]]= stoi(sp[0]);
}


void Entry::setOrigLen(int l){
  this->origLen = l;
}
void Entry::setFdataLen(int l){
  this->fdataLen = l;
}
void Entry::setLdataLen(int l){
  this->ldataLen = l;
}
int Entry::getOrigLen(){
  return this->origLen;
}
int Entry::getFdataLen(){
  return this->fdataLen;
}
int Entry::getLdataLen(){
  return this->ldataLen;
}
void Entry::setDistinctElements(int d){
  this->distinctElemenets = d;
}
int Entry::getDistinctElements(){
  return this->distinctElemenets;
}

string Entry::statToString(){
  string s="";
  s = s + "distElements:"+to_string(this->distinctElemenets) + "\n";
  s = s + "origLen:"+to_string(this->origLen) + "\n";
  s = s + "fdataLen:"+to_string(this->fdataLen) + "\n";
  s = s + "ldataLen:"+to_string(this->ldataLen) + "\n";
  return s;
}

int Entry::numOfLoops(){
  return this->ltab.size();
}

void Entry::setMaxLoopBody(int lb){
  if (lb > this->maxLoopBody){
    this->maxLoopBody = lb;
  }
}

int Entry::getMaxLoopBody(){
  return this->maxLoopBody;
}

void Entry::setMaxLC(int lc){
  if (lc > this->maxLC){
    this->maxLC = lc;
  }
}

int Entry::getMaxLC(){
  return this->maxLC;
}
