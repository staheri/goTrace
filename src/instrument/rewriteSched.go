// adapted from code by Ivan Daniluk from https://github.com/divan

package instrument

import (
	"bytes"
	_"errors"
	"strconv"
	"go/ast"
	"go/printer"
	"go/token"
	"io/ioutil"
	"path/filepath"

	"golang.org/x/tools/go/ast/astutil"

	"golang.org/x/tools/go/loader"
)

//var ErrImported = errors.New("trace already imported")

// rewriteSourceSched rewrites current source and saves
// into temporary file, returning it's path.
func rewriteSourceSched(path string, timeout,depth int) (string, error) {
	data, err := addCodeSched(path,timeout,depth)
	if err == ErrImported {
		data, err = ioutil.ReadFile(path)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	tmpDir, err := ioutil.TempDir("", "gotracer_package")
	if err != nil {
		return "", err
	}
	filename := filepath.Join(tmpDir, filepath.Base(path))
	// SAEED
	filename2 := filepath.Join(filepath.Dir(path), filepath.Base(path)+".mod")
	err = ioutil.WriteFile(filename2, data, 0666)

	err = ioutil.WriteFile(filename, data, 0666)
	if err != nil {
		return "", err
	}

	return tmpDir, nil
}

// addCode searches for main func in data, and updates AST code
// adding tracing functions.
func addCodeSched(path string, timeout,depth int) ([]byte, error) {
	var conf loader.Config
	if _, err := conf.FromArgs([]string{path}, false); err != nil {
		return nil, err
	}

	prog, err := conf.Load()
	if err != nil {
		return nil, err
	}

	// check if runtime/trace already imported
	for i, _ := range prog.Imported {
		if i == "runtime/trace" {
			return nil, ErrImported
		}
	}

	pkg := prog.Created[0]

	// TODO: find file with main func inside
	astFile := pkg.Files[0]

	// add imports
	astutil.AddImport(prog.Fset, astFile, "os")
	astutil.AddImport(prog.Fset, astFile, "runtime/trace")
	astutil.AddImport(prog.Fset, astFile, "sync")
	astutil.AddImport(prog.Fset, astFile, "runtime")
	astutil.AddImport(prog.Fset, astFile, "math/rand")
	astutil.AddImport(prog.Fset, astFile, "time")
	if timeout > 0{
		astutil.AddNamedImport(prog.Fset, astFile, "_", "net")
	}
	// add constant, struct type, global counter, function declration
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
	// add start/stop code to the main
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

func constNode(name, value string) *ast.GenDecl {
	return &ast.GenDecl{
		Tok:token.CONST,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{
					&ast.Ident{Name: name},
				},
				Values: []ast.Expr{
					&ast.BasicLit{Kind: token.INT, Value: value},
				},
			},
		},
	}
}

func structNode() *ast.GenDecl {
	return &ast.GenDecl{
		Tok:token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{ Name: "sharedInt",},
				Type: &ast.StructType{
					Fields:&ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Names: []*ast.Ident{&ast.Ident{Name: "n"}},
								Type: &ast.Ident{Name: "int"},
							},
							&ast.Field{
								Type: &ast.SelectorExpr{
									X: &ast.Ident{Name: "sync"},
									Sel: &ast.Ident{Name: "Mutex"},
								},
							},
						},
					},
				},
			},
		},
	}
}


func globalCount() *ast.GenDecl{
	return &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{&ast.Ident{Name: "cnt"}},
				Type: &ast.Ident{Name: "sharedInt"},
			},
		},
	}
}

func goMaxProcs() *ast.ExprStmt{
	return &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "runtime"},
				Sel: &ast.Ident{Name: "GOMAXPROCS"},
			},
			Args: []ast.Expr{
				&ast.BasicLit{Kind: token.INT, Value: "1"},
			},
		},
	}
}


func callFuncSched() *ast.IfStmt{
	return &ast.IfStmt{
		Cond: &ast.CallExpr{
			Fun: &ast.Ident{Name: "Reschedule"},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "runtime"},
							Sel: &ast.Ident{Name: "Gosched"},
						},
					},
				},
			},
		},
	}
}

func declFuncSched() *ast.FuncDecl{
	return &ast.FuncDecl{
		Name: &ast.Ident{Name: "Reschedule"},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Type: &ast.Ident{Name: "bool"},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{ // random seed generator
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "rand"},
							Sel: &ast.Ident{Name: "Seed"},
						},
						Args: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   &ast.Ident{Name: "time"},
											Sel: &ast.Ident{Name: "Now"},
										},
									},
									Sel: &ast.Ident{Name: "UnixNano"},
								},
							},
						},
					},
				},
				&ast.IfStmt{ // main if
					Cond: &ast.BinaryExpr{ // coint toss
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "rand"},
								Sel: &ast.Ident{Name: "Intn"},
							},
							Args: []ast.Expr{
								&ast.BasicLit{Kind: token.INT, Value: "2"},
							},
						},
						Y: &ast.BasicLit{Kind: token.INT, Value: "1"},
						Op: token.EQL,
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{ // lock
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "cnt"},
										Sel: &ast.Ident{Name: "Lock"},
									},
								},
							},
							&ast.DeferStmt{// defer unlock
								Call: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "cnt"},
										Sel: &ast.Ident{Name: "Unlock"},
									},
								},
							},
							&ast.IfStmt{// if
								Cond: &ast.BinaryExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{Name: "cnt"},
										Sel: &ast.Ident{Name: "n"},
									},
									Y: &ast.Ident{Name: "depth"},
									Op: token.LSS,
								},
								Body: &ast.BlockStmt{
									List: []ast.Stmt{
										&ast.IncDecStmt{
											X: &ast.SelectorExpr{
												X: &ast.Ident{Name: "cnt"},
												Sel: &ast.Ident{Name: "n"},
											},
											Tok: token.INC,
										},
										&ast.ReturnStmt{
											Results: []ast.Expr{
												&ast.Ident{Name: "true"},
											},
										},
									},
								},
								Else: &ast.ReturnStmt{
									Results: []ast.Expr{
										&ast.Ident{Name: "false"},
									},
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.Ident{Name: "false"},
					},
				},
			},
		},
	}
}

func newDecls(depth int) []ast.Decl {
	ret := make([]ast.Decl,4)
	ret[0] = constNode("depth",strconv.Itoa(depth))
	ret[1] = structNode()
	ret[2] = globalCount()
	ret[3] = declFuncSched()
	return ret
}
