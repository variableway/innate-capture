package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/variableway/innate/capture/internal/idea"
	"github.com/variableway/innate/capture/internal/model"
)

type ideaField int

const (
	ideaFieldTitle ideaField = iota
	ideaFieldDesc
	ideaFieldContext
)

type ideaForm struct {
	title    textinput.Model
	desc     textarea.Model
	context  textarea.Model
	focus    ideaField
	cfg      *model.Config
	result   string
	err      error
	submitted bool
}

func newIdeaForm(cfg *model.Config) *ideaForm {
	title := textinput.New()
	title.Placeholder = "Idea title"
	title.Focus()
	title.CharLimit = 128
	title.Width = 50

	desc := textarea.New()
	desc.Placeholder = "One-line description (optional)"
	desc.SetWidth(50)
	desc.SetHeight(3)
	desc.ShowLineNumbers = false
	desc.FocusedStyle.Base = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#7D56F4"))
	desc.BlurredStyle.Base = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#585858"))

	ctx := textarea.New()
	ctx.Placeholder = "Raw context (optional)"
	ctx.SetWidth(50)
	ctx.SetHeight(3)
	ctx.ShowLineNumbers = false
	ctx.FocusedStyle.Base = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#7D56F4"))
	ctx.BlurredStyle.Base = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#585858"))

	return &ideaForm{
		title:   title,
		desc:    desc,
		context: ctx,
		focus:   ideaFieldTitle,
		cfg:     cfg,
	}
}

func (f *ideaForm) Init() tea.Cmd {
	return textinput.Blink
}

func (f *ideaForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if f.submitted || f.err != nil {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "enter", "esc":
				return f, tea.Quit
			}
		}
		return f, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			f.blur()
			f.focus = (f.focus + 1) % 3
			f.focusField()
			return f, nil
		case "shift+tab", "up":
			f.blur()
			f.focus = (f.focus + 2) % 3
			f.focusField()
			return f, nil
		case "esc":
			return f, tea.Quit
		case "enter":
			if f.focus == ideaFieldTitle {
				// Enter in title submits the form
				return f.submit()
			}
			// Enter in textareas adds newline, handled by textarea
		case "ctrl+s":
			return f.submit()
		}
	}

	// Route to focused field
	var cmd tea.Cmd
	switch f.focus {
	case ideaFieldTitle:
		f.title, cmd = f.title.Update(msg)
	case ideaFieldDesc:
		f.desc, cmd = f.desc.Update(msg)
	case ideaFieldContext:
		f.context, cmd = f.context.Update(msg)
	}
	return f, cmd
}

func (f *ideaForm) blur() {
	switch f.focus {
	case ideaFieldTitle:
		f.title.Blur()
	case ideaFieldDesc:
		f.desc.Blur()
	case ideaFieldContext:
		f.context.Blur()
	}
}

func (f *ideaForm) focusField() {
	switch f.focus {
	case ideaFieldTitle:
		f.title.Focus()
	case ideaFieldDesc:
		f.desc.Focus()
	case ideaFieldContext:
		f.context.Focus()
	}
}

func (f *ideaForm) submit() (tea.Model, tea.Cmd) {
	title := strings.TrimSpace(f.title.Value())
	if title == "" {
		f.err = fmt.Errorf("title is required")
		return f, nil
	}

	desc := strings.TrimSpace(f.desc.Value())
	ctx := strings.TrimSpace(f.context.Value())

	path, err := idea.Write(f.cfg, title, desc, ctx, idea.SourceTUI)
	if err != nil {
		f.err = err
		return f, nil
	}

	f.submitted = true
	f.result = path
	return f, nil
}

func (f *ideaForm) View() string {
	var sb strings.Builder

	sb.WriteString(ideaDialogTitleStyle.Render(" New Idea "))
	sb.WriteString("\n\n")

	// Title field
	sb.WriteString(ideaLabelStyle.Render("Title"))
	sb.WriteString(" (Enter to submit)\n")
	sb.WriteString(f.title.View())
	sb.WriteString("\n\n")

	// Description field
	sb.WriteString(ideaLabelStyle.Render("Description"))
	sb.WriteString("\n")
	sb.WriteString(f.desc.View())
	sb.WriteString("\n\n")

	// Context field
	sb.WriteString(ideaLabelStyle.Render("Context"))
	sb.WriteString("\n")
	sb.WriteString(f.context.View())
	sb.WriteString("\n\n")

	if f.submitted {
		sb.WriteString(ideaSuccessStyle.Render(fmt.Sprintf("Created: %s", f.result)))
		sb.WriteString("\n")
		sb.WriteString(ideaHintStyle.Render("Press Enter or Esc to close"))
	} else if f.err != nil {
		sb.WriteString(ideaErrorStyle.Render(fmt.Sprintf("Error: %v", f.err)))
		sb.WriteString("\n")
		sb.WriteString(ideaHintStyle.Render("Press Enter or Esc to close"))
	} else {
		sb.WriteString(ideaHintStyle.Render("Tab/Shift-Tab: switch fields | Ctrl+S: submit | Esc: cancel"))
	}

	return sb.String()
}
