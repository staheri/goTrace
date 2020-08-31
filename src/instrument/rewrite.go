// adapted from code by Ivan Daniluk from https://github.com/divan

package instrument

import (
	"bytes"
	"errors"
	"go/ast"
	"go/printer"
	"go/token"
	"io/ioutil"
	"path/filepath"

	"golang.org/x/tools/go/ast/astutil"

	"golang.org/x/tools/go/loader"

	"strconv"
)

var ErrImported = errors.New("trace already imported")

// rewriteSource rewrites current source and saves
// into temporary file, returning it's path.
func rewriteSource(path string,timeout int) (string, error) {
	data, err := addCode(path,timeout)
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
func addCode(path string, timeout int) ([]byte, error) {
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
	astutil.AddImport(prog.Fset, astFile, "time")
	if timeout > 0{
		astutil.AddNamedImport(prog.Fset, astFile, "_", "net")
	}

	// add start/stop code
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
											X:  &ast.BasicLit{Kind: token.INT, Value: "1"},
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
						},
					},
					Type: &ast.FuncType{Params: &ast.FieldList{}},
				},
			},
		}
	}

	return ret
}
