// Copyright (c) 2020 Sergey Dugaev. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file in the project root for more information.

// Package ui for GUI (front end)
package ui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

func (kit *kitri) homeScreen(win fyne.Window) fyne.CanvasObject {
	lbl := widget.NewLabel(txtAboutProject)
	lbl.Wrapping = fyne.TextWrapWord
	lbl.Alignment = fyne.TextAlignCenter

	contents := widget.NewVBox(
		layout.NewSpacer(),
		lbl,
		layout.NewSpacer(),
	)

	wrapper := fyne.NewContainerWithLayout(
		layout.NewAdaptiveGridLayout(1),
		contents,
	)

	kit.containers["home"] = wrapper

	return wrapper
}
