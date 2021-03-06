package mp3mirror

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/radovskyb/watcher"
)

// TODO LOG

const (
	fileSizeLimit = 100 * 1024 * 1024
	regexFilter   = "^*(.wav)|(.aiff)$"
)

// default values
const (
	DefaultConvertTimeout  = 1 * time.Minute
	DefaultRefreshInterval = 1 * time.Second
)

// WatcherConfig ...
type WatcherConfig struct {
	Folders          []string
	FoldersRecursive []string
	MirrorFolder     string
	ConvertTimeout   time.Duration
	RefreshInterval  time.Duration
}

// Watcher ...
type Watcher struct {
	watcher   *watcher.Watcher
	interval  time.Duration
	mirror    *Mirror
	converter *Converter

	ConvertStart       chan string
	ConvertFinishOK    chan string
	ConvertFinishError chan error
}

// NewWatcher ...
func NewWatcher(cfg WatcherConfig) *Watcher {
	w := watcher.New()

	r := regexp.MustCompile(regexFilter)
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	for _, folder := range cfg.Folders {
		if err := w.Add(folder); err != nil {
			log.Fatalln(err)
		}
	}
	for _, folder := range cfg.FoldersRecursive {
		if err := w.AddRecursive(folder); err != nil {
			log.Fatalln(err)
		}
	}

	if cfg.ConvertTimeout == 0 {
		cfg.ConvertTimeout = DefaultConvertTimeout
	}
	if cfg.RefreshInterval == 0 {
		cfg.RefreshInterval = DefaultRefreshInterval
	}

	return &Watcher{
		watcher:  w,
		interval: cfg.RefreshInterval,
		mirror: &Mirror{
			Folder: cfg.MirrorFolder,
		},
		converter: &Converter{
			Timeout: cfg.ConvertTimeout,
		},
		ConvertStart:       make(chan string),
		ConvertFinishOK:    make(chan string),
		ConvertFinishError: make(chan error),
	}
}

func (w *Watcher) convert(from, to string) error {
	w.ConvertStart <- from
	err := w.converter.Convert(from, to)
	if err == nil {
		w.ConvertFinishOK <- to
	} else {
		w.ConvertFinishError <- fmt.Errorf("Convert error: %v", err)
	}
	return err
}

func (w *Watcher) handle(event watcher.Event) error {
	mp3path := w.mirror.Get(event.Path)
	if event.IsDir() || event.Size() > fileSizeLimit || len(mp3path) == 0 {
		return nil
	}
	switch event.Op {
	case watcher.Create:
		fallthrough
	case watcher.Write:
		return w.convert(event.Path, mp3path)
	case watcher.Rename:
		fallthrough
	case watcher.Move:
		oldmp3 := w.mirror.Get(event.OldPath)
		if len(oldmp3) == 0 {
			return nil
		}
		if _, err := os.Stat(oldmp3); os.IsNotExist(err) {
			return w.convert(event.Path, mp3path)
		}
		return os.Rename(oldmp3, mp3path)
	}
	return nil
}

// Run watcher loop
func (w *Watcher) Run() error {
	go func() {
		for {
			select {
			case event := <-w.watcher.Event:
				if err := w.handle(event); err != nil {
					log.Println(err)
				}
			case err := <-w.watcher.Error:
				log.Fatalln(err)
			case <-w.watcher.Closed:
				return
			}
		}
	}()

	go func() {
		w.watcher.Wait()
		fmt.Println("watching started...")
		for filepath := range w.watcher.WatchedFiles() {
			go func(file string) {
				w.convert(file, w.mirror.Get(file))
			}(filepath)
			// w.convert(filepath, w.mirror.Get(filepath))
		}
	}()

	return w.watcher.Start(w.interval)
}
