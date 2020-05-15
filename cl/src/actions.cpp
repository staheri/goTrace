/** actions.h
 * mid-level interface for communicting the front-end and back-end for trace and CL operations
 * Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2019, All rights reserved
 */
#include "actions.h"
#include <cstdlib>
#include <cstdio>

/**
 * Called from main to generate Concept Lattice.
 * It first Preprocess data (decompress, filter, detect loops, extract attributes, etc.)
 * Then it creates the lattice from generated data
 */
void genCL(string _inpath,string _filtbit,int _atrMode, int _atrFreq, int _atrOption, int k){

  // Variables
  vector<Entry> exAtrInput;

  // _outpath
  string _outpath = _inpath+"/cl/"+filtbitTranslator(_filtbit,k)+"/";

  if(mkdir(_outpath.c_str(),0777) == -1){
    perror("Error creating cl");
  }

  // Preprocess data and store results in allTraces
  printf("\nReading/Preprocessing trace entries in %s\n",_inpath.c_str());
  //unordered_map<string,vector<Entry>> allTraces = preprocess(_inpath,_filtbit, k);
  unordered_map<string,string> traceEntryPaths = preprocess(_inpath,_filtbit, k);

  // Add keys of allTraces to a vector to sort
  vector<string> allTrcKeys;
  allTrcKeys.reserve(traceEntryPaths.size()-2);
  for (auto& it : traceEntryPaths) {
    if (it.first != "ltab" && it.first != "dtab")
      allTrcKeys.push_back(it.first);
  }
  // Sort
  std::sort(allTrcKeys.begin(),allTrcKeys.end(),
		[](string a,string b){
			vector<string> tta = splitString(a,'.');
			vector<string> ttb = splitString(b,'.');
			int sss = tta.size();
			if (atoi(tta[sss-2].c_str()) == atoi(ttb[sss-2].c_str())){
				return atoi(tta[sss-1].c_str()) < atoi(ttb[sss-1].c_str());
			}
			return atoi(tta[sss-2].c_str()) < atoi(ttb[sss-2].c_str());
		});

  // Extracted attributes from allTrace (preprocessed data) stored in this
  set<string> atrSet;

  // Aux variables
  //vector<string> atrList;
  vector<string>::iterator vit;
  set<string>::iterator sit;
  //typename unordered_map<string,vector<Entry>>::iterator tit;


  // Generating CL
  string clName = clNameTranslator(_atrMode,_atrFreq,_atrOption);
  Lattice lat = Lattice(clName);
  // To hold an object of each trace and attribute for accessing their hashtables later
  Trace trc;
  Attribute<string> atr;
	set<int> attrIDs;

  clock_t t = clock();
  printf("\nExtracting Attributes & Creating CL %s\n",clName.c_str());
  for (vit = allTrcKeys.begin() ; vit != allTrcKeys.end() ; vit++){
    //printf("%s > Vector entry Size: %d \n", (tit->first).c_str(),(tit->second).size()  );

    //Creating Trace(object) and Attribute objects
		//printf("Crating Trace Object...\n\n");
		trc = Trace(*vit);

    //Extracting Attributes
    atrSet.clear();
    attrIDs.clear();
    exAtrInput.clear();
    exAtrInput = readEntryFile(traceEntryPaths[*vit]);
    atrSet = extractAttributes( exAtrInput, _atrMode, _atrFreq, _atrOption);

    //Read attributes and store their ids
    for (sit = atrSet.begin() ; sit != atrSet.end() ; sit++){
      //printf("Crating Attribute Objects...\n\n");
			atr = Attribute<string>(*sit);

      //printf("Adding Attribute Object to global ds\n\n");
			attrIDs.insert(atr.getID());
			lat.setMaxAttribute(atr.getAttributeCount());
    }

    // Making concepts and injecting to CL
    Concept c = Concept(trc.getID(),attrIDs);
		lat.addConcept(c);
		lat.addSubgraph(lat.toDotEdges(c.getID(),0),c.toString());

    /*for(vit = atrList.begin() ; vit != atrList.end() ; vit++){
      printf("\t%s\n", (*vit).c_str());
    }*/

  } // Lattice Generation finished
  t = clock()-t;

  //lat.printLatticeComponents();
	//printf("Maximum Attributes: %d\n",lat.getMaxAttribute());
	string ldot = lat.toDot(clName,0).c_str();
	//printf("%s\n",ldot.c_str());
	//printf("%s\n",trc.tableString().c_str());
	//printf("%s\n",atr.tableString().c_str());
  printf("\nFinished CL generation in %.3f seconds\nWriting CLs to %s\n",(((float)t)/CLOCKS_PER_SEC),_outpath.c_str());
  ofstream allDot;
  printf("%s\n",(_outpath+clName+".dot").c_str());
  allDot.open(_outpath+clName+".dot");
  allDot << lat.toDot(clName,0).c_str();
  allDot.close();

  ofstream ttbl;
  ofstream atbl;
  ofstream cmat; //context bit matrix
  ofstream lmat; //lattice adjacency matrix

  ttbl.open(_outpath+clName+".objTable.txt");
  ttbl << trc.tableString().c_str();
  ttbl.close();

  atbl.open(_outpath+clName+".attrTable.txt");
  atbl << atr.tableString().c_str();
  atbl.close();

  lmat.open(_outpath+clName+".latmat.txt");
  lmat << lat.toLatMat().c_str();
  lmat.close();

  cmat.open(_outpath+clName+".context.txt");
  cmat << lat.toContextBitmax().c_str();
  cmat.close();
  printf("\n############ END ############\n\n");

}

void texTraceBatch(string _inpath){

  vector<string> lot; // list of all traces within _inpath/ptrace
  vector<string>::iterator vst_it;
  vector<string>::iterator vst_it2;

  string _trace; // Holds trace full path
  string _info;
  vector<string> tmpVst; // Temporary vector string

  // For Decompression
  uint2* data ;
  string* info;
  uint8 length;
  string traceToken ;
  string infoToken ;

  //making dir Outpath to store prep files there
  string outpath = _inpath + "/texTrace/";

  if (!isDir(outpath)){
    //create directory
    if(mkdir(outpath.c_str(),0777) == -1){
      perror("Error creating texTrace");
    }
  }


  // Get the list of all traces in _inpath/ptrace
  lot = listOfTraceFiles(_inpath+"/ptrace");

  // preprocess all traces within _inpath/ptrace
  for(vst_it = lot.begin();vst_it != lot.end();vst_it++){
    _trace = _inpath+"/ptrace/" + *vst_it;
    traceToken = splitString(_trace, '/').back() ;

    // Check if trace is already preprocessed?
    ofstream fo(outpath + traceToken + ".txt");
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
    printf("decompressing %s\n", _trace.c_str() );

    data = readFile(_trace.c_str(), length);
    info = readInfo(_info.c_str());
    info[0] = "[ret]";
    printf("writing len: %llu\n",length );
    // Write data to file
    for(uint8 i = 0 ; i < length ; i++){
      fo << info[data[i]] << endl ;
    }
    fo.close();
  }
}

void fptrace(string _trace,int k){
  clock_t t;
  // For Decompression
  uint2* data ;
  vector<uint2> vdata;
  vector<uint2>::iterator vit;
  uint8 length;
  string traceToken ;
  string _info;
  vector<string>::iterator vst_it;
  Entry tmp;
  Entry ldata; // To store nlr data (loop-data)

  vector<string> traceTokenized = splitString(_trace, '/');
  //traceToken = splitString(traceTokenized.back(),'.')[0] ;
  traceToken = traceTokenized.back() ;

  string outpath = "";
  string outf;
  for (vst_it = traceTokenized.begin(); vst_it != traceTokenized.end()-1;vst_it++){
    outpath = outpath + *vst_it + "/";
  }
  //outf = outpath + traceToken + ".txt";
  _info = outpath + splitString(traceToken,'.')[0] + ".info";
  // Check if trace is already preprocessed?
  //ofstream fo(outf);

  // Decompression
  length = 0;
  t = clock();
  printf("decompressing %s\n", _trace.c_str() );
  printf("Outf: %s\nOutpath: %s\nTraceToken: %s\n", outf.c_str() , outpath.c_str() , traceToken.c_str() );
  vector<Entry> entries; // To store each trace entries
  data = readFile(_trace.c_str(), length);
  t = clock() - t;
  printf("Decompression Time: %.3f\n",((float)t)/CLOCKS_PER_SEC );
  printf("Creating CFG from data vector, size: %llu\n",length );

  // Read Info

  ifstream fi(_info);
  map<string,string> info;
  vector<string> linevector;
  string line,tline;
  while(std::getline(fi, line)){
    linevector = splitString(line,'|');
    tline = "";
    for (vst_it = linevector.begin() + 1 ; vst_it != linevector.end() ; vst_it++){
      tline = tline + (*vst_it) + '|';
    }
    info[linevector[0]]=tline;
    printf("info ID: %s\n",linevector[0].c_str() );
  }
  info["0"] = "start";

  // NLR

  // for(uint8 i = 0 ; i<length ; i++){
  //   if (i%1000000 == 0){
  //     printf("%llu/%llu\n",i,length );
  //   }
  //   if (data[i] != 0){
  //     tmp = Entry();
  //     tmp.addElement(intToString((int)data[i]));
  //     //distincts.insert(info[(*vit)]);
  //     tmp.setLC(1);
  //     entries.push_back(tmp);
  //   }
  // }
  // //free(data);
  //printf("writing len: %llu\n",length );
  //ofstream flog("log_"+outpath + traceToken + "." +intToString(k)+"nlr.txt");
  //t = clock();

  //ldata = nlr(entries,k,outpath + traceToken + "." +intToString(k)+"nlr.txt");
  //t = clock() - t;
  //flog << "time:" << t << "," << (((float)t)/CLOCKS_PER_SEC) << endl;
  //flog.close();

  // Write data to file
  //for(uint8 i = 0 ; i < length ; i++){
  //  fo << data[i] << "," ;
  //}
  //fo.close();

  // CFG Prototyping
  typedef pair<uint2,uint2> edgeID;
  edgeID newEdge;
  map<edgeID,int> edges;
  map<edgeID,int>::iterator mit;
  pair<map<edgeID,int>::iterator,bool> ret;
  uint2 prev = 0;
  set<uint2> nodes;
  set<uint2>::iterator sit;
  string cfgDot = "";
  t = clock();
  for(uint8 i = 0 ; i<length ; i++){
    if (i%10000000 == 0){
      printf("%llu/%llu\n",i,length );
    }
    if (data[i] != 0){
      nodes.insert(data[i]);
      newEdge = make_pair(prev,data[i]);
      mit = edges.find(newEdge);
      if (mit != edges.end()){
        // Key exist
        //printf("key Exist: %s\n",mit->first.first() );
        mit->second +=1;
      }else{
        edges[newEdge] = 1;
      }
      //ret = edges.insert(pair<edgeID,int>(newEdge,1));
      //if (ret.second==false) {
        //edges[newEdge] = ret.first->second + 1;
      //}
      prev = data[i];
    }
  }
  t = clock() - t;

  int countss = 0;
  cfgDot = cfgDot + "digraph g{\n\t";
  for (mit = edges.begin() ; mit != edges.end() ; mit++){
    printf("%hu -> %hu : %d \n",mit->first.first,mit->first.second,mit->second );
   cfgDot = cfgDot + "\"" +info[to_string(mit->first.first)] + "\" -> \"" + info[to_string(mit->first.second)] + "\" [label = \""+to_string(mit->second) +"\"]\n\t";
   countss += mit->second;
  }
  cfgDot = cfgDot +"\n}\n";
  printf("Counts : %d\n", countss);
  ofstream fdot(outpath + traceToken + ".dot");
  fdot << cfgDot ;
  fdot.close();
  printf("CFG Creation Time: %.3f\n",((float)t)/CLOCKS_PER_SEC );
  printf("Total Edges: %lu\nTotal Nodes: %lu\n",edges.size(),nodes.size());

}

void genFPCL(string _inpath,int _atrMode){


  // take input path
  // *.out : has all records that I need
  //  - function names
  //  - BB ids
  //  - function bit-string + frequencies
  // read in all functions
  //_inpath = *.out
  ifstream fi(_inpath);
  unordered_map<string,vector<string>> entries;
  unordered_map<string,vector<string>>::iterator umit;
  set<string> objects;
  vector<string> linevector;
  string line,tline;
  vector<string> ent;
  printf("reading from %s\n",_inpath.c_str() );
  while(std::getline(fi, line)){
    printf("line: %s\n", line.c_str());
    linevector = splitString(line,'|');
    objects.insert(linevector[0]);
    umit = entries.find(linevector[0]);
    if (umit == entries.end()){
      ent.clear();
      ent.push_back(linevector[1]+":"+linevector[2]);
      entries[linevector[0]] = ent;
    }else{
      entries[linevector[0]].push_back(linevector[1]+":"+linevector[2]);
    }
  }

  // Objects, block counts and bitsets are in entries now

  string clName = "test";
  Lattice lat = Lattice(clName);
  // To hold an object of each trace and attribute for accessing their hashtables later
  Trace trc;
  Attribute<string> atr;
	set<int> attrIDs;
  string _outpath = "";
  vector<string> temp;
  vector<string> temp2;
  vector<string>::iterator vit2;
  vector<string>::iterator vit;
  set<string>::iterator sit;
  string setArr[] = {"SSE_DATA","SSE_ARITH","SSE_OTHER","SSE2_DATA","SSE2_ARITH","SSE2_OTHER","FP_DATA","FP_ARITH","FP_OTHER","AVX","AVX2","FMA"};
  vector<string> setVec (setArr,setArr  + sizeof(setArr) / sizeof(setArr[0]) );

  set<string> atrSet;

  printf("\nExtracting Attributes & Creating CL %s\n",clName.c_str());
  for (umit = entries.begin() ; umit != entries.end() ; umit++){
    //Creating Trace(object) and Attribute objects
		printf("Crating Trace Object...\n\n");
		trc = Trace(umit->first);


    //Extracting Attributes
    atrSet.clear();
    attrIDs.clear();
    temp.clear();
    temp = umit->second;
    for (vit = temp.begin() ; vit != temp.end() ; vit++){
      temp2 = splitString(splitString(*vit,':')[1],'.');
      printf("temp entry %s\n",(*vit).c_str() );
      for (uint i = 0 ; i<temp2.size()-1 ; i++){
        printf("temp2[%d]:%s\n",i,temp2[i].c_str() );
        if (temp2[i] != "0"){
          atrSet.insert(setArr[i]);
        }
      }
    }

    //Read attributes and store their ids
    for (sit = atrSet.begin() ; sit != atrSet.end() ; sit++){
      //printf("Crating Attribute Objects...\n\n");
			atr = Attribute<string>(*sit);

      //printf("Adding Attribute Object to global ds\n\n");
			attrIDs.insert(atr.getID());
			lat.setMaxAttribute(atr.getAttributeCount());
    }

    // Making concepts and injecting to CL
    Concept c = Concept(trc.getID(),attrIDs);
		lat.addConcept(c);
		lat.addSubgraph(lat.toDotEdges(c.getID(),0),c.toString());

    /*for(vit = atrList.begin() ; vit != atrList.end() ; vit++){
      printf("\t%s\n", (*vit).c_str());
    }*/

  } // Lattice Generation finished

  //lat.printLatticeComponents();
	//printf("Maximum Attributes: %d\n",lat.getMaxAttribute());
	string ldot = lat.toDot(clName,0).c_str();
	//printf("%s\n",ldot.c_str());
	//printf("%s\n",trc.tableString().c_str());
	//printf("%s\n",atr.tableString().c_str());
  //printf("\nFinished CL generation in %.3f seconds\nWriting CLs to %s\n",(((float)t)/CLOCKS_PER_SEC),_outpath.c_str());
  ofstream allDot;
  printf("%s\n",(_outpath+clName+".dot").c_str());
  allDot.open(_outpath+clName+".dot");
  allDot << lat.toDot(clName,0).c_str();
  allDot.close();

  ofstream ttbl;
  ofstream atbl;
  ofstream cmat; //context bit matrix
  ofstream lmat; //lattice adjacency matrix

  ttbl.open(_outpath+clName+".objTable.txt");
  ttbl << trc.tableString().c_str();
  ttbl.close();

  atbl.open(_outpath+clName+".attrTable.txt");
  atbl << atr.tableString().c_str();
  atbl.close();

  lmat.open(_outpath+clName+".latmat.txt");
  lmat << lat.toLatMat().c_str();
  lmat.close();

  cmat.open(_outpath+clName+".context.txt");
  cmat << lat.toContextBitmax().c_str();
  cmat.close();
  printf("\n############ END ############\n\n");

}

void genGoCL(string _inpath, int _atrMode, int _atrFreq, int _atrOption, int k){
  // Read in files from <trace>/<pattern>/texTrace(-id)/<objectFiles>
  // for each objectFile:
  //     readLines()
  //     store them in a data vector
  //        data: create vector of Entry
  //
  //     create object/attribute
  //     create concept concept lattice
  vector<Entry> exAtrInput;
  vector<string>::iterator vit;

  // _outpath
  string _outpath = _inpath+"/cl/";
  if(mkdir(_outpath.c_str(),0777) == -1){
    perror("Error creating cl");
  }
  _outpath = _outpath + "nlr"+intToString(k) + "/";
  if(mkdir(_outpath.c_str(),0777) == -1){
    perror("Error creating cl/nlr");
  }

  // Preprocess data and store results in allTraces
  printf("\nReading/Preprocessing trace entries in %s\n",_inpath.c_str());
  unordered_map<string,string> traceEntryPaths = goPreprocess(_inpath,k);

  // Add keys of allTraces to a vector to sort
  vector<string> allTrcKeys;
  allTrcKeys.reserve(traceEntryPaths.size()-2);
  for (auto& it : traceEntryPaths) {
    if (it.first != "ltab" && it.first != "dtab"){
      allTrcKeys.push_back(it.first);
      //printf("<> %d\n",atoi(splitString(it.first,'-')[0].substr(1).c_str()));
    }
  }

  // Sort
  std::sort(allTrcKeys.begin(),allTrcKeys.end(),
		[](string a,string b){
			int ta = atoi(splitString(a,'-')[0].substr(1).c_str());
			int tb = atoi(splitString(b,'-')[0].substr(1).c_str());
			return ta < tb;
		});
    /*for (vit = allTrcKeys.begin() ; vit != allTrcKeys.end() ; vit++){
      printf("%s\n",(*vit).c_str() );
    }*/
    // Extracted attributes from allTrace (preprocessed data) stored in this
    set<string> atrSet;

    // Aux variables
    //vector<string> atrList;
  //  vector<string>::iterator vit;
    set<string>::iterator sit;
    //typename unordered_map<string,vector<Entry>>::iterator tit;


    // Generating CL
    string clName = clNameTranslator(_atrMode,_atrFreq,_atrOption);
    Lattice lat = Lattice(clName);
    // To hold an object of each trace and attribute for accessing their hashtables later
    Trace trc;
    Attribute<string> atr;
  	set<int> attrIDs;

    clock_t t = clock();
    printf("\nExtracting Attributes & Creating CL %s\n",clName.c_str());
    for (vit = allTrcKeys.begin() ; vit != allTrcKeys.end() ; vit++){
      //printf("%s > Vector entry Size: %d \n", (tit->first).c_str(),(tit->second).size()  );

      //Creating Trace(object) and Attribute objects
  		//printf("Crating Trace Object...\n\n");
  		trc = Trace(*vit);
      printf("OBJ: %s\n",(*vit).c_str() );

      //Extracting Attributes
      atrSet.clear();
      attrIDs.clear();
      exAtrInput.clear();
      exAtrInput = readEntryFile(traceEntryPaths[*vit]);
      printf("ATRINP: %s\n",traceEntryPaths[*vit].c_str() );
      atrSet = extractAttributes( exAtrInput, _atrMode, _atrFreq, _atrOption);
      printf(">>>>>>>>>>>>> AFTER ATRINP\n" );

      //Read attributes and store their ids
      for (sit = atrSet.begin() ; sit != atrSet.end() ; sit++){
        printf("Crating Attribute Objects...\n\n");
  			atr = Attribute<string>(*sit);
        printf("\tatr: %s\n",(*sit).c_str() );
        //printf("Adding Attribute Object to global ds\n\n");
  			attrIDs.insert(atr.getID());
  			lat.setMaxAttribute(atr.getAttributeCount());
      }

      // Making concepts and injecting to CL
      Concept c = Concept(trc.getID(),attrIDs);
  		lat.addConcept(c);
  		lat.addSubgraph(lat.toDotEdges(c.getID(),0),c.toString());

      /*for(vit = atrList.begin() ; vit != atrList.end() ; vit++){
        printf("\t%s\n", (*vit).c_str());
      }*/

    } // Lattice Generation finished
    t = clock()-t;

    //lat.printLatticeComponents();
  	//printf("Maximum Attributes: %d\n",lat.getMaxAttribute());
  	string ldot = lat.toDot(clName,0).c_str();
  	//printf("%s\n",ldot.c_str());
  	//printf("%s\n",trc.tableString().c_str());
  	//printf("%s\n",atr.tableString().c_str());
    printf("\nFinished CL generation in %.3f seconds\nWriting CLs to %s\n",(((float)t)/CLOCKS_PER_SEC),_outpath.c_str());
    ofstream allDot;
    printf("%s\n",(_outpath+clName+".dot").c_str());
    allDot.open(_outpath+clName+".dot");
    allDot << lat.toDot(clName,0).c_str();
    allDot.close();

    ofstream ttbl;
    ofstream atbl;
    ofstream cmat; //context bit matrix
    ofstream lmat; //lattice adjacency matrix

    ttbl.open(_outpath+clName+".objTable.txt");
    ttbl << trc.tableString().c_str();
    ttbl.close();

    atbl.open(_outpath+clName+".attrTable.txt");
    atbl << atr.tableString().c_str();
    atbl.close();

    lmat.open(_outpath+clName+".latmat.txt");
    lmat << lat.toLatMat().c_str();
    lmat.close();

    cmat.open(_outpath+clName+".context.txt");
    cmat << lat.toContextBitmax().c_str();
    cmat.close();
    printf("\n############ END ############\n\n");


}
