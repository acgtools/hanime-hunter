package progress

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"path/filepath"
	"strings"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

func (m model) View() string {
	if m.err != nil {
		return "Error downloading: " + m.err.Error() + "\n"
	}

	var sb strings.Builder

	pad := strings.Repeat(" ", padding)

	title := filepath.Base(m.fileName)
	bar := lipgloss.JoinHorizontal(lipgloss.Top, pad, m.progress.View(), pad, getDLStatus(m.pw.downloaded, m.pw.total), pad, title)

	sb.WriteString("\n")
	sb.WriteString(bar)

	sb.WriteString("\n")
	sb.WriteString(helpStyle("Press ctrl+c to quit"))

	return sb.String()
}

func getDLStatus(downloaded, total int) string {
	m := float64(1 << 20)
	return fmt.Sprintf("%.2f MiB/%.2f MiB", float64(downloaded)/m, float64(total)/m)
}
