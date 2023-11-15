package pgui

import (
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
	"time"
)

var (
	Action    *tview.TextView
	UserState *tview.TextView
	NewCall   *tview.TextView
	app       *tview.Application
)

func refresh() {
	for v := range time.Tick(time.Second * 2) {
		//fmt.Fprintf(Action, "%s", v.String()+"\n")
		//Action.ScrollToEnd()

		Action.SetText(v.String())
		UserState.SetText(v.String())
	}
}

func newTextView(title string) *tview.TextView {
	textView := tview.NewTextView()

	textView.SetTextAlign(tview.AlignLeft)
	textView.SetScrollable(true)
	textView.SetBorder(true)
	textView.SetTitle(title)

	textView.SetChangedFunc(func() {
		app.Draw()
	})
	return textView
}

func Init() {
	app = tview.NewApplication()

	header := newTextView("header")
	footer := newTextView("footer")
	Action = newTextView("Action")
	UserState = newTextView("UserState")
	NewCall = newTextView("NewCall")

	grid := tview.NewGrid().
		SetRows(1, 0, 0, 1).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(header, 0, 0, 1, 3, 0, 0, false).
		AddItem(footer, 3, 0, 1, 3, 0, 0, false)

	//Layout for screens narrower than 100 cells (menu and side bar are hidden).
	//grid.AddItem(actions, 0, 0, 0, 0, 0, 0, true).
	//	AddItem(main, 1, 0, 1, 3, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(Action, 1, 0, 2, 1, 1, 1, true).
		AddItem(UserState, 1, 1, 1, 2, 1, 1, false).
		AddItem(NewCall, 2, 1, 1, 2, 1, 1, false)

	//go refresh()

	if err := app.SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		log.Fatal().Err(err).Msg("run PguiApp")
	}

}
