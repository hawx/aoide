package song

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/nicksellen/audiotags"
	"hawx.me/code/aoide/data"
)

func Read(path string) (data.Song, error) {
	props, audioProps, err := audiotags.Read(path)
	if err != nil {
		return data.Song{}, err
	}

	disc, discTotal := splitNumber(props["discnumber"])
	if v, ok := props["disctotal"]; ok {
		if t, err := strconv.Atoi(v); err == nil {
			discTotal = t
		}
	}

	track, trackTotal := splitNumber(props["tracknumber"])
	if v, ok := props["tracktotal"]; ok {
		if t, err := strconv.Atoi(v); err == nil {
			trackTotal = t
		}
	}

	albumArtist := props["albumartist"]
	if comp, ok := props["compilation"]; ok && comp == "1" && albumArtist == "" {
		albumArtist = "Various Artists"
	}

	return data.Song{
		Path:        path,
		Album:       props["album"],
		AlbumArtist: albumArtist,
		Artist:      props["artist"],
		Composer:    props["composer"],
		Date:        props["date"],
		Disc:        disc,
		DiscTotal:   discTotal,
		Genre:       props["genre"],
		Length:      audioProps.Length,
		Title:       props["title"],
		Track:       track,
		TrackTotal:  trackTotal,
	}, nil
}

func PathFor(musicDir string, song data.Song) string {
	artist := "Unknown Artist"
	album := "Unknown Album"
	track := ""

	if song.AlbumArtist != "" {
		artist = sanitise(song.AlbumArtist)
	} else if song.Artist != "" {
		artist = sanitise(song.Artist)
	}

	if song.Album != "" {
		album = sanitise(song.Album)
	}

	if song.Disc > 0 && song.DiscTotal > 1 {
		track += strconv.Itoa(song.Disc)
		if song.Track > 0 {
			track += "-"
		} else {
			track += " "
		}
	}

	if song.Track > 0 {
		track += fmt.Sprintf("%02d", song.Track)
		track += " "
	}

	if song.Title != "" {
		track += sanitise(song.Title)
	} else {
		track += "Unknown Title"
	}

	track += filepath.Ext(song.Path)

	return musicDir + "/" + artist + "/" + album + "/" + track
}

func splitNumber(s string) (n int, total int) {
	parts := strings.Split(s, "/")
	var err error

	switch len(parts) {
	case 1:
		n, err = strconv.Atoi(parts[0])
		if err != nil {
			return -1, -1
		}
		return n, -1

	case 2:
		n, err = strconv.Atoi(parts[0])
		if err != nil {
			return -1, -1
		}

		total, err = strconv.Atoi(parts[1])
		if err != nil {
			return n, -1
		}

		return n, total
	}

	return -1, -1
}

var sanitiseRe = regexp.MustCompile("[\\/;:\\*\\?\"]")

func sanitise(name string) string {
	return sanitiseRe.ReplaceAllString(name, "_")
}
