package storage

import (
	"context"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

var (
	ErrMissingFile  = errors.New("missing file")
	ErrFileTooLarge = errors.New("file too large")
)

type TempStore struct{}

func New() *TempStore {
	return &TempStore{}
}

func (s *TempStore) Save(ctx context.Context, req *http.Request) (string, string, error) {
	mediaType, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil || mediaType != "multipart/form-data" {
		return "", "", ErrMissingFile
	}
	reader := multipart.NewReader(req.Body, params["boundary"])
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", "", err
		}
		if part.FileName() == "" {
			_ = part.Close()
			continue
		}
		path := filepath.Join(os.TempDir(), "upload-"+uuid.New().String())
		f, err := os.Create(path)
		if err != nil {
			part.Close()
			return "", "", err
		}
		_, err = io.Copy(f, part)
		f.Close()
		part.Close()
		if err != nil {
			_ = os.Remove(path)
			return "", "", err
		}
		return path, part.FileName(), nil
	}
	return "", "", ErrMissingFile
}
