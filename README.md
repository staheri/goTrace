# goTrace

goTrace is a tool that automatically:

- Instruments Go source
   * It traverses source AST tree
   * And injects trace collection API to the source
- Executes the target application and redirects its trace to ```stderr```
   * Go deadlock detector can be disabled on applications suffer from deadlock
- Inserts traces into a MySQL database
   * Now you can query the database and study the behavior

The code for creating concept lattices is also ready but many things should be done manually [to add automatic and robust functionality]


# Required Libraries
- Fine tables: [github.com/jedib0t/go-pretty/table](github.com/jedib0t/go-pretty)
  `go get github.com/jedib0t/go-pretty`
- Go MySQL driver: [github.com/go-sql-driver/mysql](github.com/go-sql-driver/mysql)
  `go get github.com/go-sql-driver/mysql`
- AST traversal: [golang.org/x/tools/go/ast/astutil](golang.org/x/tools/go/ast/astutil)
  `golang.org/x/tools/go/ast/astutil`

- There might be more libraries needed indirectly.
- [TODO] Use Go Modules to automatically detect dependencies and versions


# Instruction
