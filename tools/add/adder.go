package add

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/fhs/gompd/mpd"
	"github.com/hawx/aoide/data"
	"github.com/hawx/aoide/song"
)

func New(musicDir string, db *data.Database) *Adder {
	conn, err := mpd.Dial("tcp", ":6600")
	if err != nil {
		conn = nil
	}

	return &Adder{musicDir, db, conn}
}

type Adder struct {
	musicDir string
	db       *data.Database
	client   *mpd.Client
}

type addRecord struct {
	from, to string
	s        data.Song
}

type byPath []addRecord

func (a byPath) Len() int           { return len(a) }
func (a byPath) Less(i, j int) bool { return a[i].to < a[j].to }
func (a byPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (a *Adder) Add(paths ...string) (err error) {
	recs := []addRecord{}
	for _, path := range paths {
		s, err := song.Read(path)
		if err != nil {
			log.Println(err)
			continue
		}
		recs = append(recs, addRecord{
			from: path,
			to:   song.PathFor(a.musicDir, s),
			s:    s,
		})
	}

	sort.Sort(byPath(recs))

	lastDir := ""
	for _, rec := range recs {
		thisDir := filepath.Dir(rec.to)

		if thisDir != lastDir {
			fmt.Println(thisDir)
			lastDir = thisDir
		}

		fmt.Println(" ", filepath.Base(rec.to))
	}

	fmt.Print("Add? [yn] ")
	var response string
	_, err = fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}

	if response == "y" || response == "Y" {
		log.Println("adding...")
		a.doAdd(recs)
	} else {
		log.Println("skipping...")
	}

	return
}

func (a *Adder) doAdd(recs []addRecord) error {
	for _, rec := range recs {
		if err := os.MkdirAll(filepath.Dir(rec.to), 0755); err != nil {
			log.Printf("Could not add %s: %s\n", rec.from, err)
			continue
		}

		if err := copyFile(rec.from, rec.to); err != nil {
			log.Printf("Could not add %s: %s\n", rec.from, err)
		}

		s, err := song.Read(rec.to)
		if err != nil {
			log.Printf("Could not add %s to db: %s\n", rec.from, err)
		}

		if err := a.db.Insert(s); err != nil {
			log.Printf("Could not add %s to db: %s\n", rec.from, err)
		}

		if a.client != nil {
			p, _ := filepath.Rel(a.musicDir, rec.to)
			if _, err := a.client.Update(p); err != nil {
				log.Printf("Could not add %s to mpd: %s\n", rec.from, err)
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return out.Close()
}
