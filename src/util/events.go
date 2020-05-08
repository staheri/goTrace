// by Ivan Daniluk from https://github.com/divan

package util

import (
	"bytes"
	"errors"
	"fmt"
	"trace"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

// EventSource defines anything that can
// return tracing events.
type EventSource interface {
	Events() ([]*trace.Event, error)
}

// TraceSource implements EventSource for
// pre-made trace file
type TraceSource struct {
	// Trace is the path to the trace file.
	Trace string

	// Binary is a path to binary, needed to symbolize stacks.
	// In the future it may be dropped, see
	// https://groups.google.com/d/topic/golang-dev/PGX1H8IbhFU
	Binary string
}

// NewTraceSource inits new TraceSource.
func NewTraceSource(trace, binary string) *TraceSource {
	return &TraceSource{
		Trace:  trace,
		Binary: binary,
	}
}

// Events reads trace file from filesystem and symbolizes
// stacks.
func (t *TraceSource) Events() ([]*trace.Event, error) {
	f, err := os.Open(t.Trace)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return parseTrace(f, t.Binary)
}

func parseTrace(r io.Reader, binary string) ([]*trace.Event, error) {
	parseResult, err := trace.Parse(r,binary)
	if err != nil {
		return nil, err
	}

	err = trace.Symbolize(parseResult.Events, binary)

	return parseResult.Events, err
}

// NativeRun implements EventSource for running app locally,
// using native Go installation.
type NativeRun struct {
	OrigPath, Path string
}

// NewNativeRun inits new NativeRun source.
func NewNativeRun(path string) *NativeRun {
	return &NativeRun{
		OrigPath: path,
	}
}

// Events rewrites package if needed, adding Go execution tracer
// to the main function, then builds and runs package using default Go
// installation and returns parsed events.
func (r *NativeRun) Events() ([]*trace.Event, error) {
	// rewrite AST
	err := r.RewriteSource()
	if err != nil {
		return nil, fmt.Errorf("couldn't rewrite source code: %v", err)
	}
	defer func(tmpDir string) {
		if err := os.RemoveAll(tmpDir); err != nil {
			fmt.Println("Cannot remove temp dir:", err)
		}
	}(r.Path)

	tmpBinary, err := ioutil.TempFile("", "goTMP")
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n",tmpBinary.Name())
	defer os.Remove(tmpBinary.Name())

	// build binary
	// TODO: replace build&run part with "go run" when there is no more need
	// to keep binary
	cmd := exec.Command("go", "build", "-o", tmpBinary.Name())
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Dir = r.Path
	err = cmd.Run()
	if err != nil {
		fmt.Println("go build error", stderr.String())
		// TODO: test on most common errors, possibly add stderr to
		// error information or smth.
		return nil, err
	}

	// run
	stderr.Reset()
	cmd = exec.Command(tmpBinary.Name())
	cmd.Stderr = &stderr
	if err = cmd.Run(); err != nil {
		fmt.Println("modified program failed:", err, stderr.String())
		// TODO: test on most common errors, possibly add stderr to
		// error information or smth.
		return nil, err
	}

	if stderr.Len() == 0 {
		return nil, errors.New("empty trace")
	}

	// parse trace
	return parseTrace(&stderr, tmpBinary.Name())
}

// RewriteSource attempts to add trace-related code if needed.
// TODO: add support for multiple files and package selectors
func (r *NativeRun) RewriteSource() error {
	path, err := rewriteSource(r.OrigPath)
	if err != nil {
		return err
	}

	r.Path = path

	return nil
}
