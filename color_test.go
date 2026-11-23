package gotaskr

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

// setupColorTest sets the arguments, a controlled environment and a fake terminal setting.
// The color-related variables of the developer's environment are always cleared first,
// and all globals changed by initColor are restored after the test.
func setupColorTest(t *testing.T, args map[string]string, env map[string]string, terminal bool) {
	t.Helper()
	for _, name := range []string{colorEnvVar, "NO_COLOR", "TERM"} {
		t.Setenv(name, env[name])
	}
	oldArgs, oldIsTerminal, oldColorEnabled, oldNoColor := argumentsMap, isTerminal, colorEnabled, color.NoColor
	t.Cleanup(func() {
		argumentsMap, isTerminal, colorEnabled, color.NoColor = oldArgs, oldIsTerminal, oldColorEnabled, oldNoColor
	})
	if args == nil {
		args = map[string]string{}
	}
	argumentsMap = args
	isTerminal = func(io.Writer) bool { return terminal }
}

// captureColorOutput redirects gotaskr's colored lines (errors and the run summary) into a buffer.
func captureColorOutput(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	oldOutput := color.Output
	color.Output = &buf
	t.Cleanup(func() { color.Output = oldOutput })
	return &buf
}

func TestColorResolution(t *testing.T) {
	tests := []struct {
		name      string
		args      map[string]string
		env       map[string]string
		terminal  bool
		wantColor bool
	}{
		{"auto on terminal", nil, nil, true, true},
		{"auto not terminal", nil, nil, false, false},
		{"auto with NO_COLOR", nil, map[string]string{"NO_COLOR": "1"}, true, false},
		{"auto with TERM=dumb", nil, map[string]string{"TERM": "dumb"}, true, false},
		{"flag auto on terminal", map[string]string{"color": "auto"}, nil, true, true},
		{"flag on not terminal", map[string]string{"color": "on"}, nil, false, true},
		{"flag on beats NO_COLOR", map[string]string{"color": "on"}, map[string]string{"NO_COLOR": "1"}, false, true},
		{"flag off on terminal", map[string]string{"color": "off"}, nil, true, false},
		{"bare flag means on", map[string]string{"color": ""}, nil, false, true},
		{"flag aliases", map[string]string{"color": "ALWAYS"}, nil, false, true},
		{"env on", nil, map[string]string{colorEnvVar: "on"}, false, true},
		{"env off on terminal", nil, map[string]string{colorEnvVar: "never"}, true, false},
		{"flag beats env", map[string]string{"color": "off"}, map[string]string{colorEnvVar: "on"}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			setupColorTest(t, tt.args, tt.env, tt.terminal)

			err := initColor()

			assert.NoError(err)
			assert.Equal(tt.wantColor, ColorEnabled())
			assert.Equal(!tt.wantColor, color.NoColor, "fatih/color global switch follows the setting")
		})
	}
}

func TestColorInvalidValue(t *testing.T) {
	for name, tc := range map[string]struct {
		args map[string]string
		env  map[string]string
	}{
		"flag": {map[string]string{"color": "purple"}, nil},
		"env":  {nil, map[string]string{colorEnvVar: "purple"}},
	} {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			setupColorTest(t, tc.args, tc.env, true)

			err := initColor()

			assert.EqualError(err, "unknown color mode 'purple', supported modes: on, off, auto")
		})
	}
}

func TestInvalidColorStopsBeforeRunningTasks(t *testing.T) {
	assert := assert.New(t)
	clear()
	ran := false
	Task("default", func() error { ran = true; return nil })
	setupColorTest(t, map[string]string{"color": "purple"}, nil, false)

	exitCode := Execute()

	assert.Equal(1, exitCode)
	assert.False(ran)
}

func TestColorAppliesToRunOutput(t *testing.T) {
	for name, tc := range map[string]struct {
		args      map[string]string
		env       map[string]string
		wantColor bool
	}{
		"flag on":                {map[string]string{"color": "on"}, nil, true},
		"flag on beats NO_COLOR": {map[string]string{"color": "on"}, map[string]string{"NO_COLOR": "1"}, true},
		"flag off":               {map[string]string{"color": "off"}, nil, false},
		"env on":                 {map[string]string{}, map[string]string{colorEnvVar: "on"}, true},
	} {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			clear()
			Task("Fails", func() error { return errors.New("boom") })
			out := captureColorOutput(t)
			tc.args["target"] = "Fails"
			setupColorTest(t, tc.args, tc.env, false)

			exitCode := Execute()

			assert.Equal(1, exitCode)
			assert.Contains(out.String(), "Task error: boom")
			assert.Contains(out.String(), "Exit Code")
			assert.Equal(tc.wantColor, bytes.Contains(out.Bytes(), []byte("\x1b[")), out.String())
		})
	}
}
