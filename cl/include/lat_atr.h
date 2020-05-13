#ifndef LAT_ATR_H
#define LAT_ATR_H

#include <iostream>
#include <vector>
#include <map>
#include <vector>
#include <set>
#include <string>
#include <string.h>
#include <stdio.h>
#include <algorithm>
#include <fstream>
#include <sstream>
#include <iterator>
#include <queue>
#include <typeinfo>

#include "util.h"

using namespace std;


//! Purpose: Storing Attribute information
/*!
  This class is for instanciating new attributes from trace files or
  any other format of input. The objects from this class then will
  be used to be fed to Concept constructor to create new concepts
  and add it to the lattice
*/

template <class T>
class Attribute{
	int id; /*!< Unique id of each attribute*/
	T name; /*!< <T> Value of each attribute*/
	static int attributeCount; /*!< Total number of unique attributes*/
	static map<T,int> atable; /*!< A table to store attribute values with their key*/
public:

	//! Copy Constructor
	Attribute(){}

	//! Destructor
	~Attribute(){}

	//! Default Constructor
	/*!
	  It takes an attribute of type T and add it to the table to recieve a unique ID
	*/
	Attribute(T attr){
		this->name = attr;
		this->id = addToTable(attr);
	}

	//! Add attribute to table to receive a unique ID for that table
	int addToTable(T attr){
		typename map<T,int>::iterator it ;
		it = this->atable.find(attr);
		if (it == this->atable.end() ){
			this->attributeCount++;
			this->atable[attr] = this->attributeCount;
		}
		//printf(">inside Attr Const\n>>>ID:%d\n",this->atable[attr]);
		return this->atable[attr];
	}

	//! Return the ID of current attribute
	int getID(){
		return this->id;
	}

	//! Return current attribute count
	int getAttributeCount(){
		return this->attributeCount;
	}

	//! Initialization of the Attribute Table
	static map<T,int> initATable(){
		map<T,int> tt;
		tt = {};
		return tt;
	}

	//! Convert integer to string
	string itos(int i){
		string s = "";
		std::ostringstream oss;
		oss << i;
		s += oss.str();
		return s;
	}

	//! Returns an string-style representation of current Attribute Object
	string toString();

	//! Returns a Component-string-style representation of current Attribute Object
	string toCompString();

	//! Returns an string-style representation of current Attribute Table
	string tableCompString();

	//! Returns Attribute hashtable
	string tableString();
};




/**
 * Initializing static counter
 */
template <class T>
int Attribute<T>::attributeCount = 0;

/**
 * Initilize static Attribute Table
 */
template <class T>
map<T,int> Attribute<T>::atable = initATable();



#endif
