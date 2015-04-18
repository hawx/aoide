package cmd

import (
	"log"

	"github.com/hawx/aoide/tools/index"
	"github.com/hawx/hadfield"
)

func Index(indexer *index.Indexer) *hadfield.Command {
	var (
		force  bool
		filter string
		remove string
	)

	cmd := &hadfield.Command{
		Usage: "index",
		Short: "indexes the music directory",
		Long: `
  Indexes the set music directory, storing tag data in a local database
  for fast access.

    --filter <path>    # Only index path
    --force            # Force re-indexing of files
    --remove <path>    # Remove files from index
`,
		Run: func(cmd *hadfield.Command, args []string) {
			if remove != "" {
				if err := indexer.Remove(remove); err != nil {
					log.Fatal(err)
				}

			} else {
				opts := index.Opts{Force: force, Path: filter}

				if err := indexer.Index(opts); err != nil {
					log.Fatal(err)
				}
			}
		},
	}

	cmd.Flag.BoolVar(&force, "force", false, "")
	cmd.Flag.StringVar(&filter, "filter", "", "")
	cmd.Flag.StringVar(&remove, "remove", "", "")

	return cmd
}
