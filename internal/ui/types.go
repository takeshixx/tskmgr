package ui

import (
	"github.com/rivo/tview"
	"github.com/takeshixx/tskmgr/internal/manager"
)

// UI holds pointers to all relevant
// UI elements.
type UI struct {
	App            *tview.Application
	Manager        *manager.Manager
	OverviewLayout *tview.Grid
	Pages          *tview.Pages
	TaskLayout     *tview.Flex
	TaskList       *tview.List
	TaskInfo       *tview.TextView
	TaskFiles      *tview.TreeView
	CreateForm     *tview.Form
	DeleteModal    *tview.Modal
	ErrorModal     *tview.Modal
	Menu           *Menu
}

// NewUI creates a new UI object.
func NewUI() (ui *UI, err error) {
	ui = &UI{
		App:            tview.NewApplication(),
		OverviewLayout: tview.NewGrid(),
		Pages:          tview.NewPages(),
		TaskList:       tview.NewList(),
		TaskInfo:       tview.NewTextView(),
		TaskFiles:      tview.NewTreeView(),
		TaskLayout:     tview.NewFlex(),
		CreateForm:     tview.NewForm(),
		DeleteModal:    tview.NewModal(),
		ErrorModal:     tview.NewModal(),
		Menu:           NewMenu(),
	}
	ui.Manager, err = manager.NewManager()
	if err != nil {
		return
	}

	ui.newTaskList()
	ui.newTaskInfo()

	ui.OverviewLayout = tview.NewGrid()
	ui.OverviewLayout.SetRows(1, 0, 1)
	ui.OverviewLayout.SetColumns(25, 0, 25)
	ui.OverviewLayout.AddItem(ui.TaskList, 1, 0, 0, 0, 0, 0, true)
	ui.OverviewLayout.AddItem(ui.TaskInfo, 1, 1, 0, 3, 0, 0, false)
	ui.OverviewLayout.AddItem(ui.TaskFiles, 1, 2, 0, 0, 0, 0, false)
	ui.OverviewLayout.AddItem(ui.TaskList, 1, 0, 1, 1, 0, 100, true)
	ui.OverviewLayout.AddItem(ui.TaskInfo, 1, 1, 1, 1, 0, 100, false)
	ui.OverviewLayout.AddItem(ui.TaskFiles, 1, 2, 1, 1, 0, 100, false)
	ui.Pages.AddPage("Overview", ui.OverviewLayout, true, true)

	ui.newCreateForm()
	ui.newDeleteModal()
	ui.newErrorModal()
	ui.newMenu()

	ui.App.SetRoot(ui.Pages, true)
	ui.App.SetFocus(ui.TaskList)
	ui.registerKeyEvents()
	ui.updateUI()
	return
}

// RunUI is a wrapper that handles all the
// initialization tasks for the UI.
func (ui *UI) RunUI() (err error) {
	return ui.App.Run()
}
