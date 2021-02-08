package main

import (
	"flag"
	"log"

	"mp3mirror/mp3mirror"
)

func main() {
	folder := flag.String("f", ".", "watch folder")
	mirror := flag.String("m", "", "move files to folder (default is nearby)")
	bitrate := flag.Uint("b", mp3mirror.DefaultConvertBitrate, "mp3 bitrate")
	timeout := flag.Duration("t", mp3mirror.DefaultConvertTimeout, "convert timeout")
	interval := flag.Duration("i", mp3mirror.DefaultRefreshInterval, "rescan interval")
	flag.Parse()

	w := mp3mirror.NewWatcher(mp3mirror.WatcherConfig{
		Folders:          []string{},
		FoldersRecursive: []string{*folder},
		MirrorFolder:     *mirror,
		ConvertBitrate:   *bitrate,
		ConvertTimeout:   *timeout,
		RefreshInterval:  *interval,
	})
	if err := w.Run(); err != nil {
		log.Fatalln(err)
	}
}
