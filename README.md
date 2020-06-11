# goTrace

goTrace is a tool that automatically:

- Instruments Go source
   * It traverses source AST tree
   * And injects trace collection API to the source
- Executes the target application and redirects its trace to ```stderr```
   * Go deadlock detector can be disabled on applications that suffer from deadlock
- Inserts traces into a MySQL database
   * Now you can query the database and study the behavior

The code for creating concept lattices is also ready but many things should be done manually [to add automatic and robust functionality]


# Dependencies
GoTrace uses different libraries and drivers. [TODO] Use Go Modules/Vendors to automatically detect dependencies and versions

## Libraries

- Fine tables: [github.com/jedib0t/go-pretty/table](https://github.com/jedib0t/go-pretty)
  `go get github.com/jedib0t/go-pretty`
- Go MySQL driver: [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
  `go get github.com/go-sql-driver/mysql`
- AST traversal: [golang.org/x/tools/go/ast/astutil](https://golang.org/x/tools/go/ast/astutil)
  `golang.org/x/tools/go/ast/astutil`
- There might be more libraries needed indirectly.

## Database

- MySQL: [Install on Mac](https://dev.mysql.com/doc/mysql-osx-excerpt/5.7/en/osx-installation-pkg.html)


# Patching Runtime
`goTrace-runtime.patch` has all the needed injections to the Go runtime in order to capture additional events like channel operations, waiting groups and mutexes.

Assuming your Go installation is in `/usr/local/go`, download Go 1.14.4 and unpack it into `/usr/local/go-new`.
```
 sudo -i
 mkdir -p /usr/local/go-new
 curl https://dl.google.com/go/go1.14.4.darwin-amd64.tar.gz | tar -xz -C /usr/local/go-new
 ```

Then, copy patch and apply it:
```
sudo patch -p1 -d /usr/local/go-new/go < goTrace-runtime.patch
```

Now you can build the new runtime
```
 sudo -i
 cd /usr/local/go-new/go/src
 export GOROOT_BOOTSTRAP=/usr/local/go #or choose yours
 ./make.bash
 ```

Finally, `export PATH` or `use ln -s` command to make this Go version actual in your system:
```
 export PATH=/usr/local/go-new/go/bin:$PATH
 ```
or (assuming your PATH set to use /usr/local/go)
```
	sudo mv /usr/local/go /usr/local/go-orig
	sudo ln -nsf /usr/local/go-new/go /usr/local/go
```
NOTE: return your previous installation by `sudo ln -nsf /usr/local/go-orig /usr/local/go`

# Build GoTrace
First make sure you have set-up the Go environment variables correctly
```
export GOROOT=/usr/local/go
export GOPATH=<path-to>/goTrace
export PATH=$GOROOT/bin:$PATH
```
Then
```
cd goTrace/src
go build
```

# Run
`./src -app=test.go` would automatically take `test.go`, instruments, builds and executes it. The resulting traces are stored in the MySQL database which you can access separately.

Do `./src --help` for more information

# Case Studies

## BoltDB - Deadlock
According to the [ASPLOS'19 paper](https://dl.acm.org/doi/10.1145/3297858.3304069), there was a bug (deadlock) in BoltDB project which was fixed after [this commit](https://github.com/boltdb/bolt/commit/defdb743cdca840890fea24c3111a7bffe5cc0a3). This bug was clearly caused by different orderings of the acquiring/releasing of different data structure mutex. 
