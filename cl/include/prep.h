/**
 * Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2018, All rights reserved
 * Code: prep.h
 * Description: Preprocess trace files for Concept Lattice Generation.
 *              decompression
 *              Filter
 *              NLR (Nested Loop Recognition)
 *              Attribute Extraction
 */
#ifndef PREP_H
#define PREP_H


#include <vector>
#include <set>
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
#include <sys/stat.h>
#include <cstdlib>
#include <cstdio>
#include <cassert>
#include <cstring>
#include <utility>
#include <cxxabi.h>
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
#include <iomanip>
#include <omp.h>

#include "util.h"
#include "entry.h"
#include "nlr.h"
#include "lat_lat.h"
#include "lat_atr.h"
#include "lat_vec.h"

using namespace std;

#define INFSIZE 65536


//GLOBAL VARIABLES

//class Trace;
//template <class T> class Attribute;
//class Entry;


extern long hnum;

typedef unsigned char uint1;
typedef unsigned short uint2;
typedef unsigned int uint4;
typedef unsigned long long uint8;


/**
 * Read trace files to decompress the byte codes (By Martin Burtscher and Sindhu Devale)
 */
uint2* readFile(const char name[], uint8& length);

/**
 * Read info file to decompress traces (By Martin Burtscher and Sindhu Devale)
 */
string* readInfo(const char name[]) ;

/**
 * Filter decompressed traces based on filtbit-string
 * Input: single ParLOT trace-data and info, filtbit, vector of custome regular expression (if any)
 * Output: Vector of function ids
 */
 vector<uint2> filterData(const uint2 data[], const uint8 length, string info[], string mode,vector<regex> vreg);
 vector<uint2> filterData2(const uint2 data[], const uint8 length,int start, int end, string info[], string mode,vector<regex> vreg);

 /**
  * Preprocess trace files
  * Input: Path to ParLOT traces
  * Output: A hash map of traces and preprocessed data
  * It first checks to see if the prep files already generated.
  * if (prep exist)
  * 		load all data (decompressed, filtered, nlr, stats, etc.) from files to the hash map
  * else
  *		Decompress Trace files
  *		Filter
  *		Detect Loops
  *		Write to files
  *		return hash map
  */
unordered_map<string,string> preprocess(string _inpath,string _filtbit,int k);

/**
 * Extract attributes from vector of trace entries
 * Input: Vector of trace entry, attribute mode, attribute freq and attribute options
 * Output: Vector of strings (attributes to be injected to concept lattice)
 */
set<string> extractAttributes(vector<Entry> vec, int mode, int freq, int option);



string filtbitTranslator(string filtbit, int k);

string clNameTranslator(int m, int f, int op);


unordered_map<string,string> goPreprocess(string _inpath, int k);
#endif
