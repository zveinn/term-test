package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	menuBackground         = tcell.ColorViolet.TrueColor()
	menuForeground         = tcell.ColorWhite.TrueColor()
	menuListBackground     = tcell.ColorViolet.TrueColor()
	menuListForeground     = tcell.ColorWhite.TrueColor()
	menuSelectedBackground = tcell.ColorTeal.TrueColor()
	menuSelectedForeground = tcell.ColorWhite.TrueColor()
	footerBackground       = tcell.ColorViolet.TrueColor()
	footerForeground       = tcell.ColorWhite.TrueColor()
	paneBackground         = tcell.ColorBlue.TrueColor()
	paneForeground         = tcell.ColorWhite.TrueColor()
	paneActiveBackground   = tcell.ColorBlueViolet.TrueColor()
	paneActiveForeground   = tcell.ColorWhite.TrueColor()
	paneFooterBackground   = tcell.ColorBlueViolet.TrueColor()
	paneFooterForeground   = tcell.ColorWhite.TrueColor()
)

var (
	TUI          *tview.Application
	TUIGrid      *tview.Grid
	TUIPages     *tview.Pages
	TUIHeader    *tview.Flex
	TUIFooter    *tview.Flex
	TUILeftPane  tview.Primitive
	TUIRightPane tview.Primitive
	TUIError     *tuiError
	menus        = make([]*menu, 0)
)

func AddToLeftPane(p tview.Primitive) {
	TUILeftPane = p
	if TUIRightPane == nil {
		addTextToFooter("NR")
		TUIGrid.AddItem(p, 1, 0, 1, 2, 0, 0, true)
	} else {
		addTextToFooter("HR")
		TUIGrid.AddItem(p, 1, 0, 1, 1, 0, 0, true)
		TUIGrid.AddItem(TUIRightPane, 1, 1, 1, 1, 0, 0, true)
	}
	TUI.SetFocus(TUILeftPane)
}

func AddToRightPane(p tview.Primitive) {
	TUIRightPane = p
	if TUILeftPane != nil {
		TUIGrid.RemoveItem(TUILeftPane)
		TUIGrid.AddItem(TUILeftPane, 1, 0, 1, 1, 0, 0, false)
	}
	TUIGrid.AddItem(p, 1, 1, 1, 1, 0, 0, true)
	TUI.SetFocus(TUIRightPane)
}

func RemoveFromRightPane() {
	TUIGrid.RemoveItem(TUIRightPane)
	TUIGrid.RemoveItem(TUIRightPane)
	TUIRightPane = nil
	TUIGrid.AddItem(TUILeftPane, 1, 0, 1, 2, 0, 0, true)
	TUI.SetFocus(TUILeftPane)
}

func RemoveFromLeftPane() {
	TUIGrid.RemoveItem(TUILeftPane)
	TUIGrid.RemoveItem(TUIRightPane)
	TUILeftPane = nil
	TUIGrid.AddItem(TUIRightPane, 1, 0, 1, 2, 0, 0, true)
	TUI.SetFocus(TUIRightPane)
}

func togglePaneFocus() {
	if TUILeftPane != nil {
		if TUILeftPane.HasFocus() {
			focusRightPane()
			return
		}
	}
	if TUIRightPane != nil {
		if TUIRightPane.HasFocus() {
			focusLeftPane()
			return
		}
	}
}

func focusLeftPane() {
	if TUILeftPane != nil {
		TUI.SetFocus(TUILeftPane)
	}
}

func focusRightPane() {
	if TUIRightPane != nil {
		TUI.SetFocus(TUIRightPane)
	}
}

func main() {
	TUI = tview.NewApplication()

	TUIGrid = tview.NewGrid()
	TUIPages = tview.NewPages()
	TUIHeader = tview.NewFlex().SetDirection(tview.FlexColumn)
	TUIFooter = tview.NewFlex().SetDirection(tview.FlexColumn)

	TUIError = makeErrorModal("mainErrorModal")

	table := makeTableV2()
	table2 := makex()

	menus = append(menus, makeMenu("Menu", "mainMenu",
		[]string{"table1", "table2", "table3", "table4", "table5"},
		func(index int, mainText, secondaryText string, shortcut rune) bool {
			if mainText == "table3" {
				TUIError.showError(nil, "table 1 is not implemented")
			}
			if mainText == "table1" {
				AddToLeftPane(table)
			}
			if mainText == "table2" {
				AddToRightPane(table2)
			}
			return true
		}))

	menus = append(menus, makeMenu("Menu2", "mainMenu2",
		[]string{"m2i1"},
		func(index int, mainText, secondaryText string, shortcut rune) bool {
			if mainText == "x" {
				TUIError.showError(nil, "table 1 is not implemented")
			}
			return true
		}))

	menus = append(menus, makeMenu("Menu3", "mainMenu3",
		[]string{"m3i1"},
		func(index int, mainText, secondaryText string, shortcut rune) bool {
			if mainText == "x" {
				TUIError.showError(nil, "table 1 is not implemented")
			}
			return true
		}))

	TUIGrid.SetRows(1, 0, 1).
		SetColumns(0).
		AddItem(TUIHeader, 0, 0, 1, 2, 0, 0, false).
		AddItem(TUIFooter, 2, 0, 1, 2, 2, 1, false)

	TUIPages.AddPage("grid", TUIGrid, true, true)
	TUIPages.SetBackgroundColor(tcell.ColorBlue.TrueColor())

	TUI.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		me := handleMenuInputs(event)
		if me != nil {
			return event
		}

		switch event.Rune() {
		case 'q':
			if TUILeftPane != nil {
				if TUILeftPane.HasFocus() {
					RemoveFromLeftPane()
					return nil
				}
			}
			if TUIRightPane != nil {
				if TUIRightPane.HasFocus() {
					RemoveFromRightPane()
					return nil
				}
			}
		}

		switch event.Key() {
		case tcell.KeyEsc:
			TUI.Stop()
			return nil
		case tcell.KeyTab:
			togglePaneFocus()
		}

		return event
	})

	TUI.SetRoot(TUIPages, true)
	// The grid causes off screen scrolling if not blurred
	TUIGrid.Blur()
	err := TUI.Run()
	if err != nil {
		panic(err)
	}
}
