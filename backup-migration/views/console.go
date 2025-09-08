package views

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

// ConsoleWidget is a terminal-like log view backed by a List + data binding.
// It implements io.Writer, so you can pipe logs straight into it.
type ConsoleWidget struct {
	mu       sync.Mutex
	partial  bytes.Buffer       // holds any trailing non-terminated line
	data     binding.StringList // each list item is a log line
	List     *widget.List       // the visible list widget
	maxLines int                // 0 = unlimited, else ring buffer size
}

func NewConsoleWidget(maxLines int) (*ConsoleWidget, fyne.CanvasObject) {
	data := binding.NewStringList()
	list := widget.NewListWithData(
		data,
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(item binding.DataItem, obj fyne.CanvasObject) {
			s, _ := item.(binding.String).Get()
			obj.(*widget.Label).SetText(s)
		},
	)
	cw := &ConsoleWidget{data: data, List: list, maxLines: maxLines}
	return cw, list
}

// Write appends bytes, emitting complete lines to the list.
// (Lines are split on '\n'; the final partial line is buffered until the next write.)
func (c *ConsoleWidget) Write(p []byte) (int, error) {
	c.mu.Lock()
	c.partial.Write(p)
	s := c.partial.String()
	c.mu.Unlock()

	lastNL := strings.LastIndexByte(s, '\n')
	if lastNL >= 0 {
		lines := strings.Split(s[:lastNL], "\n")
		for _, ln := range lines {
			c.appendLine(ln)
		}
		// keep the remainder
		c.mu.Lock()
		c.partial.Reset()
		c.partial.WriteString(s[lastNL+1:])
		c.mu.Unlock()

		// Scroll to bottom on UI thread
		fyne.Do(func() {
			n := c.data.Length()
			if n > 0 {
				c.List.ScrollTo(n - 1)
			}
		})
	}
	return len(p), nil
}

func (c *ConsoleWidget) Printf(format string, args ...any) {
	_, _ = c.Write([]byte(fmt.Sprintf(format, args...)))
}

func (c *ConsoleWidget) Clear() {
	c.mu.Lock()
	c.partial.Reset()
	c.mu.Unlock()
	_ = c.data.Set([]string{})
}

func (c *ConsoleWidget) appendLine(line string) {
	if c.maxLines <= 0 {
		_ = c.data.Append(line)
		return
	}
	cur, _ := c.data.Get()
	cur = append(cur, line)
	if len(cur) > c.maxLines {
		cur = cur[len(cur)-c.maxLines:]
	}
	_ = c.data.Set(cur)
}
