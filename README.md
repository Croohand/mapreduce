# MapReduce

A distributed MapReduce framework with a Hadoop-like distributed file system (HDFS) implementation, written in Go.

## Overview

This project implements a fully functional distributed computing system featuring:

- **Distributed File System**: HDFS-like architecture with data blocks, replication, and fault tolerance
- **MapReduce Engine**: Execute custom Map and Reduce operations across distributed data
- **Coordinator-Worker Architecture**: Coordinated cluster management with leader election
- **Scheduler**: Dedicated scheduler nodes for orchestrating MapReduce jobs
- **Visualizer**: Real-time visualization of cluster activity and request flows

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                 CLIENT                                      │
│                   (read, write, mapreduce, list, etc.)                      │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           COORDINATOR NODES                                 │
│    ┌────────────┐  ┌────────────┐  ┌────────────┐                           │
│    │Coordinator2│◄►│Coordinator1│◄►│Coordinator3│  (Journal replication)    │
│    │  (Passive) │  │ (Active)   │  │ (Passive)  │                           │
│    └────────────┘  └────────────┘  └────────────┘                           │
│                           │                                                 │
│                           ├── Metadata management (BoltDB)                  │
│                           ├── Block placement & replication                 │
│                           └── Transaction coordination                      │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
               ┌─────────────────────┼─────────────────────┐
               ▼                     ▼                     ▼
┌──────────────────────┐  ┌──────────────────────┐  ┌──────────────────┐
│     SCHEDULER(S)     │  │     WORKER NODES     │  │   SIMPLE LOGGER  │
│                      │  │                      │  │                  │
│  • Map/Reduce        │  │  • Block storage     │  │  • Request       │
│    orchestration     │  │  • Map/Reduce        │  │    logging       │
│  • Status tracking   │  │    execution         │  │                  │
│                      │  │  • Data replication  │  │                  │
└──────────────────────┘  └──────────────────────┘  └──────────────────┘
```

## Components

### Coordinator

The coordinator node manages:
- **Metadata Storage**: File paths to block mappings (using BoltDB)
- **Block Management**: Block placement, replication tracking, and recovery
- **Transaction Coordination**: ACID-like transactions for file operations
- **Leader Election**: Active/passive mode with automatic failover
- **Journal Replication**: Asynchronous journal replication to passive coordinators

### Worker

Worker nodes responsible for:
- **Block Storage**: Storing and serving data blocks
- **Map Operations**: Executing user-defined Map functions on local data
- **Reduce Operations**: Executing user-defined Reduce functions
- **Block Replication**: Copying blocks to other workers for fault tolerance

### Scheduler

Dedicated nodes for MapReduce job orchestration:
- **Job Preparation**: Setting up transactions and coordinating with coordinator
- **Task Distribution**: Assigning map tasks to workers with data locality
- **Shuffle Phase**: Coordinating data movement between mappers and reducers
- **Progress Tracking**: Monitoring operation status

### Client

Command-line interface for:
- File operations: `read`, `write`, `copy`, `move`, `remove`, `list`, `exists`
- Data processing: `merge`, `mapreduce`

### Simple Logger

Optional logging service for development:
- Captures all HTTP requests between components
- Provides data for the visualizer
- Tracks component status (up/down)

### Visualizer

Web-based visualization tool:
- Real-time cluster activity visualization
- Request flow animation between components
- Recording capability for presentations

## Installation

### Prerequisites

- Go 1.21 or later
- Unix-like environment (Linux, macOS)

### Build

```bash
# Clone the repository
git clone https://github.com/Croohand/mapreduce.git
cd mapreduce

# Install all components
go install ./client ./coordinator ./worker ./simple_logger
```

The binaries will be installed to `$GOPATH/bin` (or `$HOME/go/bin` by default). Make sure this directory is in your `PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Environment Setup

Configure the required environment variables:

```bash
# Required: Data storage directory for cluster nodes
export MR_PATH=/path/to/mapreduce/data

# Required for client: Comma-separated coordinator addresses
export MR_HOSTS=http://localhost:11000,http://localhost:11001,http://localhost:11002

# Optional: Client identifier for logging/visualization
export MR_CLIENT=client1
```

## Quick Start

### Starting a Local Cluster

Use the provided run script for local development:

```bash
cd run
./run.sh 0  # Argument 0 overrides nodes configs
```

This starts:
- 3 Coordinator nodes (ports 11000-11002)
- 5 Worker nodes (ports 11003-11007)
- 2 Schedulers (ports 11008-11009)
- 1 Logger (port 11100)

The script automatically:
- Builds and installs all components
- Creates necessary data directories under `$MR_PATH`
- Copies the MapReduce template to each worker's sources directory

### Stopping the Cluster

```bash
./kill.sh
```

### Example Workflow

```bash
# 1. Set up environment
export MR_PATH=~/.mapreduce
export MR_HOSTS=http://localhost:11000,http://localhost:11001,http://localhost:11002
export PATH=$PATH:$(go env GOPATH)/bin

# 2. Start cluster
cd run && ./run.sh 0
cd ..

# 3. Write some data
echo -e "hello world\nfoo bar\nbaz bzzz" | client write -path /test/input.txt

# 4. Run word count
client mapreduce \
  -in /test/input.txt \
  -out /test/output.txt \
  -mappers 3 \
  -reducers 2 \
  -srcs ./operations/word_count/mruserlib

# 5. Read results
client read -path /test/output.txt

# 6. Clean up
client rm -path /test/input.txt
client rm -path /test/output.txt

# 7. Stop cluster
cd run && ./kill.sh
```

## Usage

### File Operations

```bash
# Write data from stdin
cat myfile.txt | client write -path /data/myfile.txt

# Append to existing file
cat more_data.txt | client write -path /data/myfile.txt -append

# Read file to stdout
client read -path /data/myfile.txt

# List files with prefix
client list -prefix /data/

# Check if file exists
client exists -path /data/myfile.txt

# Copy file
client copy -src /data/myfile.txt -dst /backup/myfile.txt
client cp -src /data/myfile.txt -dst /backup/myfile.txt  # alias

# Move file
client move -src /data/myfile.txt -dst /archive/myfile.txt
client mv -src /data/myfile.txt -dst /archive/myfile.txt  # alias

# Remove file
client remove -path /data/myfile.txt
client rm -path /data/myfile.txt  # alias

# Merge multiple files into one
client merge -in /data/part1.txt,/data/part2.txt -out /data/merged.txt
```

### MapReduce Operations

```bash
client mapreduce \
  -in /data/input.txt \
  -out /data/output.txt \
  -mappers 5 \
  -reducers 3 \
  -srcs ./operations/word_count/mruserlib
```

**Parameters:**
- `-in`: Input file path(s), comma-separated for multiple files
- `-out`: Output file path
- `-mappers`: Maximum number of concurrent map tasks
- `-reducers`: Number of reduce tasks (determines output partitioning)
- `-srcs`: Path to the `mruserlib` directory containing your Map/Reduce code
- `-detached`: Optional flag to run job in background

## Writing Custom MapReduce Operations

Create a new directory under `operations/` with a `mruserlib` subdirectory containing `lib.go`:

```go
package mruserlib

import "hash/fnv"

// Entry represents a key-value pair emitted by Map
type Entry struct {
    Key, Value string
}

// Map processes input records and emits key-value pairs
// in: channel of input lines (one per record)
// out: channel to emit Entry pairs (must close when done)
func Map(in <-chan string, out chan<- Entry) {
    for record := range in {
        // Process each record and emit key-value pairs
        out <- Entry{Key: "mykey", Value: "myvalue"}
    }
    close(out)
}

// Partition determines which reducer handles a key
// Must return a value in range [0, reducers)
func Partition(key string, reducers int) int {
    h := fnv.New32a()
    h.Write([]byte(key))
    return int(h.Sum32()) % reducers
}

// Reduce aggregates all values for a given key
// key: the key being reduced
// in: channel of values for this key
// out: channel to emit output lines (must close when done)
func Reduce(key string, in <-chan string, out chan<- string) {
    var result string
    for value := range in {
        // Aggregate values
        result += value
    }
    out <- result
    close(out)
}
```

### Included Examples

#### Word Count (`operations/word_count/mruserlib`)
Counts occurrences of each word across input files.

```bash
client mapreduce -in /input.txt -out /wordcount.txt -mappers 3 -reducers 2 \
  -srcs ./operations/word_count/mruserlib
```

#### Inverted Index (`operations/inverted_index/mruserlib`)
Builds an inverted index mapping words to document titles.

Input format: `document_title<TAB>text content`

```bash
client mapreduce -in /documents.txt -out /index.txt -mappers 3 -reducers 2 \
  -srcs ./operations/inverted_index/mruserlib
```

#### Shuffle Merge (`operations/shuffle_merge/mruserlib`)
Shuffles and merges key-value pairs by key.

## Configuration

### Coordinator

```
coordinator start [flags]
```

| Flag | Description | Default |
|------|-------------|---------|
| `-name` | Name for coordinator machine and its folder | `coordinator` |
| `-port` | Port for running coordinator on | `11000` |
| `-env` | `dev` (local development) or `prod` (production) | `dev` |
| `-coordinators` | Comma separated coordinator IP addresses | |
| `-workers` | Comma separated worker IP addresses | |
| `-schedulers` | Comma separated scheduler IP addresses | |
| `-logger` | Logger IP address | |
| `-override` | Override config.json | `false` |

### Worker

```
worker start [flags]
```

| Flag | Description | Default |
|------|-------------|---------|
| `-name` | Name for worker machine and its folder | `worker` |
| `-port` | Port for running worker on | `11001` |
| `-env` | `dev` (local development) or `prod` (production) | `dev` |
| `-coordinator` | Coordinator IP address | |
| `-logger` | Logger IP address | |
| `-scheduler` | Start worker in scheduler mode | `false` |
| `-override` | Override config.json | `false` |

### Simple Logger

```
simple_logger start [flags]
```

| Flag | Description | Default |
|------|-------------|---------|
| `-name` | Name for logger machine and its folder | `logger` |
| `-port` | Port for running logger on | `11100` |
| `-env` | `dev` (local development) or `prod` (production) | `dev` |
| `-output` | Path for logger output file | `log` |
| `-override` | Override config.json | `false` |

### Client Commands

**Environment Variables:**
- `MR_HOSTS` - Comma-separated coordinator addresses (required)
- `MR_CLIENT` - Client identifier for logging (optional)

**Commands:**

| Command | Flags | Description |
|---------|-------|-------------|
| `write` | `-path`, `-append` | Write stdin to path |
| `read` | `-path` | Read path to stdout |
| `exists` | `-path` | Check if path exists |
| `list` | `-prefix` | List files with prefix |
| `remove` / `rm` | `-path` | Remove file |
| `copy` / `cp` | `-src`, `-dst` | Copy file |
| `move` / `mv` | `-src`, `-dst` | Move file |
| `merge` | `-in`, `-out` | Merge files (comma-separated input) |
| `mapreduce` | `-in`, `-out`, `-mappers`, `-reducers`, `-srcs`, `-detached` | Run MapReduce job |

### System Parameters

The coordinator defines these system-wide parameters (hardcoded):

| Parameter | Description | Default |
|-----------|-------------|---------|
| `MaxRowLength` | Maximum length of a single row | 256 KB |
| `BlockSize` | Target size for data blocks | 4 MB |
| `ReplicationFactor` | Number of block replicas | 3 |

## Fault Tolerance

### Coordinator Failover
- Multiple coordinator nodes with active/passive roles
- Journal replication for metadata consistency
- Automatic leader election on active coordinator failure

### Block Replication
- Configurable replication factor (default: 3)
- Automatic re-replication when workers fail
- Background monitoring of block health

### Transaction Management
- Transactions expire after timeout
- Background cleanup of abandoned transactions
- Heartbeat mechanism for long-running operations

## Visualizer

The visualizer provides real-time visualization of cluster activity:

1. Start the cluster with logging: `./run.sh 0`
2. Open `visualizer/index.html` in a browser
3. Load the request log file (found at `$MR_PATH/logger/requests.log`)
4. Click "Start" to begin visualization

Features:
- Animated packet flow between nodes
- Color-coded node status (green=up, red=down)
- Live activity log
- Recording capability (PNG frames)

## Project Structure

```
mapreduce/
├── client/           # CLI client
│   ├── client.go     # Main entry point
│   └── commands/     # Command implementations
├── coordinator/      # Coordinator node
│   ├── coordinator.go
│   └── server/       # HTTP handlers and core logic
│       └── dbase/    # BoltDB wrapper and journaling
├── worker/           # Worker node
│   ├── worker.go
│   └── server/       # HTTP handlers and operations
├── common/           # Shared utilities
│   ├── fsutil/       # File system utilities
│   ├── httputil/     # HTTP client helpers
│   ├── responses/    # API response types
│   └── wrrors/       # Error wrapping
├── operations/       # Example MapReduce operations
│   ├── word_count/
│   ├── inverted_index/
│   └── shuffle_merge/
├── template/         # MapReduce job template
│   ├── build.sh      # Build script for user code
│   └── main/         # Job executor
├── simple_logger/    # Request logging service
├── visualizer/       # Web-based visualization
└── run/              # Scripts for local development
    ├── run.sh        # Start local cluster
    └── kill.sh       # Stop all processes
```

## License

This project is provided for educational and research purposes.
