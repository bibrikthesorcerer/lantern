package library

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
	"go.senan.xyz/taglib"
	_ "golang.org/x/image/webp"
)

type CoverCache struct {
	cachePath string
}

func NewCoverCache() (*CoverCache, error) {
	cachePath, err := ensureCacheDir()
	if err != nil {
		return nil, err
	}

	return &CoverCache{cachePath: cachePath}, nil
}

func ensureCacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	cacheDir = filepath.Join(cacheDir, "lantern")

	if err := os.MkdirAll(cacheDir, 0750); err != nil {
		return "", err
	}

	return cacheDir, nil
}

func (r *CoverCache) CacheAlbumCover(path string, id uint16) error {
	coverArt, err := taglib.ReadImage(path)
	if err != nil {
		return fmt.Errorf("image read: %w", err)
	}

	img, _, err := image.Decode(bytes.NewReader(coverArt))
	if err != nil {
		return err
	}
	resizedImg := resize.Thumbnail(300, 300, img, resize.Lanczos3)

	thumbFile, err := os.Create(fmt.Sprintf("%s/%d.jpg", r.cachePath, id))
	if err != nil {
		return fmt.Errorf("create thumbnail file: %w", err)
	}
	defer thumbFile.Close()

	return jpeg.Encode(thumbFile, resizedImg, nil)
}

func (r *CoverCache) GetAlbumCover(id uint16) (*AlbumCover, error) {
	thumbFile, err := os.ReadFile(fmt.Sprintf("%s/%d.jpg", r.cachePath, id))
	if err != nil {
		return nil, fmt.Errorf("open thumbnail file: %w", err)
	}

	return &AlbumCover{
		Data:      thumbFile,
		ImageDesc: taglib.ImageDesc{MIMEType: "image/jpeg"},
	}, nil
}
