package cmd

import (
	"log"

	"github.com/hawx/aoide/tools/organise"
	"github.com/hawx/hadfield"
)

func Organise(organiser *organise.Organiser) *hadfield.Command {
	var (
		printMusic     bool
		printPlaylists bool
		verbose        bool
	)

	cmd := &hadfield.Command{
		Usage: "organise [options]",
		Short: "organises the music directory",
		Long: `
  Organises the music directory into a standard "Artist/Album/00 Title.mp3"
  naming style.

    --print-music      # Print changes that will be made to musicDir
    --print-playlists  # Print changes that will be made to playlistDir
    --verbose          # Display more output when organising
`,
		Run: func(cmd *hadfield.Command, args []string) {
			if printMusic {
				organiser.Changes()
				return
			}

			if printPlaylists {
				if err := organiser.PlaylistChanges(); err != nil {
					log.Fatal(err)
				}
				return
			}

			if err := organiser.Organise(verbose); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flag.BoolVar(&printMusic, "print-music", false, "")
	cmd.Flag.BoolVar(&printPlaylists, "print-playlists", false, "")
	cmd.Flag.BoolVar(&verbose, "verbose", false, "")

	return cmd
}
