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
		/*switch x := n.(type) { // check the type of each node,
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
		if n != nil{
			fmt.Printf("%+v : %+v\n*******\n",reflect.TypeOf(n),fset.Position(n.Pos()))
		}
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
