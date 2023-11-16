package pgui

import (
	"bytes"
	"encoding/json"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/flags"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/httpconf"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

var (
	Action    *tview.Form
	UserState *tview.TextView
	NewCall   *tview.TextView
	app       *tview.Application
)

func refresh() {
	for v := range time.Tick(time.Second * 2) {
		//fmt.Fprintf(Action, "%s", v.String()+"\n")
		//Action.ScrollToEnd()

		UserState.SetText(v.String())
	}
}

func newTextView(title string) *tview.TextView {
	textView := tview.NewTextView()

	textView.SetTextAlign(tview.AlignLeft)
	textView.SetScrollable(true)
	textView.SetTitle(title).SetBorder(true)

	textView.SetChangedFunc(func() {
		app.Draw()
	})
	return textView
}

func newForm(title string, login string) *tview.Form {
	form := tview.NewForm()

	form.SetTitle(title).SetBorder(true)
	form.AddTextView("login", login, 10, 1, true, false)
	form.AddDropDown("Status", []string{"normal", "away", "dnd"}, 0, status)
	form.AddCheckbox("test", false, nil)
	form.AddButton("Answer", answer)
	form.AddButton("Cancel", cancel)
	form.MouseHandler()

	return form
}

func answer() {
	NewCall.SetText("Answer")
}

func cancel() {
	NewCall.SetText("Cansel")
}

func status(status string, index int) {
	buf := bytes.Buffer{}

	body := map[string]string{
		"ChangeUserState": status,
	}

	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(body); err != nil {
		log.Error().Err(err).Msg("Encode")
		return
	}

	url := flags.ServAddr + "/api/v1/control"

	req, _ := http.NewRequest(http.MethodPost, url, &buf)
	res, err := httpconf.HTTPClient.Do(req)
	if err != nil {
		log.Fatal().Err(err).Msg("status httpClient.Do")
	}
	defer res.Body.Close()

}

func Init() {
	// TODO добавить управление курсорами
	app = tview.NewApplication()

	header := newTextView("header")
	footer := newTextView("footer")
	UserState = newTextView("UserState")
	NewCall = newTextView("NewCall")
	Action = newForm("Action", "agent")

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

	if err := app.SetRoot(grid, true).SetFocus(grid).EnableMouse(true).Run(); err != nil {
		log.Fatal().Err(err).Msg("run PguiApp")
	}

}
