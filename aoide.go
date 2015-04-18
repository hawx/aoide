package main

import (
	"flag"
	"log"

	"github.com/hawx/aoide/cmd"
	"github.com/hawx/aoide/data"
	"github.com/hawx/aoide/tools/add"
	"github.com/hawx/aoide/tools/index"
	"github.com/hawx/aoide/tools/organise"
	"github.com/hawx/hadfield"
	"github.com/stvp/go-toml-config"
)

var (
	configFile = flag.String("config", "config.toml", "")
	help       = flag.Bool("help", false, "")

	musicDir    = config.String("musicDir", "")
	playlistDir = config.String("playlistDir", "")
	dbFile      = config.String("dbFile", "")
)

var templates = hadfield.Templates{
	Usage: `Usage: aoide [command] [arguments]

  A music library manager.

  Commands: {{range .}}
    {{.Name | printf "%-15s"}} # {{.Short}}{{end}}
`,
	Help: `usage: test {{.Usage}}
{{.Long}}
`,
}

func main() {
	flag.Parse()

	if *help {
		log.Println("Usage: aoide COMMAND [options]")
	}

	if err := config.Parse(*configFile); err != nil {
		log.Fatal("toml:", err)
	}

	if musicDir == nil || playlistDir == nil || dbFile == nil {
		log.Fatal("toml: must define musicDir, playlistDir and dbFile")
	}

	db, err := data.Open(*dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	adder := add.New(*musicDir, db)
	indexer := index.New(*musicDir, db)
	organiser := organise.New(*musicDir, *playlistDir, db)

	var commands = hadfield.Commands{
		cmd.Add(adder),
		cmd.Index(indexer),
		cmd.Organise(organiser),
	}

	hadfield.Run(commands, templates)
}
