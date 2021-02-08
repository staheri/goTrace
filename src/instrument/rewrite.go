// adapted from code by Ivan Daniluk from https://github.com/divan

package instrument

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/loader"
	"log"
	"strconv"
	_"reflect"
	"fmt"
	"strings"
)

// injects and rewrite source in app.OrigPath
// rewritten file(s) are stored in app.NewPath
// nothing returns
func (app *AppExec) RewriteSource() error {
	var data []byte
	var err error
	log.Println("RewriteSource: Add code to ",app.OrigPath)
	data, err = addCode(app.OrigPath,app.Timeout)
	if err != nil {
		fmt.Println("Error in addCode:", err)
		return err
	}
	// create files to store rewritten data
	filename := filepath.Join(app.NewPath, filepath.Base(app.OrigPath))
	toStore := filepath.Join(filepath.Dir(app.OrigPath), strings.Split(filepath.Base(app.OrigPath),".")[0]+"_mod.go")

	// write files
	err = ioutil.WriteFile(toStore, data, 0666)
	log.Println("RewriteSource: Writes data to ",toStore)
	if err != nil {
		panic(err)
		return err
	}
	err = ioutil.WriteFile(filename, data, 0666)
	log.Println("RewriteSource: Writes data to ",filename)
	if err != nil {
		panic(err)
		return err
	}
	return nil
}

// injects and rewrite source in app.OrigPath
// rewritten file(s) are stored in app.NewPath
// nothing returns
func (app *AppTest) RewriteSourceSched() error {
	var data []byte
	var err error
	log.Println("RewriteSourceSched: Add sched code to ",app.OrigPath)
	data, err = addCodeSched(app.OrigPath,app.Depth,app.ConcUsage)
	if err != nil {
		panic(err)
		return err
	}
	// create files to store rewritten data
	filename := filepath.Join(app.TestPath, app.Name+"_sched.go")


	err = ioutil.WriteFile(filename, data, 0666)
	log.Println("writes data to ",filename)
	if err != nil {
		panic(err)
		return err
	}
	return nil
}

// This function:
//    - traverses the AST
//    - finds the main package, file, function
//    - adds needed imports
//    - adds tracing mechanism (start/stop)
//    - adds constant depth, struct type, global counter and Reschedule function declaration
func addCode(path string, timeout int) ([]byte, error) {
	var conf loader.Config
	if _, err := conf.FromArgs([]string{path}, false); err != nil {
		return nil, err
	}

	prog, err := conf.Load()
	if err != nil {
		return nil, err
	}

	pkg := prog.Created[0]

	// TODO: find file with main func inside
	astFile := pkg.Files[0]

	// add imports
	log.Println("AddCode: Add imports")
	astutil.AddImport(prog.Fset, astFile, "os")
	astutil.AddImport(prog.Fset, astFile, "runtime/trace")
	astutil.AddImport(prog.Fset, astFile, "time")
	if timeout > 0{
		astutil.AddNamedImport(prog.Fset, astFile, "_", "net")
	}

	// add start/stop code
	log.Println("AddCode: Add trace")
	ast.Inspect(astFile, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			// find 'main' function
			if x.Name.Name == "main" && x.Recv == nil {
				stmts := createTraceStmts(timeout)
				stmts = append(stmts, x.Body.List...)
				x.Body.List = stmts
				return true
			}
		}
		return true
	})

	var buf bytes.Buffer
	err = printer.Fprint(&buf, prog.Fset, astFile)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// This function:
//    - traverses the AST
//    - finds the main package, file, function
//    - adds needed imports
//    - adds tracing mechanism (start/stop)
//    - adds constant depth, struct type, global counter and Reschedule function declaration
func addCodeSched(path string,depth int,concUsage map[string]int) ([]byte, error) {

	var conf loader.Config
	if _, err := conf.FromArgs([]string{path}, false); err != nil {
		return nil, err
	}
	prog, err := conf.Load()
	if err != nil {
		return nil, err
	}

	pkg := prog.Created[0]

	// TODO: find file with main func inside
	astFile := pkg.Files[0]

	// add imports
	log.Println("AddCodeSched: Add imports")
	astutil.AddImport(prog.Fset, astFile, "sync")
	astutil.AddImport(prog.Fset, astFile, "runtime")
	astutil.AddImport(prog.Fset, astFile, "math/rand")
	//fmt.Println(" >>> Added Imports")

	// add constant, struct type, global counter, function declration
	log.Println("AddCodeSched: Add constants, struct, flobal counter, func. decl.")
	ast.Inspect(astFile, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.File:
			// add constant, struct type, global counter, function declration
			decls := newDecls(depth)
			decls2 := x.Decls
			decls = append(decls2,decls...)
			x.Decls = decls
			return true
		}
		return true
	})

	log.Println("AddCodeSched: Add Gosched() invocations")
	astutil.Apply(astFile, func(cr *astutil.Cursor) bool{
		//_,ok := cr.Node().(*ast.GoStmt)
		n := cr.Node()
		if n != nil{
			t1 := n.Pos()
			t2 := prog.Fset.Position(t1)
			s := fmt.Sprintf("%v",t2)
			if !matches(n,concUsage,s) {
				return true
			}
		} else{
			return true
		}
		cr.InsertBefore(callFuncSched())
		return true
		//fmt.Println("")
	},nil)

	// add GOMACPROCS
	ast.Inspect(astFile, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			// find 'main' function
			if x.Name.Name == "main" && x.Recv == nil {
				stmts := goMaxProcs()
				stmts = append(stmts, x.Body.List...)
				x.Body.List = stmts
				return true
			}
		}
		return true
	})


	var buf bytes.Buffer
	err = printer.Fprint(&buf, prog.Fset, astFile)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// wrapper for trace statments
func createTraceStmts(timeout int) []ast.Stmt {
	ret := make([]ast.Stmt, 2)

	// trace.Start(os.Stderr)
	ret[0] = &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "trace"},
				Sel: &ast.Ident{Name: "Start"},
			},
			Args: []ast.Expr{
				&ast.SelectorExpr{
					X:   &ast.Ident{Name: "os"},
					Sel: &ast.Ident{Name: "Stderr"},
				},
			},
		},
	}
	if timeout > 0{
		// go func(){ <-time.After(5 * time.Second) trace.Stop() os.Exit(1) }()
		ret[1] = &ast.GoStmt{
			Call: &ast.CallExpr{
				Fun: &ast.FuncLit{
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "time"},
										Sel: &ast.Ident{Name: "Sleep"},
									},
									Args: []ast.Expr{
										&ast.BinaryExpr{
											X:  &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(timeout)},
											Op: token.MUL,
											Y: &ast.SelectorExpr{
												X:   &ast.Ident{Name: "time"},
												Sel: &ast.Ident{Name: "Second"},
											},
										},
									},
								},
							},
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "trace"},
										Sel: &ast.Ident{Name: "Stop"},
									},
								},
							},
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "os"},
										Sel: &ast.Ident{Name: "Exit"},
									},
									Args: []ast.Expr{
										&ast.BasicLit{Kind: token.INT, Value: "0"},
									},
								},
							},
						},
					},
					Type: &ast.FuncType{Params: &ast.FieldList{}},
				},
			},
		}
	} else{
		// defer func(){ time.Sleep(50*time.Millisecond; trace.Stop() }()
		ret[1] = &ast.DeferStmt{
			Call: &ast.CallExpr{
				Fun: &ast.FuncLit{
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "time"},
										Sel: &ast.Ident{Name: "Sleep"},
									},
									Args: []ast.Expr{
										&ast.BinaryExpr{
											X:  &ast.BasicLit{Kind: token.INT, Value: "50"},
											Op: token.MUL,
											Y: &ast.SelectorExpr{
												X:   &ast.Ident{Name: "time"},
												Sel: &ast.Ident{Name: "Millisecond"},
											},
										},
									},
								},
							},
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "trace"},
										Sel: &ast.Ident{Name: "Stop"},
									},
								},
							},
						},
					},
					Type: &ast.FuncType{Params: &ast.FieldList{}},
				},
			},
		}
	}

	return ret
}

// checks if current AST node matches with any concUsage instances
func matches(n ast.Node, conc map[string]int, location string) bool{
	if location != "-"{
		t := strings.Split(filepath.Base(location),":")[0] + ":" + strings.Split(filepath.Base(location),":")[1]
		tt := strings.Split(t,"_mod")[0] + strings.Split(t,"_mod")[1]
		if val,ok := conc[tt]; ok && val != 2 {
			conc[tt] = 2
			log.Println("ConcUsage Matches: Return True > ",tt)
			return true
		}
		return false
	}
	return false
}
