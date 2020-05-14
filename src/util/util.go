package util

import (
	_"fmt"
  _"sort"
  "trace"
	_"strconv"
)

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
    /*if _,ok := m[e.P]; ok{
      m[e.P] = append(m[e.P],e)
    } else{
      m[e.P] = append(m[e.P],e)
    }*/
		m[e.P] = append(m[e.P],e)
  }
  DispPTable(m)
}

func GroupGrtns(events []*trace.Event) {
  m := make(map[uint64][]*trace.Event)
  for _,e := range events{
    if _,ok := m[e.G]; ok{
      m[e.G] = append(m[e.G],e)
    } else{
      m[e.G] = append(m[e.G],e)
    }
  }
  DispGTable(m)
}
