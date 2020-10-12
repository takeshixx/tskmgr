package ui

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func (ui *UI) newTaskList() {
	ui.TaskList.SetBorder(true).SetTitle("Task List")
	ui.TaskList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		ui.displayTaskInfo(ui.Manager.Tasks[index])
		ui.displayTaskFiles(ui.Manager.Tasks[index])
	})
	ui.TaskList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Disable keys in list
		switch key := event.Key(); key {
		case tcell.KeyTab, tcell.KeyF5:
			return nil
		}
		return event
	})
}

func (ui *UI) newTaskInfo() {
	ui.TaskInfo.SetDynamicColors(true)
	ui.TaskInfo.SetWordWrap(true)
	ui.TaskInfo.SetBorder(true).SetTitle("Task Info")
	ui.TaskInfo.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			curTask, err := ui.getCurrentTask()
			if err != nil {
				return event
			}
			ui.App.Suspend(func() {
				ui.Manager.OpenTaskInfo(curTask)
				ui.Manager.ReloadTask(curTask)
			})
		}
		return event
	})
}

func (ui *UI) newTaskFiles() {
	ui.TaskFiles.SetBorder(true).SetTitle("Task Files")
	ui.TaskFiles.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == ' ' {
			curTask, err := ui.getCurrentTask()
			if err != nil {
				return event
			}
			ui.App.Suspend(func() {
				ui.Manager.OpenTaskFiles(curTask)
				ui.Manager.ReloadTask(curTask)
			})
		}
		return event
	})
}

func (ui *UI) newCreateForm() {
	ui.CreateForm.
		AddInputField("Title", "", 0, nil, nil).
		AddInputField("Description", "", 0, nil, nil).
		AddButton("Create", func() {
			// create a task
			newFormTitle := ui.CreateForm.GetFormItemByLabel("Title")
			newFormTitleInput := newFormTitle.(*tview.InputField)
			newTaskTitle := newFormTitleInput.GetText()

			newFormDesc := ui.CreateForm.GetFormItemByLabel("Description")
			newFormDescInput := newFormDesc.(*tview.InputField)
			newTaskDesc := newFormDescInput.GetText()

			newTask, err := ui.Manager.CreateTask(newTaskTitle)
			if err != nil {
				ui.ShowError(err)
				return
			}
			newTask.Description = newTaskDesc
			_, err = ui.Manager.SaveTask(newTask)
			if err != nil {
				ui.ShowError(err)
				return
			}

			ui.updateTaskList()
			ui.Pages.SwitchToPage("Overview")
			ui.App.SetFocus(ui.TaskList)
		}).
		AddButton("Cancel", func() {
			ui.Pages.SwitchToPage("Overview")
			ui.App.SetFocus(ui.TaskList)
		})
	ui.Pages.AddPage("Create", ui.CreateForm, true, false)
}

func (ui *UI) newDeleteModal() {
	ui.DeleteModal.
		SetText("Do you really want to delete this task?").
		AddButtons([]string{"Delete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonIndex == 0 {
				curTaskIndex := ui.TaskList.GetCurrentItem()
				if curTaskIndex-1 > len(ui.Manager.Tasks) {
					ui.ShowError(fmt.Errorf("Invalid task index"))
					return
				}
				curTask := ui.Manager.Tasks[curTaskIndex]
				err := ui.Manager.DeleteTask(curTask)
				if err != nil {
					ui.ShowError(err)
					return
				}
				ui.updateTaskList()
			}
			ui.Pages.SwitchToPage("Overview")
			ui.App.SetFocus(ui.TaskList)
		})
	ui.Pages.AddPage("Delete", ui.DeleteModal, true, false)
}

func (ui *UI) newErrorModal() {
	ui.ErrorModal.
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			ui.Pages.SwitchToPage("Overview")
			ui.App.SetFocus(ui.TaskList)
			ui.ErrorModal.SetText("")
		})
	ui.ErrorModal.SetTitle("Error")
	ui.Pages.AddPage("Error", ui.ErrorModal, true, false)
}

func (ui *UI) newMenu() {
	ui.Menu.AddItem("Overview", tcell.KeyF1, func() {
		ui.Pages.SwitchToPage("Overview")
		ui.App.SetFocus(ui.TaskList)
	})
	ui.Menu.AddItem("Create", tcell.KeyF2, func() {
		ui.Pages.SwitchToPage("Create")
		ui.App.SetFocus(ui.CreateForm)
	})
	ui.Menu.AddItem("Delete", tcell.KeyF3, func() {
		ui.Pages.SwitchToPage("Delete")
		ui.App.SetFocus(ui.DeleteModal)
	})
	ui.Menu.AddItem("Edit", tcell.KeyF4, func() {
		curTaskIndex := ui.TaskList.GetCurrentItem()
		if curTaskIndex-1 > len(ui.Manager.Tasks) {
			ui.ShowError(fmt.Errorf("Invalid task index"))
			return
		}
		curTask := ui.Manager.Tasks[curTaskIndex]
		ui.App.Suspend(func() {
			ui.Manager.OpenTaskConfig(curTask)
			ui.updateTaskList()
		})
	})
	ui.Menu.AddItem("Refresh", tcell.KeyF5, func() {
		ui.updateUI()
	})
	ui.Menu.AddItem("Ranger", tcell.KeyF9, func() {
		curTask, err := ui.getCurrentTask()
		if err != nil {
			ui.ShowError(err)
			return
		}
		ui.App.Suspend(func() {
			if err = ui.Manager.OpenTaskFiles(curTask); err != nil {
				ui.ShowError(err)
				return
			}
		})
	})
	ui.Menu.AddItem("VSCode", tcell.KeyF10, func() {
		curTask, err := ui.getCurrentTask()
		if err != nil {
			ui.ShowError(err)
			return
		}
		if err = ui.Manager.OpenTaskVSCode(curTask); err != nil {
			ui.ShowError(err)
			return
		}

	})
	ui.OverviewLayout.AddItem(ui.Menu, 2, 0, 1, 3, 0, 0, false)
}
