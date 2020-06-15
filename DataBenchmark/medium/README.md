## Naming Convention

`appName`X`execution_ID`\_l`chunk_size`\_`mode`\_`filter`.py

## Modes

- **seqAll**: Single sequence of events from all goroutines (including tracing and runtime goroutines)
- **seqAPP**: Single sequence of events from application goroutines
- **grtnAPP**: Separated sequences of events from application goroutines

## Filters

- **all**: No filter, include everything
- **CHNL**: Channel events
- **GRTN**: Goroutine events
- **MUTX**: Mutex events
- **GCMM**: Garbage collection/memory events
- **PROC**: Process events
- **WGRP**: WaitingGroup events
- **SYSC**: Syscall events
- **MISC**: Other events
