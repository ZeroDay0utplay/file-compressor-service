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
	config Config
}

func New(config Config) *Compressor {
	if config.Timeout == 0 {
		config.Timeout = 90 * time.Second
	}
	if config.DefaultPreset == "" {
		config.DefaultPreset = "/ebook"
	}
	return &Compressor{config: config}
}

func (c *Compressor) Compress(ctx context.Context, inputPath string, preset string) (string, error) {
	if preset == "" {
		preset = c.config.DefaultPreset
	}

	outputPath := filepath.Join(os.TempDir(), "compressed-"+uuid.New().String()+".pdf")

	ctx2, cancel := context.WithTimeout(ctx, c.config.Timeout)
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
		"-sOutputFile="+outputPath,
		inputPath,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		_ = os.Remove(outputPath)
		return "", fmt.Errorf("gs failed: %w (%s)", err, out)
	}

	return outputPath, nil
}
