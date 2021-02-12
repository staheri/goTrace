package schedtest

import (
	"bytes"
	_"errors"
	"fmt"
	_"trace"
	_"io"
	"io/ioutil"
	"os"
	"os/exec"
	"util"
	"db"
	"instrument"
	"strings"
	"path/filepath"
	"time"
	"log"
)

// execute sched testing
func SchedTest(app,src,x string, to,depth,iter int) *instrument.AppTest {

	// execute the base first
	fmt.Println("SchedTest: Execute base run...")
	base := instrument.NewAppExec(app,src,x,to)
	start := time.Now()
	dbn ,err := base.DBPointer()
	base.DBName = dbn
	if err != nil{
		panic(err)
	}

	log.Printf("[TIME %v: %v]\n","Total Base Run",time.Since(start))
	//fmt.Printf("***\n[TIME %v: %v]\n***\n","Total Base Run",time.Since(start))

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
	passed := 0
	failed := 0
	for i := 0 ; i<iter ; i++ {
		events, err := instrument.ExecuteTrace(test.TestPath)
		if err != nil{
			fmt.Errorf("Error in ExecuteTrace:", err)
			return nil
		}
		dbn := db.Store(events,test.Name)
		fmt.Printf("Test run %d/%d (%s):\n",i+1,iter,dbn)
		if db.Checker(dbn,false){
			passed++
		}else{
			failed++
		}
		test.DBNames[i] = dbn
	}
	fmt.Printf("Passed: %d\nFailed: %d\n",passed,failed)
	return test
}


// For measuring the native runtime
func NativeRun (app string) {
	// create tmpdir
	log.Println("NativeRun: Create tempdir ")
	tmpDir, err := ioutil.TempDir("", "GOAT")
	if err != nil {
		panic(err)
	}
	defer func(dir string){
		if err := os.RemoveAll(dir); err != nil {
			fmt.Println("Cannot remove temp dir:", err)
		}
	}(tmpDir)

	// copy app to tmpDir
	cmd := exec.Command("cp", app, tmpDir)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		panic(stderr.String())
	}

	// create binary file holder
	log.Println("NativeRun: Create tempbin ")
	tmpBinary, err := ioutil.TempFile("", "GOAT")
	if err != nil {
		panic(err)
	}
	// remove it after done
	defer os.Remove(tmpBinary.Name())

	// build binary
	log.Println("NativeRun: Build ",tmpBinary.Name()," in ", tmpDir)
	cmd = exec.Command("go", "build", "-o", tmpBinary.Name())
	stderr.Reset()
	cmd.Stderr = &stderr
	cmd.Dir = tmpDir
	start := time.Now()
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	log.Printf("[TIME %v: %v]\n","Native Build",time.Since(start))
	if util.MeasureTime{
		fmt.Printf("[TIME %v: %v]\n","Native Build",time.Since(start))
	}
	// run binary
	log.Println("ExecuteTrace: Run ",tmpBinary.Name())
	stderr.Reset()
	cmd = exec.Command(tmpBinary.Name())
	cmd.Stderr = &stderr
	start = time.Now()
	if err = cmd.Run(); err != nil {
		panic(err)
	}
	log.Printf("[TIME %v: %v]\n","Native Run",time.Since(start))
	if util.MeasureTime{
		fmt.Printf("***\n[TIME %v: %v]\n***\n","Native Run",time.Since(start))
	}

}
