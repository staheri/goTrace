package util

import (
	"strings"
)

func Contains(s []string, e string) bool {
    for _, a := range s {
        if a == "Ev"+e {
            return true
        }
    }
    return false
}

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
