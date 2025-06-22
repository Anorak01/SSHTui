package views

import (
	"anorak01.top/sshtui/views/filepicker"
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"strconv"
	"time"
	"unicode"
	"unicode/utf8"
)

type state int

const (
	StateName state = iota
	StateHost
	StatePort
	StateUser
	StateIfKey
	StateKey
	StateDone
)

var stateNames = map[state]string{
	StateName:  "Name",
	StateHost:  "Host",
	StatePort:  "Port",
	StateUser:  "User",
	StateIfKey: "Key? (y/N)",
	StateKey:   "Private Key",
	StateDone:  "Done",
}

type newConnectionViewModel struct {
	String       string
	State        state
	SshEntry     sshEntry
	filePicker   filepicker.Model
	selectedFile string
	err          error
	returnModel  *model
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func newConnectionView(returnTo *model, initialHeight int) newConnectionViewModel {
	fp := filepicker.New()
	fp.CurrentDirectory, _ = os.UserHomeDir()
	fp.AutoHeight = true
	fp.SetHeight(initialHeight - 1)
	fp.ShowHidden = true

	return newConnectionViewModel{returnModel: returnTo, filePicker: fp}
}

func (m newConnectionViewModel) Init() tea.Cmd {
	// initialize the new connection view model and file picker
	m.String = ""
	m.State = StateName
	m.SshEntry = sshEntry{}
	return m.filePicker.Init()
}

func (m newConnectionViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.filePicker.SetHeight(msg.Height - 1)
		m.filePicker.Update(msg)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			// handle Ctrl+C to exit the application
			return m, tea.Quit
		case tea.KeyEsc:
			// handle the Esc to go back to the main view
			return m.returnModel.Update(msg)
		case tea.KeyEnter:
			// process the input based on the current state
			switch m.State {
			case StateName:
				if m.String == "" {
					return m, nil
				}
				m.SshEntry.Name = m.String
				m.String = ""
				m.State = StateHost
			case StateHost:
				if m.String == "" {
					return m, nil
				}
				// no host validation done, you better know what you're doing
				m.SshEntry.Host = m.String
				m.String = ""
				m.State = StatePort
			case StatePort:
				if m.String == "" {
					return m, nil
				}
				// do port validation
				if port, err := strconv.Atoi(m.String); err != nil && port < 1 || port > 65535 {
					return m, nil
				}
				m.SshEntry.Port = m.String
				m.String = ""
				m.State = StateUser
			case StateUser:
				if m.String == "" {
					return m, nil
				}
				m.SshEntry.User = m.String
				m.String = ""
				m.State = StateIfKey
			case StateIfKey:
				if m.String == "y" || m.String == "Y" {
					// if user wants key, prompt for key
					m.State = StateKey
					m.String = ""
				} else {
					// If the user does not want to use a key, move to the done state
					m.SshEntry.Key = ""
					m.State = StateDone
				}

			case StateKey:
				// get the private key, probably use a file picker
				m.filePicker.Update(tea.KeyEnter)
			case StateDone:
				// save everything and return to the main view
				m.returnModel.data.SshEntries = append(m.returnModel.data.SshEntries, m.SshEntry)
				saveData(m.returnModel.data)
				// return to main view
				return m.returnModel.Update(nil)
			}

		case tea.KeyBackspace:
			// remove last character
			if len(m.String) > 0 {
				m.String = m.String[:len(m.String)-1]
			}
		default:
			if m.State != StateKey {
				if utf8.RuneCountInString(msg.String()) == 1 && unicode.IsPrint([]rune(msg.String())[0]) {
					m.String += msg.String()
				}
			}
		}
	}
	if m.State == StateKey {
		var cmd tea.Cmd
		m.filePicker, cmd = m.filePicker.Update(msg)

		// check if the user selected a file
		if didSelect, path := m.filePicker.DidSelectFile(msg); didSelect {
			// get path of the selected file
			m.selectedFile = path
			m.SshEntry.Key = path
			m.State = StateDone
			return m, nil
		}

		// check if user selected invalid file
		if didSelect, path := m.filePicker.DidSelectDisabledFile(msg); didSelect {
			// clear selected file and show error
			m.err = errors.New(path + " is not valid.")
			m.selectedFile = ""
			return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
		}

		return m, cmd
	}

	return m, nil
}

func (m newConnectionViewModel) View() string {
	// render the new connection view
	s := ""
	if m.State == StateKey {
		s += m.filePicker.View()
	} else if m.State == StateDone {
		// print out the summary of the SSH entry
		s += "SSH Connection Summary:\n\n"
		s += "Name: " + m.SshEntry.Name + "\n"
		s += "Host: " + m.SshEntry.Host + "\n"
		if m.SshEntry.Port != "" {
			s += "Port: " + m.SshEntry.Port + "\n"
		}
		s += "User: " + m.SshEntry.User + "\n"
		if m.SshEntry.Key != "" {
			s += "Auth: " + m.SshEntry.Key + "\n"
		} else {
			s += "Auth: Password\n"
		}
		s += "\n| Enter save | Esc cancel |\n\n"
	} else {
		s = "Enter " + stateNames[m.State] + ":\n\n"
		s += m.String + "\n"
	}
	return s
}
