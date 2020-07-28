// Copyright (c) 2020 Sergey Dugaev. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file in the project root for more information.

// Package ui for GUI (front end)
package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// Global variable for a fixed size array of input field names
var fieldNames = [...]string{"Path", "Assets", "Liabilities", "Equity", "Revenues", "Expenses"}

func showGuidance(txt string) fyne.CanvasObject {
	lbl := widget.NewLabel(txt)
	lbl.Wrapping = fyne.TextWrapWord

	scroller := widget.NewVScrollContainer(lbl)

	return scroller
}

// ***************************************************************************
// * METHODS
// ***************************************************************************

func (kit *kitri) templateScreen(win fyne.Window) fyne.CanvasObject {
	buttons := widget.NewHBox(
		kit.navigator["Load<=Review"],
		kit.navigator["Review<=Output"],
		layout.NewSpacer(),
		kit.navigator["Load=>Review"],
		kit.navigator["SaveTemplate"],
		kit.navigator["Review=>Output"],
		kit.navigator["Recalculate"],
		kit.navigator["SaveOutput"],
	)

	kit.navigator["Load<=Review"].Hide()
	kit.navigator["SaveTemplate"].Hide()
	kit.navigator["Review<=Output"].Hide()
	kit.navigator["Review=>Output"].Hide()
	kit.navigator["Recalculate"].Hide()
	kit.navigator["SaveOutput"].Hide()

	borderLayout := layout.NewBorderLayout(nil, buttons, nil, nil)
	// borderLayout := layout.NewBorderLayout(top, bottom, left, right)

	pickButton := widget.NewButtonWithIcon(
		"Select file",
		// "Browse",
		theme.SearchIcon(),
		func() {
			kit.schemaDialog(win)

			kit.containers["main"].Refresh()
		},
	)

	label := widget.NewLabel("Load configuration")
	label.TextStyle = fyne.TextStyle{Bold: true}

	box := widget.NewHBox(
		layout.NewSpacer(),
		pickButton,
		layout.NewSpacer(),
	)

	filename := widget.NewEntry()
	filename.SetPlaceHolder("full name of the file...")
	// filename.Wrapping = fyne.TextWrapBreak
	// filename.Wrapping = fyne.TextTruncate
	kit.schemaName = filename

	left := widget.NewVBox(
		// layout.NewSpacer(),
		label,
		box,
		filename,
	)

	right := showGuidance(txtGuidanceOpen)

	split := widget.NewHSplitContainer(widget.NewVScrollContainer(left), widget.NewVScrollContainer(right))

	kit.containers["1"] = fyne.NewContainerWithLayout(
		layout.NewGridLayout(1),
		split,
	)

	wrapper := fyne.NewContainerWithLayout(
		layout.NewGridLayout(1),
		kit.containers["1"],
	)

	kit.containers["main"] = wrapper

	return fyne.NewContainerWithLayout(
		borderLayout,
		buttons,
		kit.containers["main"],
	)
}

func (kit *kitri) schemaDialog(win fyne.Window) {
	fd := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
		if err == nil && file == nil {
			// Do nothing
			return
		}
		if err != nil {
			dialog.ShowError(err, win)
			// Handle error
			return
		}

		fileExt := file.URI().Extension()
		fileExt = strings.ToLower(fileExt)

		if fileExt == ".json" || fileExt == ".yaml" || fileExt == ".yml" {
			p := fmt.Sprintf("%s", file.URI())

			// Remove "file://" from file.URI() added by fyne
			p = strings.Replace(p, "file://", "", 1)

			// Note: The values returned by Split() have the property
			// that path=dir+file
			// Otherwise, it's possible to correctly join elements of the path
			// using Join() from path/filepath
			d, f := filepath.Split(p)
			kit.schemaName.SetText(d + f)
		}
	}, win)
	extFilter := storage.NewExtensionFileFilter([]string{".json", ".yaml", ".yml"})
	fd.SetFilter(extFilter)
	fd.Show()
}
