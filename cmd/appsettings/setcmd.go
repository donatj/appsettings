package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/donatj/appsettings"
	"github.com/google/subcommands"
)

type setCmd struct {
	settings *appsettings.AppSettings
}

func (*setCmd) Name() string     { return "set" }
func (*setCmd) Synopsis() string { return "Set a keys value." }
func (*setCmd) Usage() string {
	return `set [[<key>, <value>]...]:
	Set a keys value.
  `
}

func (p *setCmd) SetFlags(f *flag.FlagSet) {}

func (p *setCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	args := f.Args()

	if len(args)%2 != 0 {
		fmt.Fprintf(os.Stderr, "arguments must be an even number of key and value pairs\n")
		return subcommands.ExitFailure
	}

	for i := 0; i <= len(args)-1; i += 2 {

		parts, err := parsePath(args[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			return subcommands.ExitFailure
		}

		var tree appsettings.DataTree = p.settings
		for i := 0; i < len(parts)-1; i++ {
			tree = tree.GetTree(parts[i])
		}

		tree.SetString(parts[len(parts)-1], args[i+1])
	}

	err := p.settings.Persist()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to persist settings: %s\n", err)

		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
