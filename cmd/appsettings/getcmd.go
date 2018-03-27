package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/donatj/appsettings"
	"github.com/google/subcommands"
)

type getCmd struct {
	settings *appsettings.AppSettings
}

func (*getCmd) Name() string     { return "get" }
func (*getCmd) Synopsis() string { return "Print key value to stdout." }
func (*getCmd) Usage() string {
	return `get [-capitalize] <some text>:
	Print key value to stdout.
  `
}

func (p *getCmd) SetFlags(f *flag.FlagSet) {}

func (p *getCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	for _, arg := range f.Args() {
		parts, err := parsePath(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			return subcommands.ExitFailure
		}

		var tree appsettings.DataTree = p.settings
		for i := 0; i < len(parts)-1; i++ {
			tree = tree.GetTree(parts[i])
		}

		k, err := tree.GetString(parts[len(parts)-1])
		if err == appsettings.ErrUndefinedKey {
			fmt.Fprintf(os.Stderr, "undefined key: %s\n", arg)
			return subcommands.ExitFailure
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			return subcommands.ExitFailure
		}

		fmt.Println(k)
	}

	return subcommands.ExitSuccess
}
