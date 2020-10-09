package ui

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/takeshixx/tskmgr/internal/tasks"
)

func (ui *UI) registerKeyEvents() (err error) {
	ui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// First trigger shortcuts for menu items
		for _, i := range ui.Menu.Items {
			if event.Key() == i.Shortcut {
				if i.Selected != nil {
					i.Selected()
				}
				break
			}
		}
		switch key := event.Key(); key {
		case tcell.KeyCtrlQ:
			ui.App.Stop()
		case tcell.KeyTab:
			if ui.OverviewLayout.HasFocus() {
				if ui.TaskList.HasFocus() {
					ui.App.SetFocus(ui.TaskInfo)
				} else if ui.TaskInfo.HasFocus() {
					ui.App.SetFocus(ui.TaskFiles)
				} else if ui.TaskFiles.HasFocus() {
					ui.App.SetFocus(ui.TaskList)
				}
			}
		}
		return event
	})
	return
}

func (ui *UI) getCurrentTask() (t *tasks.Task, err error) {
	curTaskIndex := ui.TaskList.GetCurrentItem()
	if curTaskIndex-1 > len(ui.Manager.Tasks) {
		err = fmt.Errorf("Invalid task index")
	}
	t = ui.Manager.Tasks[curTaskIndex]
	return
}

func (ui *UI) updateUI() (err error) {
	err = ui.updateTaskList()
	if err != nil {
		return
	}
	return
}

func (ui *UI) updateTaskList() (err error) {
	ui.TaskList.Clear()
	if len(ui.Manager.Tasks) > 0 {
		for i, task := range ui.Manager.Tasks {
			tt := task
			if tt == nil {
				continue
			}
			ui.TaskList.AddItem(tt.Title, strings.SplitN(tt.Description, "\n", 2)[0], rune(97+i), func() {
				ui.App.Suspend(func() {
					ui.Manager.OpenTaskInfo(tt)
					ui.Manager.ReloadTask(tt)
				})
			})
		}
	}
	return
}

func (ui *UI) displayTaskInfo(t *tasks.Task) (err error) {
	// If the Task description has been rendered already
	uiText := t.GetUIText()
	if uiText != "" {
		ui.TaskInfo.SetText(uiText)
		ui.TaskInfo.ScrollToBeginning()
		return
	}
	go func() {
		t.RenderUIText(ui.Manager.Config.TasksPath + string(os.PathSeparator) + t.Directory + string(os.PathSeparator) + "task_notes.md")
		ui.TaskInfo.SetText(t.GetUIText())
		ui.TaskInfo.ScrollToBeginning()
	}()
	ui.TaskInfo.SetDoneFunc(func(key tcell.Key) {
		ui.TaskInfo.Highlight()
	})
	return
}

func (ui *UI) displayTaskFiles(t *tasks.Task) (err error) {
	taskRoot := tview.NewTreeNode(ui.Manager.Config.TasksPath + string(os.PathSeparator) + t.Directory)
	ui.TaskFiles.
		SetRoot(taskRoot).
		SetCurrentNode(taskRoot)

	add := func(target *tview.TreeNode, path string) {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			node := tview.NewTreeNode(file.Name()).
				SetReference(filepath.Join(path, file.Name())).
				SetSelectable(file.IsDir())
			if file.IsDir() {
				node.SetColor(tcell.ColorGreen)
			}
			target.AddChild(node)
		}
	}
	add(taskRoot, ui.Manager.Config.TasksPath+string(os.PathSeparator)+t.Directory)

	ui.TaskFiles.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			// Load and show files in this directory.
			path := reference.(string)
			add(node, path)
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})
	return
}
