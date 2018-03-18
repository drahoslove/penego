package main

import (
	"github.com/andlabs/ui"
)

var toolWindow *ui.Window

func init() {
	go func() {
		err := ui.Main(func() {
			input := ui.NewEntry()
			button := ui.NewButton("Greet")
			greeting := ui.NewLabel("")
			box := ui.NewVerticalBox()
			box.Append(ui.NewLabel("Enter your name:"), false)
			box.Append(input, false)
			box.Append(button, false)
			box.Append(greeting, false)
			toolWindow = ui.NewWindow("Hello", 200, 100, false)
			toolWindow.SetMargined(true)
			toolWindow.SetChild(box)
			button.OnClicked(func(*ui.Button) {
				greeting.SetText("Hello, " + input.Text() + "!")
			})
			toolWindow.OnClosing(func(*ui.Window) bool {
				toolWindow.Hide()
				return false
			})

		})
		if err != nil {
			panic(err)
		}
	}()
}

func tools() {
	toolWindow.Show()
}