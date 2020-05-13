/**
 * Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2018, All rights reserved
 * Code: prep.cpp
 * Description: Preprocess trace files for Concept Lattice Generation.
 *              decompression
 *              Filter
 *              NLR (Nested Loop Recognition)
 *              Attribute Extraction
 * (Decompression : Martin Burtscher, Texas State University, burtscher@cs.txstate.edu
 * Sindhu Devale, Texas State University)
 */
#include "prep.h"


//////////////////////////////////////////////////////////////////////////////
// traceReader(): Decompress, read info, read file, process
//////////////////////////////////////////////////////////////////////////////


static const uint2 maxstacksize = 4096;
static const uint1 hsizelg2 = 12;
static const uint2 hsize = 1 << hsizelg2;
static const uint1 psizelg2 = 16;
static const uint4 psize = 1 << psizelg2;
static const uint2 psizem1 = psize - 1;

uint2 hbuf[hsize];
uint2 pbuf[psize];

/**
 * Read info file to decompress traces (By Martin Burtscher and Sindhu Devale)
 */
string* readInfo(const char name[]){
  //printf("INSIDE READ INFO: %s\n",name);
	char img[256], fnc[256];
	string* info = new string[INFSIZE];
	FILE* f = fopen(name, "rt");  assert(f != NULL);
	uint2 i = 1;
	while (fscanf(f, "%255s\n", fnc) == 1) {
		if (fnc[0] == '+') {
			fscanf(f, "%255s\n", img);
			fscanf(f, "%255s\n", fnc);
		}
		img[255] = 0;
		fnc[255] = 0;

		char* p = fnc;
		while ((*p != 0) && (*p != '.')) p++;

		char* dname;
		int status;
		if (*p == '.') {
			*p = 0;
			dname = abi::__cxa_demangle(fnc, 0, 0, &status);
			*p = '.';
		} else {
			dname = abi::__cxa_demangle(fnc, 0, 0, &status);
		}
		if (dname != NULL) {
			info[i] = dname;
			if (*p == '.') {
				info[i] += p;
			}
			free(dname);
		} else {
			info[i] = fnc;
		}

		assert(i < INFSIZE);
		i++;
	}
	fclose(f);
	return info;
}

/**
 * Read trace files to decompress the byte codes (By Martin Burtscher and Sindhu Devale)
 */
uint2* readFile(const char name[], uint8& length){
	//printf("INSIDE READFILE NAME: %s\n",name);
	FILE* f = fopen(name, "rb");  assert(f != NULL);
	fseek(f, 0, SEEK_END);

	uint8 csize = ftell(f);
	//assert(csize > 0);
	if (csize <= 0){
		uint2* dbuf = (uint2*)malloc(1 * sizeof(uint2));
		length = 0;
		dbuf[0]=-1;
		return dbuf;
	}else{
		//printf("csize %.2f\n",csize);
		uint1* cbuf = new uint1[csize];
		fseek(f, 0, SEEK_SET);
		length = fread(cbuf, sizeof(uint1), csize, f);  assert(length == csize);
		fclose(f);

		uint8 bytes = 0;
		uint8 cpos = 0;
		while (cpos < csize) {
			uint1 bitpat = cbuf[cpos++];
			for (uint1 b = 0; b < 8; b++) {
				if (bitpat & (1 << b)) {
					cpos++;
				}
			}
			if (cpos == csize - 1) {
				bytes += cbuf[cpos++] * 2;
			} else {
				bytes += 8;
			}
		}
		uint1* bbuf = new uint1[bytes + 6];
		uint8 bpos = 0;
		cpos = 0;
		while (cpos < csize - 1) {
			uint1 bitpat = cbuf[cpos++];
			for (uint1 b = 0; b < 8; b++) {
				if (bitpat & (1 << b)) {
					bbuf[bpos++] = cbuf[cpos++];
				} else {
					bbuf[bpos++] = 0;
				}
			}
		}
		delete [] cbuf;

		uint2* dbuf = (uint2*)malloc(8192 * sizeof(uint2));  assert(dbuf != NULL);
		uint8 dcap = 8192;
		uint8 dpos = 0;

		bool iscount;
		uint2 ppos = 0;
		uint2 lpos = 0;
		uint2 hash = 0;
		uint8 tpos = 0;

		memset(hbuf, 0, hsize * sizeof(hbuf[0]));
		memset(pbuf, 0, psize * sizeof(pbuf[0]));

		uint8 words = bytes / 2;
		uint2* tbuf = (uint2*)bbuf;
		while (tpos < words) {
			lpos = hbuf[hash];
			iscount = (pbuf[(lpos - 3) & psizem1] == pbuf[(ppos - 3) & psizem1]) &&
			(pbuf[(lpos - 2) & psizem1] == pbuf[(ppos - 2) & psizem1]) &&
			(pbuf[(lpos - 1) & psizem1] == pbuf[(ppos - 1) & psizem1]);
			if (iscount) {
				uint2 count = tbuf[tpos++];
				for (uint2 i = 0; i < count; i++) {
					uint2 value = pbuf[lpos];
					if (dpos == dcap) {
						dcap *= 2;
						dbuf = (uint2*)realloc(dbuf, dcap * sizeof(uint2));  assert(dbuf != NULL);
					}
					dbuf[dpos++] = value;
					lpos = (lpos + 1) & psizem1;
					hbuf[hash] = ppos;
					hash = ((hash << (hsizelg2 / 3)) ^ value) % hsize;
					pbuf[ppos] = value;
					ppos = (ppos + 1) & psizem1;
				}
			}

			if (tpos < words) {
				uint2 value = tbuf[tpos++];
				if (dpos == dcap) {
					dcap *= 2;
					dbuf = (uint2*)realloc(dbuf, dcap * sizeof(uint2));  assert(dbuf != NULL);
				}
				dbuf[dpos++] = value;
				hbuf[hash] = ppos;
				hash = ((hash << (hsizelg2 / 3)) ^ value) % hsize;
				pbuf[ppos] = value;
				ppos = (ppos + 1) & psizem1;
			}
		}

		length = dpos;

		delete [] bbuf;
		return dbuf;
	}
}


/**
 * Filter decompressed traces based on filtbit-string
 * Input: single ParLOT trace-data and info, filtbit, vector of custome regular expression (if any)
 * Output: Vector of function ids
 */
vector<uint2> filterData(const uint2 data[], const uint8 length, string info[], string filtbit, vector<regex> vreg ){
	info[0] = "[ret]";
	vector<regex> allFilters;
	vector<regex>::iterator vreg_it;
	uint2 cur;
	int fret = filtbit[0] - '0';
	int fplt = filtbit[1] - '0';
	regex toAdd;
	bool includeEverything = false;
	allFilters.clear();
	if (filtbit[2] - '0' == 1){ //only @plt
		toAdd = "\\w*@plt$";
		allFilters.push_back(toAdd);
	}
	if (filtbit[3] - '0' == 1){ // memory
		toAdd = "\\w*mem\\w*";
		allFilters.push_back(toAdd);
	}
	if (filtbit[4] - '0' == 1){ // network
		toAdd = "\\w*network\\w*";
		allFilters.push_back(toAdd);
	}
	if (filtbit[5] - '0' == 1){ // poll
		toAdd = "\\w*poll\\w*";
		allFilters.push_back(toAdd);
	}
	if (filtbit[6] - '0' == 1){ // str
		toAdd = "\\w*str\\w*";
		allFilters.push_back(toAdd);
	}

	if (filtbit[7] - '0' == 1){ // MPI_
		toAdd = "^MPI_\\w*";
		allFilters.push_back(toAdd);
	}
	if (filtbit[8] - '0' == 1){ // MPI-related
		toAdd = "\\w*MPI\\w*";
		allFilters.push_back(toAdd);
	}
	if (filtbit[9] - '0' == 1){ // mpi_collectives
		toAdd = "^MPI_\\w*[R|r]educe\\w*";
		allFilters.push_back(toAdd);
		toAdd = "^MPI_\\w*[G|g]ather\\w*";
		allFilters.push_back(toAdd);
		toAdd = "^MPI_\\w*[b|B][c|C]ast\\w*";
		allFilters.push_back(toAdd);
		toAdd = "^MPI_\\w*[b|B]arrier\\w*";
		allFilters.push_back(toAdd);
		toAdd = "^MPI_\\w*[S|s]catter\\w*";
		allFilters.push_back(toAdd);
	}
	if (filtbit[10] - '0' == 1){ // mpi_send/rcv
		toAdd = "^MPI_Isend\\w*";
		allFilters.push_back(toAdd);
		toAdd = "^MPI_Recv\\w*";
		allFilters.push_back(toAdd);
		toAdd = "^MPI_Send\\w*";
		allFilters.push_back(toAdd);
	}
	if (filtbit[11] - '0' == 1){ // omp-critical
		toAdd = "\\w*critical\\w*";
		allFilters.push_back(toAdd);
	}
	if (filtbit[12] - '0' == 1){ // omp-mutex
		toAdd = "\\w*lock\\w*";
		allFilters.push_back(toAdd);
		toAdd = "\\w*unlock\\w*";
		allFilters.push_back(toAdd);
		toAdd = "\\w*mutex\\w*";
		allFilters.push_back(toAdd);
	}
	if (filtbit[13] - '0' == 1){ // omp-rest
		toAdd = "\\w*omp\\w*";
		allFilters.push_back(toAdd);
		toAdd = "\\w*OMP\\w*";
		allFilters.push_back(toAdd);
		toAdd = "\\w*pthread\\w*";
		allFilters.push_back(toAdd);
    toAdd = "\\w*lock\\w*";
		allFilters.push_back(toAdd);
		toAdd = "\\w*unlock\\w*";
		allFilters.push_back(toAdd);
		toAdd = "\\w*mutex\\w*";
		allFilters.push_back(toAdd);
    toAdd = "\\w*critical\\w*";
		allFilters.push_back(toAdd);
	}
	if (filtbit[14] - '0' == 1){ // custome
		for (vreg_it = vreg.begin() ; vreg_it != vreg.end() ; vreg_it++){
			allFilters.push_back(*vreg_it);
		}
	}
	if (filtbit[15] - '0' == 1){ // include Everything
		includeEverything = true;
	}

	string subject;
	smatch match;

	vector<uint2> ret;


	for (uint8 i = 0 ; i < length ; i++){
		if (i % 1000000 == 0){
			printf(">> Filtering %llu / %llu \n", i,length);
		}
		cur = data[i];
		if (cur == 0 && fret == 1){
			continue;
		}
		else{
			subject = info[cur];
			if (subject.compare(".plt") == 0 && fplt == 1){
				continue;
			}
			// iterate over all regex
			for (vreg_it = allFilters.begin() ; vreg_it != allFilters.end() ; vreg_it++){
				regex_search(subject,match,*vreg_it);
				if (match.size() >= 1){
					//printf("i: %d , %s\n",i,subject.c_str() );
					ret.push_back(cur);
					continue;
				}
			}
			//printf("F: %s\n",subject.c_str() );
			if (includeEverything){
				ret.push_back(cur);
			}else{

				continue;
			}
		}
	}
	return ret;
}

/**
 * Filter decompressed traces based on filtbit-string
 * Input: single ParLOT trace-data and info, filtbit, vector of custome regular expression (if any)
 * Output: Vector of function ids
 */
vector<uint2> filterData2(const uint2 data[], const uint8 length,int start, int end, string info[], string filtbit, vector<regex> vreg ){
	info[0] = "[ret]";
	vector<regex> allFilters;
	vector<regex>::iterator vreg_it;
	uint2 cur;
	int fret = filtbit[0] - '0';
	int fplt = filtbit[1] - '0';
	regex toAdd;
	bool includeEverything = false;
	allFilters.clear();
	if (filtbit[2] - '0' == 1){ //only @plt
		toAdd = "\\w*@plt$";
		allFilters.push_back(toAdd);
	}
	if (filtbit[3] - '0' == 1){ // memory
		toAdd = "\\w*mem\\w*";
		allFilters.push_back(toAdd);
	}
	if (filtbit[4] - '0' == 1){ // network
		toAdd = "\\w*network\\w*";
		allFilters.push_back(toAdd);
	}
	if (filtbit[5] - '0' == 1){ // poll
		toAdd = "\\w*poll\\w*";
		allFilters.push_back(toAdd);
	}
	if (filtbit[6] - '0' == 1){ // str
		toAdd = "\\w*str\\w*";
		allFilters.push_back(toAdd);
	}

	if (filtbit[7] - '0' == 1){ // MPI_
		toAdd = "^MPI_\\w*";
		allFilters.push_back(toAdd);
	}
	if (filtbit[8] - '0' == 1){ // MPI-related
		toAdd = "\\w*MPI\\w*";
		allFilters.push_back(toAdd);
	}
	if (filtbit[9] - '0' == 1){ // mpi_collectives
		toAdd = "^MPI_\\w*[R|r]educe\\w*";
		allFilters.push_back(toAdd);
		toAdd = "^MPI_\\w*[G|g]ather\\w*";
		allFilters.push_back(toAdd);
		toAdd = "^MPI_\\w*[b|B][c|C]ast\\w*";
		allFilters.push_back(toAdd);
		toAdd = "^MPI_\\w*[b|B]arrier\\w*";
		allFilters.push_back(toAdd);
		toAdd = "^MPI_\\w*[S|s]catter\\w*";
		allFilters.push_back(toAdd);
	}
	if (filtbit[10] - '0' == 1){ // mpi_send/rcv
		toAdd = "^MPI_Isend\\w*";
		allFilters.push_back(toAdd);
		toAdd = "^MPI_Recv\\w*";
		allFilters.push_back(toAdd);
		toAdd = "^MPI_Send\\w*";
		allFilters.push_back(toAdd);
	}
	if (filtbit[11] - '0' == 1){ // omp-critical
		toAdd = "\\w*critical\\w*";
		allFilters.push_back(toAdd);
	}
	if (filtbit[12] - '0' == 1){ // omp-mutex
		toAdd = "\\w*lock\\w*";
		allFilters.push_back(toAdd);
		toAdd = "\\w*unlock\\w*";
		allFilters.push_back(toAdd);
		toAdd = "\\w*mutex\\w*";
		allFilters.push_back(toAdd);
	}
	if (filtbit[13] - '0' == 1){ // omp-rest
		toAdd = "\\w*omp\\w*";
		allFilters.push_back(toAdd);
		toAdd = "\\w*OMP\\w*";
		allFilters.push_back(toAdd);
		toAdd = "\\w*pthread\\w*";
		allFilters.push_back(toAdd);
    toAdd = "\\w*lock\\w*";
		allFilters.push_back(toAdd);
		toAdd = "\\w*unlock\\w*";
		allFilters.push_back(toAdd);
		toAdd = "\\w*mutex\\w*";
		allFilters.push_back(toAdd);
    toAdd = "\\w*critical\\w*";
		allFilters.push_back(toAdd);
	}
	if (filtbit[14] - '0' == 1){ // custome
		for (vreg_it = vreg.begin() ; vreg_it != vreg.end() ; vreg_it++){
			allFilters.push_back(*vreg_it);
		}
	}
	if (filtbit[15] - '0' == 1){ // include Everything
		includeEverything = true;
	}

	string subject;
	smatch match;

	vector<uint2> ret;


	for (int i = start ; i <= end ; i++){
		// if (i % 1000000 == 0){
		// 	printf(">> Filtering %llu / %llu \n", i,length);
		// }
		cur = data[i];
		if (cur == 0 && fret == 1){
			continue;
		}
		else{
			subject = info[cur];
			if (subject.compare(".plt") == 0 && fplt == 1){
				continue;
			}
			// iterate over all regex
			for (vreg_it = allFilters.begin() ; vreg_it != allFilters.end() ; vreg_it++){
				regex_search(subject,match,*vreg_it);
				if (match.size() >= 1){
					//printf("i: %d , %s\n",i,subject.c_str() );
					ret.push_back(cur);
					continue;
				}
			}
			//printf("F: %s\n",subject.c_str() );
			if (includeEverything){
				ret.push_back(cur);
			}else{

				continue;
			}
		}
	}
	return ret;
}


/**
 * Preprocess trace files
 * Input: Path to a single ParLOT trace (PT), _filtbit, K (for NLR)
 * Output: A string path to the entries of PT (by figuring out the folder)
 * It first checks to see if the prep files already generated.
 * if (prep exist)
 * 		returns the path
 * else
 *		Decompress Trace files
 *		Filter
 *		Detect Loops
 *		Write to files
 *		Returns the path
 */
unordered_map<string,string> preprocess(string _inpath,string _filtbit, int k){
	unordered_map<string,string> ret;
	vector<string> lot; // list of all traces within _inpath/ptrace
	vector<string>::iterator vst_it;
	vector<string>::iterator vst_it2;

	string _trace; // Holds trace full path
	string _info; // Holds info full path
	vector<string> tmpVst; // Temporary vector string

	// time measurement
	clock_t t;

	vector<Entry> entries; // To store each trace entries

	// For Decompression
	string* info;
	uint2* data ;
	uint8 length;

	// Entry Initilization
	Entry tmp;

	string traceToken ;
	string infoToken ;

	vector<uint2> fdata; // To store filtered data
	vector<uint2>::iterator vit; //

	Entry ldata; // To store nlr data (loop-data)
	vector<Entry>::iterator vlit;
	set<string> distincts;

	string tmps;

  string folderName = filtbitTranslator(_filtbit,k);

	//making dir Outpath to store prep files there
	string outpath = _inpath + "/prep/"+folderName+"/";
  //printf(">>>Outpath : %s\n",outpath.c_str() );

	if (!isDir(outpath)){
		//create directory
		if(mkdir(outpath.c_str(),0777) == -1){
			perror("Error creating prep");
		}
	}

	// Up to here, the folder does exist. Now check if files do exist.
	// If they exist, read them and return
	// otherwise, do the preprocess, write to file and return


	// Get the list of all traces in _inpath/ptrace
	lot = listOfTraceFiles(_inpath+"/ptrace");

	//open log file
  ofstream flog(outpath + "log.txt");

	// preprocess all traces within _inpath/ptrace
	for(vst_it = lot.begin();vst_it != lot.end();vst_it++){
		_trace = _inpath+"/ptrace/" + *vst_it;
		traceToken = splitString(_trace, '/').back() ;

		// Write to log file
		flog << "+" << traceToken << endl;
		// Check if trace is already preprocessed?
		ifstream fi(outpath + traceToken + ".txt");
		//printf("Check if prep trace file exist \n");
		if(!fi.good()){
			// prep of current trace does not exist, so do prep, write to file
			_info = "";
			//printf(">> prep trace file does not exist \n");
			// find the info of current trace
			tmpVst = splitString(_trace,'.');
			//printf(">> 22 %d\n",tmpVst.size());
			for (vst_it2 = tmpVst.begin() ; vst_it2 != tmpVst.end() -1; vst_it2++){
				_info = _info + *vst_it2 + ".";
			}
			_info = _info + "info";

			infoToken = splitString(_info, '/').back() ;
			//printf(">> 44 %s\n",_info.c_str());
			// Decompression
			length = 0;
			t = clock();
	  	info = readInfo(_info.c_str());
	  	data = readFile(_trace.c_str(), length);
			t = clock() - t ;
			flog << "tl:" << length << endl;
			flog << "dt:" << t << "," << (((float)t)/CLOCKS_PER_SEC) << endl;
			// For custome filters
	    vector<regex> vreg;
			vreg.push_back((regex)"\\w*CPU_Exec\\w*");
			vreg.push_back((regex)"\\w*CPU_Init\\w*");
			vreg.push_back((regex)"\\w*CPU_Output\\w*");
			//vreg.push_back((regex)"^merge\\w*");
			//vreg.push_back((regex)"^qsort\\w*");
			//vreg.push_back((regex)"^copy\\w*");
			//vreg.push_back((regex)"^findPartner\\w*");





	    printf("\nFiltering data...length = %llu\n", length);
			if (length == 0){
				flog << "fl:0" << endl;
				flog << "ft:0,0.0" << endl;
				// trace file is empty
				tmp = Entry();
				distincts.clear();
				entries.clear();
		    tmp.addElement("EMPTY");
				distincts.insert("EMPTY");
		    tmp.setLC(1);
		    entries.push_back(tmp);
				printf("\nDetecting loops...length = 0\n");
			} else{
				// Filter
				fdata.clear();
				t = clock();
		    fdata = filterData(data,length,info,_filtbit,vreg);
				t = clock() - t;

				flog << "fl:" << fdata.size() << endl;
				flog << "ft:" << t << "," << (((float)t)/CLOCKS_PER_SEC) << endl;

				distincts.clear();
				entries.clear();
				// Creating entries from filtered data
		    for ( vit = fdata.begin() ; vit != fdata.end() ; vit++){
		      tmp = Entry();
		      tmp.addElement(info[(*vit)]);
					distincts.insert(info[(*vit)]);
		      tmp.setLC(1);
		      entries.push_back(tmp);
		    }
				printf("\nDetecting loops...length = %lu\n", fdata.size());
			}
			// ldata: Single entry object that holds info about dtab and ltab
			t = clock();
	    ldata = nlr(entries,k,outpath + traceToken + ".txt");
			t = clock() - t ;
			flog << "nl:" << ldata.getLdataLen() << endl;
			flog << "nt:" << t << "," << (((float)t)/CLOCKS_PER_SEC) << endl;
			//Set stats
			//ldata.setDistinctElements(distincts.size());
			//ldata.setOrigLen(length);
			//ldata.setFdataLen(fdata.size());
		}//end if
		ret[traceToken] = outpath + traceToken+".txt";
	}// end for

	flog.close();
	// write dtab to file
	ifstream fid(outpath + "dtab.txt");
	if(!fid.good()){
		ofstream fod(outpath + "dtab.txt");
		fod << ldata.dtabToString() ;
		fod.close();
	}
	ret["dtab"] = outpath + "dtab.txt";

	// write ltab to file
	ifstream fil(outpath + "ltab.txt");
	if(!fil.good()){
		ofstream fol(outpath + "ltab.txt");
		fol << ldata.ltabToString() ;
		fol.close();
	}
	ret["ltab"] = outpath + "ltab.txt";

  printf("\nOutpath : %s\n",outpath.c_str() );
	return ret;
}


/**
 * Extract attributes from vector of trace entries
 * Input: Vector of trace entry, attribute mode, attribute freq and attribute options
 * Output: Vector of strings (attributes to be injected to concept lattice)
 */
set<string> extractAttributes(vector<Entry> vec, int mode, int freq, int option){
	//printf("calling extractAttributes with vec size:%lu mode:%d freq:%d option:%d\n",vec.size(),mode,freq,option);
	set<string> ret;
	ret.clear();

	string sent;
  string sent2;
  int lc; // to hold loop count as freq

	typename vector<Entry>::iterator vit;
	unordered_map<string,int> freqTable;
	freqTable.clear();
	unordered_map<string,int>::iterator ftit;
  if (mode == 1){ // single entries
    for(vit = vec.begin() ; vit != vec.end() ; vit++){
      lc = 0;
			sent2 = (*vit).toString();
			printf("\t>>> %s\n",sent2.c_str() );
      if (splitString(sent2,'^').size() > 1){
        sent = splitString(sent2,'^')[0];
        lc = stoi(splitString(sent2,'^').back());
      }else{
        sent = sent2;
      }
			printf("\t>>> %s\n",sent.c_str() );
			if (freqTable.count(sent) > 0){
        //printf("\tif" );
        if (lc > 0){
          freqTable[sent]+=lc;
        }else{
          freqTable[sent]++;
        }
			}else{
        //printf("\telse" );
        if (lc > 0){
          freqTable[sent]=lc;
        }else{
          freqTable[sent]=1;
        }
			}
		}
    if (freq == 0){ // no freq
      // return vector of string : keys of freqTable
      for(ftit = freqTable.begin() ; ftit != freqTable.end() ; ftit++){
			//	printf("\t>added to freq" );
				ret.insert(ftit->first);
			}
			return ret;
    }else if (freq == 1){ // log 2 freq
      for(ftit = freqTable.begin() ; ftit != freqTable.end() ; ftit++){
				ret.insert(ftit->first+":"+to_string((int)trunc(log2((ftit->second) * 1.0 ))));
			}
			return ret;
    }else if (freq == 2){
      for(ftit = freqTable.begin() ; ftit != freqTable.end() ; ftit++){
				ret.insert(ftit->first+":"+to_string((int)trunc(log10((ftit->second) * 1.0 ))));
			}
      return ret;
    } else{ // actual freq (freq = 3)
      for(ftit = freqTable.begin() ; ftit != freqTable.end() ; ftit++){
				ret.insert(ftit->first+":"+to_string(ftit->second));
			}
			return ret;
    }
  } else if (mode == 2){  // set of distinct consecutive entry elements (w/ and w/o overlap)
		if (vec.size() == 0){
			//printf("Vector Size Problem\n");
		}else if (option == 0 ){
			//printf("we are going to be here. vec Empty? %d\n",vec.empty() );
      for(vit = vec.begin() ; vit != vec.end()-1 ; vit++){
				//printf("Even here?\n" );
				sent = (*vit).toString();
				//printf("OR here?\n" );
				sent = sent + " " + (*(vit+1)).toString();
				if (freqTable.count(sent) > 0){
					freqTable[sent]++;
				}else{
					freqTable[sent]=1;
				}
			}
    } else if (option == 1){
      for(vit = vec.begin() ; vit != vec.end() ; vit+=2){
				sent = (*vit).toString();
				sent = sent + " " + (*(vit+1)).toString();
				if (freqTable.count(sent) > 0){
					freqTable[sent]++;
				}else{
					freqTable[sent]=1;
				}
			}
    } else{
      printf("ERROR: Wrong Option\n");
    }
    if (freq == 0){ // no freq
      // return vector of string : keys of freqTable
      for(ftit = freqTable.begin() ; ftit != freqTable.end() ; ftit++){
				//printf("\t>added to freq" );
				ret.insert(ftit->first);
			}
			return ret;
    }else if (freq == 1){ // log 2 freq
      for(ftit = freqTable.begin() ; ftit != freqTable.end() ; ftit++){
				ret.insert(ftit->first+":"+to_string((int)trunc(log2((ftit->second) * 1.0 ))));
			}
			return ret;
    }else if (freq == 2){
      for(ftit = freqTable.begin() ; ftit != freqTable.end() ; ftit++){
				ret.insert(ftit->first+":"+to_string((int)trunc(log10((ftit->second) * 1.0 ))));
			}
      return ret;
    } else{ // actual freq (freq = 3)
      for(ftit = freqTable.begin() ; ftit != freqTable.end() ; ftit++){
				ret.insert(ftit->first+":"+to_string(ftit->second));
			}
			return ret;
    }
  } else if (mode == 3){// stats
      //orig Len
      freqTable["olen"] = vec[0].getOrigLen();
      // fdata len
      freqTable["flen"] = vec[0].getFdataLen();
      //ldata len
      freqTable["llen"] = vec[0].getLdataLen();
      //f/o
			if (vec[0].getFdataLen() != 0){
				freqTable["o/f"] = vec[0].getOrigLen() / vec[0].getFdataLen();
			}else{
				freqTable["o/f"] = 0;
			}
      //l/f
			if (vec[0].getLdataLen() != 0){
				freqTable["o/l"] = vec[0].getOrigLen() / vec[0].getLdataLen();
				freqTable["f/l"] = vec[0].getFdataLen() / vec[0].getLdataLen();
			}else{
				freqTable["o/l"] = 0;
				freqTable["f/l"] = 0;
			}
      //#loops
      freqTable["nl"] = vec[0].numOfLoops();
      // maxLoopBody
      freqTable["maxlb"] = vec[0].getMaxLoopBody();
      // MaxLC
      freqTable["maxlc"] = vec[0].getMaxLC();
      //Distinct
      freqTable["distincts"] = vec[0].getDistinctElements();
      if (freq == 0 || freq == 3) { //actual/original
        for(ftit = freqTable.begin() ; ftit != freqTable.end() ; ftit++){
  				ret.insert(ftit->first+":"+to_string(ftit->second));
  			}
  			return ret;
      } else if (freq == 1){
        for(ftit = freqTable.begin() ; ftit != freqTable.end() ; ftit++){
  				ret.insert(ftit->first+":"+to_string((int)trunc(log2((ftit->second) * 1.0 ))));
  			}
  			return ret;
      } else if (freq == 2){
        for(ftit = freqTable.begin() ; ftit != freqTable.end() ; ftit++){
					//printf("freq[%s] = %d\n",(ftit->first).c_str(),(int)ftit->second );
  				ret.insert(ftit->first+":"+to_string((int)trunc(log10((ftit->second) * 1.0 ))));
  			}
  			return ret;
      }else{
        printf("ERROR: Wrong Freq\n");
				return ret;
      }
  }
	return ret;
}






string filtbitTranslator(string filtbit,int k){
  string ret="";
  ret = ret + filtbit[0];
  ret = ret + filtbit[1];
	ret = ret + ".";
	if (filtbit[2] - '0' == 1){ //only @plt
		ret = ret + "plt.";
	}
	if (filtbit[3] - '0' == 1){ // memory
		ret = ret + "mem.";
	}
	if (filtbit[4] - '0' == 1){ // network
		ret = ret + "net.";
	}
	if (filtbit[5] - '0' == 1){ // poll
		ret = ret + "pol.";
	}
	if (filtbit[6] - '0' == 1){ // str
		ret = ret + "str.";
	}

	if (filtbit[7] - '0' == 1){ // MPI_
		ret = ret + "mpi.";
	}
	if (filtbit[8] - '0' == 1){ // MPI-related
    ret = ret + "mpiall.";
	}
	if (filtbit[9] - '0' == 1){ // mpi_collectives
		ret = ret + "mpicol.";
	}
	if (filtbit[10] - '0' == 1){ // mpi_send/rcv
		ret = ret + "mpisr.";
	}
	if (filtbit[11] - '0' == 1){ // omp-critical
		ret = ret + "ompcrit.";
	}
	if (filtbit[12] - '0' == 1){ // omp-mutex
    ret = ret + "ompmutex.";
	}
	if (filtbit[13] - '0' == 1){ // omp-rest
		ret = ret + "ompall.";
	}
	if (filtbit[14] - '0' == 1){ // custome
		ret = ret + "cust.";
	}
  ret = ret + filtbit[15];
	ret= ret + "K"+to_string(k);
  return ret;
}


string clNameTranslator(int m, int f, int op){
  //printf("%d %d %d\n",m,f,op );
  string ret;
	if (m == 1){
		ret = "sing.";
	} else if (m == 2){
		ret = "doub.";
	} else if (m == 3){
		ret = "stat.";
	}

  if (f == 0){
		ret += "orig.";
	} else if (f == 1){
		ret += "log2.";
	} else if (f == 2){
		ret += "log10.";
	} else if (f == 3){
		ret += "actual.";
	}

  if (op == 0){
		ret += "w";
	} else if (op == 1){
		ret += "wo";
	}
  //printf("ret %s\n", ret.c_str());
  return ret;
}


unordered_map<string,string> goPreprocess(string _inpath, int k){
	unordered_map<string,string> ret;
	vector<string> lot; // list of all traces within _inpath
	vector<string>::iterator vst_it;
	vector<string>::iterator vst_it2;

	string _trace; // Holds trace full path
	vector<string> tmpVst; // Temporary vector string

	// time measurement
	clock_t t;

	vector<Entry> entries; // To store each trace entries

	// Entry Initilization
	Entry tmp;

	string traceToken ;

	//vector<uint2> fdata; // To store filtered data
	//vector<uint2>::iterator vit; //

	Entry ldata; // To store nlr data (loop-data)
	vector<Entry>::iterator vlit;
	set<string> distincts;

	string tmps,line;

	string folderName = "nlr"+intToString(k);

	//making dir Outpath to store prep files there
	string outpath = _inpath + "/prep/"+folderName+"/";
	string outpath2 = _inpath + "/prep/";
	//printf(">>>Outpath : %s\n",outpath.c_str() );
	if (!isDir(outpath2)){
		//create directory
		if(mkdir(outpath2.c_str(),0777) == -1){
			perror("Error creating prep");
		}
	}

	if (!isDir(outpath)){
		//create directory
		if(mkdir(outpath.c_str(),0777) == -1){
			//printf("OUTPATH NOT: %s\n",outpath.c_str() );
			perror("Error creating prep");
		}
		else{
			printf("OUTPATH Created: %s\n",outpath.c_str() );
		}
	}
	else{
		printf("OUTPATH Existed: %s\n",outpath.c_str() );
	}

	// Up to here, the folder does exist. Now check if files do exist.
	// If they exist, read them and return
	// otherwise, do the preprocess, write to file and return


	// Get the list of all traces in _inpath/ptrace
	lot = listOfFiles(_inpath+"/","txt");

	//open log file
	ofstream flog(outpath + "log.txt");

	// preprocess all traces within _inpath/ptrace
	for(vst_it = lot.begin();vst_it != lot.end();vst_it++){
		entries.clear();
		_trace = _inpath+"/" + *vst_it;
		traceToken = splitString(_trace, '/').back() ;

		// Write to log file
		flog << "+" << traceToken << endl;
		// Check if trace is already preprocessed?
		ifstream fi(outpath + traceToken + ".txt");
		//printf("Check if prep trace file exist \n");
		if(!fi.good()){
			printf("reading from %s\n",_trace.c_str() );
			ifstream fin(_trace);
			// Read lines from each object file and create entries
			while(std::getline(fin, line)){
				//printf("line: %s\n", line.c_str());
				tmp = Entry();
				tmp.addElement(line);
				distincts.insert(line);
				tmp.setLC(1);
				entries.push_back(tmp);
			}
			printf("\nDetecting loops...length = %lu\n", entries.size());
			// ldata: Single entry object that holds info about dtab and ltab
			t = clock();
			ldata = nlr(entries,k,outpath + traceToken + ".txt");
			t = clock() - t ;
			flog << "nl:" << ldata.getLdataLen() << endl;
			flog << "nt:" << t << "," << (((float)t)/CLOCKS_PER_SEC) << endl;
			//Set stats
			//ldata.setDistinctElements(distincts.size());
			//ldata.setOrigLen(length);
			//ldata.setFdataLen(fdata.size());
		}//end if
		else{
			printf("File %s existed\n",(outpath+_trace).c_str() );
		}
		ret[traceToken] = outpath + traceToken+".txt";
	}// end for

	flog.close();
	// write dtab to file
	ifstream fid(outpath + "dtab.txt");
	if(!fid.good()){
		ofstream fod(outpath + "dtab.txt");
		fod << ldata.dtabToString() ;
		fod.close();
	}
	ret["dtab"] = outpath + "dtab.txt";

	// write ltab to file
	ifstream fil(outpath + "ltab.txt");
	if(!fil.good()){
		ofstream fol(outpath + "ltab.txt");
		fol << ldata.ltabToString() ;
		fol.close();
	}
	ret["ltab"] = outpath + "ltab.txt";

	printf("\nOutpath : %s\n",outpath.c_str() );
	return ret;
}
