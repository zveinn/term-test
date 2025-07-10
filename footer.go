package main

import (
	"fmt"

	"github.com/rivo/tview"
)

func addTextToFooter(text ...any) {
	footerText := tview.NewTextView()
	footerText.SetText(fmt.Sprint(text...))
	footerText.SetBackgroundColor(footerBackground)
	footerText.SetTextColor(footerForeground)
	ic := TUIFooter.GetItemCount()
	if ic > 0 {
		lastItem := TUIFooter.GetItem(ic - 1)
		itf, ok := lastItem.(*tview.TextView)
		if !ok {
			panic("non-text-view-in-footer")
		}
		text := itf.GetText(true)
		TUIFooter.ResizeItem(lastItem, len(text)+2, 1)
	}
	TUIFooter.AddItem(footerText, 0, 1, false)
}
