package progressbar

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	padding  = 2
	maxWidth = 80
)

type ProgressMsg struct {
	FileName string
	Ratio    float64
}

type ProgressErrMsg struct{ Err error }

type ProgressStatusMsg struct {
	FileName string
	Status   string
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:cyclop
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		return m, nil

	case tea.WindowSizeMsg:
		w := msg.Width - padding*2 - 4 //nolint:gomnd
		if w >= maxWidth {
			w = maxWidth
		}
		m.width = w
		return m, nil

	case ProgressErrMsg:
		m.err = msg.Err
		return m, tea.Quit

	case ProgressMsg:
		var cmds []tea.Cmd

		fileName, ratio := msg.FileName, msg.Ratio
		if pb, ok := m.Pbs[fileName]; ok {
			cmds = append(cmds, pb.Progress.SetPercent(ratio))
		}
		return m, tea.Batch(cmds...)

	case ProgressStatusMsg:
		fileName, status := msg.FileName, msg.Status
		if pb, ok := m.Pbs[fileName]; ok {
			pb.Status = status
		}
		return m, nil

	case progress.FrameMsg:
		var cmds []tea.Cmd

		for _, pb := range m.Pbs {
			progressModel, cmd := pb.Progress.Update(msg)
			pbm, ok := progressModel.(progress.Model)
			if ok {
				pb.Progress = pbm
			}
			cmds = append(cmds, cmd)
		}

		return m, tea.Batch(cmds...)

	default:
		return m, nil
	}
}
