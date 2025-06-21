package views

import (
	"encoding/json"
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"strconv"
)

type sshEntry struct {
	Name        string
	Host        string
	Port        string
	User        string
	Password    string // password, if not set a key is used
	Key         string // path to private key file
	KeyPassword string // password to unlock the key, if key is set
}

type data struct {
	TerminalName string     // terminal used to open
	SshEntries   []sshEntry // list of SSH entries
}

type model struct {
	cursor        int  // which to-do list item our cursor is pointing at
	data          data // data to be saved in the file
	newConnection sshEntry
}

func InitialModel() model {
	return model{
		data: data{"", []sshEntry{}},
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."

	// Check if the file exists, if not create it
	if _, err := os.Stat("./ssh_tui.json"); errors.Is(err, os.ErrNotExist) {
		fmt.Println("File ssh_tui.json does not exist, creating a new one.")
		file, err := os.Create("./ssh_tui.json")
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			return nil
		}
		dat, err := json.Marshal(m.data)
		println(string(dat))
		_, err = file.Write(dat)
		if err != nil {
			return nil
		}

		err = file.Close()
		if err != nil {
			fmt.Printf("Error closing file: %v\n", err)
			return nil
		}
	}
	dat, err := os.ReadFile("./ssh_tui.json") //UserHomeDir()
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return nil
	}
	fmt.Print(string(dat))
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.data.SshEntries)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			// open the selected SSH entry
			return newConnectionView(m).Update(msg)
		case "n":
			// Add a new SSH entry
			name := "New SSH Entry " + strconv.Itoa(len(m.data.SshEntries)+1)
			newEntry := sshEntry{
				Name:        name,
				Host:        "new_host",
				Port:        "22",
				User:        "new_user",
				Password:    "",
				Key:         "",
				KeyPassword: "",
			}
			m.data.SshEntries = append(m.data.SshEntries, newEntry)

		case "d":
			// Delete the selected SSH entry
			if len(m.data.SshEntries) > 0 && m.cursor < len(m.data.SshEntries) {
				// Remove the entry at the cursor position
				m.data.SshEntries = append(m.data.SshEntries[:m.cursor], m.data.SshEntries[m.cursor+1:]...)
				// Adjust the cursor if necessary
				if m.cursor >= len(m.data.SshEntries) {
					m.cursor = len(m.data.SshEntries) - 1
				}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "SSH connections:\n\n"

	// Iterate over our choices
	/*for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}*/

	for i, entry := range m.data.SshEntries {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected

		// Render the row
		s += fmt.Sprintf("%s [%s] %s %s@%s:%s\n", cursor, checked, entry.Name, entry.User, entry.Host, entry.Port)
	}

	// The footer
	s += "\n| Q quit | N new | D delete |\n"

	// Send the UI for rendering
	return s
}
