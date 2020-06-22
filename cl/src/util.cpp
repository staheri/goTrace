/**
 * Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2018, All rights reserved
 * Code: util.cpp
 * Description: utility functions to use in diffTrace project
 */
#include "util.h"

/**
 * Split the string<text> using sep as separator
 */

vector<string> splitString(string text, char sep) {
	//printf("Split String...%s\n",text.c_str() );
	vector<string> tokens;
	size_t start = 0, end = 0;
	while ((end = text.find(sep, start)) != string::npos) {
		tokens.push_back(text.substr(start, end - start));
		start = end + 1;
	}
	tokens.push_back(text.substr(start));
	return tokens;
}


/**
 * Convert float to string
 */
string ftos(double f, int dec){
	stringstream stream;
	stream << fixed << setprecision(dec) << f;
	return stream.str();
}
/**
 * Print help message how to use CLD
 */
void printUsage(void){
	printf("%s\n", "Two major mode [To be added more]:" );
	printf("\t%s\n", "1- Creates concept lattice out of all trace files" );
	printf("\t\t%s\n", "Usage for this mode:");
	printf("\t\t%s\n", "./exec -m 1 -p <path_to_trace> -o <output_file_name> -d [mode_options]");
	printf("\t\t%s\n", "mode options");
	printf("\t\t\t%s\n", "1: Create concept lattice from function calls");
	printf("\t\t\t%s\n", "2: Create concept lattice from function edges");
	printf("\t%s\n", "2- Creates a single text file from a trace file and info file" );
	printf("\t\t%s\n", "Usage for this mode:");
	printf("\t\t%s\n", "./exec -m 2 -i <info_file> -t <trace_file> -o <output_file_name(without extension)> -d [output_mode_options]");
	printf("\t\t%s\n", "output_mode_options");
	printf("\t\t\t%s\n", "1: Function Calls and their frequencies");
	printf("\t\t\t%s\n", "2: Function Call edges (caller-callee) and their frequencies");
	printf("\t\t\t%s\n", "3: Approximate call stack");
	printf("\t\t\t%s\n", "4: Full trace of function calls");
	printf("\t%s\n", "3- Store DECOMPRESSED data and info into text files" );
	printf("\t\t%s\n", "Usage for this mode:");
	printf("\t\t%s\n", "./exec -m 3 -i <info_file> -t <trace_file> -o <output path>");
}

/**
 * Returns a list of all files within the <path> with <ext> extension
 */
vector<string> listOfFiles(const string _path, const char* ext) {
	vector <string> l;
	string t = "";
	string path= _path+"/";
	DIR *dir;
	struct dirent *ent;
	if ((dir = opendir(path.c_str())) != NULL) {
		/* print all the files and directories within directory */
		while ((ent = readdir(dir)) != NULL) {
			if (!strcmp(ent->d_name, "."))
				continue;
			if (!strcmp(ent->d_name, ".."))
				continue;
			// in linux hidden files all start with '.'
			if (ent->d_name[0] == '.')
				continue;
			if (strstr(ent->d_name, ext)) {
				l.push_back(ent->d_name);
			}
		}
		closedir(dir);
	} else {
		/* could not open directory */
	  printf("Path: %s\n",path.c_str());
	  perror("L.O.F: could not open directory\n");
	}
	std::sort(l.begin(),l.end());
	return l;
}


/**
 * Returns a list of all trace files within the <path> (with extension .0 .1 .n)
 */
vector<string> listOfTraceFiles(const string _path) {
	vector <string> l;
	string t = "";
	string path= _path+"/";
	DIR *dir;
	struct dirent *ent;
	if ((dir = opendir(path.c_str())) != NULL) {
		/* print all the files and directories within directory */
		while ((ent = readdir(dir)) != NULL) {
			if (!strcmp(ent->d_name, "."))
				continue;
			if (!strcmp(ent->d_name, ".."))
				continue;
			// in linux hidden files all start with '.'
			if (ent->d_name[0] == '.')
				continue;
			if (!strstr(ent->d_name, ".info")) {
				l.push_back(ent->d_name);
			}

		}
		closedir(dir);
	} else {
		/* could not open directory */
	  printf("Path: %s\n",path.c_str());
	  perror("L.O.F: could not open directory\n");
	}
	std::sort(l.begin(),l.end());
	return l;
}

/**
 * Returns a list of all folders within the <path>
 */
vector<string> listOfFolders(const string _path) {
	vector <string> l;
	string path= _path+"/";
	string t = "";
	DIR *dir;
	struct dirent *ent;
	if ((dir = opendir(path.c_str())) != NULL) {
		/* print all the files and directories within directory */
		while ((ent = readdir(dir)) != NULL) {
			if (!strcmp(ent->d_name, "."))
				continue;
			if (!strcmp(ent->d_name, ".."))
				continue;
			// in linux hidden files all start with '.'
			if (ent->d_name[0] == '.')
				continue;
			// dirFile.name is the name of the file. Do whatever string comparison
			// you want here. Something like:
			if (ent->d_type == DT_DIR) {
				l.push_back(ent->d_name);
			}
		}
		closedir(dir);
	} else {
		/* could not open directory */
		perror("could not open directory");
	}
	return l;
}

/**
* Input: Set of integer, Output: A string showing all ranges of numbers within the input set
*/
string setSummary(set<int> shrinkerSet,int flag){
	if (flag) {
		int iprev,istart,iend;
		string tmps,sistart,siend;
		set<int>::iterator ito;
		iprev = -1;
		if (shrinkerSet.size() == 0) {
			tmps = "*EMPTY*";
		} else{
			for (ito = shrinkerSet.begin() ; ito != shrinkerSet.end() ; ito++){
			//printf("----<<< %d \n",int(*ito));
				if (iprev == -1){
					istart = *ito ;
					iend = *ito ;
					//printf("iprev == -1 |  %d - %d\n",istart,iend);
				}
				else{
					if (iprev == *ito - 1){
						iend = *ito ;
						//printf("iprev == current |  %d - %d\n",istart,iend);
					}
					else {
						sistart = intToString(istart);
						siend = intToString(iend);
						//printf("!! iprev == current |  %d - %d\n",istart,iend);
						//wrap previouses
						if (istart == iend){
							tmps = tmps + sistart + ",";
						}
						else{
							tmps = tmps + sistart + "-" + siend + ",";
						}
						istart = *ito;
						iend = *ito;
						//set new istart
					}
				}
				iprev = *ito ;
				//printf("val= %d \n",*ito);
			}
			//printf ("\t tmps : %s \n",tmps.c_str());
			sistart = intToString(istart);
			siend = intToString(iend);
			//printf("REACHES END |  %d - %d\n",istart,iend);
			//wrap previouses
			if (istart == iend){

				//printf ("\t tmps : %s \n",tmps.c_str());
				tmps = tmps + sistart ;
			}
			else{
				//printf ("INNNNN\n");
				tmps = tmps + sistart + "-" + siend ;
			}
		}

		return tmps;
	}
	else{
		string tmps;
		set<int>::iterator ito;
		for (ito = shrinkerSet.begin() ; ito != shrinkerSet.end() ; ito++){
			tmps = tmps + intToString(*ito) + ",";
		}
		return tmps;


	}
}

/**
 * Convert intToString to string
 */
string intToString(int i){
	string s = "";
	std::ostringstream oss;
	oss << i;
	s += oss.str();
	return s.c_str();

}

/**
 * Check if a directory exist
 */
bool isDir(string path){
	struct stat sb;
	return (stat(path.c_str(),&sb) == 0 && S_ISDIR(sb.st_mode)) ;
}
