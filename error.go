package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type tuiError struct {
	view      *tview.TextView
	pageLabel string
}

func (t *tuiError) getPageLabel() string {
	return t.pageLabel
}

func (t *tuiError) showError(err error, msg string) {
	if err != nil {
		t.view.SetText(fmt.Sprintf("%s \n Error: %s", msg, err))
	} else {
		t.view.SetText(msg)
	}
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
