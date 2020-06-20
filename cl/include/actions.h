/** actions.h
 * mid-level interface for communicting the front-end and back-end for trace and CL operations
 * Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2019, All rights reserved
 */

#ifndef ACTIONS_H
#define ACTIONS_H


#include <vector>
#include <set>
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

#include "prep.h"
#include "util.h"
#include "entry.h"
#include "nlr.h"
#include "lat_lat.h"
#include "lat_atr.h"
#include "lat_vec.h"

#define version 0.1
#define Q(x) #x
#define QUOTE(x) Q(x)
#define FILTBITSIZE 16


using namespace std;

/**
 * Called from main to generate Concept Lattice.
 * It first Preprocess data (decompress, filter, detect loops, extract attributes, etc.)
 * Then it creates the lattice from generated data
 */
void genGeneralCL(string _inpath,string _filtbit, int _atrMode, int _atrFreq, int _atrOption, int k);


#endif
