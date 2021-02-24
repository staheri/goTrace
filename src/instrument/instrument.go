package instrument

import (
	"fmt"
	"trace"
	"io"
	_"io/ioutil"
	"os"
	_"os/exec"
	"util"
	"strconv"
	"db"
	"strings"
	"log"
	"time"
	"path/filepath"
)


// holds execution info
type AppExec struct {
	App               string // DB-compatible app name
	OrigPath          string // path to the original source
	NewPath           string // path to the temp dir
	DBName            string // final DB table name
	Source            string // source of events (native,latest,x)
	X                 string // version of execution
	Timeout           int    // timeout (for DL apps)
}

// returns a pointer to new AppExec
func NewAppExec(path,src,x string, to int) *AppExec{
	return &AppExec{
		App:       util.AppName(path),
		OrigPath:  path,
		Timeout:   to,
		Source:    src,
		X:         x,
	}
}

// holds test info
type AppTest struct {
	BaseExec          *AppExec
	Name              string
	OrigPath          string
	TestPath          string
	Depth             int
	ConcUsage         map[string]int
	DBNames           map[int]string
}

// returns a pointer to new AppTest
func NewAppTest(base *AppExec,depth int) *AppTest{
	return &AppTest{
		BaseExec:  base,
		Name:      base.App+"_D"+strconv.Itoa(depth),
		Depth:     depth,
		ConcUsage: db.ConcUsage(base.DBName),
		DBNames: make(map[int]string),
	}
}

// based on app.X, retrieve a previously stored DB
// Or trace (rewrite,execute,collect) and store a new execution
func (app *AppExec) DBPointer() (dbName string, err error){

	// retrieve from db
	if app.Source != "native" && app.Source != "schedTest"{
		log.Printf("DBPointer: %s, Source: %s, X: %s", app.App,app.Source, app.X)
		fmt.Printf("DBPointer: %s, Source: %s, X: %s", app.App,app.Source, app.X)
		return db.Ops(app.Source, app.App, app.X),nil
	}

	fmt.Println("DBPointer: Native run...")
	log.Println("DBPointer: Native run...")
	// instrument, rewrite, execute, collect, store, obtain DBname
	dbName,err = app.Trace()
	if err != nil{
		return "",err
	}
	fmt.Println("DB Name:",dbName)
	app.X = strings.Split(dbName,"X")[1]
	return dbName,nil
}

// rewrite,execute,collect
func (app *AppExec) Trace() (dbName string, err error){

	// create a dir to store rewritten mod
	temp := filepath.Dir(app.OrigPath)+"/"+app.App
	app.NewPath = temp+"/"+app.App+"_mod"
	err = os.MkdirAll(app.NewPath,os.ModePerm)
	if err != nil{
		panic(err)
	}

	// writes instrumented code into app.NewPath
	fmt.Println("Trace: Rewrite...")

	// timing start
	start := time.Now()

	err = app.RewriteSource()
	if err != nil {
		return "", fmt.Errorf("couldn't rewrite source code: %v", err)
	}

	// timing end
	end := time.Since(start)
	log.Printf("[TIME %v: %v]\n","Prime Rewrite",end)
	if util.MeasureTime{
		fmt.Printf("[TIME %v: %v]\n","Prime Rewrite",end)
	}

	// exeute, capture and parse trace
	fmt.Println("Trace: Execute...")

	// timing start
	start = time.Now()
	events, err := ExecuteTrace(app.NewPath)
	if err != nil{
		return "", fmt.Errorf("Error in ExecuteTrace:", err)
	}

	// timing end
	end = time.Since(start)
	log.Printf("[TIME %v: %v]\n","Prime ExecuteTrace",end)
	if util.MeasureTime{
		fmt.Printf("[TIME %v: %v]\n","Prime ExecuteTrace",end)
	}
	// store traces
	fmt.Println("Trace: Store...")

	// timing start
	start = time.Now()

	d := db.Store(events,app.App)

	//timing end
	end = time.Since(start)
	log.Printf("[TIME %v: %v]\n","Store",end)
	if util.MeasureTime{
		fmt.Printf("[TIME %v: %v]\n","Store",end)
	}
	return d,nil

}

// reads trace from stderr (io.reader) and parse
func parseTrace(r io.Reader, binary string) ([]*trace.Event, error) {
	parseResult, err := trace.Parse(r,binary)
	if err != nil {
		return nil, err
	}

	err = trace.Symbolize(parseResult.Events, binary)

	return parseResult.Events, err
}

// checks for errors
func check(err error){
	if err != nil{
		panic(err)
	}
}

// appExec to string
func (app *AppExec) ToString() string{
	s := fmt.Sprintf("-----------\n")
	s = s + fmt.Sprintf("App: %s\n",app.App)
	s = s + fmt.Sprintf("Orig. Path: %s\n",app.OrigPath)
	s = s + fmt.Sprintf("New Path: %s\n",app.NewPath)
	s = s + fmt.Sprintf("Timeout %d\n",app.Timeout)
	s = s + fmt.Sprintf("X %s\n",app.X)
	s = s + fmt.Sprintf("-----------\n")
	return s
}

// appTest to string
func (app *AppTest) ToString() string{
	s := fmt.Sprintf("-----------\n")
	s = s + fmt.Sprintf("Base INFO\n***\n%s\n***\n",app.BaseExec.ToString())
	s = s + fmt.Sprintf("Name: %s\n",app.Name)
	s = s + fmt.Sprintf("Orig. Path: %s\n",app.OrigPath)
	s = s + fmt.Sprintf("Test Path: %s\n",app.TestPath)
	s = s + fmt.Sprintf("Depth %d\n",app.Depth)
	s = s + fmt.Sprintf("Concurrency Usage\n")
	for k,_ := range(app.ConcUsage){
		s = s + fmt.Sprintf("\t> %s\n",k)
	}
	s = s + fmt.Sprintf("-----------\n")
	return s
}
