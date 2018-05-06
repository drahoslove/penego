// other (platform native) gui - file selector, presets...

package gui

import (
	"math"
	"path/filepath"
	"time"

	"git.yo2.cz/drahoslav/penego/storage"
	"github.com/andlabs/ui"
)

const (
	// points with 72 dpi
	A4_WIDTH   = 595
	A4_HEIGHT  = 842
)

type pair struct {
	control ui.Control
	stretch bool
}

var fooWindow *ui.Window // for binding file load/save modal windows
var (
	exportWindow *ui.Window
	exportBox    ui.Control
	exportSt     *storage.Storage
	exportFunc   func(string)
)
var (
	settingsWindow *ui.Window
	settingsBox    ui.Control
	settingsSt     *storage.Storage
)

func init() {
	exportSt = storage.Of("export")
	settingsSt = storage.Of("settings")
	go func() {
		err := ui.Main(func() {
			fooWindow = ui.NewWindow("", 1, 1, false)
			exportBox = createExportBox()
			settingsBox = createSettingsBox()
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

// TODO move export gui related funcitons to single file

func createExportBox() ui.Control {
	tab := ui.NewTab()
	tab.Append("SVG", createFormatPresets("svg"))
	tab.Append("PNG", createFormatPresets("png"))
	tab.Append("PDF", createFormatPresets("pdf"))

	return tab
}

func createFormatPresets(ext string) ui.Control {
	progressBar := ui.NewProgressBar()
	progressBar.SetValue(0)
	progressBar.Disable()

	button := ui.NewButton("Export")
	button.OnClicked(func(button *ui.Button) {
		if exportFunc != nil {
			done := make(chan bool)
			setProgress := func(n int) {
				ui.QueueMain(func() {
					progressBar.SetValue(n)
				})
			}
			go func() {
				exportFunc(exportSt.Of(ext).String("filename"))
				done <- true
				time.Sleep(time.Second)
				setProgress(0)
				progressBar.Disable()
			}()
			func() {
				n := 0
				step := 1
				progressBar.Enable()
				setProgress(0)
			filling:
				for {
					select {
					case <-done:
						setProgress(100)
						break filling
					default:
						n += step
						if n < 100 {
							setProgress(n)
						}
					}
				}
			}()
		}
	})

	box := ui.NewVerticalBox()
	if ext == "pdf" {
		currentOrientation := "L"
		currentPaperFormat := "A4"
		setDimensions := func() {
			width, height := A4_WIDTH, A4_HEIGHT
			if currentOrientation == "L" {
				width, height = height, width
			}
			if currentPaperFormat == "A3" {
				width = int(float64(width) * math.Sqrt2)
				height = int(float64(height) * math.Sqrt2)
			}
			if currentPaperFormat == "A5" {
				width = int(float64(width) / math.Sqrt2 )
				height = int(float64(height) / math.Sqrt2)
			}
			exportSt.Set("width", width)
			exportSt.Set("height", height)
		} 
		orientation := func() ui.Control {
			// TODO: use radio buttons once available in ui
			portrait := ui.NewButton("Portrait")
			landscape := ui.NewButton("Landscape")
			orient := func(b *ui.Button){
				if b.Text() == "Portrait" {
					portrait.Disable()
					landscape.Enable()
					currentOrientation = "P"
				}
				if b.Text() == "Landscape"  {
					landscape.Disable()
					portrait.Enable()
					currentOrientation = "L"
				}
				setDimensions()
			}
			portrait.OnClicked(orient)
			landscape.OnClicked(orient)
			landscape.Enable()
			buttons := line(pair{portrait, false}, pair{landscape, false})
			return line(pair{ui.NewLabel("Orientation"), true}, pair{buttons, false})
		}

		paperFormat := func() ui.Control {
			// TODO: use radio buttons once available in ui
			buttons := map[string]*ui.Button{
				"A3": ui.NewButton("A3"),
				"A4": ui.NewButton("A4"),
				"A5": ui.NewButton("A5"),
			}
			buttons["A4"].Disable()
			for _, btn := range buttons {
				btn.OnClicked(func(b *ui.Button) {
					for _, btn := range buttons {
						btn.Enable()
					}
					b.Disable()
					currentPaperFormat = b.Text()
					setDimensions()
				})
			}
			return line(pair{ui.NewLabel("format"), true}, pair{buttons["A3"], false}, pair{buttons["A4"], false}, pair{buttons["A5"], false})
		}
		setDimensions()
		box.Append(orientation(), false)
		box.Append(paperFormat(), false)
	}
	box.Append(createIntInput("width", 1, math.MaxInt32), false)
	box.Append(createIntInput("height", 1, math.MaxInt32), false)
	box.Append(createIntInput("zoom", -5, +5), false)
	box.Append(createExportAs(ext), false)
	box.Append(progressBar, true)
	box.Append(button, true)
	return box
}

func createExportAs(ext string) ui.Control {
	input := ui.NewEntry()
	input.SetText(exportSt.Of(ext).String("filename"))
	input.OnChanged(func(input *ui.Entry) {
		exportSt.Of(ext).Set("filename", input.Text())
	})
	button := ui.NewButton("Browse…")
	button.OnClicked(func(*ui.Button) {
		SaveFile(func(filename string) {
			if filepath.Ext(filename) != "."+ext {
				filename += "." + ext
			}
			input.SetText(filename)
			exportSt.Of(ext).Set("filename", filename)
		})
	})
	return line(pair{input, true}, pair{button, false})
}

func createIntInput(name string, min, max int) ui.Control {
	label := ui.NewLabel(name)

	input := ui.NewSpinbox(min, max)
	input.SetValue(exportSt.Int(name))
	input.OnChanged(func(*ui.Spinbox) {
		exportSt.Set(name, input.Value())
	})
	exportSt.OnChange(func(st storage.Storage, key string) {
		input.SetValue(exportSt.Int(name))
	})

	return line(pair{label, true}, pair{input, false})
}

func createFloatInput(name string, min, max int) ui.Control {
	label := ui.NewLabel(name)

	input := ui.NewSpinbox(min, max)
	input.SetValue(int(settingsSt.Float(name)))
	input.OnChanged(func(*ui.Spinbox) {
		settingsSt.Set(name, float64(input.Value()))
	})

	return line(pair{label, true}, pair{input, false})
}

func IsExportOn() bool {
	return exportWindow != nil
}

func ToggleExport(export func(string)) {
	exportFunc = export
	ui.QueueMain(func() {
		toggleWindow(&exportWindow, "Export", exportBox)
	})
}

// TODO move settings gui related funcitons to single file

func createSettingsBox() ui.Control {
	tab := ui.NewTab()

	tab.Append("general", createFloatInput("linewidth", 1, 4))
	// tab.Append("place", nil)
	// tab.Append("transition", nil)
	// tab.Append("arc", nil)

	return tab
}

func ToggleSettings() {
	ui.QueueMain(func() {
		toggleWindow(&settingsWindow, "Settings", settingsBox)
	})
}

func toggleWindow(windowPointer **ui.Window, name string, box ui.Control) {
	if *windowPointer != nil {
		win := *windowPointer
		win.SetChild(nil)
		win.Destroy()
		*windowPointer = nil
	} else {
		win := ui.NewWindow(name, 200, 100, false)
		win.SetChild(box)
		win.OnClosing(func(win *ui.Window) bool {
			win.SetChild(nil)
			win.Destroy()
			*windowPointer = nil
			return false
		})
		win.Show()
		*windowPointer = win
	}
}

func LoadFile(cb func(string)) {
	ui.QueueMain(func() {
		filename := ui.OpenFile(fooWindow)
		if filename != "" {
			cb(filename)
		}
	})
}

func SaveFile(cb func(string)) {
	ui.QueueMain(func() {
		filename := ui.SaveFile(fooWindow)
		if filename != "" {
			cb(filename)
		}
	})
}
