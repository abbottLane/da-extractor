package main

import (
	analyzers "daextractor/pkg/analyzers"
	"errors"
	"io/ioutil"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func main() {
	var content []byte
	a := app.New()
	w := a.NewWindow("DA Extractor!")
	// Create a new label and button.
	instruction := widget.NewLabel("Select a text file to extract data from:")
	selectedFileName := widget.NewLabel("")
	textDisplayArea := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	resultDisplayArea := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	loadedTextArea := widget.NewLabel("Loaded text displays here...")
	resultTextArea := widget.NewLabel("Analysis displayed here...")
	// Create a text entry widget for the API key
	keyInput := widget.NewEntry()
	keyInput.SetPlaceHolder("Enter your API key here")

	loadfile := widget.NewButton("Choose File", func() {
		// Create a new file open dialog for loading files.
		// The dialog will filter to allow only text files to be selected.
		dialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			// Read the file into a string
			content, err = ioutil.ReadAll(reader)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			// read the filename into a string
			filename := reader.URI().Name()

			// Set the label text to the file contents
			selectedFileName.SetText(string(filename))

			//Set the content to the display area
			textDisplayArea.SetText(string(content))
		}, w)
		dialog.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
		dialog.Show()
	})

	analyzeButton := widget.NewButton("Analyze", func() {
		// Create a button that will take the value of the content variable and send it to the discourse analyzer
		if len(content) != 0 {
			println("file selected")
			//process the text content into sentences
			discourseAnalyzer := analyzers.NewDiscourseAnalyzer("openAI", keyInput.Text)
			analyzedResult := discourseAnalyzer.Analyze(string(content))
			println(analyzedResult)
			resultDisplayArea.SetText("analyzed palceholdr text")

		} else {
			println("no file selected")
			// Show an error dialog if no file is selected
			errDialog := dialog.NewError(errors.New("No file selected"), w)
			errDialog.Show()
		}
	})

	// make two columns in the window: the first column is the instruction and loadfile, the second is the display area
	w.SetContent(
		container.NewHBox(
			container.NewVBox(
				instruction,
				keyInput,
				selectedFileName,
				loadfile,
				analyzeButton,
			),
			container.NewVBox(
				loadedTextArea,
				textDisplayArea,
			),
			container.NewVBox(
				resultTextArea,
				resultDisplayArea,
			),
		),
	)

	w.Resize(fyne.NewSize(700, 700))
	w.ShowAndRun()
}
