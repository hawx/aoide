package playlist

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func Read(path string) (Playlist, error) {
	file, err := os.Open(path)
	if err != nil {
		return Playlist{}, err
	}
	return New(file)
}

const (
	ExtendedHeader = "#EXTM3U"
	InfDirective   = "#EXTINF"
)

func New(r io.Reader) (playlist Playlist, err error) {
	br := bufio.NewReader(r)
	playlist = Playlist{Tracks: []Track{}}

	header, err := br.Peek(len(ExtendedHeader))
	if err != nil && err != io.EOF {
		return
	}
	if err == io.EOF {
		return playlist, nil
	}

	if string(header) == ExtendedHeader {
		playlist.IsExtended = true
	}

	var line, inf string
	for {
		line, err = br.ReadString('\n')
		if err != nil && err != io.EOF {
			return
		}
		if err == io.EOF && len(line) == 0 {
			break
		}

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if line[0] == '#' {
			if playlist.IsExtended && line[:len(InfDirective)] == InfDirective {
				inf = line[len(InfDirective):]
			}
			continue
		}

		if inf != "" {
			playlist.Tracks = append(playlist.Tracks, Track{Val: line, Inf: inf})
			inf = ""
		} else {
			playlist.Tracks = append(playlist.Tracks, Track{Val: line})
		}
	}

	return playlist, nil
}

type Playlist struct {
	IsExtended bool
	Tracks     []Track
}

type Track struct {
	Val, Inf string // Inf only set when IsExtended is true.
}

func (p Playlist) Write(w io.Writer) error {
	if p.IsExtended {
		w.Write([]byte(ExtendedHeader))
	}

	for _, track := range p.Tracks {
		if p.IsExtended && track.Inf != "" {
			w.Write([]byte("\n" + InfDirective + track.Inf))
		}
		w.Write([]byte("\n" + track.Val))
	}

	return nil
}
