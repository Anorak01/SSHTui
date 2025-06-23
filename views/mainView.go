package views

import (
	"encoding/json"
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"os/exec"
	"syscall"
)

type sshEntry struct {
	Name string
	Host string
	Port string
	User string
	Key  string // path to private key file
}

type data struct {
	TerminalName string     // terminal used to open ssh connection
	SshEntries   []sshEntry // list of SSH entries
	HideData     bool       // hide sensitive data
	dataDir      string
}

type model struct {
	cursor        int  // which ssh entry cursor is pointing at
	data          data // data to be saved in a file
	newConnection sshEntry
	termHeight    int
}

func InitialModel() *model {
	return &model{
		data: data{"kitty", []sshEntry{}, false, ""},
	}
}

func (m *model) Init() tea.Cmd {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	m.data.dataDir = home + "/.sshtui/ssh_tui.json"
	print("Data directory: " + m.data.dataDir + "\n")

	// check if the file exists, if not create it
	if _, err := os.Stat(m.data.dataDir); errors.Is(err, os.ErrNotExist) {
		fmt.Println("File ssh_tui.json does not exist, creating a new one.")
		err := os.Mkdir(home+"/.sshtui", os.ModePerm) // make the directory just in case
		if err != nil {
			return nil
		}
		file, err := os.Create(m.data.dataDir)
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

	// read the file and unmarshal the data
	dat, err := os.ReadFile(m.data.dataDir)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return nil
	}
	err = json.Unmarshal(dat, &m.data)
	if err != nil {
		fmt.Printf("Error unmarshalling data: %v\n", err)
		return nil
	}

	return nil
}
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termHeight = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		// these should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.data.SshEntries)-1 {
				m.cursor++
			}

		case "home":
			// hide sensitive data
			m.data.HideData = !m.data.HideData
			// save the updated data to the file
			saveData(m.data)

		case "enter", " ":
			// open the selected SSH entry
			selectedEntry := m.data.SshEntries[m.cursor]
			sshArgs := []string{"ssh"}
			sshArgs = append(sshArgs, fmt.Sprintf("%s@%s", selectedEntry.User, selectedEntry.Host))
			if selectedEntry.Port != "" {
				sshArgs = append(sshArgs, "-p", selectedEntry.Port)
			}
			if selectedEntry.Key != "" { // has key
				sshArgs = append(sshArgs, "-i", selectedEntry.Key)
			}

			cmd := exec.Command(m.data.TerminalName, sshArgs...)
			cmd.SysProcAttr = &syscall.SysProcAttr{
				Setpgid: true,
				Pgid:    0,
			}
			err := cmd.Start()
			if err != nil {
				fmt.Printf("Error starting SSH command: %v\n", err)
				return m, nil
			}

		case "n": // open new connection view
			return newConnectionView(m, m.termHeight).Update(msg)

		case "d":
			// delete the selected SSH entry
			if len(m.data.SshEntries) > 0 && m.cursor < len(m.data.SshEntries) {
				// remove the entry at the cursor position
				m.data.SshEntries = append(m.data.SshEntries[:m.cursor], m.data.SshEntries[m.cursor+1:]...)
				// adjust the cursor if necessary
				if m.cursor >= len(m.data.SshEntries) {
					m.cursor = len(m.data.SshEntries) - 1
				}
				// save the updated data to the file
				saveData(m.data)
			}
		}
	}
	return m, nil
}

func saveData(m data) {
	// save the current data to the file
	dat, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling data: %v\n", err)
		return
	}

	err = os.WriteFile(m.dataDir, dat, 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}
}

func (m *model) View() string {
	s := "SSH connections: " + "\n\n"

	for i, entry := range m.data.SshEntries {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		if !m.data.HideData {
			// render the row
			s += fmt.Sprintf("%s [] %s %s@%s:%s\n", cursor, entry.Name, entry.User, entry.Host, entry.Port)
		} else {
			// render the row with hidden data
			s += fmt.Sprintf("%s [] %s %s@%s\n", cursor, entry.Name, "***", "****")
		}

	}

	s += "\n| Q quit | N new | D delete |\n"

	return s
}
