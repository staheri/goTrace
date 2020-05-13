#include <regex>
#include <stdio.h>
#include "nlr.h"


int main(int argc, char* argv[]){
  //std::string subj = "MPI_Send MPI_Recv PMPI_Reduce MPID_Poll MPI_kir";
  /*
  std::string subj = ".plt kir.plt ok@plt plt";
  std::regex re ("\\w*\\.plt");

  std::smatch match;
  std::regex_search(subj,match,re);
  printf("Size of match %d, %s\n", match.size(), match.str().c_str());
  std::sregex_iterator next(subj.begin(), subj.end(), re);
  std::sregex_iterator end;
  while (next != end){
    std::smatch match = *next;
    printf("Match: %s\n",match.str().c_str());
    next++;
  }
  */
  vector<int> v;
//  vector<Entry> container;
  vector<Entry> entries;
  vector<Entry> ret;
  Entry tmp;
  //vector<Entry>::iterator cit;
  vector<Entry>::iterator vit;
  string info[7] = {"a","b","c","d","e","f","g"};
  int fdata[60] = {0,1,3,0,1,0,1,0,1,0,1,0,1,0,1,4,3,0,1,0,1,0,1,0,1,0,1,0,1,4,3,0,1,0,1,0,1,0,1,0,1,0,1,4,3,0,1,0,1,0,1,0,1,0,1,0,1,4,0,1};

  for ( int i = 0 ; i < 60; i++){
    tmp = Entry();
    tmp.addElement(info[fdata[i]]);
    tmp.setLC(1);
    entries.push_back(tmp);
  }
  for (vit = entries.begin() ; vit != entries.end() ; vit++){
    printf("Elem: %s\n", (*vit).toString().c_str());
  }
  printf("Before NLR()\n");
  ret = nlr(entries);
  printf("After NLR()\n");
//  printf("Top of the stack: %s\n",vent.peek().toString().c_str());
//  printf("Bottom of the stack: %s\n",vent.peek_bottom().toString().c_str());
//  printf("n-th element from top of the stack: %s\n",vent.peek_n(10).toString().c_str());
//  container = vent.peek_range(5,4);
//  for (cit = container.begin();cit != container.end();cit++){
  //  printf("el> %s\n",(*cit).toString().c_str());
  //}

  //printf("Top of the stack: %s\n",vent.peek().toString().c_str());
  //printf("Ent> %s\n",vent.pop().toString().c_str());


}
 /*
 pop_ntop(int n){ //pop and return n top elements of the stack
 pop_bottom(){ // pop and return the back of the stack
 pop_nbottom(int n){ // pop and return n bottom elements of the stack
 peek(){// return the top of the stack
 peek_bottom(){ // Return the bottom of the stack
 peek_ntop(int n){//Return top N elements of the stack
 peek_n(int n){// Return the n-th element from the top of the stack
 peek_range(int n, int m){ // Return the elements [n-th to m-th] from the top of the stack
   */
