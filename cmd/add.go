package cmd

import (
	"hawx.me/code/aoide/tools/add"
	"hawx.me/code/hadfield"
)

func Add(adder *add.Adder) *hadfield.Command {
	return &hadfield.Command{
		Usage: "add FILES...",
		Short: "add tracks to the library",
		Long: `
  Adds the specified tracks to the music directory, updating the index.
`,
		Run: func(cmd *hadfield.Command, args []string) {
			adder.Add(args...)
		},
	}
}
