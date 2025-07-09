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
	TUIGrid.AddItem(p, 1, 0, 1, 1, 0, 0, true)
	TUILeftPane = p
}

func AddToRightPane(p tview.Primitive) {
	TUIGrid.AddItem(p, 1, 1, 1, 1, 0, 0, true)
	TUIRightPane = p
}

func main() {
	TUI = tview.NewApplication()

	TUIGrid = tview.NewGrid()
	TUIPages = tview.NewPages()
	TUIHeader = tview.NewFlex().SetDirection(tview.FlexColumn)
	TUIFooter = tview.NewFlex().SetDirection(tview.FlexColumn)
	addTextToFooter("Item1")
	addTextToFooter("item2")
	addTextToFooter("item3")

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
		SetColumns(0, 0).
		AddItem(TUIHeader, 0, 0, 1, 2, 0, 0, false).
		AddItem(TUIFooter, 2, 0, 1, 2, 2, 1, false)

	TUIPages.AddPage("grid", TUIGrid, true, true)
	TUIPages.SetBackgroundColor(tcell.ColorBlue.TrueColor())

	TUI.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		me := handleMenuInputs(event)
		if me != nil {
			return event
		}

		switch event.Key() {
		case tcell.KeyTab:
			if TUILeftPane.HasFocus() {
				TUI.SetFocus(TUIRightPane)
				TUILeftPane.Blur()
			} else {
				TUI.SetFocus(TUILeftPane)
				TUIRightPane.Blur()
			}
		}

		return event
	})

	TUI.SetRoot(TUIPages, true)
	TUIGrid.Blur()
	TUIPages.Blur()
	err := TUI.Run()
	if err != nil {
		panic(err)
	}
}
