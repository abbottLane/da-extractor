package main

import (
	analyzers "daextractor/pkg/analyzers"
	"errors"
	"io/ioutil"
	"strings"

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
	instruction := widget.NewLabel("1. Select a text file to extract data from, \nor simply write you text into the input box:")
	// Create a new label to display the selected file name.
	selectedFileName := widget.NewLabel("")
	// Create a new label to display the selected file contents.
	textDisplayArea := widget.NewMultiLineEntry()
	textDisplayArea.Wrapping = fyne.TextWrapWord
	textDisplayAreaScroll := container.NewScroll(textDisplayArea)
	textDisplayAreaScroll.SetMinSize(fyne.NewSize(400, 600))
	// Create a new label to display the selected file's analyzed content'.
	resultDisplayArea := widget.NewMultiLineEntry()
	resultDisplayArea.Wrapping = fyne.TextWrapWord
	resultDisplayAreaScroll := container.NewScroll(resultDisplayArea)
	resultDisplayAreaScroll.SetMinSize(fyne.NewSize(400, 600))
	// Create a new label to display the loaded text.
	loadedTextArea := widget.NewLabel("INPUT TEXT")
	// Create a new label to display the analysis result.
	resultTextArea := widget.NewLabel("ANALYSIS RESULT (json format)")
	// Create a new widget to collect a list of strings from the user
	tagsetInstructions := widget.NewLabel("2. Enter a list of discourse functions to analyze for, \nseparated by newlines:")
	tagset := widget.NewMultiLineEntry()
	tagset.SetText("statement\nquestion\nexclamation\ndirective\nappreciation\nagreement\ndisagreement\nelaboration\nbackground\ncontinuation\nconjunction\nsummary\nrestatement\nother")
	tagset.SetMinRowsVisible(15)

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

	analyzeInstructions := widget.NewLabel("3. Click the button to analyze the text:")
	analyzeButton := widget.NewButton("Analyze", func() {
		// Create a button that will take the value of the content variable and send it to the discourse analyzer
		if len(content) != 0 {
			// read the tagset into a list of strings
			tagset := tagset.Text
			// split tagset on newlines
			tagset_list := strings.Split(tagset, "\n")

			resultDisplayArea.SetText("Analyzing...(this could take up to a minute or two)")
			//process the text content into sentences
			discourseAnalyzer := analyzers.NewDiscourseAnalyzer("openAI")
			analyzedResult := discourseAnalyzer.Analyze(string(content), tagset_list)
			resultDisplayArea.SetText(analyzedResult)

		} else if len(textDisplayArea.Text) != 0 {
			// read the tagset into a list of strings
			tagset := tagset.Text
			// split tagset on newlines
			tagset_list := strings.Split(tagset, "\n")

			resultDisplayArea.SetText("Analyzing...(this could take up to a minute or two)")
			//process the text content into sentences
			discourseAnalyzer := analyzers.NewDiscourseAnalyzer("openAI")
			analyzedResult := discourseAnalyzer.Analyze(textDisplayArea.Text, tagset_list)
			resultDisplayArea.SetText(analyzedResult)
		} else {
			// Show an error dialog if no file is selected
			errDialog := dialog.NewError(errors.New("no file selected"), w)
			errDialog.Show()
		}
	})

	// make two columns in the window: the first column is the instruction and loadfile, the second is the display area
	w.SetContent(
		container.NewHBox(
			container.NewVBox(
				instruction,
				selectedFileName,
				loadfile,
				tagsetInstructions,
				tagset,
				analyzeInstructions,
				analyzeButton,
			),
			container.NewVBox(
				loadedTextArea,
				textDisplayAreaScroll,
			),
			container.NewVBox(
				resultTextArea,
				resultDisplayAreaScroll,
			),
		),
	)

	w.Resize(fyne.NewSize(700, 700))
	w.ShowAndRun()
}
