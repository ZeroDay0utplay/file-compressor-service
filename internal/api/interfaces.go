package api

import (
	"context"
	"net/http"
	"time"

	"github.com/ZeroDay0utplay/file-compressor-service/internal/compressor"
)

type Dependencies struct {
	Registry       *compressor.Registry
	Limiter        Limiter
	TempStore      TempStore
	MaxUploadBytes int64
	RequestTimeout time.Duration
}

type Limiter interface {
	Acquire(ctx context.Context) error
	Release()
}

type TempStore interface {
	Save(ctx context.Context, req *http.Request) (path string, filename string, err error)
}
