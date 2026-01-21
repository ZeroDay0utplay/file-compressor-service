package api

import (
	"context"
	"net/http"
	"time"
)

type Compressor interface {
	Compress(ctx context.Context, inputPath string, preset string) (string, error)
}

type WorkerPool interface {
	Acquire(ctx context.Context) error
	Release()
}

type TempStore interface {
	Save(ctx context.Context, req *http.Request) (string, string, error)
}

type Dependencies struct {
	Compressor     Compressor
	WorkerPool     WorkerPool
	TempStore      TempStore
	MaxUploadBytes int64
	RequestTimeout time.Duration
}
