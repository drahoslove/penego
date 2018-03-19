package gui

import (
	"fmt"
	"github.com/andlabs/ui"
)


var toolWindow *ui.Window
var box *ui.Box

func init() {
	go func() {
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
			
			box = ui.NewVerticalBox()
			box.Append(ui.NewLabel("Enter your name:"), false)
			box.Append(input, false)
			box.Append(button, false)
			box.Append(greeting, false)
			box.Append(slider, false)
		})
		if err != nil {
			panic(err)
		}
	}()
}

func IsToolsOn() bool {
	return toolWindow != nil
}

func ToggleTools() {
	ui.QueueMain(func() {
		if toolWindow != nil {
			toolWindow.SetChild(nil)
			toolWindow.Destroy()
			toolWindow = nil
		} else {
			toolWindow = ui.NewWindow("Tools", 200, 100, false)
			toolWindow.SetMargined(true)
			toolWindow.SetChild(box)
			toolWindow.OnClosing(func(*ui.Window) bool {
				toolWindow.SetChild(nil)
				toolWindow.Destroy()
				toolWindow = nil
				return false
			})
			toolWindow.Show()
		}
	})
}

func LoadFile(cb func(string)) {
	ui.QueueMain(func() {
		filename := ui.OpenFile(toolWindow)
		cb(filename)
	})
}

func SaveFile(cb func(string)) {
	ui.QueueMain(func() {
		filename := ui.SaveFile(toolWindow)
		cb(filename)
	})
}
