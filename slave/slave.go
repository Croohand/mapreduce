package main

import (
	"flag"
	"log"

	"github.com/Croohand/mapreduce/common/flagutil"
	"github.com/Croohand/mapreduce/common/osutil"
	"github.com/Croohand/mapreduce/common/wrrors"
	"github.com/Croohand/mapreduce/slave/server"
)

func main() {
	startCommand := flag.NewFlagSet("start", flag.ExitOnError)
	env := startCommand.String("env", "dev", "dev (default) — local development, prod — production")
	port := startCommand.Int("port", 11001, "Port for running slave on")
	name := startCommand.String("name", "slave", "Name for slave machine and its folder")
	masterAddr := startCommand.String("master", "", "Master IP address")
	loggerAddr := startCommand.String("logger", "", "Logger IP address")
	override := startCommand.Bool("override", false, "Override config.json")
	scheduler := startCommand.Bool("scheduler", false, "Start slave in scheduler mode")
	commandInfo := flagutil.CommandInfo{Name: "slave", Subcommands: []*flag.FlagSet{startCommand}}

	flagutil.CheckArgs(commandInfo)
	switch flagutil.Parse(commandInfo) {
	case startCommand:
		wrrors.SetSubject(*name)
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		server.Config = server.SlaveConfig{
			Env:        *env,
			Name:       *name,
			Port:       *port,
			MasterAddr: *masterAddr,
			LoggerAddr: *loggerAddr,
			Scheduler:  *scheduler,
		}
		osutil.Init(*name, *override, &server.Config)
		if *name != server.Config.Name {
			panic("Name in config doesn't match with folder name")
		}
		server.Run()
	}
}
