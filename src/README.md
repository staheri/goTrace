
# Run

This repository is under heavy construction so to run the tool, you have to modify `main.go` file directly. Currently, here are the available functionalities that can be activated through main:

## Instrument, Execute, Store
Pass the Go source as argument `./src -app=<APP.GO>`. GoTrace then instruments, executes and stores the whole program (end to end) sequence of events in a database name `app_nameX?` where `?` is the id of execution (in case you execute the same app multiple times)

## Generate reports, etc.
Let us assume that we already stored traces in a database, now we want to study traces from different aspects. Here are the available functionalities that you can achieve by passing the DB name as argument `./src -dbName=<DB_NAME_X0>`:

- Display events grouped by goroutines (through `cl.GroupGrtns()`)
- Generate data for word2vec ideas (through `db.WriteData()`) **NEED SOME WORK**
- Generate formal context, concept lattice, jaccard matrix, dendogram and flat cluster of goroutines (through `db.FormalContext()`)
- Generate **Resource Reports** where resources are (a sample is provided in Case Studies section):
   * Channels
   * Mutexes
   * WaitingGroups



Do `./src --help` for more information
