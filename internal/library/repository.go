package library

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

type LibraryRepository struct {
	db *sql.DB
}

func NewTrackRepository(path string) (*LibraryRepository, error) {
	var err error
	r := LibraryRepository{}
	if !Exists(path) {
		r.db, err = createSchema()
		if err != nil {
			return &r, fmt.Errorf("create schema: %w", err)
		}
	}

	r.db, err = sql.Open("sqlite", "./library.db")
	if err != nil {
		return nil, fmt.Errorf("open db connection: %w", err)
	}

	return &r, nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return err == nil
}

func createSchema() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./library.db")
	if err != nil {
		return nil, fmt.Errorf("open db connection: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("transaction start: %w", err)
	}
	defer tx.Rollback()

	// TODO: autoincrement on id - could break????
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS tracks(
		id INTEGER PRIMARY KEY,
		title TEXT,
		artist TEXT,
		album TEXT,
		album_artist TEXT,
		track_num INTEGER,
		year INTEGER,
		mtime INTEGER,
		size INTEGER,
		path TEXT UNIQUE NOT NULL
	)
	`)
	if err != nil {
		return nil, fmt.Errorf("create tracks table: %w", err)
	}

	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS albums(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		album_artist TEXT,
		year INTEGER,
		track_count INTEGER
	)
	`)
	if err != nil {
		return nil, fmt.Errorf("create albums table: %w", err)
	}

	_, err = tx.Exec(`
	CREATE UNIQUE INDEX idx_album_unique
	ON albums(title, album_artist);
	`)
	if err != nil {
		return nil, fmt.Errorf("create unique albums index: %w", err)
	}

	return db, tx.Commit()
}

func (r *LibraryRepository) GetTrackByID(id uint16) (Track, error) {
	res := Track{}
	var mtime int64
	err := r.db.QueryRow(`
		SELECT 
			id, title,
			artist, album,
			mtime, size, path
		FROM tracks
		WHERE id=?`,
		id).Scan(
		&res.ID,
		&res.Title,
		&res.Artist,
		&res.Album,
		&mtime,
		&res.Size,
		&res.Path,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Track{}, ErrTrackNotFound
	}
	if err != nil {
		return Track{}, fmt.Errorf("get track: %w", err)
	}

	res.ModTime = time.Unix(mtime, 0)

	return res, nil
}

func (r *LibraryRepository) GetAlbumList() ([]Album, error) {
	albums := []Album{}
	return albums, nil
}

func (r *LibraryRepository) GetAllTracks() ([]Track, error) {
	tracks := []Track{}
	return tracks, nil
}

func (r *LibraryRepository) ImportLibrary(path string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("trasaction start: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO tracks (
	id, title, artist, album, album_artist, track_num, year, mtime, size, path
	)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("prepare insert stmt: %w", err)
	}

	err = ScanLibrary(path, func(track Track) error {
		_, err := stmt.Exec(
			track.ID, track.Title, track.Artist,
			track.Album, track.AlbumArtist, track.TrackNum,
			track.Year, track.ModTime.Unix(), track.Size, track.Path,
		)
		return err
	})

	if err != nil {
		return err
	}

	if err := r.RebuildAlbums(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *LibraryRepository) RebuildAlbums(tx *sql.Tx) error {
	if _, err := tx.Exec(`DELETE FROM albums`); err != nil {
		return fmt.Errorf("clear albums: %w", err)
	}

	_, err := tx.Exec(`
		INSERT INTO albums(
			title, album_artist,
			year, track_count
		)
		SELECT
			album, artist,
			MIN(year) AS year,
			COUNT(*) AS track_count
		FROM track
		GROUP BY album, album_artist
	`)

	if err != nil {
		return fmt.Errorf("albums rebuild: %w", err)
	}

	return nil
}
