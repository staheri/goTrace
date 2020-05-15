package analyze

import (
	"fmt"
  "trace"
	"util"
	_"reflect"
	"path"
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

const (
  AtrMode_StkNone           = 0 // no stack
  AtrMode_StkTopAll         = 1 // Top element of stack (immediate parent) - File, Function, Line
  AtrMode_StkTopFlFn        = 2 // Top element of stack (immediate parent) - File, Function
  AtrMode_StkTopFnLn        = 3 // Top element of stack (immediate parent) - Function, Line
  AtrMode_StkTopFn          = 4 // Top element of stack (immediate parent) - Function
  AtrMode_StkBotAll         = 5 // Bottom element of stack (great ancesstor) - File, Function, Line
  AtrMode_StkBotFlFn        = 6 // Bottom element of stack (great ancesstor) - File, Function
  AtrMode_StkBotFnLn        = 7 // Bottom element of stack (great ancesstor) - Function, Line
  AtrMode_StkBotFn          = 8 // Bottom element of stack (great ancesstor) - Function
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
// Event types in the trace.
// Verbatim copy from src/runtime/trace.go with the "trace" prefix removed.
const (
	EvNone              = 0  // unused
	EvBatch             = 1  // start of per-P batch of events [pid, timestamp]
	EvFrequency         = 2  // contains tracer timer frequency [frequency (ticks per second)]
	EvStack             = 3  // stack [stack id, number of PCs, array of {PC, func string ID, file string ID, line}]
	EvGomaxprocs        = 4  // current value of GOMAXPROCS [timestamp, GOMAXPROCS, stack id]
	EvProcStart         = 5  // start of P [timestamp, thread id]
	EvProcStop          = 6  // stop of P [timestamp]
	EvGCStart           = 7  // GC start [timestamp, seq, stack id]
	EvGCDone            = 8  // GC done [timestamp]
	EvGCSTWStart        = 9  // GC mark termination start [timestamp, kind]
	EvGCSTWDone         = 10 // GC mark termination done [timestamp]
	EvGCSweepStart      = 11 // GC sweep start [timestamp, stack id]
	EvGCSweepDone       = 12 // GC sweep done [timestamp, swept, reclaimed]
	EvGoCreate          = 13 // goroutine creation [timestamp, new goroutine id, new stack id, stack id]
	EvGoStart           = 14 // goroutine starts running [timestamp, goroutine id, seq]
	EvGoEnd             = 15 // goroutine ends [timestamp]
	EvGoStop            = 16 // goroutine stops (like in select{}) [timestamp, stack]
	EvGoSched           = 17 // goroutine calls Gosched [timestamp, stack]
	EvGoPreempt         = 18 // goroutine is preempted [timestamp, stack]
	EvGoSleep           = 19 // goroutine calls Sleep [timestamp, stack]
	EvGoBlock           = 20 // goroutine blocks [timestamp, stack]
	EvGoUnblock         = 21 // goroutine is unblocked [timestamp, goroutine id, seq, stack]
	EvGoBlockSend       = 22 // goroutine blocks on chan send [timestamp, stack]
	EvGoBlockRecv       = 23 // goroutine blocks on chan recv [timestamp, stack]
	EvGoBlockSelect     = 24 // goroutine blocks on select [timestamp, stack]
	EvGoBlockSync       = 25 // goroutine blocks on Mutex/RWMutex [timestamp, stack]
	EvGoBlockCond       = 26 // goroutine blocks on Cond [timestamp, stack]
	EvGoBlockNet        = 27 // goroutine blocks on network [timestamp, stack]
	EvGoSysCall         = 28 // syscall enter [timestamp, stack]
	EvGoSysExit         = 29 // syscall exit [timestamp, goroutine id, seq, real timestamp]
	EvGoSysBlock        = 30 // syscall blocks [timestamp]
	EvGoWaiting         = 31 // denotes that goroutine is blocked when tracing starts [timestamp, goroutine id]
	EvGoInSyscall       = 32 // denotes that goroutine is in syscall when tracing starts [timestamp, goroutine id]
	EvHeapAlloc         = 33 // memstats.heap_live change [timestamp, heap_alloc]
	EvNextGC            = 34 // memstats.next_gc change [timestamp, next_gc]
	EvTimerGoroutine    = 35 // denotes timer goroutine [timer goroutine id]
	EvFutileWakeup      = 36 // denotes that the previous wakeup of this goroutine was futile [timestamp]
	EvString            = 37 // string dictionary entry [ID, length, string]
	EvGoStartLocal      = 38 // goroutine starts running on the same P as the last event [timestamp, goroutine id]
	EvGoUnblockLocal    = 39 // goroutine is unblocked on the same P as the last event [timestamp, goroutine id, stack]
	EvGoSysExitLocal    = 40 // syscall exit on the same P as the last event [timestamp, goroutine id, real timestamp]
	EvGoStartLabel      = 41 // goroutine starts running with label [timestamp, goroutine id, seq, label string id]
	EvGoBlockGC         = 42 // goroutine blocks on GC assist [timestamp, stack]
	EvGCMarkAssistStart = 43 // GC mark assist start [timestamp, stack]
	EvGCMarkAssistDone  = 44 // GC mark assist done [timestamp]
	EvUserTaskCreate    = 45 // trace.NewContext [timestamp, internal task id, internal parent id, stack, name string]
	EvUserTaskEnd       = 46 // end of task [timestamp, internal task id, stack]
	EvUserRegion        = 47 // trace.WithRegion [timestamp, internal task id, mode(0:start, 1:end), stack, name string]
	EvUserLog           = 48 // trace.Log [timestamp, internal id, key string id, stack, value string]
	EvChSend            = 49 // goTrace: chan send [timestamp, stack, event id, channel id, value]
	EvChRecv            = 50 // goTrace: chan recv [timestamp, stack, event id, channel id, value]
	EvChMake            = 51 // goTrace: chan make [timestamp, stack, channel id]
	EvChClose           = 52 // goTrace: chan close [timestamp, stack, channel id]
	EvCount             = 53
)

var EventDescriptions = [EvCount]struct {
	Name       string
	minVersion int
	Stack      bool
	Args       []string
	SArgs      []string // string arguments
}{
	EvNone:              {"None", 1005, false, []string{}, nil},
	EvBatch:             {"Batch", 1005, false, []string{"p", "ticks"}, nil}, // in 1.5 format it was {"p", "seq", "ticks"}
	EvFrequency:         {"Frequency", 1005, false, []string{"freq"}, nil},   // in 1.5 format it was {"freq", "unused"}
	EvStack:             {"Stack", 1005, false, []string{"id", "siz"}, nil},
	EvGomaxprocs:        {"Gomaxprocs", 1005, true, []string{"procs"}, nil},
	EvProcStart:         {"ProcStart", 1005, false, []string{"thread"}, nil},
	EvProcStop:          {"ProcStop", 1005, false, []string{}, nil},
	EvGCStart:           {"GCStart", 1005, true, []string{"seq"}, nil}, // in 1.5 format it was {}
	EvGCDone:            {"GCDone", 1005, false, []string{}, nil},
	EvGCSTWStart:        {"GCSTWStart", 1005, false, []string{"kindid"}, []string{"kind"}}, // <= 1.9, args was {} (implicitly {0})
	EvGCSTWDone:         {"GCSTWDone", 1005, false, []string{}, nil},
	EvGCSweepStart:      {"GCSweepStart", 1005, true, []string{}, nil},
	EvGCSweepDone:       {"GCSweepDone", 1005, false, []string{"swept", "reclaimed"}, nil}, // before 1.9, format was {}
	EvGoCreate:          {"GoCreate", 1005, true, []string{"g", "stack"}, nil},
	EvGoStart:           {"GoStart", 1005, false, []string{"g", "seq"}, nil}, // in 1.5 format it was {"g"}
	EvGoEnd:             {"GoEnd", 1005, false, []string{}, nil},
	EvGoStop:            {"GoStop", 1005, true, []string{}, nil},
	EvGoSched:           {"GoSched", 1005, true, []string{}, nil},
	EvGoPreempt:         {"GoPreempt", 1005, true, []string{}, nil},
	EvGoSleep:           {"GoSleep", 1005, true, []string{}, nil},
	EvGoBlock:           {"GoBlock", 1005, true, []string{}, nil},
	EvGoUnblock:         {"GoUnblock", 1005, true, []string{"g", "seq"}, nil}, // in 1.5 format it was {"g"}
	EvGoBlockSend:       {"GoBlockSend", 1005, true, []string{}, nil},
	EvGoBlockRecv:       {"GoBlockRecv", 1005, true, []string{}, nil},
	EvGoBlockSelect:     {"GoBlockSelect", 1005, true, []string{}, nil},
	EvGoBlockSync:       {"GoBlockSync", 1005, true, []string{}, nil},
	EvGoBlockCond:       {"GoBlockCond", 1005, true, []string{}, nil},
	EvGoBlockNet:        {"GoBlockNet", 1005, true, []string{}, nil},
	EvGoSysCall:         {"GoSysCall", 1005, true, []string{}, nil},
	EvGoSysExit:         {"GoSysExit", 1005, false, []string{"g", "seq", "ts"}, nil},
	EvGoSysBlock:        {"GoSysBlock", 1005, false, []string{}, nil},
	EvGoWaiting:         {"GoWaiting", 1005, false, []string{"g"}, nil},
	EvGoInSyscall:       {"GoInSyscall", 1005, false, []string{"g"}, nil},
	EvHeapAlloc:         {"HeapAlloc", 1005, false, []string{"mem"}, nil},
	EvNextGC:            {"NextGC", 1005, false, []string{"mem"}, nil},
	EvTimerGoroutine:    {"TimerGoroutine", 1005, false, []string{"g"}, nil}, // in 1.5 format it was {"g", "unused"}
	EvFutileWakeup:      {"FutileWakeup", 1005, false, []string{}, nil},
	EvString:            {"String", 1007, false, []string{}, nil},
	EvGoStartLocal:      {"GoStartLocal", 1007, false, []string{"g"}, nil},
	EvGoUnblockLocal:    {"GoUnblockLocal", 1007, true, []string{"g"}, nil},
	EvGoSysExitLocal:    {"GoSysExitLocal", 1007, false, []string{"g", "ts"}, nil},
	EvGoStartLabel:      {"GoStartLabel", 1008, false, []string{"g", "seq", "labelid"}, []string{"label"}},
	EvGoBlockGC:         {"GoBlockGC", 1008, true, []string{}, nil},
	EvGCMarkAssistStart: {"GCMarkAssistStart", 1009, true, []string{}, nil},
	EvGCMarkAssistDone:  {"GCMarkAssistDone", 1009, false, []string{}, nil},
	EvUserTaskCreate:    {"UserTaskCreate", 1011, true, []string{"taskid", "pid", "typeid"}, []string{"name"}},
	EvUserTaskEnd:       {"UserTaskEnd", 1011, true, []string{"taskid"}, nil},
	EvUserRegion:        {"UserRegion", 1011, true, []string{"taskid", "mode", "typeid"}, []string{"name"}},
	EvUserLog:           {"UserLog", 1011, true, []string{"id", "keyid"}, []string{"category", "message"}},
	EvChSend:            {"ChSend", 1011, true, []string{"eid","cid","val"}, nil}, // goTrace
	EvChRecv:            {"ChRecv", 1011, true, []string{"eid","cid","val"}, nil}, // goTrace
	EvChMake:            {"ChMake", 1011, true, []string{"cid"}, nil}, // goTrace
	EvChClose:           {"ChClose", 1011, true, []string{"cid"}, nil}, // goTrace
}
