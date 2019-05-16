package main

import (
	"flag"
	"log"
	"strings"

	"github.com/Croohand/mapreduce/common/flagutil"
	"github.com/Croohand/mapreduce/common/osutil"
	"github.com/Croohand/mapreduce/common/wrrors"
	"github.com/Croohand/mapreduce/master/server"
)

func main() {
	startCommand := flag.NewFlagSet("start", flag.ExitOnError)
	port := startCommand.Int("port", 11000, "Port for running master on")
	name := startCommand.String("name", "master", "Name for master machine and its folder")
	slaveAddrs := startCommand.String("slaves", "", "Comma separated slaves IP addresses")
	override := startCommand.Bool("override", false, "Override config.json")
	commandInfo := flagutil.CommandInfo{Name: "master", Subcommands: []*flag.FlagSet{startCommand}}

	flagutil.CheckArgs(commandInfo)
	switch flagutil.Parse(commandInfo) {
	case startCommand:
		wrrors.SetSubject(*name)
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		server.Config = server.MasterConfig{Name: *name, Port: *port, SlaveAddrs: strings.Split(*slaveAddrs, ",")}
		osutil.Init(*name, *override, &server.Config)
		if *name != server.Config.Name {
			panic("Name in config doesn't match with folder name")
		}
		server.Run()
	}
}
