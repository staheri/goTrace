package schedtest

import (
	_"errors"
	"fmt"
	_"trace"
	_"io"
	"db"
	"os"
	"strings"
	"strconv"
	"github.com/jedib0t/go-pretty/table"
	"sort"
)

type CoverageID struct{
	f          string // function
	res        string
	//event      string
}

type Resources struct{
	Mutexes    map[string]string
	Channels   map[string]string
	Waitgroups map[string]string
	Condvars   map[string]string
	Selects    map[string]string
	Gs         map[string]string
	//events     []string
}





func InitCoverageTable(cu []*db.ConUse){
	// take the concusage slice
	// create a data structure
	sort.Slice(cu,func(i,j int) bool {return cu[i].Line < cu[j].Line })
	resMap   := make(map[string]map[string]string)
	resTable := make(map[string][]string)
	//eventMap := make(map[CoverageID]string)
	locMap := make(map[CoverageID]string)
	var funcs  []string
	var res,ev string
	for _,c := range(cu){
		ev = c.Event[4:]
		f := c.Funct
		//cnt := 0
		//if val,ok := resTable[f],
		if val,ok := resMap[f];ok{
			if v,ok2 := val[c.Rid];ok2{
				res = v
				if strings.HasPrefix(res,"G"){
					ev=c.Event[2:]
				}
			}else{
				if strings.HasPrefix(c.Rid,"M"){
					res = "Mu"
				} else if strings.HasPrefix(c.Rid,"CV"){
					res = "CV"
				} else if strings.HasPrefix(c.Rid,"C"){
					res = "Ch"
				}else if strings.HasPrefix(c.Rid,"W"){
					res = "WG"
				} else{
					res = "G"
					ev=c.Event[2:]
				}
				res = res + strconv.Itoa(len(val)+1)
				val[c.Rid] = res
			}
		}else{
			if strings.HasPrefix(c.Rid,"M"){
				res = "Mu1"
			} else if strings.HasPrefix(c.Rid,"CV"){
				res = "CV1"
			} else if strings.HasPrefix(c.Rid,"C"){
				res = "Ch1"
			}else if strings.HasPrefix(c.Rid,"W"){
				res = "WG1"
			} else{
				res = "G1"
				ev=c.Event[2:]
			}
			mm := make(map[string]string)
			mm[c.Rid]=res
			resMap[f]=mm
		}
		res = res + "("+ev+")"
		resTable[f] = append(resTable[f],res)
		//eventMap[CoverageID{f,res}]=c.Event
		locMap[CoverageID{f,res}]=c.File+":"+c.Line
	}

	for k,_ := range(resTable){
		funcs = append(funcs,k)
	}

	sort.Strings(funcs)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Function","Resources (Event)","Loc"})
	for _,fu := range(funcs){
		k := fu
		v := resTable[fu]
		for i,el := range(v){
			var row []interface{}
			if i == 0{
				row = append(row,k)
			} else{
				row = append(row,"")
			}
			row = append(row,el)
			//row = append(row,eventMap[CoverageID{k,el}])
			row = append(row,locMap[CoverageID{k,el}])
			t.AppendRow(row)
		}
		t.AppendSeparator()

	}
	t.Render()
}

func InitCoverageTable2(cu []*db.ConUse){
	// take the concusage slice
	// create a data structure
	sort.Slice(cu,func(i,j int) bool {return cu[i].Line < cu[j].Line })
	resMap   := make(map[string]*Resources)
	resTable := make(map[string][]string)
	//eventMap := make(map[CoverageID]string)
	locMap := make(map[CoverageID]string)
	var funcs  []string
	var res,ev string
	for _,c := range(cu){
		fmt.Println(">>>>>>",c.Rid,c.Funct,c.Event)
		ev = c.Event[4:]
		f := c.Funct

		if _,ok := resMap[f];!ok{
			// seeing the function for the first time
			// lets make a resource for it
			r := new(Resources)
			r.Mutexes   = make(map[string]string)
			r.Channels  = make(map[string]string)
			r.Waitgroups= make(map[string]string)
			r.Condvars  = make(map[string]string)
			r.Selects   = make(map[string]string)
			r.Gs        = make(map[string]string)
			resMap[f] = r
		}
		rs := resMap[f]
		if strings.HasPrefix(c.Rid,"G"){
			ev = c.Event[2:]
			if _,ok2 := rs.Gs[c.Rid];!ok2{
				res = "G"+strconv.Itoa(len(rs.Gs)+1)
				rs.Gs[c.Rid]=res
			} else{
				res = rs.Gs[c.Rid]
			}
		} else if strings.HasPrefix(c.Rid,"M"){
			if _,ok2 := rs.Mutexes[c.Rid];!ok2{
				res = "Mu"+strconv.Itoa(len(rs.Mutexes)+1)
				rs.Mutexes[c.Rid]=res
			}else{
				res = rs.Mutexes[c.Rid]
			}
		} else if strings.HasPrefix(c.Rid,"CV"){
			if _,ok2 := rs.Condvars[c.Rid];!ok2{
				res = "CV"+strconv.Itoa(len(rs.Condvars)+1)
				rs.Condvars[c.Rid]=res
			}else{
				res = rs.Condvars[c.Rid]
			}
		} else if strings.HasPrefix(c.Rid,"C"){
			if _,ok2 := rs.Channels[c.Rid];!ok2{
				res = "Ch"+strconv.Itoa(len(rs.Channels)+1)
				rs.Channels[c.Rid]=res
			}else{
				res = rs.Channels[c.Rid]
			}
		} else if strings.HasPrefix(c.Rid,"W"){
			if _,ok2 := rs.Waitgroups[c.Rid];!ok2{
				res = "Wg"+strconv.Itoa(len(rs.Waitgroups)+1)
				rs.Waitgroups[c.Rid]=res
			}else{
				res = rs.Waitgroups[c.Rid]
			}
		} else if strings.HasPrefix(c.Rid,"S"){
			ev = c.Event[2:]
			if _,ok2 := rs.Selects[c.Rid];!ok2{
				res = "Sel"+strconv.Itoa(len(rs.Selects)+1)
				rs.Selects[c.Rid]=res
			}else{
				res = rs.Selects[c.Rid]
			}
		}
		resMap[f] = rs
		res = res + "("+ev+")"+c.Rid
		resTable[f] = append(resTable[f],res)
		//eventMap[CoverageID{f,res}]=c.Event
		locMap[CoverageID{f,res}]=c.File+":"+c.Line
	}

	for k,_ := range(resTable){
		funcs = append(funcs,k)
	}

	sort.Strings(funcs)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Function","Resources (Event)","Loc"})
	for _,fu := range(funcs){
		k := fu
		v := resTable[fu]
		for i,el := range(v){
			var row []interface{}
			if i == 0{
				row = append(row,k)
			} else{
				row = append(row,"")
			}
			row = append(row,el)
			//row = append(row,eventMap[CoverageID{k,el}])
			row = append(row,locMap[CoverageID{k,el}])
			t.AppendRow(row)
		}
		t.AppendSeparator()

	}
	t.Render()
}
