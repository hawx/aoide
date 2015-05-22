package cmd

import (
	"hawx.me/code/aoide/tools/list"
	"hawx.me/code/hadfield"
)

func List(lister *list.Lister) *hadfield.Command {
	var (
		artist string
		album  string
		track  string
		after  string
		before string
	)

	cmd := &hadfield.Command{
		Usage: "list",
		Short: "list items in the library",
		Long: `
  List items in the library.

    --artist PRED   # Show if artist matches PRED
    --album PRED    # Show if album matches PRED
    --track PRED    # Show if track matches PRED
    --after DUR     # Show if added after duration
    --before DUR    # Show if added before duration
`,
		Run: func(cmd *hadfield.Command, args []string) {
			lister.List(artist, album, track, after, before)
		},
	}

	cmd.Flag.StringVar(&artist, "artist", "", "")
	cmd.Flag.StringVar(&album, "album", "", "")
	cmd.Flag.StringVar(&track, "track", "", "")
	cmd.Flag.StringVar(&after, "after", "", "")
	cmd.Flag.StringVar(&before, "before", "", "")

	return cmd
}
