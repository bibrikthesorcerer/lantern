package library

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

type TrackRepository struct {
	db *sql.DB
}

func NewTrackRepository(path string) (*TrackRepository, error) {
	var err error
	r := TrackRepository{}
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

	// TODO: autoincrement on id - could break????
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS tracks(
		id INTEGER PRIMARY KEY,
		title TEXT,
		artist TEXT,
		album TEXT,
		mtime INTEGER,
		size INTEGER,
		path TEXT UNIQUE NOT NULL
	)
	`)
	if err != nil {
		return nil, fmt.Errorf("create tracks table: %w", err)
	}

	return db, nil
}

func (r *TrackRepository) ImportLibrary(path string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("trasaction start: %w", err)
	}

	stmt, err := tx.Prepare("INSERT INTO tracks (id, title, artist, album, mtime, size, path) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("prepare insert stmt: %w", err)
	}

	err = ScanLibrary(path, func(track Track) error {
		_, err := stmt.Exec(track.ID, track.Title, track.Artist, track.Album, track.ModTime.Unix(), track.Size, track.Path)
		return err
	})

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
