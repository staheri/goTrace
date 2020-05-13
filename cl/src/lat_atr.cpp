#include "lat_atr.h"

// SPECIALIZATION FOR INT

template <> string Attribute<int>::toCompString(){
	string s="\t---------------------";
	s = s + "\n\tAttribut ID:" + this->itos(this->id);
	//s = s + "\n\tAttribut ID:" + attrIntToString(1);
	s = s + "\n\t----------------------";
	s = s + "\n\tAttribute Name: " + this->itos(this->name);
	s = s + "\n\tattrCount: " + this->itos(this->attributeCount);
	s = s + "\n\t-----------end---------\n";
	return s;
}



template <> string Attribute<int>::tableCompString(){
	map<int,int>::iterator it;
	string s="\t--------------------- Attribute Table ---------------------";
	s = s + "\n\tAttribute ID:" + this->itos(this->id);
	s = s + "\n\t-------------------------------------------------------";
	s = s + "\n\tTrace Name | Trace ID";
	s = s + "\n\t_____________________";
	for (it = this->atable.begin() ; it != this->atable.end() ; it++){
		s = s + "\n\t" + this->itos(it->first) + " | " + this->itos(it->second) ; 
	}
	s = s + "\n\t_____________________";
	s = s + "\n\t--------------------------end--------------------------\n";
	return s;
}

 template <> string Attribute<int>::tableString(){
	map<int,int>::iterator it;
	stringstream ss;
	string s;
	
	const char separator    = ' ';
    const int nameWidth     = 20;
    const int numWidth      = 5;
	
	ss << "--------------------------" << endl;
    ss << left << setw(numWidth) << setfill(separator) << "ID" << "|";
    ss << left << setw(nameWidth) << setfill(separator) << "Attribute(int)" << "|" << endl;
	ss << "--------------------------" << endl;
    s = s + ss.str();
	
	for (it = this->atable.begin() ; it != this->atable.end() ; it++){
		stringstream sss;
		sss << left << setw(numWidth) << setfill(separator) << it->second << "|";
		sss << left << setw(nameWidth) << setfill(separator) << it->first << "|" << endl;
		s = s + sss.str(); 
	}
	s = s + "--------------------------\n";
	return s;
}






// SPECIALIZATION FOR STRING


template <> string Attribute<string>::toCompString(){
	string s="\t---------------------";
	s = s + "\n\tAttribut ID:" + this->itos(this->id);
	s = s + "\n\t----------------------";
	s = s + "\n\tAttribute Name: " + this->name;
	s = s + "\n\tattrCount: " + this->itos(this->attributeCount);
	s = s + "\n\t-----------end---------\n";
	return s;
}

template <> string Attribute<string>::tableCompString(){
	map<string,int>::iterator it;
	string s="\t***********************************************************\n";
	s = s + "\t*                   Attribute Table                       *\n";
	s = s + "\t***********************************************************\n";
	s = s + "\n\tAttribute ID:" + this->itos(this->id);
	s = s + "\n\t-------------------------------------------------------";
	s = s + "\n\tTrace Name | Trace ID";
	s = s + "\n\t_____________________";
	for (it = this->atable.begin() ; it != this->atable.end() ; it++){
		s = s + "\n\t" + it->first + " | " + this->itos(it->second) ; 
	}
	s = s + "\n\t_____________________";
	s = s + "\n\t--------------------------end--------------------------\n";
	return s;
}

template <> string Attribute<string>::tableString(){
	map<string,int>::iterator it;
	stringstream ss;
	string s;
	
	const char separator    = ' ';
    const int nameWidth     = 30;
    const int numWidth      = 5;
	
	ss << "------------------------------------" << endl;
    ss << left << setw(numWidth) << setfill(separator) << "ID" << "|";
    ss << left << setw(nameWidth) << setfill(separator) << "Attribute(str)" << "|" << endl;
	ss << "------------------------------------" << endl;
    s = s + ss.str();
	
	for (it = this->atable.begin() ; it != this->atable.end() ; it++){
		stringstream sss;
		sss << left << setw(numWidth) << setfill(separator) << it->second << "|";
		sss << left << setw(nameWidth) << setfill(separator) << it->first << "|" << endl;
		s = s + sss.str(); 
	}
	s = s + "------------------------------------\n";
	return s;
}





