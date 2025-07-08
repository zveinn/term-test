package main

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	TUI          *tview.Application
	TUIGrid      *tview.Grid
	TUIPages     *tview.Pages
	TUIHeader    *tview.Flex
	TUIFooter    *tview.Flex
	TUILeftPane  *tview.TextView
	TUIRightPane *tview.TextView
	TUIError     *tuiError
	menus        = make([]*menu, 10)
)

func main() {
	TUI = tview.NewApplication()
	TUILeftPane = tview.NewTextView().SetText("left..")
	TUIRightPane = tview.NewTextView().SetText("right..")
	TUIGrid = tview.NewGrid()
	TUIPages = tview.NewPages()
	TUIHeader = tview.NewFlex().SetDirection(tview.FlexColumn)
	TUIFooter = tview.NewFlex().SetDirection(tview.FlexColumn)
	TUIError = makeErrorModal("mainErrorModal")

	menus[0] = makeMenu("Menu", "mainMenu", '1',
		[]string{"table1", "table2", "table3", "table4", "table5"},
		func(index int, mainText, secondaryText string, shortcut rune) bool {
			if mainText == "table1" {
				TUIError.showError(nil, "table 1 is not implemented")
			}
			return true
		})

	menus[1] = makeMenu("Menu2", "mainMenu2", '2',
		[]string{"m2i1"},
		func(index int, mainText, secondaryText string, shortcut rune) bool {
			if mainText == "table1" {
				TUIError.showError(nil, "table 1 is not implemented")
			}
			return true
		})

	menus[2] = makeMenu("Menu3", "mainMenu3", '3',
		[]string{"m3i1"},
		func(index int, mainText, secondaryText string, shortcut rune) bool {
			if mainText == "table1" {
				TUIError.showError(nil, "table 1 is not implemented")
			}
			return true
		})

	TUIGrid.SetRows(1, 0, 1).
		SetColumns(0, 0).
		SetBorders(true).
		AddItem(TUIFooter, 2, 0, 1, 2, 0, 0, false).
		AddItem(TUIHeader, 0, 0, 1, 2, 0, 0, false)

	TUIPages.AddPage("grid", TUIGrid, true, true)

	TUIGrid.AddItem(TUILeftPane, 1, 0, 1, 1, 0, 0, false)
	TUIGrid.AddItem(TUIRightPane, 1, 1, 1, 1, 0, 0, false)

	TUIGrid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if handleMenuInputs(event.Rune()) {
			return event
		}
		switch event.Rune() {
		case 'j', 'k', 'l', 'h':
			hasOpenMenu := false
			for _, v := range menus {
				if v == nil {
					continue
				}
				if v.isOpen {
					hasOpenMenu = true
				}
			}
			if !hasOpenMenu {
				TUILeftPane.SetText("no open menus")
				return nil
			}
		case 'q', 'Q':
			// Allow quitting with 'q'
			TUI.Stop()
			return nil
		}

		switch event.Key() {
		case tcell.KeyPgUp, tcell.KeyPgDn, tcell.KeyDown, tcell.KeyUp, tcell.KeyLeft, tcell.KeyRight:
			hasOpenMenu := false
			for _, v := range menus {
				if v == nil {
					continue
				}
				if v.isOpen {
					hasOpenMenu = true
				}
			}
			if !hasOpenMenu {
				TUILeftPane.SetText("no open menus")
				return nil
			}
		}

		return event
	})

	TUI.SetRoot(TUIPages, true).SetFocus(TUIGrid)
	err := TUI.Run()
	if err != nil {
		panic(err)
	}
}

type menu struct {
	isOpen    bool
	labelText string
	list      *tview.List
	label     *tview.TextView
	pageLabel string
	shortcut  rune
}

func (m *menu) close() {
	TUIPages.HidePage(m.getPageListLabel())
	TUIPages.SendToBack(m.getPageListLabel())
	TUI = TUI.SetFocus(TUIGrid)
	m.isOpen = false
}

func (m *menu) open() {
	x, y, width, height := m.label.GetRect()
	_, _, mw, mh := m.list.GetRect()
	m.list.SetRect(x, y+1, mw, mh)
	TUILeftPane.SetText(fmt.Sprintf("%s %d %d %d %d", m.labelText, x, y, width, height))
	m.isOpen = true
	TUIPages.ShowPage(m.getPageListLabel())
	TUIPages.SendToFront(m.getPageListLabel())
	TUI = TUI.SetFocus(m.list)
}

func (m *menu) toggleMenu() {
	if !m.isOpen {
		m.open()
	} else {
		m.close()
	}
}

func (m *menu) getPageListLabel() string {
	return m.pageLabel + "_list"
}

func (m *menu) getPageLabel() string {
	return m.pageLabel
}

func makeMenu(label string, pageLabel string, shortcut rune, options []string, selectFunc func(index int, mainText, secondaryText string, shortcut rune) bool) (m *menu) {
	m = new(menu)
	m.pageLabel = pageLabel
	m.shortcut = shortcut
	m.labelText = label
	m.label = tview.NewTextView().
		SetText(label).
		SetTextColor(tcell.ColorWhite)

	m.list = tview.NewList()
	m.list.SetUseStyleTags(false, false)
	itemStyle := tcell.StyleDefault
	itemStyle = itemStyle.Background(tcell.ColorBlue.TrueColor())
	itemStyle = itemStyle.Foreground(tcell.ColorWhite)
	m.list.SetMainTextStyle(itemStyle)
	m.list.SetBackgroundColor(tcell.ColorBlue.TrueColor())
	m.list.SetBorderPadding(1, 1, 1, 1)
	sstyle := tcell.StyleDefault
	sstyle = sstyle.Background(tcell.ColorNone)
	sstyle = sstyle.Foreground(tcell.ColorBisque)
	m.list.SetSecondaryTextStyle(sstyle)
	// m.list.SetBorder(true)
	// bStyle := tcell.StyleDefault
	// bStyle = bStyle.Background(tcell.ColorBlue.TrueColor())
	// m.list.SetBorderStyle(bStyle)
	// m.list.SetBorderStyle(bStyle)
	m.list.ShowSecondaryText(false)
	lw := 0
	lh := 0
	for _, option := range options {
		if lw < len(option)+2 {
			lw = len(option) + 2
		}
		lh++
		m.list.AddItem(option, "", 0, nil)
	}

	m.list.SetRect(1, 1, lw, lh+2)

	m.list.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		ok := selectFunc(index, mainText, secondaryText, shortcut)
		if ok {
			TUIPages.HidePage(m.getPageListLabel())
			TUIPages.SendToBack(m.getPageListLabel())
			TUI = TUI.SetFocus(TUIGrid)
			m.isOpen = false
		}
	})

	m.list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if handleMenuInputs(event.Rune()) {
			return event
		}

		if event.Key() == tcell.KeyEscape {
			TUIPages.HidePage(m.getPageListLabel())
			TUIPages.SendToBack(m.getPageListLabel())
			TUI = TUI.SetFocus(TUIGrid)
			m.isOpen = false
			return nil
		}
		return event
	})

	TUIHeader.AddItem(m.label, len(label)+1, 1, false)
	TUIPages.AddPage(m.getPageListLabel(), m.list, false, false)
	return
}

type tuiError struct {
	originalError error
	customMsg     string
	view          *tview.TextView
	pageLabel     string
}

func (t *tuiError) getPageLabel() string {
	return t.pageLabel
}

func (t *tuiError) showError(_ error, msg string) {
	t.view.SetText(msg)
	TUIPages.ShowPage(t.getPageLabel())
	TUIPages.SendToFront(t.getPageLabel())
	TUI = TUI.SetFocus(t.view)
}

func makeErrorModal(pageLabel string) (e *tuiError) {
	e = new(tuiError)
	e.pageLabel = pageLabel
	e.view = tview.NewTextView()
	e.view.SetText("").
		SetTextColor(tcell.ColorRed).
		SetBorder(true).
		SetBorderColor(tcell.ColorDarkRed)
	e.view.SetBackgroundColor(tcell.ColorBlack)

	e.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyEnter {
			TUIPages.HidePage(e.getPageLabel())
			TUIPages.SendToBack(e.getPageLabel())
			TUI = TUI.SetFocus(TUIGrid)
			return nil
		}
		return event
	})

	TUIPages.AddPage(e.getPageLabel(), e.view, false, false)
	return
}

func handleMenuInputs(key rune) (wasMenuTrigger bool) {
	switch key {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		for i := range menus {
			if menus[i] == nil {
				continue
			}
			menus[i].close()
		}
		TUIRightPane.SetText(string(key))
		ri, err := strconv.Atoi(string(key))
		if err != nil {
			TUIError.showError(err, "menu item not initialized")
			return true
		}
		ri -= 1
		if menus[ri] == nil {
			TUIError.showError(err, "menu item not initialized")
			return true
		}
		menus[ri].toggleMenu()
	default:
	}

	return false
}
