package index

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hawx/aoide/data"
	"github.com/hawx/aoide/song"
)

func New(root string, db *data.Database) *Indexer {
	return &Indexer{root: root, db: db}
}

type Indexer struct {
	root string
	db   *data.Database
}

type Opts struct {
	Path  string
	Force bool
}

func (i *Indexer) Index(opts Opts) error {
	indexPath := i.root
	if opts.Path != "" {
		indexPath = opts.Path
	}

	return filepath.Walk(indexPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if !opts.Force && !i.db.IsStale(path, info.ModTime()) {
			fmt.Printf("\rskipping: %s", padRight(truncate(path, 80), 80))
			return nil
		}

		s, err := song.Read(path)
		if err != nil {
			return err
		}

		fmt.Printf("\radding: %s\n", padRight(truncate(path, 80), 80))
		return i.db.Insert(s)
	})
}

func (i *Indexer) Remove(path string) error {
	toRemove := []string{}

	err := i.db.Each(func(song data.Song) error {
		toRemove = append(toRemove, song.Path)
		return nil
	}, data.PrefixClause{"Path", path})

	if err != nil {
		return err
	}

	for _, t := range toRemove {
		fmt.Printf("removing: %s\n", t)
		if err := i.db.Delete(t); err != nil {
			return err
		}
	}

	return nil
}

func padRight(str string, length int) string {
	for {
		str += " "
		if len(str) > length {
			return str[0:length]
		}
	}
}

func truncate(str string, length int) string {
	if len(str) > 80 {
		return "..." + str[len(str)-77:]
	}

	return str
}
