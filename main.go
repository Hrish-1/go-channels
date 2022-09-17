package main

import (
	"channels/pets"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	initialModel := model{
		spinner: s,
		out:     petInfo(""),
		loading: true,
	}
	err := tea.NewProgram(initialModel).Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type petInfo string

type model struct {
	out     petInfo
	spinner spinner.Model
	loading bool
}

// Init implements tea.Model
func (m model) Init() tea.Cmd {
	return tea.Batch(spinner.Tick, fetchPetInfo)
}

// Update implements tea.Model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case petInfo:
		m.loading = false
		m.out = msg
		return m, tea.Quit
	}
	if m.loading {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

// View implements tea.Model
func (m model) View() string {
	if m.loading {
		return fmt.Sprintf("%s fetching pet info... please wait", m.spinner.View())
	}
	return string(m.out)
}

func fetchPetInfo() tea.Msg {
	return petInfo(pets.PetInfo())
}
