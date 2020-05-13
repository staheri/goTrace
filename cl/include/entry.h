/**
 * Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2019, All rights reserved
 * Code: entry.h
 * Description: decompression of ParLOT traces
 */
#ifndef ENTRY_H
#define ENTRY_H


#include <vector>
#include <map>
#include <unordered_map>
#include <string>
#include <string.h>
#include <stdio.h>
#include <algorithm>
#include <sstream>
#include <iterator>
#include <dirent.h>
#include <unistd.h>
#include <sstream>
#include <iostream>
#include <fstream>
#include <iterator>
#include <regex>
#include <assert.h>
#include "util.h"


using namespace std;

class Entry{
  vector<int> elements;
  int lc;
  static int distinctElemenets;
  static int origLen;
  static int ldataLen;
  static int fdataLen;
  static int maxLC;
  static int maxLoopBody;
  static map<int,string> dtab;
  static unordered_map<string,int> rdtab;
  static map<int,string> ltab;
  static unordered_map<string,int> rltab;
public:
  //! Constructor
  Entry();

  //! Destructor
  ~Entry();

  //! Add Element
	/*!
	  Add element to the vector of current entry
	*/
  void addElement(string el);

  //! Set loop coint to lc
  void setLC(int lc);

  //! Increment loop count by one
  void incLC();

  //! Return string representation of the entry
  string toString();

  string toLString();

  //! Overloading == operator for entries
  bool operator==(const Entry& b);

  //! Return loop count of current entry
  int getLC();

  //! Return size of elements in current entry
  int getElementLen();

  //! add new entry (element) to ltab-rltab
  int addToDtab(string s);
  int addToLtab(string s);



  //! return ltab elements
  string retFromDtab(int i);
  string retFromLtab(int i);
  vector<int> getElements();

  string ltabToString();
  string dtabToString();
  void loadLtab(string s);
  void loadDtab(string s);

  void setOrigLen(int l);
  void setFdataLen(int l);
  void setLdataLen(int l);
  int getOrigLen();
  int getFdataLen();
  int getLdataLen();
  void setDistinctElements(int d);
  int getDistinctElements();
  string statToString();
  int numOfLoops();
  void setMaxLoopBody(int lb);
  int getMaxLoopBody();
  void setMaxLC(int lc);
  int getMaxLC();

  // Return the body of the loop in string format (if the entry is a loop)
  string LBtoString();
};


#endif
