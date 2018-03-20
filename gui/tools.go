// other (platform native) gui - file selector, presets...

package gui

import (
	_ "fmt"
	"math"
	"path/filepath"
	"github.com/andlabs/ui"
	"git.yo2.cz/drahoslav/penego/storage"
)

type pair struct {
	control ui.Control
	stretch bool
}

var fooWindow *ui.Window // for binding file load/save modal windows
var toolWindow *ui.Window
var box ui.Control
var store storage.Storage

func Init(st storage.Storage) {
	store = st
}

func init() {
	go func() {
		err := ui.Main(func() {
			fooWindow = ui.NewWindow("", 0, 0, false)
			box = createToolsBox()
		})
		if err != nil {
			panic(err)
		}
	}()
}

func line(stuff ...pair) ui.Control {
	box := ui.NewHorizontalBox()
	box.SetPadded(true)
	for _, p := range stuff {
		box.Append(p.control, p.stretch)
	}
	return box
}

func createToolsBox() ui.Control {
	exportTab := ui.NewTab()
	exportTab.Append("SVG", createSvgPresets())
	exportTab.Append("PNG", createSvgPresets())
	exportTab.Append("PDF", createSvgPresets())

	return exportTab
}

func createSvgPresets () ui.Control {
	box := ui.NewVerticalBox()
	w := store.Get("export.width").(int)
	h := store.Get("export.height").(int)
	z := store.Get("export.zoom").(int)
	box.Append(createIntInput("width", w, 1, math.MaxInt32), false)
	box.Append(createIntInput("height", h, 1, math.MaxInt32), false)
	box.Append(createIntInput("zoom", z, -5, +5), false)
	box.Append(createExportAs(".png"), false)
	return box
}

func createExportAs(ext string) ui.Control {
	input := ui.NewEntry()
	button := ui.NewButton("export as...")
	button.OnClicked(func(*ui.Button) {
		SaveFile(func(filename string) {
			if filepath.Ext(filename) != ext {
				filename += ext
			}
			input.SetText(filename)
		})
	})
	return line(pair{input, true}, pair{button, false})
}

func createIntInput(name string, value, max, min int) ui.Control {
	input := ui.NewSpinbox(max, min)
	input.SetValue(value)
	label := ui.NewLabel(name)

	return line(pair{label, true}, pair{input, false})
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
			toolWindow = ui.NewWindow("Export", 200, 100, false)
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
		filename := ui.OpenFile(fooWindow)
		cb(filename)
	})
}

func SaveFile(cb func(string)) {
	ui.QueueMain(func() {
		filename := ui.SaveFile(fooWindow)
		cb(filename)
	})
}
