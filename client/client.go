package main

import (
	"flag"

	"github.com/Croohand/mapreduce/client/commands"
	"github.com/Croohand/mapreduce/common/flagutil"
)

func main() {
	commands.Init()

	writeCommand := flag.NewFlagSet("write", flag.ExitOnError)
	path := writeCommand.String("path", "", "MR path to write")

	commandInfo := flagutil.CommandInfo{Name: "client", Subcommands: []*flag.FlagSet{writeCommand}}

	flagutil.CheckArgs(commandInfo)
	switch flagutil.Parse(commandInfo) {
	case writeCommand:
		commands.Write(*path)
	}
}
