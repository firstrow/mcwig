package pipe

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/firstrow/mcwig"
	"github.com/firstrow/mcwig/testutils"
)

func TestPipe(t *testing.T) {
	e := mcwig.NewEditor(
		testutils.Viewport,
		nil,
	)

	buf := e.BufferGetByName("test-1")
	e.ActiveWindow().Buffer = buf
	e.ActiveBuffer().AppendStringLine(`echo "%s"`)

	p := New(e, Options{IsPrompt: false})
	assert.Equal(t, "echo", p.getCommand())

	p.send(`ping pong`)
	p.cmd.Wait()

	outBuffer := e.BufferGetByName("[output]")
	assert.Equal(t, "ping pong", outBuffer.String())

	args := p.buildArgs("ping pong")
	assert.Equal(t, 1, len(args))
	assert.Equal(t, `ping pong`, args[0])
}

func TestPipeLongRunningProcess(t *testing.T) {
	e := mcwig.NewEditor(
		testutils.Viewport,
		nil,
	)

	buf := e.BufferGetByName("test-1")
	e.ActiveWindow().Buffer = buf
	e.ActiveBuffer().AppendStringLine(`python -i`)

	p := New(e, Options{IsPrompt: true})
	assert.Equal(t, "python", p.getCommand())

	p.send(`help`)
	// TODO: figure out how to Wait properly
	time.Sleep(100 * time.Millisecond)

	outBuffer := e.BufferGetByName("[output]")
	assert.Contains(t, outBuffer.String(), "Type help() for interactive help, or help(object) for help about object.")
}
