package prompt_test

import (
	"errors"
	"os"
	"testing"

	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/assert"

	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/AlecAivazis/survey.v1/terminal"

	"github.com/petrovskiborislav/docker-cli/prompt"

	pseudotty "github.com/creack/pty"
)

func TestSelectPrompt_WhenNoProblem(t *testing.T) {
	// Arrange
	msg := "Select items"
	items := []string{"item1", "item2", "item3"}

	procedure := func(c expectConsole) {
		// select item1
		c.Send(" ")
		// select item3
		c.Send(string(terminal.KeyArrowDown))
		c.Send(string(terminal.KeyArrowDown))
		c.SendLine(" ")
	}

	pty, tty, err := pseudotty.Open()
	assert.NoError(t, err)

	var result []string
	exec := func(stdio terminal.Stdio) error {
		opt := survey.WithStdio(stdio.In, stdio.Out, stdio.Err)
		result, err = prompt.NewPrompt().SelectPrompt(msg, items, opt)
		return err
	}

	// Act
	runTest(t, pty, tty, procedure, exec)

	// Assert
	want := []string{"item1", "item3"}

	assert.NoError(t, err)
	assert.EqualValues(t, want, result)
}

func TestSelectPrompt_WhenErrorOccursOnCreationOfPrompt(t *testing.T) {
	// Arrange
	msg := "Select items"
	items := []string{"item1", "item2", "item3"}

	opt := func(options *survey.AskOptions) error {
		return errors.New("error")
	}

	// Act
	var result []string
	result, err := prompt.NewPrompt().SelectPrompt(msg, items, opt)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
}

// Helpers
type expectConsole interface {
	Send(string)
	SendLine(string)
}

type consoleWithErrorHandling struct {
	console *expect.Console
	t       *testing.T
}

func (c *consoleWithErrorHandling) SendLine(s string) {
	if _, err := c.console.SendLine(s); err != nil {
		c.t.Helper()
		c.t.Fatalf("SendLine(%q) = %v", s, err)
	}
}

func (c *consoleWithErrorHandling) Send(s string) {
	if _, err := c.console.Send(s); err != nil {
		c.t.Helper()
		c.t.Fatalf("Send(%q) = %v", s, err)
	}
}

func runTest(t *testing.T, pty, tty *os.File, procedure func(expectConsole), exec func(terminal.Stdio) error) {
	term := vt10x.New(vt10x.WithWriter(tty))
	c, err := expect.NewConsole(expect.WithStdin(pty), expect.WithStdout(term), expect.WithCloser(pty, tty))
	if err != nil {
		t.Fatalf("failed to create console: %v", err)
	}
	defer c.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		procedure(&consoleWithErrorHandling{console: c, t: t})
	}()

	stdio := terminal.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()}
	if err := exec(stdio); err != nil {
		t.Error(err)
	}

	if err := c.Tty().Close(); err != nil {
		t.Errorf("error closing Tty: %v", err)
	}
	<-done
}
