package main

import (
	"flag"
	"log"

	"github.com/Croohand/mapreduce/common/flagutil"
	"github.com/Croohand/mapreduce/common/osutil"
	"github.com/Croohand/mapreduce/common/wrrors"
	"github.com/Croohand/mapreduce/worker/server"
)

func main() {
	startCommand := flag.NewFlagSet("start", flag.ExitOnError)
	env := startCommand.String("env", "dev", "dev — local development, prod — production")
	port := startCommand.Int("port", 11001, "Port for running worker on")
	name := startCommand.String("name", "worker", "Name for worker machine and its folder")
	coordinatorAddr := startCommand.String("coordinator", "", "Coordinator IP address")
	loggerAddr := startCommand.String("logger", "", "Logger IP address")
	override := startCommand.Bool("override", false, "Override config.json")
	scheduler := startCommand.Bool("scheduler", false, "Start worker in scheduler mode")
	commandInfo := flagutil.CommandInfo{Name: "worker", Subcommands: []*flag.FlagSet{startCommand}}

	flagutil.CheckArgs(commandInfo)
	switch flagutil.Parse(commandInfo) {
	case startCommand:
		wrrors.SetSubject(*name)
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		server.Config = server.WorkerConfig{
			Env:        *env,
			Name:       *name,
			Port:       *port,
			CoordinatorAddr: *coordinatorAddr,
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
