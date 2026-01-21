package compressor

import (
	"context"
	"errors"
)

var ErrNotSupported = errors.New("not supported")

type Backend interface {
	Compress(ctx context.Context, inputPath string, preset string) (string, error)
}

type Registry struct {
	backends map[string]Backend
}

func NewRegistry() *Registry {
	return &Registry{backends: make(map[string]Backend)}
}

func (r *Registry) Register(mimePrefix string, b Backend) {
	r.backends[mimePrefix] = b
}

func (r *Registry) findBackend(mime string) (Backend, bool) {
	if b, ok := r.backends[mime]; ok {
		return b, true
	}
	for prefix, b := range r.backends {
		if len(prefix) > 0 && prefix[len(prefix)-1] == '/' {
			if len(mime) >= len(prefix) && mime[:len(prefix)] == prefix {
				return b, true
			}
		}
	}
	return nil, false
}

func (r *Registry) Compress(ctx context.Context, inputPath, mime, preset string) (string, error) {
	if b, ok := r.findBackend(mime); ok {
		return b.Compress(ctx, inputPath, preset)
	}
	return "", ErrNotSupported
}
