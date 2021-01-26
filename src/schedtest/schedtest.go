package schedtest

import (
	_"bytes"
	_"errors"
	"fmt"
	_"trace"
	_"io"
	_"io/ioutil"
	_"os"
	_"os/exec"
	_"util"
	"db"
)

func SchedTest(dbName string, depth int) {
	// find concurrency usage
	concurrencyUsage := db.ConcUsage(dbName)

	for k,v := range(concurrencyUsage){
		fmt.Println(k)
		fmt.Println(v)
	}
	// instrument the app based on the usage
	// execute/store the app
	// check for error (either here or back in main)
}
