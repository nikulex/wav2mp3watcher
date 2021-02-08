package main

import (
	"flag"
	"fmt"
	"log"

	"mp3mirror/mp3mirror"
)

func main() {
	folder := flag.String("f", ".", "watch folder")
	mirror := flag.String("m", "", "move files to folder (default is nearby)")
	timeout := flag.Duration("t", mp3mirror.DefaultConvertTimeout, "convert timeout")
	interval := flag.Duration("i", mp3mirror.DefaultRefreshInterval, "rescan interval")
	flag.Parse()

	w := mp3mirror.NewWatcher(mp3mirror.WatcherConfig{
		Folders:          []string{},
		FoldersRecursive: []string{*folder},
		MirrorFolder:     *mirror,
		ConvertTimeout:   *timeout,
		RefreshInterval:  *interval,
	})

	go func() {
		for {
			select {
			case file := <-w.ConvertStart:
				fmt.Printf("Start convert: %s\n", file)
			case file := <-w.ConvertFinishOK:
				fmt.Printf("Finish convert: %s\n", file)
			case err := <-w.ConvertFinishError:
				fmt.Printf("Finish with Error: %v\n", err)
			}
		}
	}()

	if err := w.Run(); err != nil {
		log.Fatalln(err)
	}
}
