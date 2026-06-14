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

	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS albums(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		album_artist TEXT
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
		album_id INTEGER,
		mtime INTEGER,
		size INTEGER,
		path TEXT UNIQUE NOT NULL,
		FOREIGN KEY (album_id) REFERENCES albums(id)
	)
	`)
	if err != nil {
		return nil, fmt.Errorf("create tracks table: %w", err)
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

func (r *LibraryRepository) GetAlbums() ([]Album, error) {
	albums := []Album{}
	rows, err := r.db.Query(`
	SELECT
		a.id, a.title,
		a.album_artist,
		MIN(year) AS year,
		COUNT(*) AS total_tracks
	FROM albums a
	LEFT JOIN tracks t ON a.id = t.album_id 
	GROUP BY a.title, a.album_artist
	ORDER BY a.id
	`)
	if err != nil {
		return nil, fmt.Errorf("query execution: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		al := Album{}
		err = rows.Scan(
			&al.ID, &al.Title,
			&al.AlbumArtist,
			&al.Year, &al.TotalTracks,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan err: %w", err)
		}
		albums = append(albums, al)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration err: %s", err)
	}

	return albums, nil
}

func (r *LibraryRepository) GetAlbumByID(id uint16) (Album, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return Album{}, fmt.Errorf("transaction begin: %w", err)
	}
	defer tx.Rollback()

	album := Album{}
	err = tx.QueryRow(`
		SELECT id, title, album_artist
		FROM albums
		WHERE id = ?
	`, id).Scan(&album.ID, &album.Title, &album.AlbumArtist)
	if errors.Is(err, sql.ErrNoRows) {
		return Album{}, ErrAlbumNotFound
	}
	if err != nil {
		return Album{}, fmt.Errorf("album query: %w", err)
	}

	fields := "id, title, artist, album, album_artist, track_num, album_id"
	q := fmt.Sprintf("SELECT %s FROM tracks WHERE album_id = ?", fields)
	rows, err := tx.Query(q, id)
	if err != nil {
		return Album{}, fmt.Errorf("bulk select tracks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		tr := Track{}
		err = rows.Scan(
			&tr.ID, &tr.Title,
			&tr.Artist, &tr.Album,
			&tr.AlbumArtist, &tr.TrackNum,
			&tr.AlbumID,
		)
		if err != nil {
			return Album{}, fmt.Errorf("row scan err: %w", err)
		}

		album.Tracks = append(album.Tracks, tr)
	}
	if err := rows.Err(); err != nil {
		return Album{}, fmt.Errorf("row iteration err: %w", err)
	}

	return album, tx.Commit()
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
	id, title, artist, album, album_artist, track_num, year, album_id, mtime, size, path
	)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("prepare insert stmt: %w", err)
	}

	err = ScanLibrary(path, func(track Track) error {
		albumID, err := getOrCreateAlbum(tx, track)
		if err != nil {
			return err
		}

		track.AlbumID = albumID

		_, err = stmt.Exec(
			track.ID, track.Title, track.Artist,
			track.Album, track.AlbumArtist, track.TrackNum,
			track.Year, track.AlbumID, track.ModTime.Unix(),
			track.Size, track.Path,
		)
		return err
	})

	if err != nil {
		return err
	}

	return tx.Commit()
}

func getOrCreateAlbum(tx *sql.Tx, track Track) (uint16, error) {
	var id uint16

	err := tx.QueryRow(`
		INSERT INTO albums (title, album_artist)
		VALUES (?, ?)
		ON CONFLICT (title, album_artist)
		DO UPDATE SET album_artist = excluded.album_artist
		RETURNING id
	`, track.Album, track.AlbumArtist).Scan(&id)
	return id, err
}
