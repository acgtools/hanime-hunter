package progressbar

import (
	"io"
	"os"
	"sync"
	"sync/atomic"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	err   error
	width int
	pbs   []*ProgressBar // sorted pbs, cache for View()
	Mux   sync.Mutex     // protect fields below
	Pbs   map[string]*ProgressBar
}

var _ tea.Model = (*Model)(nil)

func (m *Model) AddPb(pb *ProgressBar) {
	m.Mux.Lock()
	defer m.Mux.Unlock()
	m.Pbs[pb.FileName] = pb
}

type ProgressBar struct {
	Pw       *ProgressWriter
	Pc       *ProgressCounter
	Progress progress.Model
	FileName string
	Status   string
}

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
		p.Send(ProgressErrMsg{err})
	}
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	pw.Downloaded += int64(len(p))
	if pw.Total > 0 && pw.OnProgress != nil {
		pw.OnProgress(pw.FileName, float64(pw.Downloaded)/float64(pw.Total))
	}
	return len(p), nil
}

type ProgressCounter struct {
	Total      int64
	Downloaded atomic.Int64
	FileName   string
	Onprogress func(string, float64)
}

func (pc *ProgressCounter) Increase() {
	d := pc.Downloaded.Add(1)
	if pc.Total > 0 && pc.Onprogress != nil {
		pc.Onprogress(pc.FileName, float64(d)/float64(pc.Total))
	}
}
