package logger

import (
	"io"
	"sync"
	"time"
)

const (
	defaultTimeLayout = "2006-02-01 15:04:05 -0700"
)

type Handler interface {
	Handle(r Record) error
	Enabled(l Level) bool
}

type TextHandler struct {
	level      Level
	mu         sync.Mutex
	w          io.Writer
	timeLayout string
}

type Formatter struct {
	buf *Buffer
	sep string
}

func NewTextHandler(level Level, w io.Writer) Handler {
	return &TextHandler{
		level:      level,
		w:          w,
		timeLayout: defaultTimeLayout,
	}
}

func NewFormatter(buf *Buffer, sep string) *Formatter {
	return &Formatter{
		buf: buf,
		sep: sep,
	}
}

func (h *TextHandler) Enabled(l Level) bool {
	return l >= h.level
}

func (h *TextHandler) Handle(r Record) error {
	if !h.Enabled(r.Lvl) {
		return nil
	}
	buf := make(Buffer, 0, 256)
	f := NewFormatter(&buf, " ")
	if !r.Time.IsZero() {
		f.AppendTime(r.Time, h.timeLayout)
	}
	f.AppendLevel(r.Lvl)
	f.AppendString(r.Msg)
	for _, a := range r.Fields {
		f.AppendArg(a)
	}
	f.AppendStringNL()
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(buf)

	return err
}

func (f *Formatter) AppendLevel(lvl Level) error {
	f.buf.WriteString(lvl.String())
	f.buf.WriteString(f.sep)

	return nil
}

func (f *Formatter) AppendTime(t time.Time, layout string) {
	f.buf.WriteString(t.Format(layout))
	f.buf.WriteString(f.sep)
}

func (f *Formatter) AppendString(s string) {
	f.buf.WriteString(s)
	f.buf.WriteString(f.sep)
}

func (f *Formatter) AppendArg(a Field) {
	f.buf.WriteByte('{')
	f.buf.WriteByte('"')
	f.buf.WriteString(a.Key)
	f.buf.WriteByte('"')
	f.buf.WriteByte('}')
	f.buf.WriteByte(',')
	f.buf.WriteByte('{')
	f.buf.WriteByte('"')
	f.buf.WriteString(a.Value)
	f.buf.WriteByte('"')
	f.buf.WriteByte('}')

	f.buf.WriteString(f.sep)
}

func (f *Formatter) AppendStringNL() {
	f.buf.WriteByte('\n')
}
