package data

type Song struct {
	Path string

	Album,
	AlbumArtist,
	Artist,
	Composer,
	Date,
	Genre,
	Title string

	Disc,
	DiscTotal,
	Length,
	Track,
	TrackTotal int
}

func (s Song) DisplayArtist() string {
	if s.AlbumArtist != "" {
		return s.AlbumArtist
	}
	return s.Artist
}
