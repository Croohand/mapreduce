package main

import (
	"flag"
	"log"
	"strings"

	"github.com/Croohand/mapreduce/client/commands"
	"github.com/Croohand/mapreduce/common/flagutil"
	"github.com/Croohand/mapreduce/common/wrrors"
)

var aliases = map[string]string{
	"mv": "move",
	"cp": "copy",
	"ls": "list",
	"rm": "remove"}

func main() {
	// Avoid printing stacktrace while panicking
	defer func() {
		recover()
	}()

	wrrors.SetSubject("client")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	commands.Init()

	writeCommand := flag.NewFlagSet("write", flag.ExitOnError)
	writePath := writeCommand.String("path", "", "HDFS path to write")
	doAppend := writeCommand.Bool("append", false, "Append to file")

	readCommand := flag.NewFlagSet("read", flag.ExitOnError)
	readPath := readCommand.String("path", "", "HDFS path to read")

	existsCommand := flag.NewFlagSet("exists", flag.ExitOnError)
	existsPath := existsCommand.String("path", "", "HDFS path to check")

	listCommand := flag.NewFlagSet("list", flag.ExitOnError)
	listPrefix := listCommand.String("prefix", "", "List files with given prefix")

	removeCommand := flag.NewFlagSet("remove", flag.ExitOnError)
	removePath := removeCommand.String("path", "", "HDFS path to remove")

	mergeCommand := flag.NewFlagSet("merge", flag.ExitOnError)
	mergeInPaths := mergeCommand.String("in", "", "HDFS input paths, comma separated")
	mergeOutPath := mergeCommand.String("out", "", "HDFS output path")

	mapReduceCommand := flag.NewFlagSet("mapreduce", flag.ExitOnError)
	mapReduceInPaths := mapReduceCommand.String("in", "", "HDFS input paths, comma separated")
	mapReduceOutPath := mapReduceCommand.String("out", "", "HDFS output path")
	mapReduceReducersNum := mapReduceCommand.Int("reducers", 0, "Number of reducers for HDFS operation")
	mapReduceMappersNum := mapReduceCommand.Int("mappers", 0, "Number of maximum running mappers for HDFS operation")
	mapReduceSrcsPath := mapReduceCommand.String("srcs", "", "Path to mruserlib")
	mapReduceDetached := mapReduceCommand.Bool("detached", false, "Detach from scheduler when operation starts")

	copyCommand := flag.NewFlagSet("copy", flag.ExitOnError)
	copyInPath := copyCommand.String("src", "", "HDFS source path")
	copyOutPath := copyCommand.String("dst", "", "HDFS destination path")

	moveCommand := flag.NewFlagSet("move", flag.ExitOnError)
	moveInPath := moveCommand.String("src", "", "HDFS source path")
	moveOutPath := moveCommand.String("dst", "", "HDFS destination path")

	commandInfo := flagutil.CommandInfo{Name: "client", Aliases: aliases, Subcommands: []*flag.FlagSet{
		existsCommand,
		readCommand,
		writeCommand,
		mergeCommand,
		copyCommand,
		removeCommand,
		moveCommand,
		listCommand,
		mapReduceCommand}}

	flagutil.CheckArgs(commandInfo)
	switch flagutil.Parse(commandInfo) {
	case writeCommand:
		commands.Write(*writePath, *doAppend)
	case readCommand:
		commands.Read(*readPath)
	case existsCommand:
		commands.Exists(*existsPath)
	case mergeCommand:
		commands.Merge(strings.Split(*mergeInPaths, ","), *mergeOutPath)
	case mapReduceCommand:
		commands.MapReduce(strings.Split(*mapReduceInPaths, ","), *mapReduceOutPath, *mapReduceSrcsPath, *mapReduceMappersNum, *mapReduceReducersNum, *mapReduceDetached)
	case copyCommand:
		commands.Copy(*copyInPath, *copyOutPath)
	case moveCommand:
		commands.Move(*moveInPath, *moveOutPath)
	case removeCommand:
		commands.Remove(*removePath)
	case listCommand:
		commands.List(*listPrefix)
	}
}
