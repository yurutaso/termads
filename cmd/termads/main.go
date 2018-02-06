package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/yurutaso/termads"
	"log"
)

const (
	statusExit     = -1
	statusContinue = 1
)

type Panel struct {
	x                 int
	y                 int
	width             int
	text              string
	cursor            *Cursor
	modifiable        bool
	storeOverflowText bool
	name              string
}

func NewPanelForm(x, y, width int, text string, name string) *Panel {
	return &Panel{x: x, y: y, width: width, modifiable: true, text: text, cursor: NewCursor(x, y), storeOverflowText: false, name: name}
}

func NewPanel(x, y int, text string) *Panel {
	return &Panel{x: x, y: y, width: len(text), modifiable: false, text: text, cursor: NewCursor(x, y), storeOverflowText: false, name: ""}
}

func (panel *Panel) DrawText() {
	if panel.width > len(panel.text) {
		drawLine(panel.x, panel.y, fmt.Sprintf(panel.text))
	} else {
		drawLine(panel.x, panel.y, fmt.Sprintf(panel.text[0:panel.width]))
	}
}

func (panel *Panel) DrawCursor() {
	if panel.x+panel.width < panel.cursor.x {
		termbox.SetCursor(panel.x+panel.width, panel.cursor.y)
	} else {
		termbox.SetCursor(panel.cursor.x, panel.cursor.y)
	}
}

func (panel *Panel) CheckCursorPosition() error {
	if panel.x > panel.cursor.x || panel.GetLastXpos() < panel.cursor.x {
		return fmt.Errorf("cursor x position out of range")
	}
	return nil
}

func (panel *Panel) InsertText(s string) error {
	err := panel.CheckCursorPosition()
	if err != nil {
		return err
	}
	if panel.modifiable {
		if !panel.storeOverflowText && len(panel.text) == panel.width {
			return nil
		}
		cx := panel.cursor.x - panel.x
		switch {
		case cx == 0:
			panel.text = s + panel.text
		case cx == len(panel.text):
			panel.text += s
		default:
			panel.text = panel.text[0:cx] + s + panel.text[cx:]
		}
		panel.MoveCursorRight()
	}
	return nil
}

func (panel *Panel) Backspace() error {
	if len(panel.text) == 0 {
		return nil
	}
	err := panel.CheckCursorPosition()
	if err != nil {
		return err
	}
	if panel.modifiable {
		cx := panel.cursor.x - panel.x
		lx := len(panel.text)
		switch {
		case cx == lx:
			panel.text = panel.text[0 : lx-1]
			panel.MoveCursorLeft()
		case cx == panel.x:
			break
		default:
			panel.text = panel.text[0:cx-1] + panel.text[cx:]
			panel.MoveCursorLeft()
		}
	}
	return nil
}

func (panel *Panel) RemoveText() {
	panel.text = ""
	panel.MoveCursorFirst()
}

func (panel *Panel) GetLastXpos() int {
	return panel.x + len(panel.text)
}

func (panel *Panel) MoveCursorFirst() {
	panel.cursor.Set(panel.x, panel.y)
}

func (panel *Panel) MoveCursorLast() {
	panel.cursor.Set(panel.GetLastXpos(), panel.y)
}

func (panel *Panel) MoveCursorLeft() {
	if panel.cursor.x > panel.x {
		panel.cursor.Left()
	}
}

func (panel *Panel) MoveCursorRight() {
	if panel.cursor.x < panel.GetLastXpos() {
		panel.cursor.Right()
	}
}

/* Window */
type Window struct {
	panels []*Panel
	papers []*termads.Paper
	data   map[string]string
	active int
}

func NewWindow(panels []*Panel) *Window {
	return &Window{panels: panels, data: map[string]string{}}
}

func (window *Window) ActivePanel() *Panel {
	return window.panels[window.active]
}

func (window *Window) RedrawAll() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for _, panel := range window.panels {
		panel.DrawText()
	}
	window.ActivePanel().DrawCursor()
	yoff := 10
	for i, paper := range window.papers {
		if i > 10 {
			break
		}
		drawLine(0, yoff+i, paper.GetBibcode()+" "+paper.AvailableLinkTypes()+" "+paper.GetTitle())
	}
	termbox.Flush()
}

func (window *Window) FocusNext() {
	if window.active == len(window.panels)-1 {
		window.active = 0
	} else {
		window.active += 1
	}
}

func (window *Window) FocusPrev() {
	if window.active == 0 {
		window.active = len(window.panels) - 1
	} else {
		window.active -= 1
	}
}

func (window *Window) FocusNextForm() error {
	for i := 0; i < len(window.panels); i++ {
		window.FocusNext()
		if window.ActivePanel().modifiable {
			return nil
		}
	}
	return fmt.Errorf("No modifiable window found")
}

func (window *Window) FocusPrevForm() error {
	for i := 0; i < len(window.panels); i++ {
		window.FocusPrev()
		if window.ActivePanel().modifiable {
			return nil
		}
	}
	return fmt.Errorf("No modifiable window found")
}

func (window *Window) GetForms() {
	window.active = 0
	window.FocusNextForm()
	i := window.active
	for key, _ := range window.data {
		window.data[key] = ""
	}
	for {
		key := window.ActivePanel().name
		val := window.ActivePanel().text
		window.data[key] = val
		window.FocusNextForm()
		if i == window.active {
			break
		}
	}
}

func (window *Window) GetPapersFromADS() error {
	form := termads.NewForm()
	for key, val := range window.data {
		form.Set(key, val)
	}
	papers, err := termads.GetLinks(form)
	window.papers = papers
	if err != nil {
		return err
	}
	return nil
}

func (window *Window) FormIsEmpty() bool {
	for _, val := range window.data {
		if len(val) > 0 {
			return false
		}
	}
	return true
}

func pollEvent() {
	panels := []*Panel{
		NewPanel(0, 0, "Input seach forms, then press <Enter> to get links from ADS."),
		NewPanel(0, 1, "Press <TAB>/<Ctrl-N> or <Ctrl-P> to move between forms."),
		NewPanel(0, 2, "------------------------------------------------------------"),
		// Row 1
		NewPanel(0, 3, "   Authors:"),
		NewPanelForm(12, 3, 100, "", "author"),
		// Row 2
		NewPanel(0, 4, "start year:"),
		NewPanelForm(12, 4, 4, "", "start_year"),
		NewPanel(18, 4, "month:"),
		NewPanelForm(25, 4, 2, "", "end_mon"),
		// Row 3
		NewPanel(0, 5, "  end year:"),
		NewPanelForm(12, 5, 4, "", "end_year"),
		NewPanel(18, 5, "month:"),
		NewPanelForm(25, 5, 2, "", "end_mon"),
		// Row 4
		NewPanel(0, 6, "     Title:"),
		NewPanelForm(12, 6, 100, "", "title"),
		// Row 5
		NewPanel(0, 7, "  Abstract:"),
		NewPanelForm(12, 7, 100, "", "text"),
		// Row 6
		NewPanel(0, 8, "------------------------------------------------------------"),
	}

	window := NewWindow(panels)
	window.FocusNextForm()
	window.RedrawAll()
	for {
		status := window.HandleKeyEventInForm(termbox.PollEvent())
		if status == statusExit {
			return
		}
		window.RedrawAll()
	}
}

/*
func (window *Window) HandleKeyEventInResult(ev termbox.Event) int {
	switch ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		// Terminate
		case termbox.KeyCtrlC:
			return statusExit
		// Send forms to ADS
		case termbox.KeyEnter:
			return statusExit
		// Motions
		case termbox.KeyArrowDown:
			window.NextResult()
		case termbox.KeyArrowUp:
			window.PrevResult()
		case termbox.KeyArrowRight:
			window.MoveCursorRight()
		case termbox.KeyArrowLeft:
			window.MoveCursorLeft()
		case termbox.KeySpace:
			window.SelectResult()
		default:
			if ev.Ch != 0 {
				s = string(ev.Ch)
				switch s {
				case 'h':
					window.MoveCursorLeft()
				case 'j':
					window.NextResult()
				case 'k':
					window.PrevResult()
				case 'l':
					window.MoveCursorRight()
				case 'g':
					window.FirstResult()
				case 'G':
					window.LastResult()
				}
			}
		}
	}
	return statusContinue
}
*/

func (window *Window) HandleKeyEventInForm(ev termbox.Event) int {
	switch ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		// Terminate
		case termbox.KeyEsc, termbox.KeyCtrlC:
			return statusExit
		// Send forms to ADS
		case termbox.KeyEnter:
			window.GetForms()
			err := window.GetPapersFromADS()
			if err != nil {
				return statusExit
			}
		// Motions
		case termbox.KeyTab, termbox.KeyCtrlN, termbox.KeyArrowDown:
			window.FocusNextForm()
		case termbox.KeyCtrlP, termbox.KeyArrowUp:
			window.FocusPrevForm()
		case termbox.KeyCtrlF, termbox.KeyArrowRight:
			window.ActivePanel().MoveCursorRight()
		case termbox.KeyCtrlB, termbox.KeyArrowLeft:
			window.ActivePanel().MoveCursorLeft()
		case termbox.KeyCtrlE:
			window.ActivePanel().MoveCursorLast()
		case termbox.KeyCtrlA:
			window.ActivePanel().MoveCursorFirst()
		// Edit text
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			window.ActivePanel().Backspace()
		case termbox.KeySpace:
			window.ActivePanel().InsertText(" ")
		case termbox.KeyCtrlU:
			window.ActivePanel().RemoveText()
		default:
			if ev.Ch != 0 {
				s := string(ev.Ch)
				window.ActivePanel().InsertText(s)
			}
		}
	}
	return statusContinue
}

func main() {
	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()
	pollEvent()
}

func drawLine(x, y int, str string) {
	color := termbox.ColorDefault
	bgrcolor := termbox.ColorDefault
	runes := []rune(str)
	for i := 0; i < len(runes); i++ {
		termbox.SetCell(x+i, y, runes[i], color, bgrcolor)
	}
}
