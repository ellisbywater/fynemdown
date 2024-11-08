package main

import (
	"io/ioutil"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type config struct {
	EditWidget    *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile   fyne.URI
	SaveMenuItem  *fyne.MenuItem
}

func (app *config) makeUI() (*widget.Entry, *widget.RichText) {
	edit := widget.NewMultiLineEntry()
	preview := widget.NewRichTextFromMarkdown("")
	app.EditWidget = edit
	app.PreviewWidget = preview
	edit.OnChanged = preview.ParseMarkdown

	return edit, preview
}

var cfg config

func main() {
	// create a fyne app
	a := app.New()

	// create a window for the app
	win := a.NewWindow("Markdown")
	cfg.createMenuItems(win)
	// get the user interface
	edit, preview := cfg.makeUI()
	// set the content of the window
	win.SetContent(container.NewHSplit(edit, preview))
	// show window and run app
	win.Resize(fyne.NewSize(800, 800))
	win.CenterOnScreen()
	win.ShowAndRun()
}

func (app *config) createMenuItems(win fyne.Window) {
	openMenuItem := fyne.NewMenuItem("Open...", app.openFile(win))
	saveMenuItem := fyne.NewMenuItem("Save", app.saveFile(win))
	app.SaveMenuItem = saveMenuItem
	app.SaveMenuItem.Disabled = true
	saveAsMenuItem := fyne.NewMenuItem("Save as...", app.saveAs(win))
	fileMenu := fyne.NewMenu("File", openMenuItem, saveMenuItem, saveAsMenuItem)
	menu := fyne.NewMainMenu(fileMenu)

	win.SetMainMenu(menu)
}

var filter = storage.NewExtensionFileFilter([]string{".md", ".MD"})

func (app *config) saveFile(win fyne.Window) func() {
	return func() {
		if app.CurrentFile != nil {
			write, err := storage.Writer(app.CurrentFile)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			write.Write([]byte(app.EditWidget.Text))
			defer write.Close()
		}
	}
}

func (app *config) openFile(win fyne.Window) func() {
	return func() {
		openDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
			}
			if reader == nil {
				return
			}
			defer reader.Close()

			data, err := ioutil.ReadAll(reader)
			if err != nil {
				dialog.ShowError(err, win)
			}

			app.EditWidget.SetText(string(data))
			app.CurrentFile = reader.URI()
			win.SetTitle(win.Title() + "-" + reader.URI().Name())
			app.SaveMenuItem.Disabled = false
		}, win)
		openDialog.SetFilter(filter)
		openDialog.Show()
	}
}

func (app *config) saveAs(win fyne.Window) func() {
	return func() {
		saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if writer == nil { // user canceled
				return
			}

			if !strings.HasSuffix(strings.ToLower(writer.URI().String()), ".md") {
				dialog.ShowInformation("Error", "Please name your file with an .md extension", win)
				return
			}
			// Save file
			writer.Write([]byte(app.EditWidget.Text))
			app.CurrentFile = writer.URI()
			defer writer.Close()

			win.SetTitle(win.Title() + "-" + writer.URI().Name())
			app.SaveMenuItem.Disabled = false
		}, win)
		saveDialog.SetFileName("untitled.md")
		saveDialog.SetFilter(filter)
		saveDialog.Show()
	}
}
