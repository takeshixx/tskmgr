package ui

import (
	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

type menuItem struct {
	Text     string
	Shortcut tcell.Key
	Selected func()
}

type Menu struct {
	*tview.Box
	Items []menuItem
}

func NewMenu() *Menu {
	m := &Menu{
		Box: tview.NewBox(),
	}
	return m
}

func (m *Menu) Draw(screen tcell.Screen) {
	x, y, width, _ := m.GetInnerRect()
	x++
	for _, i := range m.Items {
		tview.Print(screen, tcell.KeyNames[i.Shortcut], x, y, width-2, tview.AlignLeft, tcell.GetColor("yellow"))
		x += 3

		tview.Print(screen, i.Text, x, y, width, tview.AlignLeft, tcell.GetColor("grey"))
		x += runewidth.StringWidth(i.Text) + 2
	}
}

func (m *Menu) AddItem(text string, shortcut tcell.Key, selected func()) *Menu {
	m.Items = append(m.Items, menuItem{
		Text:     text,
		Shortcut: shortcut,
		Selected: selected,
	})
	return m
}

func (m *Menu) Clear() *Menu {
	m.Items = nil
	return m
}
