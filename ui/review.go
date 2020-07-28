// Copyright (c) 2020 Sergey Dugaev. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file in the project root for more information.

// Package ui for GUI (front end)
package ui

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"

	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/serdug/kitri/conti"
	"github.com/serdug/kitri/handlers"
)

const recPrefix string = "rec-"

var sectionKeys = [...]string{"assets", "liabls", "equity", "revenues", "expenses"}

// ***************************************************************************
// * METHODS
// ***************************************************************************
func (kit *kitri) reviewInput(source string, win fyne.Window) {
	var s conti.Schema

	// Note: clean up old records before loading a config
	kit.recEntry = make(map[string]*recordMenuButton)

	fileExt := filepath.Ext(kit.schemaName.Text)

	switch fileExt {
	case ".json":
		s = handlers.ReadSchemaJSON(kit.schemaName.Text)

	case ".yaml":
		s = handlers.ReadSchemaYAML(kit.schemaName.Text)

	case ".yml":
		s = handlers.ReadSchemaYAML(kit.schemaName.Text)
	}

	left := kit.makeChartGroup(s, win)

	recrdGroup := kit.arrangeRecords(s, win)

	addButton := &widget.Button{
		// Alignment:     widget.ButtonAlignLeading,
		IconPlacement: widget.ButtonIconLeadingText,
		Icon:          theme.ContentAddIcon(),
		Text:          "Add",
		OnTapped: func() {
			scount := strconv.Itoa(len(kit.recEntry) + 1)
			recKey := recPrefix + scount
			rec := newRecordMenuButton()
			rec.Alignment = widget.ButtonAlignLeading

			rec.Text = kit.newRecFile.Text
			// fmt.Println("Add record:", rec.Text)

			rec.Icon = theme.CheckButtonIcon()
			kit.recEntry[recKey] = rec
			kit.recGroup.Append(rec)

			// Note: the order of the following refreshes is important!!!
			if kit.source == "2" {
				// Note: important, refreshes screen 3
				kit.containers["3"].Refresh()
			} else {
				// Note: important, refreshes screen 2 after returning from screen 3
				kit.containers["2"].Refresh()
			}

			// fmt.Println("Add a record")
		},
	}

	pickButton := widget.NewButtonWithIcon("Browse",
		theme.SearchIcon(),
		func() {
			kit.recordDialog(win)
		},
	)

	filename := widget.NewLabel("")
	kit.newRecFile = filename

	addRecord := widget.NewHBox(
		pickButton,
		filename,
		layout.NewSpacer(),
		addButton,
	)

	noteAddRecords := widget.NewLabel("Files (.csv) with records may be added to or erased from the input")
	noteAddRecords.Wrapping = fyne.TextWrapWord

	records := widget.NewVBox(
		recrdGroup,
		noteAddRecords,
		addRecord,
	)

	right := widget.NewVScrollContainer(records)

	kit.containers["2.2"] = fyne.NewContainerWithLayout(
		layout.NewGridLayout(1),
		right,
	)

	split := widget.NewHSplitContainer(widget.NewVScrollContainer(left), kit.containers["2.2"])

	kit.containers["2"] = fyne.NewContainerWithLayout(
		layout.NewGridLayout(1),
		split,
	)

	kit.containers["main"].AddObject(kit.containers["2"])

	kit.containers["main"].Refresh()
}

// makeChartGroup renders chart sections (and a path to working directory)
// with file names for review / edit
func (kit *kitri) makeChartGroup(s conti.Schema, win fyne.Window) fyne.Widget {
	var sectionBox *widget.Box

	pathEntry := widget.NewEntry()
	kit.wDirInput = pathEntry
	kit.wDirInput.SetText(s.Path)

	slist := make([]fyne.CanvasObject, len(sectionKeys))
	for i, skey := range sectionKeys {
		kit.sectionButton(skey, win)

		switch skey {
		case "assets":
			sectionBox = widget.NewVBox(
				widget.NewLabel("Assets:"),
				kit.sectEntry["assets"],
			)
			kit.section["assets"] = s.Chart.Assets
			kit.sectEntry["assets"].Text = s.Chart.Assets

		case "liabls":
			sectionBox = widget.NewVBox(
				widget.NewLabel("Liabilities:"),
				kit.sectEntry["liabls"],
			)
			kit.section["liabls"] = s.Chart.Liabilities
			kit.sectEntry["liabls"].Text = s.Chart.Liabilities

		case "equity":
			sectionBox = widget.NewVBox(
				widget.NewLabel("Equity:"),
				kit.sectEntry["equity"],
			)
			kit.section["equity"] = s.Chart.Equity
			kit.sectEntry["equity"].Text = s.Chart.Equity

		case "revenues":
			sectionBox = widget.NewVBox(
				widget.NewLabel("Revenues:"),
				kit.sectEntry["revenues"],
			)
			kit.section["revenues"] = s.Chart.Revenues
			kit.sectEntry["revenues"].Text = s.Chart.Revenues

		case "expenses":
			sectionBox = widget.NewVBox(
				widget.NewLabel("Expenses:"),
				kit.sectEntry["expenses"],
			)
			kit.section["expenses"] = s.Chart.Expenses
			kit.sectEntry["expenses"].Text = s.Chart.Expenses

		default:
			// fmt.Printf("Chart section key %s not found\n", skey)
			dialog.ShowInformation("Warning", "Chart section key '"+skey+"'' not found!", win)
		}
		slist[i] = sectionBox
	}

	chartGroup := widget.NewGroup("Chart of Accounts", slist...)

	labelWorkDir := widget.NewLabel("Working directory:")

	return widget.NewVBox(
		chartGroup,
		labelWorkDir,
		pathEntry,
	)
}

// sectionButton creates and describes a section entry button for a section
// by its idenifier 'skey'
func (kit *kitri) sectionButton(skey string, win fyne.Window) {
	kit.sectEntry[skey] = &widget.Button{
		Alignment: widget.ButtonAlignLeading,
	}

	kit.sectEntry[skey].OnTapped = func() {
		var d, f string
		// fmt.Println("Section Dialog", kit.sectEntry[skey].Text)

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

			if fileExt == ".csv" {
				p := fmt.Sprintf("%s", file.URI())

				// Remove "file://" from file.URI() added by fyne
				p = strings.Replace(p, "file://", "", 1)

				// Note: The values returned by Split() have the property
				// that path=dir+file
				// Otherwise, it's possible to correctly join elements of the path
				// using Join() from path/filepath
				d, f = filepath.Split(p)

				switch skey {
				case "assets":
					kit.sectEntry["assets"].SetText(f)
					// fmt.Println("Assets - replaced:", "(file '", f, "')")

				case "liabls":
					kit.sectEntry["liabls"].SetText(f)
					// fmt.Println("Liabilities - replaced:", "(file '", f, "')")

				case "equity":
					kit.sectEntry["equity"].SetText(f)
					// fmt.Println("Equity - replaced:", "(file '", f, "')")

				case "revenues":
					kit.sectEntry["revenues"].SetText(f)
					// fmt.Println("Revenues - replaced:", "(file '", f, "')")

				case "expenses":
					kit.sectEntry["expenses"].SetText(f)
					// fmt.Println("Expenses - replaced:", "(file '", f, "')")
				}

				// TO DO: WARN of a wrong directory!!!
				kit.wDirInput.SetText(d)
			}

		}, win)
		extFilter := storage.NewExtensionFileFilter([]string{".csv"})
		fd.SetFilter(extFilter)
		fd.Show()
	}
}

// arrangeRecords renders files with records for review
func (kit *kitri) arrangeRecords(s conti.Schema, win fyne.Window) fyne.CanvasObject {
	var (
		scount string
		recKey string
	)

	rlist := make([]fyne.CanvasObject, len(s.Records))

	remove := fyne.NewMenuItem("Remove", nil)

	remove.ChildMenu = fyne.NewMenu(
		"",
		fyne.NewMenuItem("Erase from input template", func() { fmt.Println("Erase record") }),
	)

	// fmt.Println("Records (loaded schema):", len(s.Records))

	for i := range s.Records {
		scount = strconv.Itoa(i + 1)
		recKey = recPrefix + scount

		rec := newRecordMenuButton()
		rec.Alignment = widget.ButtonAlignLeading
		rec.Text = s.Records[i].Id
		if s.Records[i].Include != 1 {
			rec.Icon = theme.CheckButtonIcon()

			rec.menu = fyne.NewMenu(
				"",
				fyne.NewMenuItem(
					"Include in calculation",
					func() {
						// fmt.Println("Include record", recKey)
					},
				),
				remove,
			)

			kit.record[recKey] = conti.Record{
				Include: 0,
				Id:      s.Records[i].Id,
			}
		} else {
			rec.Icon = theme.CheckButtonCheckedIcon()

			rec.menu = fyne.NewMenu(
				"",
				fyne.NewMenuItem(
					"Exclude from calculation",
					func() {
						// fmt.Println("Exclude record", recKey)
					},
				),
				remove,
			)

			kit.record[recKey] = conti.Record{
				Include: 1,
				Id:      s.Records[i].Id,
			}
		}

		kit.recEntry[recKey] = rec

		rlist[i] = rec
	}

	group := widget.NewGroup("Transactions", rlist...)
	kit.recGroup = group

	return group
}

func (kit *kitri) recordDialog(win fyne.Window) {
	var d, f string

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

		if fileExt == ".csv" {
			p := fmt.Sprintf("%s", file.URI())

			// Remove "file://" from file.URI() added by fyne
			p = strings.Replace(p, "file://", "", 1)

			// Note: The values returned by Split() have the property
			// that path=dir+file
			// Otherwise, it's possible to correctly join elements of the path
			// using Join() from path/filepath
			d, f = filepath.Split(p)
			if d != kit.wDirInput.Text {
				// fmt.Println("Add record->Wrong directory")
				return
			}
			kit.newRecName = f
			kit.newRecFile.SetText(f)
			// fmt.Println("recordDialog:", kit.newRecFile.Text, "(file '", f, "')")
		}

	}, win)
	extFilter := storage.NewExtensionFileFilter([]string{".csv"})
	fd.SetFilter(extFilter)
	fd.Show()
}

type recordMenuButton struct {
	widget.Button
	menu *fyne.Menu
}

func newRecordMenuButton() *recordMenuButton {
	b := &recordMenuButton{}
	b.ExtendBaseWidget(b)
	return b
}

func (b *recordMenuButton) Tapped(e *fyne.PointEvent) {
	// fmt.Println("Tapped: record", b.Text)
	remove := fyne.NewMenuItem("Remove", nil)

	remove.ChildMenu = fyne.NewMenu(
		"",
		fyne.NewMenuItem(
			"Erase from input template",
			func() {
				// fmt.Println("Erase record", b.Text)
				b.Hide()
			},
		),
	)
	if b.Icon == theme.CheckButtonIcon() {
		b.menu = fyne.NewMenu(
			"",
			fyne.NewMenuItem(
				"Include in calculation",
				func() {
					// fmt.Println("Include record", b.Text)
					b.Icon = theme.CheckButtonCheckedIcon()
					b.Refresh()
				},
			),
			remove,
		)
	} else {
		b.menu = fyne.NewMenu(
			"",
			fyne.NewMenuItem(
				"Exclude from calculation",
				func() {
					// fmt.Println("Exclude record", b.Text)
					b.Icon = theme.CheckButtonIcon()
					b.Refresh()
				},
			),
			remove,
		)
	}
	widget.ShowPopUpMenuAtPosition(b.menu, fyne.CurrentApp().Driver().CanvasForObject(b), e.AbsolutePosition)
}
