package flagutil

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type CommandInfo struct {
	Name        string
	Subcommands []*flag.FlagSet
}

func usage(info CommandInfo) {
	fmt.Print("Usage: ", info.Name)
	names := make([]string, len(info.Subcommands))
	for i, com := range info.Subcommands {
		names[i] = com.Name()
	}
	switch len(names) {
	case 0:
		fmt.Println(" [options]")
	case 1:
		fmt.Printf(" %s [options]\n", names[0])
	default:
		fmt.Printf(" {%s} [options]\n", strings.Join(names, ","))
	}
}

func Parse(info CommandInfo) *flag.FlagSet {
	for _, com := range info.Subcommands {
		if os.Args[1] == com.Name() {
			com.Parse(os.Args[2:])
			return com
		}
	}
	usage(info)
	os.Exit(1)
	return nil
}

func CheckArgs(info CommandInfo) {
	if len(info.Subcommands) > 0 && len(os.Args) < 2 {
		usage(info)
		os.Exit(1)
	}
}
