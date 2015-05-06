package playlist_test

import (
	"bytes"
	"strings"
	"testing"

	playlist "."
	"github.com/stretchr/testify/assert"
)

func TestPlaylistNotExtended(t *testing.T) {
	for _, c := range []string{"", "\n"} {
		p, err := playlist.New(strings.NewReader(c))

		assert.Nil(t, err)
		assert.False(t, p.IsExtended)
		assert.Equal(t, []playlist.Track{}, p.Tracks)
	}
}

func TestPlaylistExtended(t *testing.T) {
	for _, c := range []string{"#EXTM3U", "#EXTM3U\n"} {
		p, err := playlist.New(strings.NewReader(c))

		assert.Nil(t, err)
		assert.True(t, p.IsExtended)
		assert.Equal(t, []playlist.Track{}, p.Tracks)
	}
}

func TestPlaylist(t *testing.T) {
	p, err := playlist.New(strings.NewReader("my_song.mp3\n"))

	assert.Nil(t, err)
	assert.False(t, p.IsExtended)
	assert.Equal(t, []playlist.Track{{Val: "my_song.mp3"}}, p.Tracks)
}

func TestPlaylistWithMultipleTracks(t *testing.T) {
	text := `# my cool list
my_song.mp3

/some/file/song.m4a

http://example.com/stream.mp3
/some/other/file.mp3
song.mp3

# comment
  # comment
`

	p, err := playlist.New(strings.NewReader(text))

	assert.Nil(t, err)
	assert.False(t, p.IsExtended)
	assert.Equal(t, []playlist.Track{
		{Val: "my_song.mp3"},
		{Val: "/some/file/song.m4a"},
		{Val: "http://example.com/stream.mp3"},
		{Val: "/some/other/file.mp3"},
		{Val: "song.mp3"},
	}, p.Tracks)
}

func TestPlaylistExtendedWithMultipleTracks(t *testing.T) {
	text := `#EXTM3U
my_song.mp3

#EXTINF this is the best
/some/file/song.m4a

#NOTEXTINF
http://example.com/stream.mp3

#EXTINF:123;Song - Band
/some/other/file.mp3
song.mp3

# comment
  # comment
`

	p, err := playlist.New(strings.NewReader(text))

	assert.Nil(t, err)
	assert.True(t, p.IsExtended)
	assert.Equal(t, []playlist.Track{
		{Val: "my_song.mp3"},
		{Val: "/some/file/song.m4a", Inf: " this is the best"},
		{Val: "http://example.com/stream.mp3"},
		{Val: "/some/other/file.mp3", Inf: ":123;Song - Band"},
		{Val: "song.mp3"},
	}, p.Tracks)
}

func TestPlaylistWrite(t *testing.T) {
	p := playlist.Playlist{
		IsExtended: false,
		Tracks: []playlist.Track{
			{Val: "my_song.mp3"},
			{Val: "other.m4a", Inf: ":456;song - artist"},
		},
	}

	expectedText := `
my_song.mp3
other.m4a`

	buf := new(bytes.Buffer)
	assert.Nil(t, p.Write(buf))
	assert.Equal(t, expectedText, buf.String())
}

func TestPlaylistWriteExtended(t *testing.T) {
	p := playlist.Playlist{
		IsExtended: true,
		Tracks: []playlist.Track{
			{Val: "my_song.mp3"},
			{Val: "other.m4a", Inf: ":456;song - artist"},
		},
	}

	expectedText := `#EXTM3U
my_song.mp3
#EXTINF:456;song - artist
other.m4a`

	buf := new(bytes.Buffer)
	assert.Nil(t, p.Write(buf))
	assert.Equal(t, expectedText, buf.String())
}
