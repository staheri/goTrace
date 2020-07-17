
# Dinning Philosophers in Go

## Three implementations
- **dine3.go**: One channel per fork - Adjustable timing (from source) - N=3
- **dineN.go**: One channel per fork - Adjustable N (from source) - Random(fair) timing
- **dinphil-promela.go**: Implementation from Ganesh's book - data-dependent - two channels per fork (*WARNING: IT IS BUGGY. DO NOT USE FOR EXPERIMENTS*)

## Generating Data:

(*for instructions on how to run goTrace and options, please refer to the repo README*)

To generate data:

```
./src -app=<path-to-target-app> -cmd=dineData -src=<native/latest/x> [-x=X] -outdir=<path-to-data-folder> -n=<#-of-philosophers>
```
Example:
```
./src -app=../CodeBenchmark/dinPhil/dineN.go -cmd=dineData -src=native -outdir=../DataBenchmark/dineData/n5 -n=5
```

## Data Format
The above command generates four sets of data (in separate folders) under `outdir` folder:
- **all**: All events
- **all-chid**: All events + channel ID (for channel events)
- **ch**: channel events only
- **ch-chid**: channel events only + channel ID

For examples, please refer to `DataBenchmark/dineData/`

## Hint for prototyping/testing
If your target application generates too long traces, you can limit your analysis to TOP-X elements of the global trace. In order to do so, just uncomment lines 108, 109, 110 in `src/db/inops.go` under `Store` function and rebuild:
```
108   if cnt > TOPX{
109    	break
110   }
```

```
?> cd src
?> cd build
?> ./src --help
```

You can set `TOPX` from `src/db/constants.go` (default TOPX=2000)
