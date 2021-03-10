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
	"strconv"
)

// execute sched testing
func SchedTest(app,src,x string, to,depth,iter int) *instrument.AppTest {

	var maxConcUsage int

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

	/////////////////////////////////////////////////////////////////////////////
	//  Concurrency Usage
	/////////////////////////////////////////////////////////////////////////////
	baseConcUsage := db.ConcUsageStruct(dbn)
	fmt.Println("Concurrency Usage:")
	db.DisplayConcUsageTable(baseConcUsage)
	maxConcUsage = len(baseConcUsage)
	//InitCoverageTable2(baseConcUsage)

	// initilize a table to store coverage metrics


	fmt.Println("SchedTest: Initialize new test run...")
	test := instrument.NewAppTest(base,depth)
	// obtain the base rewritten version
	test.OrigPath = filepath.Join(base.NewPath, strings.Split(filepath.Base(base.OrigPath),".")[0]+"_mod.go")

	// create a dir to store rewritten schedtest version permanently
	temp := filepath.Dir(base.OrigPath)+"/"+base.App+"/schedTests"
	test.TestPath = temp+"/"+test.Name+"_S0"
	err = os.MkdirAll(test.TestPath,os.ModePerm)
	if err != nil{
		panic(err)
	}

	log.Println("SchedTest: Rewrites permanent schedTest version: ",test.TestPath)
	// rewrite the schedTest based on the base rewritten and concusage
	err = test.RewriteSourceSched(0)
	if err != nil{
		panic(err)
	}

	// for loop:
	//    trace (execute,collect,store) the permanent version schedTest
	//    executeTrace(app.NewPath)
	//    add dbnames to test object

	//fmt.Println(test.ToString())
	//fmt.Println("SchedTest: ///////////////////////////")
	//fmt.Println("SchedTest: Testing iterations begin...")
	//fmt.Println("SchedTest: ///////////////////////////")

	var passed,failed,latest int

	for i := 0 ; i<iter ; i++ {
		//fmt.Println("SchedTest: Executing ",test.TestPath)
		events, err := instrument.ExecuteTrace(test.TestPath)
		if err != nil{
			fmt.Errorf("Error in ExecuteTrace:", err)
			return nil
		}
		//fmt.Println("SchedTest: Storing ",test.TestPath," in ", test.Name)
		dbn := db.Store(events,test.Name)
		fmt.Printf("Test run %d/%d (%s):\n",i+1,iter,dbn)
		if db.Checker(dbn,false){
			passed++
		}else{
			failed++
		}
		test.DBNames[i] = dbn

		// if concurrency usage changes, do the re-write
		testConcUsage := db.ConcUsageStruct(dbn)
		//db.DisplayConcUsageTable(testConcUsage)
		//InitCoverageTable2(testConcUsage)

		// we can find a better way to see if concusage is updated
		if len(testConcUsage) > maxConcUsage{
			fmt.Println("more concurrency usage is found, Update concurrency table!")
			maxConcUsage = len(testConcUsage)
			db.DisplayConcUsageTable(testConcUsage)
			// obtain the base rewritten version
			test.OrigPath = filepath.Join(test.TestPath,test.Name+"_s"+strconv.Itoa(latest)+"_sched.go")
			latest = i

			// create a dir to store rewritten schedtest version permanently
			temp := filepath.Dir(base.OrigPath)+"/"+base.App+"/schedTests"
			test.TestPath = temp+"/"+test.Name+"_S"+strconv.Itoa(i+1)
			err = os.MkdirAll(test.TestPath,os.ModePerm)
			if err != nil{
				panic(err)
			}

			//fmt.Println(test.ToString())
			test.ConcUsage = db.ConcUsage(dbn)
			//fmt.Println("SchedTest: New concurrency usages are added!")
			//fmt.Println(test.ToString())

			//fmt.Println("SchedTest: Rewrites permanent schedTest version: ",test.TestPath)
			//fmt.Println("SchedTest: Test Origpath: ",test.OrigPath)
			err = test.RewriteSourceSched(i+1)
			if err != nil{
				panic(err)
			}
		}
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
