package util

import (
	"strings"
	"fmt"
  "os"
  "path"
  "github.com/jedib0t/go-pretty/table"
  "sort"
  "trace"
	"strconv"
)

// If s contains e
func Contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

// Returns appName from long paths (omitting forbidden chars for database)
func AppName(app string) string{
	a := strings.Split(app,"/")
	b := strings.Split(a[len(a)-1],".")
	s := ""
	ret := ""
	for i:=0;i<len(b)-1;i++{
		if i == len(b) - 2{
			s = s + b[i]
		}else{
			s = s + b[i]+"_"
		}
	}
	for _,b := range s{
		if string(b) == "-"{
			ret = ret + "_"
		} else{
			ret = ret + string(b)
		}
	}
	return ret
}


// Display trace.Events grouped by Goroutines
func dispGTable(m map[uint64][]*trace.Event) (){
  t := table.NewWriter()
  var keys []uint64
  for k,_ := range m{
    keys = append(keys,k)
  }
  sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

  t.SetOutputMirror(os.Stdout)
  t.AppendHeader(table.Row{"GoRoutine","Seq ID", "Event", "Args","SArgs", "Caller"})
  //w := new(bytes.Buffer)
  var w string
  var prev int
  prev = -1
  for _,k := range keys{
    for _,ev := range m[k]{
      var row []interface{}
      if int(k) != prev{
        row = append(row,k)
        prev = int(k)
      } else{
        row = append(row,"")
      }
      row = append(row,ev.Ts)
      desc := EventDescriptions[ev.Type]
      row = append(row,desc.Name)
      w = ""
      for i, a := range desc.Args {
    		w = w + fmt.Sprintf(" %v=%v", a, int64(ev.Args[i]))
    	}
      row = append(row,w)

      w = ""
    	for i, a := range desc.SArgs {
    		w = w + fmt.Sprintf(" %v=%v", a, ev.SArgs[i])
    	}
      row = append(row,w)

      w = ""
      for _,a := range ev.Stk{
        w = w + fmt.Sprintf(" %v-%v:%v\n", path.Base(a.File),a.Fn, a.Line)
      }
      /*if len(ev.Stk) != 0{
    	   w = w + fmt.Sprintf(" %v-%v:%v", path.Base(ev.Stk[len(ev.Stk)-1].File),ev.Stk[len(ev.Stk)-1].Fn, ev.Stk[len(ev.Stk)-1].Line)
    	}*/
      row = append(row,w)

      t.AppendRow(row)
    }
    t.AppendSeparator()
  }
  t.Render()
}


// Display trace.Events grouped by Processes
func dispPTable(m map[int][]*trace.Event) (){
  t := table.NewWriter()
  var keys []int
  for k,_ := range m{
    keys = append(keys,k)
  }
  sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

  t.SetOutputMirror(os.Stdout)
  t.AppendHeader(table.Row{"Process", "Event", "Args","SArgs", "Caller"})
  //w := new(bytes.Buffer)
  var w string
  var prev int
  prev = -1
  for _,k := range keys{
    for _,ev := range m[k]{
      var row []interface{}
      if k != prev{
        row = append(row,k)
        prev = k
      } else{
        row = append(row,"")
      }
      desc := EventDescriptions[ev.Type]
      row = append(row,desc.Name)
      w = ""
      for i, a := range desc.Args {
    		w = w + fmt.Sprintf(" %v=%v", a, ev.Args[i])
    	}
      row = append(row,w)

      w = ""
    	for i, a := range desc.SArgs {
    		w = w + fmt.Sprintf(" %v=%v", a, ev.SArgs[i])
    	}
      row = append(row,w)

      w = ""
      for _,a := range ev.Stk{
        w = w + fmt.Sprintf(" %v-%v:%v\n", path.Base(a.File),a.Fn, a.Line)
      }
      /*if len(ev.Stk) != 0{
    	   w = w + fmt.Sprintf(" %v-%v:%v", path.Base(ev.Stk[len(ev.Stk)-1].File),ev.Stk[len(ev.Stk)-1].Fn, ev.Stk[len(ev.Stk)-1].Line)
    	}*/
      row = append(row,w)

      t.AppendRow(row)
    }
    t.AppendSeparator()
  }
  t.Render()
}

// Display CL.attributes grouped by Goroutines as objects
func DispGAttribute(m map[uint64][]*trace.Event) (){
	t := table.NewWriter()
  var keys []uint64
  for k,_ := range m{
    keys = append(keys,k)
  }
  sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

  t.SetOutputMirror(os.Stdout)
  t.AppendHeader(table.Row{"Object(GoRoutine)","Attribute"})
  //w := new(bytes.Buffer)
  //var w string
  var prev int
  prev = -1
  for _,k := range keys{
    for _,ev := range m[k]{
      var row []interface{}
      if int(k) != prev{
        row = append(row,"G"+strconv.Itoa(int(k)))
        prev = int(k)
      } else{
        row = append(row,"")
      }
      desc := EventDescriptions[ev.Type]
      row = append(row,desc.Name)
      t.AppendRow(row)
    }
    t.AppendSeparator()
  }
  t.Render()
}

// Display CL.attributes map
func DispAtrMap(m map[int][]string, obj string) {
	t := table.NewWriter()
	var objPrefix string
	switch obj{
	case "grtn":
		objPrefix = "G"
	case "proc":
		objPrefix = "P"
	case "chan":
		objPrefix = "C"
	default:
		objPrefix = "X"
	}
  var keys []int
  for k,_ := range m{
    keys = append(keys,k)
  }
  sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

  t.SetOutputMirror(os.Stdout)
  t.AppendHeader(table.Row{fmt.Sprintf("Object(%v)",obj),"Attribute"})
  //w := new(bytes.Buffer)
  //var w string
  var prev int
  prev = -1
  for _,k := range keys{
    for _,s := range m[k]{
      var row []interface{}
      if int(k) != prev{
        row = append(row,objPrefix+strconv.Itoa(int(k)))
        prev = int(k)
      } else{
        row = append(row,"")
      }
      row = append(row,s)
      t.AppendRow(row)
    }
    t.AppendSeparator()
  }
  t.Render()
}


func DispByProcs(events []*trace.Event) {
  m := make(map[int][]*trace.Event)
  for _,e := range events{
		m[e.P] = append(m[e.P],e)
  }
  dispPTable(m)
}

func DispByGrtns(events []*trace.Event) {
  m := make(map[uint64][]*trace.Event)
  for _,e := range events{
		m[e.G] = append(m[e.G],e)
  }
  dispGTable(m)
}

func AttributeModesDescription() string {
	s := "Include stack snapshots\n\t\t"
	s = s + "0: no stack\n\t\t"
	s = s + "1: Top element of stack (immediate parent) - File, Function, Line\n\t\t"
	s = s + "2: Top element of stack (immediate parent) - File, Function\n\t\t"
	s = s + "3: Top element of stack (immediate parent) - Function, Line\n\t\t"
	s = s + "4: Top element of stack (immediate parent) - Function\n\t\t"
	s = s + "5: Bottom element of stack (great ancesstor) - File, Function, Line\n\t\t"
	s = s + "6: Bottom element of stack (great ancesstor) - File, Function\n\t\t"
	s = s + "7: Bottom element of stack (great ancesstor) - Function, Line\n\t\t"
	s = s + "8: Bottom element of stack (great ancesstor) - Function\n\t\t"
	return s
}


func PrintUsage() {
	fmt.Println("Failed")
}
