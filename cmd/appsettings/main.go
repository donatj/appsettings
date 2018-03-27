package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/donatj/appsettings"
	"github.com/google/subcommands"
)

func main() {
	x, err := appsettings.NewAppSettings("appsettings.json")
	if err != nil {
		log.Fatal(err)
	}

	subcommands.Register(&getCmd{settings: x}, "")
	subcommands.Register(&setCmd{settings: x}, "")
	subcommands.Register(&deleteCmd{settings: x}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))

}
