package progressbar

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"sort"
	"strings"
)

var (
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
	pbStyle   = lipgloss.NewStyle().MaxHeight(1).Render
)

func (m *Model) View() string {
	if m.err != nil {
		return "Error downloading: " + m.err.Error() + "\n"
	}

	var sb strings.Builder
	pad := strings.Repeat(" ", padding)

	if m.pbs == nil {
		m.pbs = make([]*ProgressBar, 0)
	}
	if len(m.pbs) != len(m.Pbs) {
		m.pbs = pbMapToSortedSlice(m.Pbs, m.width)
	}

	for _, pb := range m.pbs {
		bar := lipgloss.JoinHorizontal(lipgloss.Top, pad, pb.Progress.View(),
			pad, getDLStatus(pb.Pw.Downloaded, pb.Pw.Total),
			pad, pb.FileName)

		sb.WriteString("\n\n")
		sb.WriteString(pbStyle(bar))
	}

	sb.WriteString("\n\n")
	sb.WriteString(helpStyle("Press ctrl+c to quit"))

	return sb.String()
}

func pbMapToSortedSlice(m map[string]*ProgressBar, w int) []*ProgressBar {
	res := make([]*ProgressBar, 0, len(m))
	for _, v := range m {
		v.Progress.Width = w
		res = append(res, v)
	}

	sort.SliceStable(res, func(i, j int) bool {
		return res[i].FileName < res[j].FileName
	})

	return res
}

func getDLStatus(downloaded, total int64) string {
	m := float64(1 << 20)
	return fmt.Sprintf("%.2f MiB/%.2f MiB", float64(downloaded)/m, float64(total)/m)
}
