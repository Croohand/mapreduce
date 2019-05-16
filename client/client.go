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
	writePath := writeCommand.String("path", "", "MR path to write")
	doAppend := writeCommand.Bool("append", false, "Append to file")

	readCommand := flag.NewFlagSet("read", flag.ExitOnError)
	readPath := readCommand.String("path", "", "MR path to read")

	existsCommand := flag.NewFlagSet("exists", flag.ExitOnError)
	existsPath := existsCommand.String("path", "", "MR path to check")

	listCommand := flag.NewFlagSet("list", flag.ExitOnError)
	listPrefix := listCommand.String("prefix", "", "List files with given prefix")

	removeCommand := flag.NewFlagSet("remove", flag.ExitOnError)
	removePath := removeCommand.String("path", "", "MR path to remove")

	mergeCommand := flag.NewFlagSet("merge", flag.ExitOnError)
	mergeInPaths := mergeCommand.String("in", "", "MR input paths, comma separated")
	mergeOutPath := mergeCommand.String("out", "", "MR output path")

	copyCommand := flag.NewFlagSet("copy", flag.ExitOnError)
	copyInPath := copyCommand.String("src", "", "MR source path")
	copyOutPath := copyCommand.String("dst", "", "MR destination path")

	moveCommand := flag.NewFlagSet("move", flag.ExitOnError)
	moveInPath := moveCommand.String("src", "", "MR source path")
	moveOutPath := moveCommand.String("dst", "", "MR destination path")

	commandInfo := flagutil.CommandInfo{Name: "client", Aliases: aliases, Subcommands: []*flag.FlagSet{
		existsCommand,
		readCommand,
		writeCommand,
		mergeCommand,
		copyCommand,
		removeCommand,
		moveCommand,
		listCommand}}

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
