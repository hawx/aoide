package cmd

import (
	"log"

	"hawx.me/code/aoide/tools/autotag"
	"hawx.me/code/hadfield"
)

func Autotag(tagger *autotag.Tagger) *hadfield.Command {
	var filter string

	cmd := &hadfield.Command{
		Usage: "autotag",
		Short: "autotag files in library",
		Long: `
  Autotag files in library.

    --filter <path>     # Only autotag path
`,
		Run: func(cmd *hadfield.Command, args []string) {
			opts := autotag.Opts{Path: filter}

			if err := tagger.Autotag(opts); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flag.StringVar(&filter, "filter", "", "")

	return cmd
}
