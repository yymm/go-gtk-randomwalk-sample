package main

import (
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"os"
	"unsafe"
)

type point struct {
	x int
	y int
}

func main() {
    gtk.Init(&os.Args)
    window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
    window.SetTitle("Random Walking")
    window.Connect("destroy", gtk.MainQuit)

	var p1, p2 point
	var gdkwin *gdk.Window
	var pixmap *gdk.Pixmap
	var gc *gdk.GC
	p1.x = -1
	p1.y = -1

	// Box
	vbox := gtk.NewVBox(false, 0)
	drawingarea := gtk.NewDrawingArea()

	drawingarea.Connect("configure-event", func() {
		if pixmap != nil {
			pixmap.Unref()
		}
		allocation := drawingarea.GetAllocation()
		pixmap = gdk.NewPixmap(drawingarea.GetWindow().GetDrawable(), allocation.Width, allocation.Height, 24)
		gc = gdk.NewGC(pixmap.GetDrawable())
		gc.SetRgbFgColor(gdk.NewColor("black"))
		pixmap.GetDrawable().DrawRectangle(gc, true, 0, 0, -1, -1)
		gc.SetRgbFgColor(gdk.NewColor("white"))
	})

	drawingarea.Connect("motion-notify-event", func(ctx *glib.CallbackContext) {
		if gdkwin == nil {
			gdkwin = drawingarea.GetWindow()
		}
		arg := ctx.Args(0)
		mev := *(**gdk.EventMotion)(unsafe.Pointer(&arg))
		var mt gdk.ModifierType
		if mev.IsHint != 0 {
			gdkwin.GetPointer(&p2.x, &p2.y, &mt)
		} else {
			p2.x, p2.y = int(mev.X), int(mev.Y)
		}
		if p1.x != -1 && p2.x != -1 && (gdk.EventMask(mt)&gdk.BUTTON_PRESS_MASK) != 0 {
			pixmap.GetDrawable().DrawLine(gc, p1.x, p1.y, p2.x, p2.y)
			drawingarea.GetWindow().Invalidate(nil, false)
		}
		p1 = p2
	})

	drawingarea.Connect("expose-event", func() {
		if pixmap != nil {
			drawingarea.GetWindow().GetDrawable().DrawDrawable(gc, pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
		}
	})

	drawingarea.SetEvents(int(gdk.POINTER_MOTION_MASK | gdk.POINTER_MOTION_HINT_MASK | gdk.BUTTON_PRESS_MASK))
	vbox.Add(drawingarea)

	btns := gtk.NewHBox(false, 0)
	startbtn := gtk.NewButtonWithLabel("start")
	stopbtn := gtk.NewButtonWithLabel("stop")

	btns.Add(startbtn)
	btns.Add(stopbtn)

	vbox.PackEnd(btns, false, false, 0)




	window.Add(vbox)
    window.SetSizeRequest(800, 800)
    window.ShowAll()
    gtk.Main()
}
