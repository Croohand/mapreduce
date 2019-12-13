package main

import (
	"flag"
	"log"

	"github.com/Croohand/mapreduce/common/flagutil"
	"github.com/Croohand/mapreduce/common/osutil"
	"github.com/Croohand/mapreduce/common/wrrors"
	"github.com/Croohand/mapreduce/simple_logger/server"
)

func main() {
	startCommand := flag.NewFlagSet("start", flag.ExitOnError)
	env := startCommand.String("env", "dev", "dev — local development, prod — production")
	port := startCommand.Int("port", 11100, "Port for running logger on")
	name := startCommand.String("name", "logger", "Name for logger machine and its folder")
	outputFile := startCommand.String("output", "log", "Path for logger output file")
	override := startCommand.Bool("override", false, "Override config.json")
	commandInfo := flagutil.CommandInfo{Name: "simple_logger", Subcommands: []*flag.FlagSet{startCommand}}

	flagutil.CheckArgs(commandInfo)
	switch flagutil.Parse(commandInfo) {
	case startCommand:
		wrrors.SetSubject(*name)
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		server.Config = server.LoggerConfig{
			Env:        *env,
			Name:       *name,
			Port:       *port,
			OutputFile: *outputFile,
		}
		osutil.Init(*name, *override, &server.Config)
		if *name != server.Config.Name {
			panic("Name in config doesn't match with folder name")
		}
		server.Run()
	}
}
