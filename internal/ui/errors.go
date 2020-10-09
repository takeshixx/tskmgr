package ui

// ShowError displays an error in the ErrorModal.
func (ui *UI) ShowError(err error) {
	ui.ErrorModal.SetText(err.Error())
	ui.Pages.SwitchToPage("Error")
}
