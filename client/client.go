package main

import (
	"flag"
	"log"

	"github.com/Croohand/mapreduce/client/commands"
	"github.com/Croohand/mapreduce/common/flagutil"
	"github.com/Croohand/mapreduce/common/wrrors"
)

func main() {
	wrrors.SetSubject("client")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	commands.Init()

	writeCommand := flag.NewFlagSet("write", flag.ExitOnError)
	writePath := writeCommand.String("path", "", "MR path to write")

	readCommand := flag.NewFlagSet("read", flag.ExitOnError)
	readPath := readCommand.String("path", "", "MR path to read")

	existsCommand := flag.NewFlagSet("exists", flag.ExitOnError)
	existsPath := existsCommand.String("path", "", "MR path to check")

	commandInfo := flagutil.CommandInfo{Name: "client", Subcommands: []*flag.FlagSet{existsCommand, readCommand, writeCommand}}

	flagutil.CheckArgs(commandInfo)
	switch flagutil.Parse(commandInfo) {
	case writeCommand:
		commands.Write(*writePath)
	case readCommand:
		commands.Read(*readPath)
	case existsCommand:
		commands.Exists(*existsPath)
	}
}
