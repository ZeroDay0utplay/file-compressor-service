package gs

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type Config struct {
	DefaultPreset string
	Timeout       time.Duration
}

type Compressor struct {
	cfg Config
}

func New(cfg Config) *Compressor {
	if cfg.Timeout == 0 {
		cfg.Timeout = 90 * time.Second
	}
	if cfg.DefaultPreset == "" {
		cfg.DefaultPreset = "/ebook"
	}
	return &Compressor{cfg: cfg}
}

func (c *Compressor) Compress(ctx context.Context, inputPath string, preset string) (string, error) {
	if preset == "" {
		preset = c.cfg.DefaultPreset
	}
	outPath := filepath.Join(os.TempDir(), "compressed-"+uuid.New().String()+".pdf")
	ctx2, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()
	cmd := exec.CommandContext(
		ctx2,
		"gs",
		"-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.4",
		"-dPDFSETTINGS="+preset,
		"-dNOPAUSE",
		"-dQUIET",
		"-dBATCH",
		"-sOutputFile="+outPath,
		inputPath,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		_ = os.Remove(outPath)
		return "", fmt.Errorf("gs failed: %w (%s)", err, out)
	}
	return outPath, nil
}
