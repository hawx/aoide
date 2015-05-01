package organise

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	"hawx.me/code/aoide/data"
	"hawx.me/code/aoide/diff"
	"hawx.me/code/aoide/playlist"
	"hawx.me/code/aoide/song"
)

func New(musicDir, playlistDir string, db *data.Database) *Organiser {
	return &Organiser{musicDir, playlistDir, db}
}

type Organiser struct {
	musicDir    string
	playlistDir string
	db          *data.Database
}

func (o *Organiser) Organise(verbose bool) (err error) {
	changes := map[string]string{}

	log.Println("Moving songs...")
	for oldPath, newPath := range o.changes() {
		if verbose {
			log.Printf("Move '%s' to '%s'", oldPath, newPath)
		}

		if err = os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
			log.Printf("Error moving '%s': %v", oldPath, err)
			continue
		}
		if err = os.Rename(oldPath, newPath); err != nil {
			log.Printf("Error moving '%s': %v", oldPath, err)
			continue
		}
		if err = o.db.Move(oldPath, newPath); err != nil {
			log.Println("Error moving '%s': %v", oldPath, err)
			os.Rename(newPath, oldPath) // put back to be consistent, if this errors who knows?
			continue
		}
		changes[oldPath] = newPath
	}

	log.Println("Calculating playlist changes...")
	// if there is an error at this point bad things!
	playlistChanges, _ := o.playlistChanges(changes)

	log.Println("Updating playlists...")
	for path, p := range playlistChanges {
		if verbose {
			log.Printf("Updating '%s'\n", path)
		}

		file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		p.Write(file)
		if err := file.Close(); err != nil {
			return err
		}
		if verbose {
			log.Printf("Wrote '%s'\n", path)
		}
	}

	log.Println("Done")

	return nil
}

func (o *Organiser) Changes() {
	changes := o.changes()
	keys := []string{}
	for k, _ := range changes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Println(diff.Path(k, changes[k]))
	}
}

func (o *Organiser) PlaylistChanges() error {
	changes, err := o.playlistChanges(o.changes())
	if err != nil {
		return err
	}

	for path, p := range changes {
		old, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		buf := new(bytes.Buffer)
		p.Write(buf)
		fmt.Printf("%s:\n%s\n\n", path, diff.Diff(string(old), buf.String()))
	}

	return nil
}

func (o *Organiser) playlistChanges(changes map[string]string) (map[string]playlist.Playlist, error) {
	m := map[string]playlist.Playlist{}

	err := filepath.Walk(o.playlistDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if matched, err := filepath.Match("*.m3u", filepath.Base(path)); !matched || err != nil {
			return nil
		}

		p, err := playlist.Read(path)
		if err != nil {
			return err
		}

		changed := false
		for i, track := range p.Tracks {
			if n, ok := changes[track.Val]; ok {
				p.Tracks[i].Val = n
				changed = true
			}
		}

		if !changed {
			return nil
		}

		m[path] = p
		return nil
	})

	return m, err
}

func (o *Organiser) changes() map[string]string {
	m := map[string]string{}

	o.db.Each(func(s data.Song) error {
		p := s.Path
		n := song.PathFor(o.musicDir, s)

		if p == n {
			return nil
		}

		m[p] = n
		return nil
	})

	return m
}
