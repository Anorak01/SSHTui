package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"unicode"
	"unicode/utf8"
)

type newConnectionViewModel struct {
	String      string
	returnModel model
}

func newConnectionView(returnTo model) newConnectionViewModel {
	return newConnectionViewModel{returnModel: returnTo}
}

func (m newConnectionViewModel) Init() tea.Cmd {
	// Initialize the new connection view
	// This could include setting up any necessary state or UI elements
	return nil
}

func (m newConnectionViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			// Handle Ctrl+C to exit the application
			return m, tea.Quit
		case tea.KeyEsc:
			// Handle the Esc to go back to the main view
			return m.returnModel.Update(msg)
		case tea.KeyEnter:
			// Handle the Enter key to save the connection name
		case tea.KeyBackspace:
			// Handle the Backspace key to remove the last character
			if len(m.String) > 0 {
				m.String = m.String[:len(m.String)-1]
			}
		default:
			if utf8.RuneCountInString(msg.String()) == 1 && unicode.IsPrint([]rune(msg.String())[0]) {
				m.String += msg.String()
			}
		}
	}
	return m, nil
}

func (m newConnectionViewModel) View() string {
	// Render the new connection view
	s := "Enter Connection Name: "
	s += m.String + "\n"
	return s
}
