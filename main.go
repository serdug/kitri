// Copyright (c) 2020 Sergey Dugaev. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file in the project root for more information.

package main

import (
	// "fmt"

	"fyne.io/fyne/app"

	"github.com/serdug/kitri/ui"
)

func main() {
	a := app.NewWithID("io.kitri") // for preferences
	// a := app.New()

	ui.Show(a)
	// a.Run()

	tidyUp()
}

func tidyUp() {
	// a.Preferences().SetInt(preferenceCurrentTab, tabs.CurrentTabIndex())
	// fmt.Println("Exited")
}
