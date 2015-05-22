package main

import (
	"flag"
	"log"

	"github.com/stvp/go-toml-config"
	"hawx.me/code/aoide/cmd"
	"hawx.me/code/aoide/data"
	"hawx.me/code/aoide/tools/add"
	"hawx.me/code/aoide/tools/autotag"
	"hawx.me/code/aoide/tools/index"
	"hawx.me/code/aoide/tools/list"
	"hawx.me/code/aoide/tools/organise"
	"hawx.me/code/hadfield"
	"hawx.me/code/xdg"
)

var (
	help = flag.Bool("help", false, "")

	musicDir    = config.String("musicDir", "")
	playlistDir = config.String("playlistDir", "")
	dbFile      = config.String("dbFile", "")
)

var templates = hadfield.Templates{
	Help: `Usage: aoide [command] [arguments]

  A music library manager.

  By default will search for config at $XDG_CONFIG_HOME/aoide/rc, or in one of
  $XDG_CONFIG_DIRS called aoide/rc; otherwise tries to load ./config.toml. This
  configuartion file must include entries like:

    musicDir = "/path/to/music"
    playlistDir = /path/to/playlists"
    dbFile = "/path/to/db"

  Commands: {{range .}}
    {{.Name | printf "%-15s"}} # {{.Short}}{{end}}
`,
	Command: `usage: aoide {{.Usage}}
{{.Long}}
`,
}

func main() {
	flag.Parse()

	if *help {
		log.Println("Usage: aoide COMMAND [options]")
	}

	if err := config.Parse(xdg.Config("aoide/rc")); err != nil {
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
	tagger := autotag.New(*musicDir, db)
	lister := list.New(db)

	var commands = hadfield.Commands{
		cmd.Add(adder),
		cmd.Autotag(tagger),
		cmd.Index(indexer),
		cmd.List(lister),
		cmd.Organise(organiser),
	}

	hadfield.Run(commands, templates)
}
