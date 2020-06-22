/**
 * Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2018, All rights reserved
 * Code: gen_convert.h
 * Description: Part of traceToText decompression application
 */
#ifndef UTIL_H
#define UTIL_H

#include <cstdlib>
#include <cstdio>


#include <vector>
#include <set>
#include <unordered_map>
#include <string>
#include <string.h>
#include <stdio.h>
#include <time.h>
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
#include <map>
#include <utility>

#include "lat_lat.h"
#include "lat_atr.h"
#include "lat_vec.h"

using namespace std;

//GLOBAL VARIABLES

class Concept;
class Trace;
template <class T> class Attribute;
class Entry;

/**
 * Takes a string and returns a vector of strings which is the splitter based on char 'sep' as seperator
 */
vector<string> splitString(string text, char sep);

/**
 * Pairing each trace file with its corresponding info file.
 */
pair<string,string> infoPair(const string info,const string path);

/**
 * Returns a vector of all files within path with ext extension
 */
vector<string> listOfFiles(const string path,const char* ext);

/**
 * Returns a vector of all folders within path
 */
vector<string> listOfFolders(const string path);

/**
 * Returns a vector of all ptrace files within path
 */
vector<string> listOfTraceFiles(const string _path);

/**
 * Prints help message
 */
void printUsage(void);

/**
 * Summarize long range of ints (1,2,3,...,n) to (1-n)
 */
string setSummary(set<int>,int flag);

/**
 * Converts integer to string
 */
string intToString(int i);

/**
 * Check if a directory exists
 */
bool isDir(string path);



#endif
