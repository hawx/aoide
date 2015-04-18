package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/hawx/aoide/cmd"
	"github.com/hawx/aoide/data"
	"github.com/hawx/aoide/tools/add"
	"github.com/hawx/aoide/tools/index"
	"github.com/hawx/aoide/tools/organise"
	"github.com/hawx/hadfield"
	"github.com/stvp/go-toml-config"
)

var (
	help = flag.Bool("help", false, "")

	musicDir    = config.String("musicDir", "")
	playlistDir = config.String("playlistDir", "")
	dbFile      = config.String("dbFile", "")
)

var templates = hadfield.Templates{
	Usage: `Usage: aoide [command] [arguments]

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
	Help: `usage: test {{.Usage}}
{{.Long}}
`,
}

func findConfigFile() string {
	const (
		configSuffix         = "/aoide/rc"
		defaultXdgConfigDirs = "/etc/xdg"
	)

	var (
		defaultXdgConfigHome = os.ExpandEnv("$HOME/.config")
		xdgConfigDirs        = os.Getenv("XDG_CONFIG_DIRS")
		xdgConfigHome        = os.Getenv("XDG_CONFIG_HOME")
	)

	if xdgConfigDirs == "" {
		xdgConfigDirs = defaultXdgConfigDirs
	}
	if xdgConfigHome == "" {
		xdgConfigHome = defaultXdgConfigHome
	}

	if _, err := os.Stat(xdgConfigHome + configSuffix); err == nil {
		return xdgConfigHome + configSuffix
	}

	for _, path := range strings.Split(xdgConfigDirs, ":") {
		if _, err := os.Stat(path + configSuffix); err == nil {
			return path + configSuffix
		}
	}

	return ""
}

func main() {
	flag.Parse()

	if *help {
		log.Println("Usage: aoide COMMAND [options]")
	}

	if err := config.Parse(findConfigFile()); err != nil {
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
