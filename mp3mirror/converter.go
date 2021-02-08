package mp3mirror

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// TODO goav wrapper

// Converter to mp3
type Converter struct {
	Timeout time.Duration
	Bitrate uint
}

// Convert from wav/aiff to mp3
func (c *Converter) Convert(from string, to string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", from, "-b:a", fmt.Sprintf("%vk", c.Bitrate), to)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("convert error: %v", err)
	}
	return nil
}
