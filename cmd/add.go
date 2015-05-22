package cmd

import (
	"hawx.me/code/aoide/tools/add"
	"hawx.me/code/hadfield"
)

func Add(adder *add.Adder) *hadfield.Command {
	var queue bool

	cmd := &hadfield.Command{
		Usage: "add [options] FILES...",
		Short: "add tracks to the library",
		Long: `
  Adds the specified tracks to the music directory, updating the index.

    --queue      # Add tracks to Playlist
`,
		Run: func(cmd *hadfield.Command, args []string) {
			adder.Add(queue, args...)
		},
	}

	cmd.Flag.BoolVar(&queue, "queue", false, "")

	return cmd
}
