package main

import (
	"fmt"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func saveToClipboardAndWait(pixbuf *gdk.Pixbuf, onOwnerChanged func()) {
	clip, err := gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD)
	if err != nil {
		fmt.Println(err)
		return
	}

	var count = 0

	_ = clip.Connect("owner-change", func(clip *gtk.Clipboard, event *gdk.Event) {
		count += 1
		if count == 2 {
			onOwnerChanged()
		}
	})
	clip.SetImage(pixbuf)
	clip.Store()
}
