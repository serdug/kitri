// Copyright (c) 2020 Sergey Dugaev. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file in the project root for more information.

// Package ui for GUI (front end)
package ui

import (
	// "fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/serdug/kitri/conti"
)

const preferenceCurrentTab = "currentTab"

type kitri struct {
	// The name of loaded schema JSON file
	schemaName *widget.Entry

	recGroup *widget.Group

	containers map[string]*fyne.Container

	// Working directory input
	wDirInput *widget.Entry

	// Buttons for input chart sections
	// Examples: 'assets', 'liabls', 'equity', 'revenues', 'expenses'
	sectEntry map[string]*widget.Button

	// Back and Next navigation buttons
	navigator map[string]*widget.Button

	// Buttons for input records
	// Example: 'rec-1'
	recEntry map[string]*recordMenuButton

	// Previous screen (source)
	source string

	// New record file name
	newRecFile *widget.Label
	newRecName string

	// Names of input chart sections
	section map[string]string

	// Input records
	record map[string]conti.Record
}

// newKitri initiates a new Kitri app struct
func newKitri() *kitri {
	return &kitri{
		sectEntry:  make(map[string]*widget.Button),
		recEntry:   make(map[string]*recordMenuButton),
		navigator:  make(map[string]*widget.Button),
		containers: make(map[string]*fyne.Container),
		section:    make(map[string]string),
		record:     make(map[string]conti.Record),
	}
}

// newSchema initiates a new conti.Schema struct
func newSchema() conti.Schema {
	return conti.Schema{}
}

// Show loads a Kitri window for the specified app context
func Show(app fyne.App) {
	kit := newKitri()
	app.Settings().SetTheme(theme.LightTheme())

	w := app.NewWindow("Kitri")
	// w.SetIcon(icon)

	// Main Menu
	/*
		newItem := fyne.NewMenuItem("New", nil)
		newItem.ChildMenu = fyne.NewMenu("",
			fyne.NewMenuItem("File", func() {fmt.Println("Menu New->File")}),
		)
		settingsItem := fyne.NewMenuItem("Settings", func() {fmt.Println("Menu Settings")})
	*/

	helpMenu := fyne.NewMenu("Help") /*
		fyne.NewMenuItem("Kitri Help", func() {
			widget.NewHyperlink("web", func() {
				link, err := url.Parse("https://kitribalance.com/help")
					if err!= nil {
						fyne.LogError("Could not parse URL", err)
					}
				return link
			}),

			// fmt.Println("Menu->Help")
		}),
	*/
	/*
		fyne.NewMenuItem("About Kitri", func() {
			widget.NewHyperlink("web", func() {
				link, err := url.Parse("https://kitribalance.com/about")
					if err!= nil {
						fyne.LogError("Could not parse URL", err)
					}
				return link
			}),

			// fmt.Println("Menu->About")
		}),
	*/

	mainMenu := fyne.NewMainMenu(
		// A Quit... item appended to the first menu

		// A Settings item appended to the first menu
		// fyne.NewMenu("Template", newItem, fyne.NewMenuItemSeparator(), settingsItem),

		helpMenu,
	)
	w.SetMainMenu(mainMenu)
	w.SetMaster()

	// Set navigation buttons
	kit.navigation(w)

	// A layout with separate tabs for Load and Create
	tabs := widget.NewTabContainer(
		widget.NewTabItemWithIcon("Home", theme.HomeIcon(), kit.homeScreen(w)),
		// widget.NewTabItemWithIcon("Home", theme.SettingsIcon(), kit.homeScreen(w)),
		widget.NewTabItemWithIcon("Calculate", theme.ViewRefreshIcon(), kit.templateScreen(w)),
		// widget.NewTabItemWithIcon("Open", theme.SearchIcon(), kit.templateScreen(w)),
		// widget.NewTabItemWithIcon("Open", theme.FileIcon(), kit.templateScreen(w)),
		// widget.NewTabItemWithIcon("Open", theme.FileApplicationIcon(), kit.templateScreen(w)),
	)

	tabs.SetTabLocation(widget.TabLocationLeading)

	// Load the latest tab from the previous session
	// Const: preferenceCurrentTab
	tabs.SelectTabIndex(app.Preferences().Int(preferenceCurrentTab))

	w.SetContent(tabs)

	kit.navigator["Review=>Output"].Hide()

	w.Resize(fyne.NewSize(1400, 800))
	w.ShowAndRun()
	// w.Show()

	// Save tab at exit
	// Const: preferenceCurrentTab
	app.Preferences().SetInt(preferenceCurrentTab, tabs.CurrentTabIndex())
}

// ***************************************************************************
// * METHODS
// ***************************************************************************
func (kit *kitri) navigation(win fyne.Window) {
	kit.navigator["Create=>Review"] = &widget.Button{
		// Alignment:     widget.ButtonAlignTrailing,
		IconPlacement: widget.ButtonIconTrailingText,
		Icon:          theme.NavigateNextIcon(),
		Text:          "Next",
		OnTapped: func() {
			// fmt.Println("Create->Next")
		},
	}

	kit.navigator["Load=>Review"] = &widget.Button{
		// Alignment:     widget.ButtonAlignTrailing,
		IconPlacement: widget.ButtonIconTrailingText,
		Icon:          theme.NavigateNextIcon(),
		Text:          "Next",
		OnTapped: func() {
			kit.containers["1"].Hide()

			kit.navigator["Load=>Review"].Hide()
			kit.navigator["Recalculate"].Hide()
			kit.navigator["Load<=Review"].Show()
			kit.navigator["SaveTemplate"].Show()
			kit.navigator["Review=>Output"].Show()

			kit.source = "1"

			kit.reviewInput("1", win)

			// fmt.Println("Load->Next")
		},
	}

	kit.navigator["Load<=Review"] = &widget.Button{
		// Alignment:     widget.ButtonAlignLeading,
		IconPlacement: widget.ButtonIconLeadingText,
		Icon:          theme.NavigateBackIcon(),
		Text:          "Back",
		OnTapped: func() {
			kit.containers["2"].Hide()
			kit.containers["1"].Show()

			kit.navigator["Load<=Review"].Hide()
			kit.navigator["SaveTemplate"].Hide()
			kit.navigator["Recalculate"].Hide()
			kit.navigator["Review=>Output"].Hide()
			kit.navigator["Load=>Review"].Show()

			// fmt.Println("Load<-Back")
		},
	}

	kit.navigator["SaveTemplate"] = &widget.Button{
		// Alignment:     widget.ButtonAlignLeading,
		IconPlacement: widget.ButtonIconLeadingText,
		Icon:          theme.DocumentSaveIcon(),
		Text:          "Save Template",
		OnTapped: func() {
			// fmt.Println("Save Template")
			dialog.ShowFileSave(
				func(writer fyne.URIWriteCloser, err error) {
					if err != nil {
						dialog.ShowError(err, win)
						return
					}
					if writer == nil {
						// fmt.Println("Save cancelled")
						return
					}
					schemaWriter(*kit, fileNamed(writer))
				},
				win,
			)
		},
	}

	kit.navigator["Review=>Output"] = &widget.Button{
		// Alignment:     widget.ButtonAlignTrailing,
		IconPlacement: widget.ButtonIconTrailingText,
		Icon:          theme.NavigateNextIcon(),
		Text:          "Calculate",
		OnTapped: func() {
			kit.containers["2"].Hide()

			kit.navigator["Load<=Review"].Hide()
			kit.navigator["SaveTemplate"].Hide()
			kit.navigator["Review=>Output"].Hide()
			kit.navigator["Review<=Output"].Show()

			kit.navigator["Recalculate"].Show()
			kit.navigator["SaveOutput"].Show()

			kit.source = "2"

			kit.showOutput(win)

			// fmt.Println("Review->Next")
		},
	}

	kit.navigator["Review<=Output"] = &widget.Button{
		// Alignment:     widget.ButtonAlignLeading,
		IconPlacement: widget.ButtonIconLeadingText,
		Icon:          theme.NavigateBackIcon(),
		Text:          "Back",
		OnTapped: func() {
			kit.containers["3"].Hide()
			kit.containers["2"].Show()

			kit.navigator["Review<=Output"].Hide()
			kit.navigator["SaveOutput"].Hide()
			kit.navigator["Recalculate"].Hide()
			kit.navigator["SaveTemplate"].Show()
			kit.navigator["Review=>Output"].Show()
			kit.navigator["Load<=Review"].Show()

			kit.source = "3"

			// fmt.Println("Review<-Back")
		},
	}

	kit.navigator["Recalculate"] = &widget.Button{
		// Alignment:     widget.ButtonAlignLeading,
		IconPlacement: widget.ButtonIconLeadingText,
		Icon:          theme.ViewRefreshIcon(),
		Text:          "Recalculate",
		OnTapped: func() {
			kit.refreshOutput(win)
			// fmt.Println("Recalculate")
		},
	}

	kit.navigator["SaveOutput"] = &widget.Button{
		// Alignment:     widget.ButtonAlignLeading,
		IconPlacement: widget.ButtonIconLeadingText,
		Icon:          theme.DocumentSaveIcon(),
		Text:          "Recalculate and Save",
		OnTapped: func() {
			// fmt.Println("Save Output")
			dialog.ShowFileSave(
				func(writer fyne.URIWriteCloser, err error) {
					if err != nil {
						dialog.ShowError(err, win)
						return
					}
					if writer == nil {
						// fmt.Println("Save cancelled")
						return
					}
					outputWriter(*kit, fileNamed(writer))
				},
				win,
			)
		},
	}
}
