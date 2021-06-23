package main

import (
	"fmt"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/KrystianD/screentool/main/utils"
)

func CropPixbuf(pb *gdk.Pixbuf, rect utils.Rectangle) *gdk.Pixbuf {
	var surface = cairo.CreateImageSurface(cairo.FORMAT_ARGB32, rect.Width(), rect.Height())
	fmt.Println(rect.Width(), rect.Height())
	var ctx = cairo.Create(surface)

	gtk.GdkCairoSetSourcePixBuf(ctx, pb, float64(-rect.X()), float64(-rect.Y()))
	fmt.Println(float64(rect.X()), float64(rect.Y()))
	ctx.Paint()

	pixbuf, _ := gdk.PixbufGetFromSurface(surface, 0, 0, rect.Width(), rect.Height())
	fmt.Println(0, 0, rect.Width(), rect.Height())

	return pixbuf
}
