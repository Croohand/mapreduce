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
	port := startCommand.Int("port", 11001, "Port for running slave on")
	name := startCommand.String("name", "slave", "Name for slave machine and its folder")
	masterAddr := startCommand.String("master", "", "Master IP address")
	override := startCommand.Bool("override", false, "Override config.json")
	scheduler := startCommand.Bool("scheduler", false, "Start slave in scheduler mode")
	commandInfo := flagutil.CommandInfo{Name: "slave", Subcommands: []*flag.FlagSet{startCommand}}

	flagutil.CheckArgs(commandInfo)
	switch flagutil.Parse(commandInfo) {
	case startCommand:
		wrrors.SetSubject(*name)
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		server.Config = server.SlaveConfig{Name: *name, Port: *port, MasterAddr: *masterAddr, Scheduler: *scheduler}
		osutil.Init(*name, *override, &server.Config)
		if *name != server.Config.Name {
			panic("Name in config doesn't match with folder name")
		}
		server.Run()
	}
}
