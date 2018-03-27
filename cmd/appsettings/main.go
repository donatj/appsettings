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
	cget := &getCmd{}
	cset := &setCmd{}
	cdelete := &deleteCmd{}

	subcommands.Register(cget, "")
	subcommands.Register(cset, "")
	subcommands.Register(cdelete, "")

	f := flag.String("file", "appsettings.json", "Appsetting Database to Query")

	flag.Parse()

	settings, err := appsettings.NewAppSettings(*f)
	if err != nil {
		log.Fatal(err)
	}

	cget.settings = settings
	cset.settings = settings
	cdelete.settings = settings

	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
