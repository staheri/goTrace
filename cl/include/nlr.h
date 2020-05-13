/** nlr.h
 * Implementing NLR algorithm (nested loop detection)
 * Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2019, All rights reserved
 */
#ifndef NLR_H
#define NLR_H


#include <vector>
#include <stack>
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
#include "entry.h"
#include "prep.h"

using namespace std;

#define SMAX 10000
#define NFLUSH 1000

typedef unsigned char uint1;
typedef unsigned short uint2;
typedef unsigned int uint4;
typedef unsigned long long uint8;


// Template Stack class for NLR implementation
template <class T> class VecStack {
  static vector<T> stk;
public:
  VecStack(){
    stk.clear();
  }
  ~VecStack(){}

  static vector<T> initStack(){
		vector<T> v;
		return v;
	}
  string toString(){
    string s = "[";
    typename vector<T>::iterator vit;
    for(vit = this->stk.begin(); vit != this->stk.end() ; vit++){
      s = s + (*vit).toString();
      s = s + " , ";
    }
    s = s + "]";
    return s;
  }
  int len(){
    return this->stk.size();
  }
  bool isEmpty(){
    if (this->stk.size() == 0){
      return true;
    } else{
      return false;
    }
  }
  void clear(){
    this->stk.clear();
  }
  void push(T a){
    this->stk.push_back(a);
  }
  T pop(){// pop and return the top of the stack
    T ret = this->stk.back();
    this->stk.pop_back();
    return ret;
  }
  vector<T> pop_ntop(int n){ //pop and return n top elements of the stack
    vector<T> ret;
    typename vector<T>::iterator sit;
    if ((unsigned int)n <= this->stk.size()){
      for (sit = this->stk.begin()+this->stk.size()-n ; sit != this->stk.end(); sit++){
        ret.push_back(*sit);
      }
      for(int i = 0 ; i < n ; i++){
        this->stk.pop_back();
      }
      return ret;
    } else{
      for (sit = this->stk.begin(); sit != this->stk.end(); sit++){
        ret.push_back(*sit);
      }
      this->stk.clear();
      return ret;
    }
  }
  T pop_bottom(){ // pop and return the back of the stack
    T ret= this->stk.front();
    this->stk.erase(this->stk.begin());
    return ret;
  }
  vector<T> pop_nbottom(int n){ // pop and return n bottom elements of the stack
    vector<T> ret;
    typename vector<T>::iterator sit;
    if ((unsigned int)n <= this->stk.size()){
      for (sit = this->stk.begin() ; sit != this->stk.begin()+n; sit++){
        ret.push_back(*sit);
      }
      this->stk.erase(this->stk.begin(),this->stk.begin()+n);
      return ret;
    } else{
      for (sit = this->stk.begin(); sit != this->stk.end(); sit++){
        ret.push_back(*sit);
      }
      this->stk.clear();
      return ret;
    }
  }
  T peek(){// return the top of the stack
    return this->stk.back();
  }
  T peek_bottom(){ // Return the bottom of the stack
    return this->stk.front();
  }
  vector<T> peek_ntop(int n){//Return top N elements of the stack
    vector<T> ret;
    typename vector<T>::iterator sit;
    if (n <= this->stk.size()){
      for (sit = this->stk.begin()+this->stk.size()-n ; sit != this->stk.end(); sit++){
        ret.push_back(*sit);
      }
      return ret;
    } else{
      for (sit = this->stk.begin(); sit != this->stk.end(); sit++){
        ret.push_back(*sit);
      }
      return ret;
    }
  }
  T peek_n(int n){// Return the n-th element from the top of the stack
    if ((unsigned int)n <= this->stk.size()){
      return this->stk[this->stk.size()-n];
    }
    else{
      printf("ERROR in VecSTACK (peek_n) - n is greater than stack size\n\t...Returning the top of the stack");
      return this->peek();
    }
  }
  vector<T> peek_range(int n, int m){ // Return the elements [n-th to m-th] from the top of the stack
    //printf("\t\t\t\tWithin peek_range with(%d,%d)\n",n,m);
    string s = "";
    vector<T> ret;
    if ((unsigned int)n > this->stk.size()){
      return ret;
    }
    typename vector<T>::iterator sit;
    for (sit = this->stk.begin()+this->stk.size()-n ; sit != this->stk.begin()+this->stk.size()-m+1 ; sit++){
      s = s + (*sit).toString();
      s = s + " ";
      ret.push_back(*sit);
    }
    //printf("\t\t\t\treturn %s\n",s.c_str());
    return ret;
  }
};


template <class T>
 vector<T> VecStack<T>::stk = initStack();

extern VecStack<Entry> vent;
extern int K1;

void reduce();
Entry nlr(vector<Entry> seq,int k,string outFile);
Entry triplet(vector<Entry> u,vector<Entry> v,vector<Entry> w);
bool follows(Entry u, vector<Entry> v );

Entry serialNLR(vector<Entry> &seq,int start, int end, int k, int tid);
void serialReduce(VecStack<Entry> &vent);

#endif
