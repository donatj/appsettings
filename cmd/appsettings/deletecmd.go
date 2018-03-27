package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/donatj/appsettings"
	"github.com/google/subcommands"
)

type deleteCmd struct {
	settings *appsettings.AppSettings
}

func (*deleteCmd) Name() string     { return "delete" }
func (*deleteCmd) Synopsis() string { return "Delete a key." }
func (*deleteCmd) Usage() string {
	return `delete [<key>...]:
	Delete a key.
  `
}

func (p *deleteCmd) SetFlags(f *flag.FlagSet) {}

func (p *deleteCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

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

		tree.Delete(parts[len(parts)-1])
	}

	err := p.settings.Persist()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to persist settings: %s\n", err)

		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
