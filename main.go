package main

import (
	"fmt"
	"time"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gtk"
	"os"
	"math/rand"
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

	drawingarea.Connect("expose-event", func() {
		if pixmap != nil {
			drawingarea.GetWindow().GetDrawable().DrawDrawable(gc, pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
		}
	})

	drawingarea.SetEvents(int(gdk.POINTER_MOTION_MASK | gdk.POINTER_MOTION_HINT_MASK | gdk.BUTTON_PRESS_MASK))
	vbox.Add(drawingarea)

	walkable := make(chan bool)
	counter := 0

	btns := gtk.NewHBox(false, 0)
	startbtn := gtk.NewButtonWithLabel("start")
	resetbtn := gtk.NewButtonWithLabel("reset")

	startbtn.Clicked(func() {
		fmt.Println("button clicked: ", startbtn.GetLabel())
		walkable <- true
	})

	resetbtn.Clicked(func() {
		fmt.Println("button clicked: ", resetbtn.GetLabel())
		counter = 0
		drawingarea.GetWindow().Invalidate(nil, false)
	})

	go func() {
		if gdkwin == nil {
			gdkwin = drawingarea.GetWindow()
		}
		p1.x = 400
		p1.y = 400
		p2.x = 400
		p2.y = 400
		rand.Seed(time.Now().Unix())
		for {
			wa := <-walkable
			if wa {
				for ; counter < 1000 ; {
					time.Sleep(40*time.Millisecond)
					//fmt.Println(counter)
					//fmt.Println(p1, p2)
					counter += 1
					r := rand.Float64()
					if r < 0.33333333 { p2.x += -5 } else if r < 0.66666666 { p2.x += 5 } else { p2.x += 0}
					r = rand.Float64()
					if r < 0.33333333 { p2.y += -5 } else if r < 0.66666666 { p2.y += 5 } else { p2.y += 0}
					pixmap.GetDrawable().DrawLine(gc, p1.x, p1.y, p2.x, p2.y)
					drawingarea.GetWindow().Invalidate(nil, false)
					p1 = p2
				}
			}
		}
	}()

	btns.Add(startbtn)
	btns.Add(resetbtn)

	vbox.PackEnd(btns, false, false, 0)

	window.Add(vbox)
    window.SetSizeRequest(800, 800)
    window.ShowAll()
    gtk.Main()
}
