package api

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ZeroDay0utplay/file-compressor-service/internal/storage"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	deps Dependencies
}

func New(deps Dependencies) *Handler {
	return &Handler{deps: deps}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) Compress(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.deps.RequestTimeout)
	defer cancel()

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.deps.MaxUploadBytes)

	inputPath, filename, err := h.deps.TempStore.Save(ctx, c.Request)
	if err != nil {
		if errors.Is(err, storage.ErrMissingFile) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing file"})
			return
		}
		if errors.Is(err, storage.ErrFileTooLarge) {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file too large"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid upload"})
		return
	}
	defer os.Remove(inputPath)

	if err := h.deps.WorkerPool.Acquire(ctx); err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "service busy"})
		return
	}
	defer h.deps.WorkerPool.Release()

	preset := choosePreset(inputPath)

	outputPath, err := h.deps.Compressor.Compress(ctx, inputPath, preset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "compression failed"})
		return
	}
	defer os.Remove(outputPath)

	f, err := os.Open(outputPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot read output"})
		return
	}
	defer f.Close()

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", `attachment; filename="`+filepath.Base(filename)+`.compressed"`)

	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, f)
}

func choosePreset(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return "/ebook"
	}
	switch size := info.Size(); {
	case size < 5*1024*1024:
		return "/printer"
	case size < 30*1024*1024:
		return "/ebook"
	default:
		return "/screen"
	}
}
