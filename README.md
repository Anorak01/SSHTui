# SSHTui
A small TUI for keeping SSH connections.

I made this tool for myself and as such there are some hardcoded values that you might want (need) to change. There is also little to no input validation as of now.

I make no guarantees that it will work for you, it's just here for anyone who needs something similar.

# Installation
```bash
git clone https://github.com/Anorak01/SSHTui
cd SSHTui
go install
```
If you dont use `kitty` terminal, go into the config file and change the terminal to what you want

## Usage
Open the application with:
```bash
sshtui
```
*make sure you have your `$HOME/go/bin` in your `$PATH`

```
N to make a new connection
D to delete a connection
Q to quit
```

Connections are stored in `$HOME/.sshtui/sshtui.json`, go there if you need to edit them

## Thanks to
Bubble Tea for the TUI framework
Bubbles - filepicker has been pulled in for local testing, its license in inside the views/filepicker folder

### Known issues
* file picker doesn't properly load the initial directory, starts working after you go up the tree