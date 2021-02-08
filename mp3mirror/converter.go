package mp3mirror

import (
	"context"
	"os/exec"
	"time"
)

// TODO goav wrapper

// Converter to mp3
type Converter struct {
	Timeout time.Duration
}

// Convert from wav/aiff to mp3
func (c *Converter) Convert(from, to string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	//cmd := exec.CommandContext(ctx, "ffmpeg", "-i", from, "-b:a", "320k", to)
	cmd := exec.CommandContext(ctx, "/usr/local/bin/lame", "--silent", "-h", "-V2", from, to)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
