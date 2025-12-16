package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/Croohand/mapreduce/common/flagutil"
	"github.com/Croohand/mapreduce/common/osutil"
	"github.com/Croohand/mapreduce/common/wrrors"
	"github.com/Croohand/mapreduce/coordinator/server"
)

func main() {
	startCommand := flag.NewFlagSet("start", flag.ExitOnError)
	env := startCommand.String("env", "dev", "dev — local development, prod — production")
	port := startCommand.Int("port", 11000, "Port for running coordinator on")
	name := startCommand.String("name", "coordinator", "Name for coordinator machine and its folder")
	workerAddrs := startCommand.String("workers", "", "Comma separated worker IP addresses")
	coordinatorAddrs := startCommand.String("coordinators", "", "Comma separated coordinator IP addresses")
	schedulerAddrs := startCommand.String("schedulers", "", "Comma separated scheduler IP addresses")
	loggerAddr := startCommand.String("logger", "", "Logger IP address")
	override := startCommand.Bool("override", false, "Override config.json")
	commandInfo := flagutil.CommandInfo{Name: "coordinator", Subcommands: []*flag.FlagSet{startCommand}}

	flagutil.CheckArgs(commandInfo)
	switch flagutil.Parse(commandInfo) {
	case startCommand:
		wrrors.SetSubject(*name)
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		server.Config = server.CoordinatorConfig{
			Env:            *env,
			Name:           *name,
			Port:           *port,
			WorkerAddrs:     strings.Split(*workerAddrs, ","),
			CoordinatorAddrs:    strings.Split(*coordinatorAddrs, ","),
			SchedulerAddrs: strings.Split(*schedulerAddrs, ","),
			LoggerAddr:     *loggerAddr,
			LastJournalTs:  time.Now(),
		}
		osutil.Init(*name, *override, &server.Config)
		if *name != server.Config.Name {
			panic("Name in config doesn't match with folder name")
		}
		server.Run()
	}
}
