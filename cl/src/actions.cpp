/** actions.h
* mid-level interface for communicting the front-end and back-end for trace and CL operations
* Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2019, All rights reserved
*/
#include "actions.h"
#include <cstdlib>
#include <cstdio>



void genGeneralCL(string _inpath){

  vector<string> lot; // list of all traces within _inpath
  vector<string>::iterator vst_it;

  string _trace,line;
  // _outpath
  string _outpath = _inpath+"/cl/";
  if(mkdir(_outpath.c_str(),0777) == -1){
    perror("Error creating cl");
  }

  // Preprocess data and store
  printf("\nReading trace entries in %s\n",_inpath.c_str());

  // Get the list of all traces in _inpath/
  lot = listOfFiles(_inpath+"/","txt");

  // Sort
  std::sort(lot.begin(),lot.end(),
  [](string a,string b){
    int ta = atoi(splitString(a,'.')[0].substr(1).c_str());
    int tb = atoi(splitString(b,'.')[0].substr(1).c_str());
    return ta < tb;
  });


  //string clName = clNameTranslator(_atrMode,_atrFreq,_atrOption);
  string clName = "test";
  Lattice lat = Lattice(clName);
  // To hold an object of each trace and attribute for accessing their hashtables later
  Trace trc;
  Attribute<string> atr;
  set<int> attrIDs;

  printf("\nExtracting Attributes & Creating CL %s\n",clName.c_str());


  // preprocess all traces within _inpath/ptrace
  for(vst_it = lot.begin();vst_it != lot.end();vst_it++){
    _trace = _inpath+"/" + *vst_it;

    ifstream fin(_trace);

    trc = Trace(*vst_it);
    printf("OBJ: %s\n",(*vst_it).c_str() );

    //Extracting Attributes
    attrIDs.clear();
    //Read attributes and store their ids
    while(std::getline(fin, line)){
      printf("Crating Attribute Objects...\n\n");
      atr = Attribute<string>(line);
      printf("\tatr: %s\n",line.c_str() );
      //printf("Adding Attribute Object to global ds\n\n");
      attrIDs.insert(atr.getID());
      lat.setMaxAttribute(atr.getAttributeCount());

      // Making concepts and injecting to CL
      Concept c = Concept(trc.getID(),attrIDs);
      lat.addConcept(c);
      lat.addSubgraph(lat.toDotEdges(c.getID(),0),c.toString());
    }
  }
  string ldot = lat.toDot(clName,0).c_str();
  //printf("%s\n",ldot.c_str());
  //printf("%s\n",trc.tableString().c_str());
  //printf("%s\n",atr.tableString().c_str());

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
