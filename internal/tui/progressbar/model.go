package progressbar

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"io"
	"os"
	"sync"
)

type ProgressMsg struct {
	FileName string
	Ratio    float64
}

type progressErrMsg struct{ err error }

type Model struct {
	err   error
	width int
	pbs   []*ProgressBar // sorted pbs, cache for View()
	Mux   sync.Mutex     // protect fields bellow
	Pbs   map[string]*ProgressBar
}

type ProgressBar struct {
	Pw       *ProgressWriter
	Progress progress.Model
	FileName string
}

var _ tea.Model = (*Model)(nil)

type ProgressWriter struct {
	Total      int64
	Downloaded int64
	File       *os.File
	Reader     io.Reader
	FileName   string
	OnProgress func(string, float64)
}

func (pw *ProgressWriter) Start(p *tea.Program) {
	// TeeReader calls PW.Write() each time a new response is received
	_, err := io.Copy(pw.File, io.TeeReader(pw.Reader, pw))
	if err != nil {
		p.Send(progressErrMsg{err})
	}
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	pw.Downloaded += int64(len(p))
	if pw.Total > 0 && pw.OnProgress != nil {
		pw.OnProgress(pw.FileName, float64(pw.Downloaded)/float64(pw.Total))
	}
	return len(p), nil
}
