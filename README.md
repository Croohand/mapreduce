# MapReduce

A distributed MapReduce framework with a Hadoop-like distributed file system (HDFS) implementation, written in Go.

## Overview

This project implements a fully functional distributed computing system featuring:

- **Distributed File System**: HDFS-like architecture with data blocks, replication, and fault tolerance
- **MapReduce Engine**: Execute custom Map and Reduce operations across distributed data
- **Master-Slave Architecture**: Coordinated cluster management with leader election
- **Scheduler**: Dedicated scheduler nodes for orchestrating MapReduce jobs
- **Visualizer**: Real-time visualization of cluster activity and request flows

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                 CLIENT                                       │
│                   (read, write, mapreduce, list, etc.)                      │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              MASTER NODES                                    │
│    ┌──────────┐    ┌──────────┐    ┌──────────┐                            │
│    │ Master 1 │◄──►│ Master 2 │◄──►│ Master 3 │   (Journal replication)    │
│    │ (Active) │    │(Passive) │    │(Passive) │                            │
│    └──────────┘    └──────────┘    └──────────┘                            │
│         │                                                                    │
│         ├── Metadata management (BoltDB)                                    │
│         ├── Block placement & replication                                   │
│         └── Transaction coordination                                        │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
               ┌─────────────────────┼─────────────────────┐
               ▼                     ▼                     ▼
┌──────────────────────┐  ┌──────────────────────┐  ┌──────────────────┐
│     SCHEDULER(S)     │  │      SLAVE NODES     │  │   SIMPLE LOGGER  │
│                      │  │                      │  │                  │
│  • Job orchestration │  │  • Block storage     │  │  • Request       │
│  • Map/Reduce        │  │  • Map execution     │  │    logging       │
│    coordination      │  │  • Reduce execution  │  │  • Visualization │
│  • Status tracking   │  │  • Data replication  │  │    support       │
└──────────────────────┘  └──────────────────────┘  └──────────────────┘
```

## Components

### Master

The master node manages:
- **Metadata Storage**: File paths to block mappings (using BoltDB)
- **Block Management**: Block placement, replication tracking, and recovery
- **Transaction Coordination**: ACID-like transactions for file operations
- **Leader Election**: Active/passive mode with automatic failover
- **Journal Replication**: Write-ahead logging replicated to passive masters

### Slave

Worker nodes responsible for:
- **Block Storage**: Storing and serving data blocks
- **Map Operations**: Executing user-defined Map functions on local data
- **Reduce Operations**: Executing user-defined Reduce functions
- **Block Replication**: Copying blocks to other slaves for fault tolerance

### Scheduler

Dedicated nodes for MapReduce job orchestration:
- **Job Preparation**: Setting up transactions and coordinating with master
- **Task Distribution**: Assigning map tasks to slaves with data locality
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
go install ./client
go install ./master
go install ./slave
go install ./simple_logger
```

The binaries will be installed to `$GOPATH/bin` (or `$HOME/go/bin` by default).

### Environment Setup

Set the `MR_PATH` environment variable for data storage:

```bash
export MR_PATH=/path/to/mapreduce/data
export MR_HOSTS=http://localhost:11000  # Comma-separated master addresses
export MR_CLIENT=client1  # Client identifier for logging
```

## Quick Start

### Starting a Local Cluster

Use the provided run script for local development:

```bash
cd run
./run.sh 0  # Start with logging enabled
```

This starts:
- 1 Master node (port 11000)
- 5 Slave nodes (ports 11001-11005)
- 1 Scheduler (port 11006)
- 1 Logger (port 11100)

### Stopping the Cluster

```bash
./kill.sh
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

# List files
client list -prefix /data/

# Check if file exists
client exists -path /data/myfile.txt

# Copy file
client copy -src /data/myfile.txt -dst /backup/myfile.txt
# Or use alias: client cp -src /data/myfile.txt -dst /backup/myfile.txt

# Move file
client move -src /data/myfile.txt -dst /archive/myfile.txt
# Or use alias: client mv -src /data/myfile.txt -dst /archive/myfile.txt

# Remove file
client remove -path /data/myfile.txt
# Or use alias: client rm -path /data/myfile.txt

# Merge multiple files
client merge -in /data/part1.txt,/data/part2.txt -out /data/merged.txt
```

### MapReduce Operations

```bash
client mapreduce \
  -in /data/input.txt \
  -out /data/output.txt \
  -mappers 5 \
  -reducers 3 \
  -srcs /path/to/operations/word_count/mruserlib

# Run in detached mode (returns immediately)
client mapreduce \
  -in /data/input.txt \
  -out /data/output.txt \
  -mappers 5 \
  -reducers 3 \
  -srcs /path/to/operations/word_count/mruserlib \
  -detached
```

## Writing Custom MapReduce Operations

Create a new directory under `operations/` with a `mruserlib` subdirectory containing `lib.go`:

```go
package mruserlib

// Entry represents a key-value pair emitted by Map
type Entry struct {
    Key, Value string
}

// Map processes input records and emits key-value pairs
// in: channel of input lines (one per record)
// out: channel to emit Entry pairs
func Map(in <-chan string, out chan<- Entry) {
    for record := range in {
        // Process each record
        // Emit key-value pairs:
        out <- Entry{Key: "mykey", Value: "myvalue"}
    }
    close(out)
}

// Partition determines which reducer handles a key
// Must return a value in range [0, reducers)
func Partition(key string, reducers int) int {
    // Hash-based partitioning example
    h := fnv.New32a()
    h.Write([]byte(key))
    return int(h.Sum32()) % reducers
}

// Reduce aggregates values for a key
// key: the key being reduced
// in: channel of values for this key
// out: channel to emit final output lines
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

Input format: `document_title\ttext content`

```bash
client mapreduce -in /documents.txt -out /index.txt -mappers 3 -reducers 2 \
  -srcs ./operations/inverted_index/mruserlib
```

#### Shuffle Merge (`operations/shuffle_merge/mruserlib`)
Shuffles and merges key-value pairs.

## Configuration

### Master Configuration

```bash
master start \
  -name master1 \
  -port 11000 \
  -env dev \
  -slaves http://localhost:11001,http://localhost:11002 \
  -masters http://localhost:11000,http://localhost:11010 \
  -schedulers http://localhost:11006 \
  -logger http://localhost:11100 \
  -override  # Override existing config
```

### Slave Configuration

```bash
slave start \
  -name slave1 \
  -port 11001 \
  -env dev \
  -master http://localhost:11000 \
  -logger http://localhost:11100 \
  -override
```

### Scheduler Configuration

```bash
slave start \
  -name scheduler1 \
  -port 11006 \
  -env dev \
  -master http://localhost:11000 \
  -logger http://localhost:11100 \
  -scheduler  # Enable scheduler mode
  -override
```

### System Parameters (MrConfig)

The system uses these configurable parameters:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `MaxRowLength` | Maximum length of a single row | Configured on master |
| `BlockSize` | Target size for data blocks (bytes) | Configured on master |
| `ReplicationFactor` | Number of block replicas | Configured on master |

## API Endpoints

### Master Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /IsAlive` | Check master status (active/passive) |
| `GET /GetMrConfig` | Get system configuration |
| `GET /GetAvailableSlaves` | Get list of available slave nodes |
| `POST /GetAvailableReducers` | Get reducer nodes for MapReduce |
| `POST /GetAvailableScheduler` | Get an available scheduler |
| `POST /Transaction/Start` | Start a new transaction |
| `POST /Transaction/Update` | Ping/update transaction |
| `POST /Transaction/Close` | Close transaction |
| `POST /File/Write` | Register file with blocks |
| `POST /File/Read` | Get block info for file |
| `POST /File/Remove` | Remove a file |
| `POST /File/Merge` | Merge multiple files |
| `GET /File/List` | List files by prefix |
| `GET /File/IsExists` | Check if file exists |

### Slave Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /IsAlive` | Health check |
| `POST /Block/Write` | Write block data |
| `POST /Block/Read` | Read block data |
| `POST /Block/Copy` | Replicate block |
| `POST /Block/Remove` | Remove block |
| `POST /Block/Check` | Check block existence |
| `POST /Block/Validate` | Validate block integrity |
| `POST /Operation/Map` | Execute map operation |
| `POST /Operation/Reduce` | Execute reduce operation |
| `POST /Source/Write` | Upload user code |
| `POST /Source/Build` | Build user code |

### Scheduler Endpoints

| Endpoint | Description |
|----------|-------------|
| `POST /Operation/PrepareMapReduce` | Prepare MapReduce job |
| `POST /Operation/StartMapReduce` | Start MapReduce execution |
| `POST /Operation/GetStatus` | Get job status |

## Fault Tolerance

### Master Failover
- Multiple master nodes with active/passive roles
- Journal replication for metadata consistency
- Automatic leader election on active master failure

### Block Replication
- Configurable replication factor
- Automatic re-replication when slaves fail
- Background monitoring of block health

### Transaction Management
- Transactions expire after timeout
- Background cleanup of abandoned transactions
- Heartbeat mechanism for long-running operations

## Visualizer

The visualizer provides real-time visualization of cluster activity:

1. Start the cluster with logging enabled
2. Open `visualizer/index.html` in a browser
3. Load the request log file (`requests.log`)
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
├── master/           # Master node
│   ├── master.go     # Main entry point
│   └── server/       # HTTP handlers and core logic
│       └── dbase/    # BoltDB wrapper and journaling
├── slave/            # Slave node
│   ├── slave.go      # Main entry point
│   └── server/       # HTTP handlers and operations
├── common/           # Shared utilities
│   ├── fsutil/       # File system utilities
│   ├── httputil/     # HTTP helpers
│   ├── responses/    # API response types
│   └── wrrors/       # Error wrapping
├── operations/       # Example MapReduce operations
│   ├── word_count/
│   ├── inverted_index/
│   └── shuffle_merge/
├── template/         # MapReduce job template
│   └── main/         # Job executor (do_map, do_reduce)
├── simple_logger/    # Request logging service
├── visualizer/       # Web-based visualization
└── run/              # Scripts for local development
```

## Adding New Operations

1. Create `operations/my_operation/mruserlib/lib.go`
2. Implement `Map`, `Partition`, and `Reduce` functions
3. Use with: `client mapreduce -srcs ./operations/my_operation/mruserlib`

## License

This project is provided for educational and research purposes.
