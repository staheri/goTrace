package main

import (
	"fmt"
  "log"
  "os"
	"go/ast"
  _ "strconv"
	_ "go/format"
	"go/parser"
	"go/token"
)


func printImportInfo(f *ast.ImportSpec){
  //if len(f.Doc.List) > 0:
  //fmt.Printf("*********<Import>**********\nDoc: %s\nName: %s\n",f.Doc.List[0].Text, f.Name.Name)
  fmt.Printf("Path: %s\nComment: %s\n*********</Import>**********\n",f.Path.Value,f.Comment.Text())
}


func channelAnalyis(source string){
  fset := token.NewFileSet() // A fileset to store AST
  node, err := parser.ParseFile(fset, source, nil, parser.ParseComments) // Reads and parses the source, stores the AST root in node
  if err != nil {
      log.Fatal(err)
  }
  fmt.Println("Imports:")
  for _, i := range node.Imports {
    printImportInfo(i)
  }

  //fmt.Println("Comments:")
  //for _, i := range node.Comments {
  //  fmt.Println(i.Text())
  //}
  ast.Inspect(node, func(n ast.Node) bool {
		var s string

		/*switch x := n.(type) {
		case *ast.ChanType:
			//s = x.Value
      //t = x.Value
      _ = x
      tt := n.(*ast.ChanType)
      fmt.Printf("%s: \tChanTYPE\ntt: %+v\n",fset.Position(n.Pos()),tt)
		case *ast.SendStmt:
      _ = x
      qq := n.(*ast.SendStmt)
      fmt.Printf("%s: \tSend\nSend Struct: %+v\n",fset.Position(n.Pos()),qq)
			//s = strconv.Itoa(n)
      //fmt.Printf("ChanDIR")
    case *ast.CommClause:
      _ = x
      g := n.(*ast.CommClause).
      fmt.Printf("%s: \tComm\nCommClause Struct: %+v\n",fset.Position(n.Pos()),g)
			//s = strconv.Itoa(n)
      //fmt.Printf("ChanDIR")
		}
		if s != "" {
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), s)
		}*/
		fmt.Printf("NODE:\n%+v\n*******\n",n.(type))
		return true
	})
}

func main(){
  if len(os.Args) != 2 {
    log.Fatalf("specify the input source")
  }
  src := os.Args[1]

  fmt.Printf("AST-analysis of %s\n",src)
  channelAnalyis(src)

}
