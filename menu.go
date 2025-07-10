package main

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type menu struct {
	isOpen    bool
	labelText string
	list      *tview.List
	label     *tview.TextView
	pageLabel string
}

func (m *menu) close() {
	TUIPages.HidePage(m.getPageListLabel())
	TUIPages.SendToBack(m.getPageListLabel())
	TUI.SetFocus(TUIGrid)
	m.isOpen = false
}

func (m *menu) open() {
	x, y, _, _ := m.label.GetRect()
	_, _, mw, mh := m.list.GetRect()
	m.list.SetRect(x, y+1, mw, mh)
	// TUILeftPane.SetText(fmt.Sprintf("%s %d %d %d %d", m.labelText, x, y, width, height))
	m.isOpen = true
	TUIPages.ShowPage(m.getPageListLabel())
	TUIPages.SendToFront(m.getPageListLabel())
	TUI.SetFocus(m.list)
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

func makeMenu(label string, pageLabel string, options []string, selectFunc func(index int, mainText, secondaryText string, shortcut rune) bool) (m *menu) {
	count := len(menus)
	m = new(menu)
	m.pageLabel = pageLabel
	m.labelText = label
	m.label = tview.NewTextView().
		SetText(strconv.Itoa(count+1) + "|" + label)

	m.label.SetTextColor(menuForeground)
	m.label.SetBackgroundColor(menuBackground)

	m.list = tview.NewList()
	m.list.SetUseStyleTags(false, false)
	m.list.SetBorderColor(menuListBackground)
	m.list.SetBackgroundColor(menuListBackground)
	itemStyle := tcell.StyleDefault
	itemStyle = itemStyle.Background(menuListBackground)
	itemStyle = itemStyle.Foreground(menuListForeground)
	m.list.SetMainTextStyle(itemStyle)
	// m.list.SetBorderPadding(1, 1, 1, 1)

	sstyle := tcell.StyleDefault
	sstyle = sstyle.Foreground(menuSelectedForeground)
	sstyle = sstyle.Background(menuSelectedBackground)
	m.list.SetSelectedStyle(sstyle)
	m.list.ShowSecondaryText(false)
	lw := 0
	lh := 0
	for _, option := range options {
		if lw < len(option) {
			lw = len(option)
		}
		lh++
		m.list.AddItem(option, "", 0, nil)
	}
	m.list.SetRect(1, 1, lw, lh)

	m.list.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		ok := selectFunc(index, mainText, secondaryText, shortcut)
		if ok {
			TUIPages.HidePage(m.getPageListLabel())
			TUIPages.SendToBack(m.getPageListLabel())
			m.isOpen = false
		}
	})

	// enable wim style navigation in menus
	m.list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, 0, 0)
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, 0, 0)
		}
		return event
	})

	// This is a hack to resize the last item to fill the flex box
	ic := TUIHeader.GetItemCount()
	if ic > 0 {
		lastItem := TUIHeader.GetItem(ic - 1)
		itf, ok := lastItem.(*tview.TextView)
		if !ok {
			panic("non-text-view-in-menu")
		}
		text := itf.GetText(true)
		TUIHeader.ResizeItem(lastItem, len(text)+2, 1)
	}
	TUIHeader.AddItem(m.label, 0, 1, false)
	TUIPages.AddPage(m.getPageListLabel(), m.list, false, false)
	return
}

func handleMenuInputs(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		for i := range menus {
			if menus[i] == nil {
				continue
			}
			menus[i].close()
		}
		ri, err := strconv.Atoi(string(event.Rune()))
		if err != nil {
			TUIError.showError(err, "menu item not initialized")
			return nil
		}
		ri -= 1
		if menus[ri] == nil {
			TUIError.showError(err, "menu item not initialized")
			return nil
		}
		menus[ri].toggleMenu()
	default:
	}

	switch event.Key() {
	case tcell.KeyPgUp, tcell.KeyPgDn, tcell.KeyDown, tcell.KeyUp:
		return event
	case tcell.KeyEscape:
		for i := range menus {
			if menus[i] != nil {
				menus[i].close()
			}
		}
	}

	return nil
}
