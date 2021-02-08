package main

import (
	"mp3mirror/mp3mirror"
	"mp3mirror/trayapp/icon"

	"github.com/getlantern/systray"
	"github.com/martinlindhe/notify"
	"github.com/sqweek/dialog"
)

func main() {
	systray.Run(onReady, onExit)
}

// CONFIG: ~/Library/Application Support/<appname>/<appname>.cfg

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTooltip("mp3 mirror")

	watchFolder, err := dialog.Directory().Title("Select watch folder").Browse()
	if err != nil {
		systray.Quit()
		return
	}
	watcher := mp3mirror.NewWatcher(mp3mirror.WatcherConfig{
		Folders:          []string{},
		FoldersRecursive: []string{watchFolder},
		MirrorFolder:     "",
		ConvertBitrate:   mp3mirror.DefaultConvertBitrate,
		ConvertTimeout:   mp3mirror.DefaultConvertTimeout,
		RefreshInterval:  mp3mirror.DefaultRefreshInterval,
	})

	quitItem := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		for {
			select {
			case <-quitItem.ClickedCh:
				systray.Quit()
			}
		}
	}()

	go func() {
		if err := watcher.Run(); err != nil {
			systray.Quit()
			dialog.Message("Watcher error: %v", err).Title("Watcher error").Error()
		}
	}()

	go func() {
		for {
			select {
			case convertedfile := <-watcher.Converted:
				notify.Notify("MP3mirror", "", convertedfile, "")
			}
		}
	}()
}

func onExit() {
	// clean up here
}
