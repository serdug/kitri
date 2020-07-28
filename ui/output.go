// Copyright (c) 2020 Sergey Dugaev. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file in the project root for more information.

// Package ui for GUI (front end)
package ui

import (
	"fmt"
	"math"
	"strconv"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"

	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/serdug/kitri/conti"
)

const symbolsInDescription int = 35
const splitOffset float64 = 0.37

// templateSchema creates a schema request
func templateSchema(kit kitri) conti.Schema {
	var scount string
	var recKey string
	var s conti.Schema
	var ent recordMenuButton

	recs := len(kit.recEntry)
	s.Records = make([]conti.Record, recs)

	s.Path = kit.wDirInput.Text

	s.Chart.Assets = kit.sectEntry["assets"].Text
	s.Chart.Liabilities = kit.sectEntry["liabls"].Text
	s.Chart.Equity = kit.sectEntry["equity"].Text
	s.Chart.Revenues = kit.sectEntry["revenues"].Text
	s.Chart.Expenses = kit.sectEntry["expenses"].Text

	for i := 0; i < recs; i++ {
		scount = strconv.Itoa(i + 1)
		recKey = recPrefix + scount
		ent = *kit.recEntry[recKey]
		if !ent.Visible() {
			continue
		}

		s.Records[i].Id = ent.Text
		if ent.Icon == theme.CheckButtonCheckedIcon() {
			s.Records[i].Include = 1
		} else {
			s.Records[i].Include = 0
		}
	}

	return s
}

// arrangeOutput creates an object
func arrangeOutput(cats []conti.Categories) fyne.CanvasObject {
	var catTxt, sectTxt, nameTxt, staTxt, difTxt, endTxt *widget.Label
	var cat, sect, name, end, dif, sta string

	// Note: print using localized formatting with golang.org/x/text/message
	p := message.NewPrinter(language.English)

	catTitle := widget.NewLabelWithStyle("Cat", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	sectTitle := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	nameTitle := widget.NewLabelWithStyle("Description", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	staTitle := widget.NewLabelWithStyle("Starting Value", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})
	difTitle := widget.NewLabelWithStyle("Change", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})
	endTitle := widget.NewLabelWithStyle("Ending Value", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})

	// Add another row for column titles
	catCol := make([]fyne.CanvasObject, len(cats)+1)
	sectCol := make([]fyne.CanvasObject, len(cats)+1)
	nameCol := make([]fyne.CanvasObject, len(cats)+1)
	staCol := make([]fyne.CanvasObject, len(cats)+1)
	difCol := make([]fyne.CanvasObject, len(cats)+1)
	endCol := make([]fyne.CanvasObject, len(cats)+1)

	// The first row contains column titles
	catCol[0] = catTitle
	sectCol[0] = sectTitle
	nameCol[0] = nameTitle
	staCol[0] = staTitle
	difCol[0] = difTitle
	endCol[0] = endTitle

	for i, one := range cats {
		cat = one.Cat
		sect = one.Sect

		// Truncate the name
		if len(one.Name) <= symbolsInDescription {
			name = one.Name
		} else {
			name = one.Name[0:symbolsInDescription] + "..."
		}

		sta = p.Sprintf("%.2f", one.Bal.Sta)

		if math.Abs(math.Round(one.Bal.Dif*100)/100) < 0.01 {
			dif = "0.00"
		} else {
			dif = p.Sprintf("%.2f", one.Bal.Dif)
		}

		if math.Abs(math.Round(one.Bal.End*100)/100) < 0.01 {
			end = "0.00"
		} else {
			end = p.Sprintf("%.2f", one.Bal.End)
		}

		catTxt = widget.NewLabel(cat)
		sectTxt = widget.NewLabel(sect)
		nameTxt = widget.NewLabel(name)
		staTxt = widget.NewLabel(sta)
		difTxt = widget.NewLabel(dif)
		endTxt = widget.NewLabel(end)

		staTxt.Alignment = fyne.TextAlignTrailing
		difTxt.Alignment = fyne.TextAlignTrailing
		endTxt.Alignment = fyne.TextAlignTrailing

		catCol[i+1] = catTxt
		sectCol[i+1] = sectTxt
		nameCol[i+1] = nameTxt
		staCol[i+1] = staTxt
		difCol[i+1] = difTxt
		endCol[i+1] = endTxt
	}

	catCont := widget.NewVBox(catCol...)
	sectCont := widget.NewVBox(sectCol...)
	nameCont := widget.NewVBox(nameCol...)
	staCont := widget.NewVBox(staCol...)
	difCont := widget.NewVBox(difCol...)
	endCont := widget.NewVBox(endCol...)

	return widget.NewHBox(catCont, sectCont, nameCont, staCont, difCont, endCont)
}

// ***************************************************************************
// * METHODS
// ***************************************************************************
// showOutput renders calculation results in a two-column grid container
func (kit *kitri) showOutput(win fyne.Window) {
	s := templateSchema(*kit)

	cats, alert := conti.Accounts(s)
	if len(alert.Code) != 0 {
		dialog.ShowInformation("Information", alert.Code+"\n"+alert.Hint, win)
		if alert.Error != nil {
			dialog.ShowError(alert.Error, win)
		}
		fmt.Println(alert.Code)
	}

	contents := arrangeOutput(cats)

	right := widget.NewVScrollContainer(contents)

	kit.containers["3.2"] = fyne.NewContainerWithLayout(
		layout.NewGridLayout(1),
		right,
	)

	split := widget.NewHSplitContainer(kit.containers["2.2"], right)
	split.SetOffset(splitOffset)

	kit.containers["3"] = fyne.NewContainerWithLayout(
		layout.NewGridLayout(1),
		split,
	)

	kit.containers["main"].AddObject(kit.containers["3"])

	kit.containers["main"].Refresh()
}

// refreshOutput refreshes calculation results
func (kit *kitri) refreshOutput(win fyne.Window) {
	s := templateSchema(*kit)

	cats, alert := conti.Accounts(s)
	if len(alert.Code) != 0 {
		dialog.ShowInformation("Information", alert.Code+"\n"+alert.Hint, win)
		if alert.Error != nil {
			dialog.ShowError(alert.Error, win)
		}
		fmt.Println(alert.Code)
	}

	contents := arrangeOutput(cats)

	right := widget.NewVScrollContainer(contents)

	kit.containers["3"].Hide()

	split := widget.NewHSplitContainer(kit.containers["2.2"], right)
	split.SetOffset(splitOffset)

	kit.containers["3"] = fyne.NewContainerWithLayout(
		layout.NewAdaptiveGridLayout(1),
		split,
	)

	kit.containers["main"].AddObject(kit.containers["3"])
	kit.containers["3"].Show()

	kit.containers["main"].Refresh()
}
