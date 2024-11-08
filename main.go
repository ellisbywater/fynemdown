package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
	openMenuItem := fyne.NewMenuItem("Open...", func() {

	})
	saveMenuItem := fyne.NewMenuItem("Save", func() {})
	saveAsMenuItem := fyne.NewMenuItem("Save as...", app.saveAs(win))
	fileMenu := fyne.NewMenu("File", openMenuItem, saveMenuItem, saveAsMenuItem)
	menu := fyne.NewMainMenu(fileMenu)

	win.SetMainMenu(menu)
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
			// Save file
			writer.Write([]byte(app.EditWidget.Text))
			app.CurrentFile = writer.URI()
			defer writer.Close()

			win.SetTitle(win.Title() + "-" + writer.URI().Name())
			app.SaveMenuItem.Disabled = false
		}, win)
		saveDialog.Show()
	}
}
