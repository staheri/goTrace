/** nlr.cpp
 * Implementing NLR algorithm (nested loop detection)
 * Author: Saeed Taheri, University of Utah, staheri@cs.utah.edu, 2019, All rights reserved
 */
#include "nlr.h"

VecStack<Entry> vent;
int K1;


//Entry serialNLR(vector<Entry> seq, int k){
Entry serialNLR(vector<Entry> &seq,int start,int end, int k, int tid){
  VecStack<Entry> serialVent;
  int counter=0;
  int cnt = 0;
  Entry ret;
  vector<Entry> rett;
  vector<Entry>::iterator vit;
  vector<Entry>::iterator vit2;
  vector<Entry>::const_iterator first = seq.begin() + start;
  vector<Entry>::const_iterator last = seq.begin() + end;
  vector<Entry> newVec(first, last);
  //vent.clear();
  K1=k;

  // /printf("Thread %d: Serial NLR - Input Length :%lu\n",id,seq.size());
  printf("\tNLR: before first loop TID: %d, Start: %d, End: %d, VENT SIZE: %d\n",tid,start,end,serialVent.len());
  serialVent.push(newVec[0]);
  serialReduce(serialVent);
  serialVent.push(newVec[1]);
  serialReduce(serialVent);
  serialVent.push(newVec[2]);
  serialReduce(serialVent);
  // serialVent.push(newVec[3]);
  // serialReduce(serialVent);
  // serialVent.push(newVec[4]);
  // serialVent.push(newVec[5]);
  // serialVent.push(newVec[6]);
  // serialVent.push(newVec[7]);
  // serialVent.push(newVec[8]);
  // serialVent.push(newVec[9]);
  // serialVent.push(newVec[10]);
  // serialVent.push(newVec[11]);
  // serialVent.push(newVec[12]);
  // serialReduce(serialVent);
  printf("\tNLR: before first loop TID: %d, Start: %d, End: %d, VENT SIZE: %d\n",tid,start,end,serialVent.len());
  // for(vit = newVec.begin() ; vit != newVec.end() ; vit++){
  //   serialVent.push(*vit);
  //   //serialReduce(serialVent);
  // }
  //printf("\tNLR: before first loop TID: %d, Start: %d, End: %d, VENT SIZE: %d\n",tid,start,end,serialVent.len());
  /*for(vit = newVec.begin() ; vit != newVec.end() ; vit++){
    // if (cnt % 10000 == 0){
    //   printf("Inserting element %d/%lu\n",cnt,seq.size());
    // }
    serialVent.push(*vit);
    //printf("\t before reduce\n");
    serialReduce(serialVent);
    //printf("\t After reduce\n");
    if (serialVent.len() > SMAX){
      ret = serialVent.peek();
      // FLUSH N elements from bottom of the stack
      //printf("\tFLUSHINGggg\n");
      rett = serialVent.pop_nbottom(NFLUSH);
      counter+=rett.size();
    //  for (vit2 = rett.begin() ; vit2 != rett.end() ; vit2++){
        //fo << (*vit2).toString() << endl;
        //printf("%s\n",(*vit2).toString().c_str() );
        //counter++;
    //  }
      rett.clear();
    }
    else{
      //printf("\tPEEEEKING\n");
      ret = serialVent.peek();
    }
    //printf("\tContinue\n");
    cnt++;
  }
  //printf("\t>>>After first loop NLR()\n");
  while (! serialVent.isEmpty()){
    ret = serialVent.peek();
    //fo << (serialVent.pop_bottom()).toString() << endl;
    //printf("%s\n",(serialVent.pop_bottom()).toString().c_str() );
    serialVent.pop_bottom();
    counter++;
  }
  //fo.close();
  ret.setLdataLen(counter);*/
  return ret;
}

Entry nlr(vector<Entry> seq, int k,string outFile){
  int counter=0;
  int cnt = 0;
  Entry ret;
  vector<Entry> rett;
  vector<Entry>::iterator vit;
  vector<Entry>::iterator vit2;
  vent.clear();
  K1=k;

  //write to file
  ofstream fo(outFile);
  if (!fo.is_open()){
    printf("NLR:I/O Error\nUnable to open/write %s\n",(outFile).c_str());
  }
  printf("NLR - Input Length :%lu\n",seq.size());
  //printf("\t before first loop NLR()\n");
  for(vit = seq.begin() ; vit != seq.end() ; vit++){
    if (cnt % 10000 == 0){
      printf("Inserting element %d/%lu\n",cnt,seq.size());
    }
    vent.push(*vit);
    //printf("\t before reduce\n");
    reduce();
    if (vent.len() > SMAX){
      ret = vent.peek();
      // FLUSH N elements from bottom of the stack
      rett = vent.pop_nbottom(NFLUSH);
      for (vit2 = rett.begin() ; vit2 != rett.end() ; vit2++){
        fo << (*vit2).toLString() << endl;
        counter++;
      }
      rett.clear();
    }
    else{
      ret = vent.peek();
    }
    cnt++;
  }
  while (! vent.isEmpty()){
    ret = vent.peek();
    fo << (vent.pop_bottom()).toLString() << endl;
    counter++;
  }
  fo.close();
  ret.setLdataLen(counter);
  return ret;
}

void reduce(){
  // global maxLC

  // assert stack is not empty
  //printf("\t\t Within reduce\n");

  Entry tmp;
  int i,b;
  for(i = 2; i < 3*K1 ; i++){
    //printf("\t\t K LOOP i=%d\n",i);
    if(i%3 ==0){
      b = i/3 ;
      //printf("\t\t\tbefore triplet\n");
      //printf("\t\t current Stack: %s\n",vent.toString().c_str());
      //Entry tmp = triplet(vent.peek_range(3*b,2*b+1),vent.peek_range(2*b,b+1),vent.peek_range(b,1));
      tmp = triplet(vent.peek_range(3*b,2*b+1),vent.peek_range(2*b,b+1),vent.peek_range(b,1));
      //printf("\t\t\tafter triplet\n");
      if (tmp.getLC() != 0){
        vent.pop_ntop(i);
        vent.push(tmp);
        reduce();
        return;
      }
    }
    //printf("\t\t After first if\n");
    if (i > vent.len()){
      //printf("\t\t Here?\n");
      return;
    } else if (i <= K1+1){
      //printf("\t\t or Here?\n");
      //printf("\t\t current Stack: %s\n",vent.toString().c_str());
      if (vent.peek_n(i).getLC() != 1 && follows(vent.peek_n(i),vent.peek_range(i-1,1))){
        //printf("\t\t or even Here?\n");
        vent.pop_ntop(i-1);
        Entry tmp1 = vent.pop();
        tmp1.incLC();
        tmp1.setMaxLC(tmp1.getLC());
        vent.push(tmp1);
        reduce();
        return;
      }
    }
  }
}

bool follows(Entry u, vector<Entry> v ){
  // assert len(u.elements) == 1
  // elements of u == [x for x in elements of v]
  //printf("\tIN FOLLOWS\n");
  vector<int> velemenets;
  vector<int> vtmp;
  vector<int>::iterator vit;
  vector<Entry>::iterator vec_it;
  //printf("\tbefore first for\n");
  for (vec_it = v.begin() ; vec_it != v.end() ; vec_it++){
    //printf("\t\tbefore second for %s\n",(*vec_it).toString().c_str());
    vtmp = (*vec_it).getElements();
    //printf("\t\tbefore second for %s , len: %d\n",(*vec_it).toString().c_str(),vtmp.size());
    for (vit = vtmp.begin() ; vit != vtmp.end() ; vit++){
      //printf("\t\t\twithin second for %d\n",*vit);
      velemenets.push_back(*vit);
    }
  }
  //printf("\tafter fors\n");
  if (velemenets.size() != (unsigned int)u.getElementLen()){
    return false;
  }
  for (int i = 0 ; i < u.getElementLen() ; i++){
    if (u.getElements()[i] != velemenets[i]){
      return false;
    }
  }
  return true;
}

Entry triplet(vector<Entry> u,vector<Entry> v,vector<Entry> w){
  //printf("\t\t\t\twithin triplet\n");
  Entry ret = Entry();
  ret.setLC(0);
  unsigned int i;
  vector<int> newEntryElements;
  //<int>::iterator vit;
  string s,tmp;
  //typename vector<Entry>::iterator vit;
  //printf("\t\t\t\tBefore check\n");
  if (! (u.size() == v.size() && u.size() == w.size() )){
    return ret;
  } else{
    for (i = 0 ; i < u.size() ; i++){
      if (! (u[i] == v[i] && u[i] == w[i] ) ){
        return ret;
      }
    }
  }
  //printf("\t\t\t\tAfter check\n");
  ret.setLC(3);
  ret.setMaxLC(3);
  ret.setMaxLoopBody(u.size());
  s = "";
  for (i = 0 ; i < u.size() ; i++){
    // add element.toString to new entry elements (consequently to dtab)
    tmp = u[i].toString();
    ret.addElement(tmp);
    // add element.toString to S to be added to ltab
    s = s + u[i].toString();
    if (i < u.size() -1 ){
      s = s + " - ";
    }
  }
  ret.addToLtab(s);
  return ret;
}

void serialReduce(VecStack<Entry> &vent2){
  // global maxLC

  // assert stack is not empty
  //printf("\t\t Within reduce\n");

  Entry tmp;
  int i,b;
  for(i = 2; i < 3*K1 ; i++){
    //printf("\t\t K LOOP i=%d\n",i);
    if(i%3 ==0){
      b = i/3 ;
      //printf("\t\t\tbefore triplet\n");
      //printf("\t\t current Stack: %s\n",vent.toString().c_str());
      //Entry tmp = triplet(vent.peek_range(3*b,2*b+1),vent.peek_range(2*b,b+1),vent.peek_range(b,1));
      tmp = triplet(vent2.peek_range(3*b,2*b+1),vent2.peek_range(2*b,b+1),vent2.peek_range(b,1));
      //printf("\t\t\tafter triplet\n");
      if (tmp.getLC() != 0){
        vent2.pop_ntop(i);
        vent2.push(tmp);
        serialReduce(vent2);
        return;
      }
    }
    //printf("\t\t After first if\n");
    if (i > vent2.len()){
      //printf("\t\t Here?\n");
      return;
    } else if (i <= K1+1){
      //printf("\t\t or Here?\n");
      //printf("\t\t current Stack: %s\n",vent.toString().c_str());
      if (vent2.peek_n(i).getLC() != 1 && follows(vent2.peek_n(i),vent2.peek_range(i-1,1))){
        //printf("\t\t or even Here?\n");
        vent.pop_ntop(i-1);
        Entry tmp1 = vent2.pop();
        tmp1.incLC();
        tmp1.setMaxLC(tmp1.getLC());
        vent.push(tmp1);
        serialReduce(vent2);
        return;
      }
    }
  }
}
