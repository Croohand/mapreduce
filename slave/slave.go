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
	port := startCommand.Int("port", 11001, "port for running master on")
	name := startCommand.String("name", "slave", "name for slave machine and its folder")
	masterAddr := startCommand.String("master", "", "master IP address")
	override := startCommand.Bool("override", false, "override config.json")
	commandInfo := flagutil.CommandInfo{Name: "slave", Subcommands: []*flag.FlagSet{startCommand}}

	flagutil.CheckArgs(commandInfo)
	switch flagutil.Parse(commandInfo) {
	case startCommand:
		wrrors.SetSubject(*name)
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		server.Config = server.SlaveConfig{Name: *name, Port: *port, MasterAddr: *masterAddr}
		osutil.Init(*name, *override, &server.Config)
		if *name != server.Config.Name {
			panic("name in config doesn't match with folder name")
		}
		server.Run()
	}
}
