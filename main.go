package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
)

const fileSizeLimit = 100 * 1024 * 1024

func mp3(wav string) string {
	return wav[:len(wav)-4] + ".mp3"
}

type mp3Converter struct {
	Bitrate uint
	Timeout time.Duration
}

func (c *mp3Converter) handle(event watcher.Event) error {
	if event.IsDir() || event.Size() > fileSizeLimit || !strings.HasSuffix(event.Name(), ".wav") {
		fmt.Printf("skip: %s\n", event.Path)
		return nil
	}
	switch event.Op {
	case watcher.Create:
		fallthrough
	case watcher.Write:
		return c.convert(event.Path)
	case watcher.Remove:
		fmt.Printf("remove: %s\n", event.Path)
		return os.Remove(mp3(event.Path))
	case watcher.Rename:
		fallthrough
	case watcher.Move:
		if _, err := os.Stat(mp3(event.OldPath)); os.IsNotExist(err) {
			return c.convert(event.Path)
		}
		fmt.Printf("move: %s -> %s\n", mp3(event.OldPath), mp3(event.Path))
		return os.Rename(mp3(event.OldPath), mp3(event.Path))
	}
	return nil
}

func (c *mp3Converter) convert(wav string) error {
	if len(wav) == 0 {
		return fmt.Errorf("file name \"%s\" is empty", wav)
	}
	if !strings.HasSuffix(wav, ".wav") {
		return fmt.Errorf("file \"%s\" has no \".wav\" extention", wav)
	}
	fmt.Printf("convert: %s\n", mp3(wav))

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", wav, "-b:a", fmt.Sprintf("%vk", c.Bitrate), mp3(wav))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("convert error: %v", err)
	}
	return nil
}

func main() {
	folder := flag.String("f", ".", "watch folder")
	bitrate := flag.Uint("b", 320, "mp3 bitrate")
	interval := flag.Duration("i", 1*time.Second, "rescan interval")
	timeout := flag.Duration("t", 1*time.Second, "convert timeout")
	//daemon := flag.Bool("d", false, "run as daemon") // TODO
	//quiet := flag.Bool("q", false, "quiet") // TODO
	flag.Parse()

	w := watcher.New()

	r := regexp.MustCompile("^*.wav$")
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	c := &mp3Converter{
		Bitrate: *bitrate,
		Timeout: *timeout,
	}

	go func() {
		for {
			select {
			case event := <-w.Event:
				if err := c.handle(event); err != nil {
					log.Println(err)
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.AddRecursive(*folder); err != nil {
		log.Fatalln(err)
	}

	go func() {
		w.Wait()
		fmt.Println("watching started...")
		for path, _ := range w.WatchedFiles() {
			c.convert(path)
		}
	}()

	if err := w.Start(*interval); err != nil {
		log.Fatalln(err)
	}
}
