// adapted from code by Ivan Daniluk from https://github.com/divan

package instrument

import (
	"go/ast"
	"go/token"
	"strconv"
)


// returns a general declration representing a constant node
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

// returns sharedInt type structure node
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

// returns a global instance of sharedInt node
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

// returns GOMAXPROCS line node
func goMaxProcs() []ast.Stmt{
	ret := make([]ast.Stmt, 1)
	ret[0] = &ast.ExprStmt{
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
	return ret
}

// returns "if Reschedule then Gosched()" line node
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

// returns Reschedule() delration node
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

// wrapper for new declrations
func newDecls(depth int) []ast.Decl {
	ret := make([]ast.Decl,4)
	ret[0] = constNode("depth",strconv.Itoa(depth))
	ret[1] = structNode()
	ret[2] = globalCount()
	ret[3] = declFuncSched()
	return ret
}
