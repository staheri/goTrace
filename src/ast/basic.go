package main

import (
	"fmt"
  "log"
	"bytes"
  "os"
	"go/ast"
  _ "strconv"
	"go/format"
	"go/parser"
	"go/token"
	"reflect"
	"golang.org/x/tools/go/ast/astutil"
)

func channelAnalyis(source string){
  fset := token.NewFileSet() // A fileset to store AST
  node, err := parser.ParseFile(fset, source, nil, parser.ParseComments) // Reads and parses the source, stores the AST root in node
  if err != nil {
      log.Fatal(err)
  }
  ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) { // check the type of each node,

		case *ast.File:
			_ = x // blank identifier
		file := n.(*ast.File)
		fmt.Printf("||||\nFile: %+v\npos:%+v\n||||\n",file,fset.Position(n.Pos()))
/*
		case *ast.CommClause:
			_ = x // blank identifier
			commClause := n.(*ast.CommClause)
			fmt.Printf("||||\nCommClause: %+v\npos:%+v\n||||\n",commClause,fset.Position(n.Pos()))

		case *ast.GenDecl:
			_ = x // blank identifier
		gendecl := n.(*ast.GenDecl)
		fmt.Printf("||||\nGenDecl: %+v\npos:%+v\n||||\n",gendecl,fset.Position(n.Pos()))

		case *ast.ValueSpec:
			_ = x // blank identifier
			valueSpec := n.(*ast.ValueSpec)
			fmt.Printf("||||\nValueSpec: %+v\npos:%+v\n||||\n",valueSpec,fset.Position(n.Pos()))

		case *ast.TypeSpec:
			_ = x // blank identifier
			typeSpec := n.(*ast.TypeSpec)
			fmt.Printf("||||\ntypeSpec: %+v\npos:%+v\n||||\n",typeSpec,fset.Position(n.Pos()))

		case *ast.Ident:
			_ = x // blank identifier
			ident := n.(*ast.Ident)
			fmt.Printf("||||\nIdent: %+v\npos:%+v\n||||\n",ident,fset.Position(n.Pos()))

		case *ast.IfStmt:
			_ = x // blank identifier
			ifstmt := n.(*ast.IfStmt)
			fmt.Printf("||||\nifstmt: %+v\npos:%+v\n||||\n",ifstmt,fset.Position(n.Pos()))

		case *ast.ExprStmt:
			_ = x // blank identifier
			exprStmt := n.(*ast.ExprStmt)
			fmt.Printf("||||\nexprStmt: %+v\npos:%+v\n||||\n",exprStmt,fset.Position(n.Pos()))

		case *ast.CallExpr:
			_ = x // blank identifier
			callExpr := n.(*ast.CallExpr)
			fmt.Printf("||||\ncallExpr: %+v\npos:%+v\n||||\n",callExpr,fset.Position(n.Pos()))

		case *ast.SelectorExpr:
			_ = x // blank identifier
			selectorExpr := n.(*ast.SelectorExpr)
			fmt.Printf("||||\nselectorExpr: %+v\npos:%+v\n||||\n",selectorExpr,fset.Position(n.Pos()))

		case *ast.StructType:
			_ = x // blank identifier
			structType := n.(*ast.StructType)
			fmt.Printf("||||\nStructType: %+v\npos:%+v\n||||\n",structType,fset.Position(n.Pos()))

		case *ast.FieldList:
			_ = x // blank identifier
			fieldList := n.(*ast.FieldList)
			fmt.Printf("||||\nFieldList: %+v\npos:%+v\n||||\n",fieldList,fset.Position(n.Pos()))

		case *ast.Field:
			_ = x // blank identifier
			field := n.(*ast.Field)
			fmt.Printf("||||\nField: %+v\npos:%+v\n||||\n",field,fset.Position(n.Pos()))

		case *ast.FuncDecl:
			_ = x // blank identifier
			funcDecl := n.(*ast.FuncDecl)
			fmt.Printf("||||\nFuncDecl: %+v\npos:%+v\n||||\n",funcDecl,fset.Position(n.Pos()))

*/
		default:
			//fmt.Printf(".")
			_ = x // blank identifier
			if n != nil{
				fmt.Printf("%+v : %+v\n*******\n",reflect.TypeOf(n),fset.Position(n.Pos()))
			}
			//*/
		}
		/*
		case *ast.ChanType:
      _ = x // blank identifier
      ct := n.(*ast.ChanType) // extract the structs := type assertion
      fmt.Printf("%s: \tChanTYPE\n%+v\n\n",fset.Position(n.Pos()),ct)
		case *ast.SendStmt:
      _ = x
      ss := n.(*ast.SendStmt)
      fmt.Printf("%s: \tSend\n%+v\n\n",fset.Position(n.Pos()),ss)
    case *ast.CommClause:
      _ = x
      cc := n.(*ast.CommClause)
      fmt.Printf("%s: \tSelect statment\n%+v\n\n",fset.Position(n.Pos()),cc)
		case *ast.CommClause:
      _ = x
      cc := n.(*ast.CommClause)
      fmt.Printf("%s: \tSelect statment\n%+v\n\n",fset.Position(n.Pos()),cc)
		}*/

		/*ident,ok := n.(*ast.Ident)
		if ok{
			fmt.Printf("\tIdent Name: %s\n\tPosition: %s\n",ident.Name,fset.Position(n.Pos()))
		}
		//functype,ok := n.(*ast.FuncType)
		_,ok = n.(*ast.FuncType)
		if ok{
			fmt.Printf("\t\"func\" keyword detected at Position: %s\n",fset.Position(n.Pos()))
		}*/
		return true
	})
}

func ex(source string){
	fset := token.NewFileSet() // A fileset to store AST
  node, err := parser.ParseFile(fset, source, nil, parser.ParseComments) // Reads and parses the source, stores the AST root in node
  if err != nil {
      log.Fatal(err)
  }
	for _, d := range node.Decls {
		fmt.Println("#### ", reflect.TypeOf(d))
		funcdecl,ok := d.(*ast.FuncDecl)
		if ok{
			fmt.Println("\t","Name (ident):",funcdecl.Name)
			fmt.Println("\t","Receiver (FieldList):",funcdecl.Recv)
			fmt.Println("\t","Type (FuncType):",funcdecl.Type)
			//functype, ok2 := funcdecl.Type.(*ast.FuncType)
			//if ok2{
				//fmt.Printf("\tFunc:%v\n\tParams: %v\n\tResults: %v\n",functype.Func,functype.Params,functype.Results)
			//}
			fmt.Printf("\tFunc:%v\n\tParams: %v\n\tResults: %v\n",funcdecl.Type.Func,funcdecl.Type.Params,funcdecl.Type.Results)
			fmt.Println("\t","Body (BlockStmt):",funcdecl.Body)
		}
		gendecl,ok := d.(*ast.GenDecl)
		if ok{
			fmt.Println("\t","TOK:",gendecl.Tok)
			impspec, ok2 := gendecl.Specs[0].(*ast.ImportSpec)
			if ok2{
				fmt.Println("\t","Spec ([]Spec):",impspec.Name)
			}

			fmt.Println("\t","( (Tok.Pos):",fset.Position(gendecl.Lparen))
			fmt.Println("\t",") (Tok.Pos):",fset.Position(gendecl.Rparen))
		}
		// print funcDecl fields
		// print GenDecl fields
	}
	for _, s := range node.Imports {
		fmt.Println("**** ", reflect.TypeOf(s))
		fmt.Println("\t","Name:",s.Name)
	}
	fmt.Println("Name: ", node.Name)
	fmt.Println("Package: ", node.Package)
	fmt.Println("Scope: ", node.Scope)
}

func extractImports(source string){
	fset := token.NewFileSet() // A fileset to store AST
  node, err := parser.ParseFile(fset, source, nil, parser.ParseComments) // Reads and parses the source, stores the AST root in node
  if err != nil {
      log.Fatal(err)
  }

	for _, s := range node.Imports {
		fmt.Println("**** ", reflect.TypeOf(s))
		fmt.Println("\t","Name:",s.Name)
		fmt.Println("\t\t","Path (BasicLit):",s.Path)
	}
	newImport := ast.ImportSpec{Doc: nil, Name: nil, Path: &ast.BasicLit{Kind: token.STRING, Value: "tracer"}}
	node.Imports = append(node.Imports, &newImport)
	for _, s := range node.Imports {
		fmt.Println("**** ", reflect.TypeOf(s))
		fmt.Println("\t","Name:",s.Name)
		fmt.Println("\t\t","Path (BasicLit):",s.Path)
		fmt.Println("\t\t","POS:",fset.Position(s.Path.Pos()))
	}
	for _, d := range node.Decls {
		gendecl,ok := d.(*ast.GenDecl)
		if ok{
			fmt.Println("\t","TOK:",gendecl.Tok)
			fmt.Println("\t","( loc:",gendecl.Lparen)
			fmt.Println("\t",") loc:",gendecl.Rparen)
			for _,elem := range gendecl.Specs{
				impspec, ok2 := elem.(*ast.ImportSpec)
				if ok2{
					fmt.Println("\t\t","Name (Ident)):",impspec.Name)
					fmt.Println("\t\t","Path (BasicLit):",impspec.Path)
				} else{
					fmt.Println("\t\t >> ",reflect.TypeOf(elem))
				}
			}
		}
	}
	astutil.AddImport(fset, node, "tracer")
	var buf bytes.Buffer
	err1 := format.Node(&buf, fset, node)
	if err1 != nil {
		panic(err1)
	}
	fmt.Println("---")
	fmt.Println(buf.String())
	fmt.Println("---")
}

func main(){
  if len(os.Args) != 2 {
    log.Fatalf("specify the input source")
  }
  src := os.Args[1]

  fmt.Printf("AST-analysis of %s\n",src)
  channelAnalyis(src)
	//ex(src)
	//extractImports(src)
}
