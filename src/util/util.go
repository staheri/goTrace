package util

import (
	"fmt"
  "sort"
  "trace"
	"strconv"
	"os"
	"os/exec"
	"strings"
)

const num_of_ctgs = 6
const num_of_atrmodes = 9

var ctgDescriptions = [num_of_ctgs]struct {
	Category      string
	Members    []string
}{
	0:  {"GRTN", []string{"EvGoCreate","EvGoStart","EvGoEnd","EvGoStop","EvGoSched","EvGoPreempt","EvGoSleep","EvGoBlock","EvGoUnblock","EvGoBlockSend","EvGoBlockRecv","EvGoBlockSelect","EvGoBlockSync","EvGoBlockCond","EvGoBlockNet","EvGoWaiting","EvGoInSyscall","EvGoStartLocal","EvGoUnblockLocal","EvGoSysExitLocal","EvGoStartLabel","EvGoBlockGC"}},
  1:  {"CHNL",[]string{"EvChSend","EvChRecv","EvChMake","EvChClose"}},
  2:  {"PROC",[]string{"EvNone","EvBatch","EvFrequency","EvStack","EvGomaxprocs","EvProcStart","EvProcStop"}},
  3:  {"GCMM",[]string{"EvGCStart","EvGCDone","EvGCSTWStart","EvGCSTWDone","EvGCSweepStart","EvGCSweepDone","EvHeapAlloc","EvNextGC","EvGCMarkAssistStart","EvGCMarkAssistDone"}},
  4:  {"SYSC",[]string{"EvGoSysCall","EvGoSysExit","EvGoSysBlock"}},
  5:  {"MISC",[]string{"EvUserTaskCreate","EvUserTaskEnd","EvUserRegion","EvUserLog","EvTimerGoroutine","EvFutileWakeup","EvString"}},
}

func Contains(s []string, e string) bool {
    for _, a := range s {
        if a == "Ev"+e {
            return true
        }
    }
    return false
}


func GroupProcs(events []*trace.Event) {
  m := make(map[int][]*trace.Event)
  for _,e := range events{
		m[e.P] = append(m[e.P],e)
  }
  DispPTable(m)
}

func GroupGrtns(events []*trace.Event) {
  m := make(map[uint64][]*trace.Event)
  for _,e := range events{
		m[e.G] = append(m[e.G],e)
  }
  DispGTable(m)
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

func bitstrTranslate(bitstr string) string{
	s := ""
	for i,b := range bitstr{
		if string(b) == "1"{
			s = s + ctgDescriptions[i].Category +"_"
		}
	}
	return s
}

func AppName(app string) string{
	a := strings.Split(app,"/")
	b := strings.Split(a[len(a)-1],".")
	s := ""
	for i:=0;i<len(b)-1;i++{
		if i == len(b) - 2{
			s = s + b[i]
		}else{
			s = s + b[i]+"-"
		}
	}
	return s
}
