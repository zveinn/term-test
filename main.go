package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}
	app := tview.NewApplication()
	mainLeft := newPrimitive("Main content left side")
	mainRight := newPrimitive("Main content left side")

	// Create custom menu system
	menuLabel := tview.NewTextView().
		SetText("Menu").
		SetTextColor(tcell.ColorWhite)

	// Create the dropdown list (initially not visible)
	menuOptions := []string{"Option 1", "Option 2", "Option 3", "Option 4", "Option 5"}
	menuList := tview.NewList()
	for _, option := range menuOptions {
		menuList.AddItem(option, "", 0, nil)
	}
	menuList.SetBackgroundColor(tcell.ColorBlack)
	menuList.SetMainTextColor(tcell.ColorWhite)
	menuList.SetSelectedTextColor(tcell.ColorWhite)
	menuList.SetSelectedBackgroundColor(tcell.ColorBlue)
	menuList.ShowSecondaryText(false)

	// Use the menu label as the header element

	var menuOpen bool = false

	footer := newPrimitive("footer")

	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(0, 0).
		SetBorders(true).
		AddItem(footer, 2, 0, 1, 2, 0, 0, false).
		AddItem(menuLabel, 0, 0, 1, 2, 0, 0, false)

	// Set up menu list selection handler
	menuList.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		// Handle menu selection here - close dropdown and return focus to grid
		app.SetRoot(grid, true).SetFocus(grid)
		menuOpen = false
	})

	// Add input capture to menu list for escape key
	menuList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			// Close menu without selection
			app.SetRoot(grid, true).SetFocus(grid)
			menuOpen = false
			return nil
		}
		return event
	})

	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'M', 'm':
			if !menuOpen {
				// Open the dropdown menu using Pages for proper overlay
				menuOpen = true

				// Create pages for overlay
				pages := tview.NewPages()
				pages.AddPage("main", grid, true, true)

				// Create a positioned container for the menu at top-left
				menuContainer := tview.NewFlex().SetDirection(tview.FlexColumn).
					AddItem(menuList, 20, 0, true).
					AddItem(nil, 0, 1, false) // Spacer to push menu to left

				// Position at top with spacer below
				topContainer := tview.NewFlex().SetDirection(tview.FlexRow).
					AddItem(menuContainer, len(menuOptions)+2, 0, true).
					AddItem(nil, 0, 1, false) // Spacer to push menu to top

				pages.AddPage("menu", topContainer, true, true)

				app.SetRoot(pages, true).SetFocus(menuList)
			}
			return nil
		case 'q', 'Q':
			// Allow quitting with 'q'
			app.Stop()
			return nil
		}
		return event
	})

	grid.AddItem(mainLeft, 1, 0, 1, 1, 0, 0, false)
	grid.AddItem(mainRight, 1, 1, 1, 1, 0, 0, false)

	app.SetRoot(grid, true).SetFocus(grid)
	err := app.Run()
	if err != nil {
		panic(err)
	}
}
