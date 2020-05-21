package cl

import (
	"fmt"
  "trace"
	"util"
	"path"
	"strconv"
	"os"
	"os/exec"
	"sort"
)

func Convert(events []*trace.Event, obj string, bitstr string, atrmode int) (m map[int][]string, err error){
  m = make(map[int][]string)
  var objkey int // hold object key

  if !validateBitstr(bitstr, obj){
    return nil, fmt.Errorf("Conversion Failed: invalid bitsrting: %v\n",bitstr)
  }
	//fmt.Println("Len(events): ",len(events))
  processedEvents := filter(events,bitstr)
	//fmt.Println("Len(Processed Events): ",len(processedEvents))

  for _,e := range processedEvents{
    // finding key to group with
    switch obj{
    case "grtn":
      objkey = int(e.G)
    case "proc":
      objkey = e.P
    case "chan":
			if e.Type == 51 || e.Type == 52 { // event type is channel make or close
				objkey = int(e.Args[0]) // channel id in Args {cid}
			} else{ // event type is channel or send/recv
				objkey = int(e.Args[1]) // channel id in Args {eid, cid, val}
			}
    default:
      return nil, fmt.Errorf("Conversion Failed: wrong obj to group: %v\n",obj)
    }

    // append attributes to keys
    m[objkey] = append(m[objkey],getAttribute(e,atrmode))
  }

  return m, nil
}

func filter(events []*trace.Event, bitstr string) []*trace.Event{
  ret := []*trace.Event{}
	//fmt.Printf("** FILTER\n")
  for _,e := range events{
		//fmt.Printf("IN FOR\n")
    desc := EventDescriptions[e.Type]
    for i,bit := range bitstr{
			//fmt.Printf("IN SECOND FOR bitstr[%v]=%v - type: %v\n",i,strconv.QuoteRune(bit), reflect.TypeOf(strconv.QuoteRune(bit)))
			//fmt.Printf("IN SECOND FOR bitstr[%v]=RAW: %v - type: %v\n",i,bit,reflect.TypeOf(bit))
			//fmt.Printf("IN SECOND FOR bitstr[%v]=FMT: %v - type: %v\n",i,fmt.Sprintf("%b",bit),reflect.TypeOf(fmt.Sprintf("%b",bit)))
			//fmt.Printf("IN SECOND FOR bitstr[%v]=FMT: %v - type: %v\n",i,string(bit),reflect.TypeOf(string(bit)))
			//fmt.Printf("bit(type %v): %v - comp(type %v) %v, Cond: %v\n",reflect.TypeOf(strconv.QuoteRune(bit)),strconv.QuoteRune(bit),reflect.TypeOf(strconv.Itoa(1)),strconv.Itoa(1),strconv.QuoteRune(bit) == strconv.Itoa(1))
			//fmt.Printf("bit(type %v): %v - comp(type %v) %v, Cond: %v\n",reflect.TypeOf(bit),bit,reflect.TypeOf(strconv.Itoa(1)),strconv.Itoa(1),bit == strconv.Itoa(1))
			if string(bit) == "1"{
				//fmt.Printf("bitstr[%v] is enabled: %v\n",i,strconv.QuoteRune(bit))
				if util.Contains(ctgDescriptions[i].Members,desc.Name){
					ret = append(ret,e)
				}
			}
      /*if strconv.QuoteRune(bit) == "1" && util.Contains(ctgDescriptions[i].Members,desc.Name){
        ret = append(ret,e)
      }*/
    }
  }
  return ret
}

func getAttribute(e *trace.Event, atrmode int) string{
  desc := EventDescriptions[e.Type]
  if len(e.Stk) != 0{
    switch atrmode{
    case AtrMode_StkTopAll:
      return fmt.Sprintf("%v:%v:%v:%v",desc.Name,path.Base(e.Stk[len(e.Stk)-1].File),e.Stk[len(e.Stk)-1].Fn, e.Stk[len(e.Stk)-1].Line)
    case AtrMode_StkTopFlFn:
      return fmt.Sprintf("%v:%v:%v:",desc.Name,path.Base(e.Stk[len(e.Stk)-1].File),e.Stk[len(e.Stk)-1].Fn)
    case AtrMode_StkTopFnLn:
      return fmt.Sprintf("%v::%v:%v",desc.Name,e.Stk[len(e.Stk)-1].Fn,e.Stk[len(e.Stk)-1].Line)
    case AtrMode_StkTopFn:
      return fmt.Sprintf("%v::%v:",desc.Name,e.Stk[len(e.Stk)-1].Fn)
    case AtrMode_StkBotAll:
      return fmt.Sprintf("%v:%v:%v:%v",desc.Name,path.Base(e.Stk[0].File),e.Stk[0].Fn, e.Stk[0].Line)
    case AtrMode_StkBotFlFn:
      return fmt.Sprintf("%v:%v:%v:",desc.Name,path.Base(e.Stk[0].File),e.Stk[0].Fn)
    case AtrMode_StkBotFnLn:
      return fmt.Sprintf("%v::%v:%v",desc.Name,e.Stk[0].Fn, e.Stk[0].Line)
    case AtrMode_StkBotFn:
      return fmt.Sprintf("%v::%v:",desc.Name,e.Stk[len(e.Stk)-1].Fn)
    default:
      return desc.Name
    }
  }
  return desc.Name
}

func validateBitstr(bitstr, obj string) bool{
  if len(bitstr) != num_of_ctgs{
    return false
  }
  for i,b := range bitstr{
    if i > 1 && string(b) == "1" && obj == "chan"{
      return false
    }
  }
  return true
}

func bitstrTranslate(bitstr string) string{
	s := ""
	for i,b := range bitstr{
		if string(b) == "1"{
			s = s + ctgDescriptions[i].Category +"_"
		}
	}
	return s
}

func WriteContext(path , obj , bitstr  string, m map[int][]string, atrmode int){
	// path must include app
	folderName := path+"/"+obj+"-"+bitstrTranslate(bitstr)+strconv.Itoa(atrmode)+"/"
	cmd := exec.Command("mkdir","-p",folderName)
	outCmd,err := cmd.Output()
	if err != nil{
		fmt.Println(outCmd)
		panic(err)
	}
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

  for _,k := range keys{
		kk := strconv.Itoa(k)
		fmt.Println("Writing content of ",objPrefix+kk," to ",folderName+objPrefix+kk+".txt")
		f, err := os.Create(folderName+objPrefix+kk+".txt")
		if err != nil{
			panic(err)
		}
    for _,s := range m[k]{
      f.WriteString(s+"\n")
			fmt.Printf("\t%v\n",s)
    }
    f.Close()
  }


}
