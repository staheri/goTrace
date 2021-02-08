package schedtest

import (
	_"bytes"
	_"errors"
	"fmt"
	_"trace"
	_"io"
	_"io/ioutil"
	"os"
	_"os/exec"
	_"util"
	"db"
	"instrument"
	"strings"
	"path/filepath"
)

func SchedTest(app,src,x string, to,depth,iter int) *instrument.AppTest {

	// execute the base first
	fmt.Println("SchedTest: Execute base run...")
	base := instrument.NewAppExec(app,src,x,to)
	dbn ,err := base.DBPointer()
	base.DBName = dbn
	if err != nil{
		panic(err)
	}

	fmt.Println("SchedTest: Initialize new test run...")
	test := instrument.NewAppTest(base,depth)
	// obtain the base rewritten version
	test.OrigPath = strings.Split(base.OrigPath,".go")[0]+"_mod.go"

	// create a dir to store rewritten schedtest version permanently
	temp := filepath.Dir(test.OrigPath)+"/"+base.App
	test.TestPath = temp+"/"+test.Name
	err = os.MkdirAll(test.TestPath,os.ModePerm)
	if err != nil{
		panic(err)
	}

	fmt.Println("SchedTest: Rewrites permanent schedTest version: ",test.TestPath)
	// rewrite the schedTest based on the base rewritten and concusage
	err = test.RewriteSourceSched()
	if err != nil{
		panic(err)
	}

	// for loop:
	//    trace (execute,collect,store) the permanent version schedTest
	//    executeTrace(app.NewPath)
	//    add dbnames to test object
	fmt.Println("SchedTest: Testing iterations begin...")
	for i := 0 ; i<iter ; i++ {
		events, err := instrument.ExecuteTrace(test.TestPath)
		if err != nil{
			fmt.Errorf("Error in ExecuteTrace:", err)
			return nil
		}
		dbn := db.Store(events,test.Name)
		fmt.Printf("Test run %d/%d (%s):\n",i+1,iter,dbn)
		db.Checker(dbn,false)
		test.DBNames[i] = dbn
	}
	return test
}
