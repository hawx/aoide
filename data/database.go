package data

import (
	"database/sql"
	"time"

	"code.google.com/p/go-sqlite/go1/sqlite3"
)

type Database struct {
	db *sql.DB
}

func Open(path string) (*Database, error) {
	sqlite, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	db := &Database{sqlite}
	return db, db.setup()
}

func (d *Database) setup() error {
	_, err := d.db.Exec(`CREATE TABLE IF NOT EXISTS songs (
    Id          INTEGER PRIMARY KEY,
    Path        TEXT UNIQUE,
    Updated     DATETIME,
    Album       TEXT,
    AlbumArtist TEXT,
    Artist      TEXT,
    Composer    TEXT,
    Date        TEXT,
    Disc        INTEGER,
    DiscTotal   INTEGER,
    Genre       TEXT,
    Length      INTEGER,
    Title       TEXT,
    Track       INTEGER,
    TrackTotal  INTEGER,
    Mbid        CHARACTER(36),
    AlbumMbid   CHARACTER(36),
    ArtistMbid  CHARACTER(36)
  );
`)

	return err
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) IsStale(path string, modtime time.Time) bool {
	row := d.db.QueryRow("SELECT Updated FROM songs WHERE Path=?",
		path)

	var comptime time.Time
	if err := row.Scan(&comptime); err != nil {
		return true
	}

	return comptime.Before(modtime)
}

func (d *Database) Insert(song Song) error {
	_, err := d.db.Exec("INSERT INTO songs VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		nil,
		song.Path,
		time.Now(),
		song.Album,
		song.AlbumArtist,
		song.Artist,
		song.Composer,
		song.Date,
		song.Disc,
		song.DiscTotal,
		song.Genre,
		song.Length,
		song.Title,
		song.Track,
		song.TrackTotal,
		"",
		"",
		"")

	if serr, ok := err.(*sqlite3.Error); ok && serr.Code() == 2067 {
		_, err := d.db.Exec(`UPDATE songs SET
    Updated=?,
    Album=?, AlbumArtist=?, Artist=?,
    Composer=?, Date=?, Disc=?,
    DiscTotal=?, Genre=?, Length=?,
    Title=?, Track=?, TrackTotal=?
  WHERE Path=?
`,
			time.Now(),
			song.Album,
			song.AlbumArtist,
			song.Artist,
			song.Composer,
			song.Date,
			song.Disc,
			song.DiscTotal,
			song.Genre,
			song.Length,
			song.Title,
			song.Track,
			song.TrackTotal,
			song.Path)

		return err
	}

	return err
}

func (d *Database) Each(f func(Song) error, qs ...Clause) error {
	stmt := `SELECT
  Path,
  Album, AlbumArtist, Artist,
  Composer, Date, Disc, DiscTotal,
  Genre, Length, Title, Track, TrackTotal
  FROM songs`

	var (
		rows *sql.Rows
		err  error
	)

	if len(qs) == 0 {
		rows, err = d.db.Query(stmt)
	} else {
		stmt += " WHERE"
		args := []interface{}{}
		for i, q := range qs {
			if i > 0 {
				stmt += " AND"
			}
			stmt += " " + q.Clause()
			args = append(args, q.Arg())
		}

		rows, err = d.db.Query(stmt, args...)
	}

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var song Song
		if err := rows.Scan(&song.Path,
			&song.Album, &song.AlbumArtist, &song.Artist,
			&song.Composer, &song.Date, &song.Disc, &song.DiscTotal,
			&song.Genre, &song.Length, &song.Title, &song.Track, &song.TrackTotal,
		); err != nil {
			return err
		}

		if err := f(song); err != nil {
			return err
		}
	}

	return rows.Err()
}

func (d *Database) Delete(path string) error {
	_, err := d.db.Exec("DELETE FROM songs WHERE Path=?", path)
	return err
}

func (d *Database) Move(oldPath, newPath string) error {
	_, err := d.db.Exec("UPDATE songs SET Path=? WHERE Path=?",
		newPath,
		oldPath)
	return err
}
