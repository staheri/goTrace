/**
 * Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2018, All rights reserved
 * Code: main.cpp
 * Description: The skeleton of this code is basically traceReader by Martin Burtscher and Sindhu Devale.
 * It reads in compressed trace files one by one, decompress them and write the sequence of function IDs to text files. Actual functions are in .info files(mode to run with: -m 3).
 */
#include "actions.h"
#include <typeinfo>

using namespace std;

//set<pair<string,pair<int,string> > > ftable;

struct globalArgs_t{
	char *path;
	char *info;
	char *trace;
	int mode;
  int atrMode;
  int atrFreq;
  int atrOption;
	int k;
	char *output;
	string filtbit;
}globalArgs;


int main(int argc, char* argv[]) {
	printf("diffTrace v%s (%s)\n", QUOTE(version), __FILE__);
	//printUsage();

	static const char *optstring = "m:p:o:i:t:f:a:k:q:n:h?";
	globalArgs.path = NULL;
	globalArgs.info = NULL;
	globalArgs.trace = NULL;
	globalArgs.mode = 0;
	globalArgs.output = NULL;
  globalArgs.atrMode = 0;
  globalArgs.atrFreq = 0;
  globalArgs.atrOption = 0;
	globalArgs.filtbit = "00000";

	int opt;

	while ((opt = getopt(argc,argv,optstring)) != -1){
		switch(opt){
			case 'm': // Mode
				globalArgs.mode = stoi(optarg);
			case 'p': // Path
				globalArgs.path = optarg;
				break;
      case 'a': //Attribute Mode
  			globalArgs.atrMode = stoi(optarg);
  			break;
      case 'q': // Atr Freq Mode
  			globalArgs.atrFreq = stoi(optarg);
  			break;
      case 'n': // Atr Option Mode
  			globalArgs.atrOption = stoi(optarg);
  			break;
			case 'o': // Outpath
				globalArgs.output = optarg;
				break;
			case 'i': // single info path
				globalArgs.info = optarg;
				break;
			case 't': // single trace path
				globalArgs.trace = optarg;
				break;
			case 'k': // single trace path
				globalArgs.k = stoi(optarg);
					break;
			case 'f': // filtbit
				globalArgs.filtbit = optarg;
				if (strlen(optarg) != FILTBITSIZE){
					fprintf (stderr, "Short filtBit, supposed to be %d but got %lu.\n",FILTBITSIZE, strlen(optarg));
					abort();
				}
				break;
			case 'h':
			case '?':
				if (optopt == 'm' || optopt == 'p' || optopt == 'o' || optopt == 'i' || optopt == 't' || optopt == 'f' || optopt == 'a' || optopt == 'q' || optopt == 'n' || optopt == 'k' ) {
					fprintf (stderr, "Option -%c requires an argument.\n", optopt);
					printUsage();
				} else if (isprint (optopt)){
					fprintf (stderr, "Unknown option `-%c'.\n", optopt);
					printUsage();
				} else {
					fprintf (stderr,"Unknown option character `\\x%x'.\n",optopt);
					printUsage();
				}
        		return 1;
				break;
			default:
				printUsage();
				abort();
		}
	}
	if ( globalArgs.mode == 1 ){
		texTraceBatch(globalArgs.path);
	}
	if ( globalArgs.mode == 2 ){
		genCL(globalArgs.path,globalArgs.filtbit,globalArgs.atrMode,globalArgs.atrFreq,globalArgs.atrOption,globalArgs.k);
	}
	if ( globalArgs.mode == 3 ){
		fptrace(globalArgs.path,globalArgs.k);
	}
	if ( globalArgs.mode == 4 ){
		genFPCL(globalArgs.path,globalArgs.atrMode);
	}
	if ( globalArgs.mode == 5 ){
		genGoCL(globalArgs.path,globalArgs.atrMode,globalArgs.atrFreq,globalArgs.atrOption,globalArgs.k);
	}
}
