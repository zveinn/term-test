package main

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	madmin "github.com/minio/madmin-go/v4"
	"github.com/rivo/tview"
)

type tuiTable struct {
	Flex        *tview.Flex
	Table       *tview.Table
	Footer      *tview.TextView
	TotalItem   int
	CurrentItem int
	Offset      int
	Limit       int
	Filter      string
}

type tableColumn struct {
	key    string
	format func(string) string
}

var defaultFormat = func(x string) string { return x }

func newTableColumn(key string, format func(string) string) (tc *tableColumn) {
	tc = new(tableColumn)
	tc.format = format
	tc.key = key
	if tc.format == nil {
		tc.format = defaultFormat
	}
	return
}

func makeTableV2() (flex *tview.Flex) {
	table := tview.NewTable()
	table.SetBackgroundColor(tcell.ColorNone)

	columns := make([]*tableColumn, 4)
	columns[0] = newTableColumn("Path", nil)
	columns[1] = newTableColumn("PoolIndex", nil)
	columns[2] = newTableColumn("State", nil)
	columns[3] = newTableColumn("NodeID", nil)

	for i, v := range columns {
		table.SetCell(0, i, tview.NewTableCell(v.key))
	}

	data := getData()
	for i, v := range data.Results {
		for ii, vv := range columns {
			table.SetCell(i+1, ii, tview.NewTableCell(fmt.Sprint(printStructField(v, vv.key))))
		}
	}

	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			table.SetSelectable(true, true)
		}
	}).SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).SetTextColor(tcell.ColorRed)
		table.SetSelectable(false, false)
	})

	footer := tview.NewTextView()
	footer.SetText("footer!!")
	footer.SetBackgroundColor(paneFooterBackground)
	footer.SetTextColor(paneFooterForeground)

	flex = tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(table, 0, 1, false)
	flex.AddItem(footer, 1, 0, false)

	// Hijack focus from flex and force focus on the table
	flex.SetFocusFunc(func() {
		TUI.SetFocus(table)
	})

	return flex
}

func getData() *madmin.PaginatedDrivesResponse {
	client, err := madmin.New("localhost:9001", "minioadmin", "minioadmin", false)
	if err != nil {
		fmt.Println("Error initializing client:", err)
		return nil
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	nodes, err := client.DrivesQuery(ctx, &madmin.DrivesResourceOpts{
		Limit:  9999,
		Offset: 0,
		Filter: "",
	})
	if err != nil {
		panic(err)
	}

	return nodes
}

func printStructField(item interface{}, fieldName string) interface{} {
	v := reflect.ValueOf(item)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	field := v.FieldByName(fieldName)

	if !field.IsValid() {
		return nil
	}

	return field.Interface()
}

func makex() *tview.Table {
	table := tview.NewTable().
		SetBorders(true)
	lorem := strings.Split("Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.", " ")
	cols, rows := 10, 40
	word := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
			if c < 1 || r < 1 {
				color = tcell.ColorYellow
			}
			table.SetCell(r, c,
				tview.NewTableCell(lorem[word]).
					SetTextColor(color).
					SetAlign(tview.AlignCenter))
			word = (word + 1) % len(lorem)
		}
	}
	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			table.SetSelectable(true, true)
		}
	}).SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).SetTextColor(tcell.ColorRed)
		table.SetSelectable(false, false)
	})
	return table
}
