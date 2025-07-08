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
	header := newPrimitive("header")
	footer := newPrimitive("footer")

	grid := tview.NewGrid().
		SetRows(1, 0, 2).
		SetColumns(0, 0).
		SetBorders(true).
		AddItem(footer, 2, 0, 1, 2, 0, 0, false).
		AddItem(header, 0, 0, 1, 2, 0, 0, false)

	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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
