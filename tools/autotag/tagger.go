package autotag

import (
	"errors"
	"fmt"
	"log"

	"github.com/bmatsuo/goline"
	"github.com/michiwend/gomusicbrainz"
	"hawx.me/code/aoide/data"
)

func New(musicDir string, db *data.Database) *Tagger {
	client, _ := gomusicbrainz.NewWS2Client(
		"https://musicbrainz.org/ws/2",
		"A GoMusicBrainz example",
		"0.0.1-beta",
		"http://github.com/michiwend/gomusicbrainz")

	return &Tagger{musicDir, db, client}
}

type Tagger struct {
	musicDir string
	db       *data.Database
	client   *gomusicbrainz.WS2Client
}

type Opts struct {
	Path string
}

func (t *Tagger) Autotag(opts Opts) error {
	var clauses []data.Clause
	if opts.Path != "" {
		clauses = []data.Clause{data.PrefixClause{"Path", opts.Path}}
	}

	var lastAlbum, lastArtist string

	return t.db.Each(func(s data.Song) error {
		artist := s.AlbumArtist
		if artist == "" {
			artist = s.Artist
		}

		if artist == lastArtist && s.Album == lastAlbum {
			return nil // skip
		}

		lastAlbum = s.Album
		lastArtist = artist

		resp, err := t.client.SearchRelease(`artist:"`+artist+`" AND release:"`+s.Album+`"`, 3, 0)
		if err != nil {
			return err
		}

		chosenIdx, _ := goline.Choose(func(menu *goline.Menu) {
			for i, release := range resp.Releases {
				itemname := fmt.Sprintf("%s by %s released in %s on %s [%d]\n",
					release.Title,
					release.ArtistCredit.NameCredit.Artist.Name,
					release.Date.Format("2006"),
					release.Mediums[0].Format,
					i)

				menu.Choice(itemname, func(name string, arg string) {
					log.Printf("got %s with %s\n", name, arg)
				})
			}
		})

		chosen := resp.Releases[chosenIdx]
		log.Printf("Chose: %s by %s released in %s on %s\n",
			chosen.Title,
			chosen.ArtistCredit.NameCredit.Artist.Name,
			chosen.Date.Format("2006"),
			chosen.Mediums[0].Format)

		if len(resp.ResultsWithScore(100)) == 0 {
			log.Printf("Could not find match for %s by %s\n", s.Album, artist)
			return nil
		}

		return errors.New("break")
	}, clauses...)
}
