package main

import (
	"fmt"
	"github.com/andlabs/ui"
)

var toolWindow *ui.Window

func tools() {
	go func() {
		if toolWindow != nil {
			c := make(chan bool)
			ui.QueueMain(func() {
				toolWindow.Destroy()
				c <- true
			})
			<- c // wait for destroy
		}
		err := ui.Main(func() {
			input := ui.NewEntry()
			greeting := ui.NewLabel("")
			
			slider := ui.NewSlider(-10, 10)
			slider.SetValue(0)
			slider.OnChanged(func(slider *ui.Slider) {
				fmt.Println("value", slider.Value())
			})

			button := ui.NewButton("Greet")
			button.OnClicked(func(*ui.Button) {
				greeting.SetText("Hello, " + input.Text() + "!")
			})
			
			box := ui.NewVerticalBox()
			box.Append(ui.NewLabel("Enter your name:"), false)
			box.Append(input, false)
			box.Append(button, false)
			box.Append(greeting, false)
			box.Append(slider, false)

			toolWindow = ui.NewWindow("Hello", 200, 100, false)
			toolWindow.SetMargined(true)
			toolWindow.SetChild(box)
			toolWindow.OnClosing(func(*ui.Window) bool {
				return true
			})
			toolWindow.Show()
		})
		if err != nil {
			panic(err)
		}
	}()
}

func toolsLoadFile(cb func(string)) {
	ui.QueueMain(func() {
		filename := ui.OpenFile(toolWindow)
		cb(filename)
	})
}

func toolsSaveFile(cb func(string)) {
	ui.QueueMain(func() {
		filename := ui.SaveFile(toolWindow)
		cb(filename)
	})
}
