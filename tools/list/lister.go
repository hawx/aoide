package list

import (
	"fmt"
	"time"

	"hawx.me/code/aoide/data"
)

func New(db *data.Database) *Lister {
	return &Lister{db}
}

type Lister struct {
	db *data.Database
}

func (l *Lister) List(artist, album, track, after, before string) {
	var clauses []data.Clause
	if artist != "" {
		clauses = append(clauses, data.PrefixClause{"Artist", artist})
	}
	if album != "" {
		clauses = append(clauses, data.PrefixClause{"Album", album})
	}
	if track != "" {
		clauses = append(clauses, data.PrefixClause{"Title", track})
	}

	if after != "" {
		dur, _ := time.ParseDuration(after)
		t := time.Now().Add(dur)
		clauses = append(clauses, data.GreaterThanClause{"Updated", t.Unix()})
	}
	if before != "" {
		dur, _ := time.ParseDuration(before)
		t := time.Now().Add(dur)
		clauses = append(clauses, data.LessThanClause{"Updated", t.Unix()})
	}

	l.db.Each(func(s data.Song) error {
		str := s.Title

		if s.Artist != "" {
			str += " by " + s.Artist
		} else if s.AlbumArtist != "" {
			str += " by " + s.AlbumArtist
		}

		if s.Album != "" {
			str += " on " + s.Album
		}

		if s.Date != "" {
			str += " [" + s.Date + "]"
		}

		fmt.Println(str)

		return nil
	}, clauses...)
}
