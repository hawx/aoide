package add

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/fhs/gompd/mpd"
	"hawx.me/code/aoide/data"
	"hawx.me/code/aoide/song"
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

func (a *Adder) Add(queue bool, paths ...string) (err error) {
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
		a.doAdd(queue, recs)
	} else {
		log.Println("skipping...")
	}

	return
}

func (a *Adder) doAdd(queue bool, recs []addRecord) error {
	var toQueue []string

	for _, rec := range recs {
		if err := os.MkdirAll(filepath.Dir(rec.to), 0755); err != nil {
			log.Printf("Could not add %s: %s\n", rec.from, err)
			continue
		}

		if err := copyFile(rec.from, rec.to); err != nil {
			log.Printf("Could not add %s: %s\n", rec.from, err)
			continue
		}

		s, err := song.Read(rec.to)
		if err != nil {
			log.Printf("Could not add %s to db: %s\n", rec.from, err)
			continue
		}

		if err := a.db.Insert(s); err != nil {
			log.Printf("Could not add %s to db: %s\n", rec.from, err)
			continue
		}

		if a.client != nil {
			p, _ := filepath.Rel(a.musicDir, rec.to)
			if _, err := a.client.Update(p); err != nil {
				log.Printf("Could not add %s to mpd: %s\n", rec.from, err)
				continue
			}

			if queue {
				toQueue = append(toQueue, p)
			}
		}
	}

	if len(toQueue) > 0 {
		time.Sleep(time.Second) // ugh

		for _, p := range toQueue {
			if err := a.client.Add(p); err != nil {
				log.Printf("Could not add %s to playlist: %s\n", p, err)
				continue
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	if _, err := os.Stat(dst); err == nil || os.IsExist(err) {
		return errors.New("destination already exists")
	}

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
