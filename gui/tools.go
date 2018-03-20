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
var exportWindow *ui.Window
var exportBox ui.Control
var exportSt storage.Storage
var exportFunc func(string)

func Init(st storage.Storage) {
	exportSt = st.Of("export")
}

func init() {
	go func() {
		err := ui.Main(func() {
			fooWindow = ui.NewWindow("", 0, 0, false)
			exportBox = createExportBox()
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

func createExportBox() ui.Control {
	exportTab := ui.NewTab()
	exportTab.Append("SVG", createSvgPresets())
	exportTab.Append("PNG", createSvgPresets())
	exportTab.Append("PDF", createSvgPresets())

	return exportTab
}

func createSvgPresets () ui.Control {
	button := ui.NewButton("Export")
	button.OnClicked(func(button *ui.Button) {
		if exportFunc != nil {
			exportFunc(exportSt.String("filename"))
		}
	})

	box := ui.NewVerticalBox()
	box.Append(createIntInput("width", 1, math.MaxInt32), false)
	box.Append(createIntInput("height", 1, math.MaxInt32), false)
	box.Append(createIntInput("zoom", -5, +5), false)
	box.Append(createExportAs(".png"), false)
	box.Append(button, true)
	return box
}

func createExportAs(ext string) ui.Control {
	input := ui.NewEntry()
	input.SetText(exportSt.String("filename"))
	input.OnChanged(func(input *ui.Entry) {
		exportSt.Set("filename", input.Text())
	})
	button := ui.NewButton("browser...")
	button.OnClicked(func(*ui.Button) {
		SaveFile(func(filename string) {
			if filepath.Ext(filename) != ext {
				filename += ext
			}
			input.SetText(filename)
			exportSt.Set("filename", filename)
		})
	})
	return line(pair{input, true}, pair{button, false})
}

func createIntInput(name string, max, min int) ui.Control {
	label := ui.NewLabel(name)

	input := ui.NewSpinbox(max, min)
	input.SetValue(exportSt.Int(name))
	input.OnChanged(func(*ui.Spinbox) {
		exportSt.Set(name, input.Value())
	})
	
	return line(pair{label, true}, pair{input, false})
}

func IsExportOn() bool {
	return exportWindow != nil
}

func ToggleExport(export func(string)) {
	exportFunc = export
	ui.QueueMain(func() {
		if exportWindow != nil {
			exportWindow.SetChild(nil)
			exportWindow.Destroy()
			exportWindow = nil
		} else {
			exportWindow = ui.NewWindow("Export", 200, 100, false)
			exportWindow.SetMargined(true)
			exportWindow.SetChild(exportBox)
			exportWindow.OnClosing(func(*ui.Window) bool {
				exportWindow.SetChild(nil)
				exportWindow.Destroy()
				exportWindow = nil
				return false
			})
			exportWindow.Show()
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