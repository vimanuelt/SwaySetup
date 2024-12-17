package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styling
var (
	// General colors
	baseStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("249"))
	titleStyle      = lipgloss.NewStyle().Padding(0, 2).Bold(true).Foreground(lipgloss.Color("117")).Background(lipgloss.Color("236"))
	statusStyle     = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("255"))
	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true)
	successStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("41")).Bold(true)
	listItemStyle   = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("250"))
	listSelectedStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Background(lipgloss.Color("32")).
				Foreground(lipgloss.Color("230")).
				Bold(true)
)

// Timeout duration for system commands
const commandTimeout = 10 * time.Second

// Item structure for the list
type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// Model to manage TUI
type model struct {
	list    list.Model
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			return m, executeAction(m.list.SelectedItem().(item).title)
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Sway & Wayland Setup for FreeBSD"),
		m.list.View(),
		statusStyle.Render("Press q to quit, ↑/↓ to navigate, Enter to select an option"),
	)
}

// Helper function to execute system commands with a timeout
func runCommandWithTimeout(name string, args ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	output, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("command timed out: %s", name)
	}

	if err != nil {
		return fmt.Errorf("command failed: %s, output: %s", err, string(output))
	}

	return nil
}

// Actions
func installPackages() error {
	return runCommandWithTimeout("pkg", "install", "-y", "sway", "swaylock", "swayidle", "seatd", "wayland")
}

func configureSway() error {
	homeDir, _ := os.UserHomeDir()
	configDir := fmt.Sprintf("%s/.config/sway", homeDir)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	configPath := fmt.Sprintf("%s/config", configDir)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return runCommandWithTimeout("cp", "/usr/local/etc/sway/config", configPath)
	}
	return nil
}

func setupSeatd() error {
	if err := runCommandWithTimeout("sysrc", "seatd_enable=YES"); err != nil {
		return fmt.Errorf("failed to enable seatd: %w", err)
	}

	if err := runCommandWithTimeout("service", "seatd", "status"); err == nil {
		return nil
	}

	if err := runCommandWithTimeout("service", "seatd", "start"); err != nil {
		return fmt.Errorf("failed to start seatd: %w", err)
	}

	return nil
}

func setEnvironment() error {
	homeDir, _ := os.UserHomeDir()
	profilePath := fmt.Sprintf("%s/.profile", homeDir)
	envVar := "\nexport XDG_RUNTIME_DIR=/tmp/xdg-runtime-$(id -u)\n"
	file, err := os.OpenFile(profilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(envVar)
	return err
}

func executeAction(action string) tea.Cmd {
	var err error
	switch action {
	case "Install Packages":
		err = installPackages()
	case "Configure Sway":
		err = configureSway()
	case "Setup seatd":
		err = setupSeatd()
	case "Set Environment":
		err = setEnvironment()
	}

	if err != nil {
		return tea.Printf(errorStyle.Render(fmt.Sprintf("Error: %v", err)))
	}
	return tea.Printf(successStyle.Render(fmt.Sprintf("Success: %s completed!", action)))
}

func main() {
	items := []list.Item{
		item{title: "Install Packages", desc: "Install Sway, Wayland, and dependencies"},
		item{title: "Configure Sway", desc: "Set up initial Sway configuration"},
		item{title: "Setup seatd", desc: "Enable and start seatd service"},
		item{title: "Set Environment", desc: "Set necessary environment variables"},
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = listSelectedStyle
	delegate.Styles.NormalTitle = listItemStyle

	list := list.New(items, delegate, 0, 0)
	list.Title = "Choose an action"

	p := tea.NewProgram(model{list: list})
	if err := p.Start(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

