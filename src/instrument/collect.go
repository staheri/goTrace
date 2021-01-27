package instrument

import (
	"bytes"
	"errors"
	"fmt"
	"trace"
	"io/ioutil"
	"os"
	"os/exec"
)

// - Compile and executes the modified source
// - Parse collected trace
func ExecuteTrace(path string) ([]*trace.Event, error){
	// create binary file holder
	tmpBinary, err := ioutil.TempFile("", "GOAT")
	if err != nil {
		//fmt.Println("Error creating binary file:",err)
		return nil, fmt.Errorf("Error creating binary file:",err)
	}

	// remove it after done
	defer os.Remove(tmpBinary.Name())

	// build binary
	cmd := exec.Command("go", "build", "-o", tmpBinary.Name())
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Dir = path
	err = cmd.Run()
	if err != nil {
		fmt.Println("go build error", stderr.String())
		return nil, err
	}

	// run
	stderr.Reset()
	cmd = exec.Command(tmpBinary.Name())
	cmd.Stderr = &stderr
	if err = cmd.Run(); err != nil {
		fmt.Println("modified program failed:", err, stderr.String())
		return nil, err
	}

	// check length of stderr
	if stderr.Len() == 0 {
		return nil, errors.New("empty trace")
	}

	// parse
	return parseTrace(&stderr, tmpBinary.Name())
}

// removes dir
func removeDir(dir string) {
	if err := os.RemoveAll(dir); err != nil {
		fmt.Println("Cannot remove temp dir:", err)
	}
}
