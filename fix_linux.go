package main

import (
	_ "runtime"
	_ "time"
	// "github.com/mattn/go-gtk/gtk"
)

// hack to make dialog work under linux
func init() {
	// var dlg *gtk.MessageDialog
	// dlg = gtk.NewMessageDialog(nil, 0, gtk.MESSAGE_OTHER, gtk.BUTTONS_NONE, "Loading Penego . . .") // this is a lie!
	// go func() {
	// 	time.Sleep(time.Second)
	// 	dlg.Response(gtk.RESPONSE_NONE)
	// }()
	// dlg.Run()
	// dlg.Destroy()
	// for gtk.EventsPending() {
	// 	gtk.MainIteration()
	// }
}
