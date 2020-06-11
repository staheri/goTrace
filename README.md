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

### Sample query after the fix
```
mysql>
  SELECT t1.type, t1.ts, t4.arg, t4.value, t3.func
  FROM sample_1.Events t1
  INNER JOIN global.catMUTX t2 ON t1.type = t2.eventName
  INNER JOIN sample_1.StackFrames t3 ON t1.id = t3.eventId
  INNER JOIN sample_1.Args t4 ON t1.id = t4.eventID
  WHERE t3.func="github.com/boltdb/bolt.(*DB).beginTx"
        OR
        t3.func="github.com/boltdb/bolt.(*DB).removeTx" ;
```

|-------------|-----------|------|-------|---------------------------------------|
| type        | ts        | arg  | value | func                                  |
|-------------|-----------|------|-------|---------------------------------------|
| EvMuLock    | 113889903 | muid |     7 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvRWMrLock  | 113895317 | rwid |     3 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuUnlock  | 113904913 | muid |     7 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvRWMLock   | 113908300 | rwid |     4 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuLock    | 113911379 | muid |     9 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvRWMUnlock | 113914946 | rwid |     4 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuUnlock  | 113918230 | muid |     9 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuLock    | 113963337 | muid |     7 | github.com/boltdb/bolt.(*DB).removeTx |
| EvMuUnlock  | 113969162 | muid |     7 | github.com/boltdb/bolt.(*DB).removeTx |
| EvRWMLock   | 113982427 | rwid |     4 | github.com/boltdb/bolt.(*DB).removeTx |
| EvMuLock    | 113986712 | muid |     9 | github.com/boltdb/bolt.(*DB).removeTx |
| EvRWMUnlock | 113991177 | rwid |     4 | github.com/boltdb/bolt.(*DB).removeTx |
| EvMuUnlock  | 113995513 | muid |     9 | github.com/boltdb/bolt.(*DB).removeTx |
| EvMuLock    | 114000311 | muid |     7 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvRWMrLock  | 114003390 | rwid |     3 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuUnlock  | 114008983 | muid |     7 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvRWMLock   | 114012499 | rwid |     4 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuLock    | 114015886 | muid |     9 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvRWMUnlock | 114019195 | rwid |     4 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuUnlock  | 114022274 | muid |     9 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuLock    | 114042493 | muid |     7 | github.com/boltdb/bolt.(*DB).removeTx |
| EvMuUnlock  | 114047086 | muid |     7 | github.com/boltdb/bolt.(*DB).removeTx |
| EvRWMLock   | 114051397 | rwid |     4 | github.com/boltdb/bolt.(*DB).removeTx |
| EvMuLock    | 114055425 | muid |     9 | github.com/boltdb/bolt.(*DB).removeTx |
| EvRWMUnlock | 114059684 | rwid |     4 | github.com/boltdb/bolt.(*DB).removeTx |
| EvMuUnlock  | 114063892 | muid |     9 | github.com/boltdb/bolt.(*DB).removeTx |
| EvMuLock    | 114067972 | muid |     7 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvRWMrLock  | 114071179 | rwid |     3 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuUnlock  | 114075772 | muid |     7 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvRWMLock   | 114079364 | rwid |     4 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuLock    | 114082597 | muid |     9 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvRWMUnlock | 114085958 | rwid |     4 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuUnlock  | 114089037 | muid |     9 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuLock    | 114124882 | muid |     7 | github.com/boltdb/bolt.(*DB).removeTx |
| EvMuUnlock  | 114129475 | muid |     7 | github.com/boltdb/bolt.(*DB).removeTx |
| EvRWMLock   | 114133580 | rwid |     4 | github.com/boltdb/bolt.(*DB).removeTx |
| EvMuLock    | 114137250 | muid |     9 | github.com/boltdb/bolt.(*DB).removeTx |
| EvRWMUnlock | 114141483 | rwid |     4 | github.com/boltdb/bolt.(*DB).removeTx |
| EvMuUnlock  | 114145691 | muid |     9 | github.com/boltdb/bolt.(*DB).removeTx |
|-------------|-----------|------|-------|---------------------------------------|


### Sample query (deadlock)
```
mysql>
  SELECT t1.type, t1.ts, t4.arg, t4.value, t3.func
  FROM sample_2.Events t1
  INNER JOIN global.catMUTX t2 ON t1.type = t2.eventName
  INNER JOIN sample_2.StackFrames t3 ON t1.id = t3.eventId
  INNER JOIN sample_2.Args t4 ON t1.id = t4.eventID
  WHERE t3.func="github.com/boltdb/bolt.(*DB).beginTx"
        OR
        t3.func="github.com/boltdb/bolt.(*DB).removeTx" ;
```

|-------------|----------|------|-------|---------------------------------------|
| type        | ts       | arg  | value | func                                  |
|-------------|----------|------|-------|---------------------------------------|
| EvRWMrLock  | 76503916 | rwid |     3 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuLock    | 76506789 | muid |     7 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuUnlock  | 76513255 | muid |     7 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvRWMLock   | 76515103 | rwid |     4 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuLock    | 76516770 | muid |     9 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvRWMUnlock | 76518772 | rwid |     4 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuUnlock  | 76528291 | muid |     9 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuLock    | 76556977 | muid |     7 | github.com/boltdb/bolt.(*DB).removeTx |
| EvMuUnlock  | 76560287 | muid |     7 | github.com/boltdb/bolt.(*DB).removeTx |
| EvRWMLock   | 76562340 | rwid |     4 | github.com/boltdb/bolt.(*DB).removeTx |
| EvMuLock    | 76564367 | muid |     9 | github.com/boltdb/bolt.(*DB).removeTx |
| EvRWMUnlock | 76566573 | rwid |     4 | github.com/boltdb/bolt.(*DB).removeTx |
| EvMuUnlock  | 76568677 | muid |     9 | github.com/boltdb/bolt.(*DB).removeTx |
| EvRWMrLock  | 76571140 | rwid |     3 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuLock    | 76572783 | muid |     7 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuUnlock  | 76575066 | muid |     7 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvRWMLock   | 76576760 | rwid |     4 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuLock    | 76578325 | muid |     9 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvRWMUnlock | 76579993 | rwid |     4 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuUnlock  | 76581583 | muid |     9 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuLock    | 76601443 | muid |     7 | github.com/boltdb/bolt.(*DB).removeTx |
| EvMuUnlock  | 76603855 | muid |     7 | github.com/boltdb/bolt.(*DB).removeTx |
| EvRWMLock   | 76605831 | rwid |     4 | github.com/boltdb/bolt.(*DB).removeTx |
| EvMuLock    | 76607806 | muid |     9 | github.com/boltdb/bolt.(*DB).removeTx |
| EvRWMUnlock | 76609936 | rwid |     4 | github.com/boltdb/bolt.(*DB).removeTx |
| EvMuUnlock  | 76611963 | muid |     9 | github.com/boltdb/bolt.(*DB).removeTx |
| EvRWMrLock  | 76614298 | rwid |     3 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuLock    | 76615889 | muid |     7 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuUnlock  | 76618044 | muid |     7 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvRWMLock   | 76619737 | rwid |     4 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuLock    | 76621303 | muid |     9 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvRWMUnlock | 76622970 | rwid |     4 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuUnlock  | 76624535 | muid |     9 | github.com/boltdb/bolt.(*DB).beginTx  |
| EvMuLock    | 76655736 | muid |     7 | github.com/boltdb/bolt.(*DB).removeTx |
| EvMuUnlock  | 76658122 | muid |     7 | github.com/boltdb/bolt.(*DB).removeTx |
| EvRWMLock   | 76660098 | rwid |     4 | github.com/boltdb/bolt.(*DB).removeTx |
| EvMuLock    | 76662176 | muid |     9 | github.com/boltdb/bolt.(*DB).removeTx |
| EvRWMUnlock | 76664280 | rwid |     4 | github.com/boltdb/bolt.(*DB).removeTx |
| EvMuUnlock  | 76666282 | muid |     9 | github.com/boltdb/bolt.(*DB).removeTx |
|-------------|----------|------|-------|---------------------------------------|
