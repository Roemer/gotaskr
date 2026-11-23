package gotaskr

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
)

// colorEnvVar is the environment variable that sets the color mode when --color is not given.
const colorEnvVar = "GOTASKR_COLOR"

// colorEnabled is the resolved color setting for all gotaskr output.
// Until Execute resolves --color, it follows fatih/color's own detection.
var colorEnabled = !color.NoColor

// ColorEnabled reports whether gotaskr output uses color, after applying --color, GOTASKR_COLOR,
// NO_COLOR and the terminal check. Tasks can use it to decide whether to color their own output.
func ColorEnabled() bool {
	return colorEnabled
}

// colorMode is the color mode requested by the user.
type colorMode int

const (
	colorAuto colorMode = iota
	colorOn
	colorOff
)

// parseColorMode accepts on/off/auto and the common aliases always/never/true/false.
func parseColorMode(value string) (colorMode, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "auto":
		return colorAuto, nil
	case "on", "always", "true":
		return colorOn, nil
	case "off", "never", "false":
		return colorOff, nil
	default:
		return colorAuto, fmt.Errorf("unknown color mode '%s', supported modes: on, off, auto", value)
	}
}

// requestedColorMode reads the color mode: the --color argument, then the environment variable, then auto.
// A bare --color (no value) means on.
func requestedColorMode() (colorMode, error) {
	if value, ok := GetArgument("color"); ok {
		if value == "" {
			return colorOn, nil
		}
		return parseColorMode(value)
	}
	if value := os.Getenv(colorEnvVar); value != "" {
		return parseColorMode(value)
	}
	return colorAuto, nil
}

// isTerminal reports whether w is a terminal. It is a variable so tests can replace it.
var isTerminal = func(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
}

// initColor resolves the color setting and applies it to gotaskr and to fatih/color's global switch.
func initColor() error {
	mode, err := requestedColorMode()
	if err != nil {
		return err
	}
	switch mode {
	case colorOn:
		colorEnabled = true
	case colorOff:
		colorEnabled = false
	default:
		// See https://no-color.org: a non-empty NO_COLOR disables color unless the user explicitly asks for it
		noColor := os.Getenv("NO_COLOR") != ""
		colorEnabled = !noColor && os.Getenv("TERM") != "dumb" && isTerminal(os.Stdout)
	}
	color.NoColor = !colorEnabled
	return nil
}

// newColor returns a color that is on or off only according to enabled.
// fatih/color disables each new color when NO_COLOR is set, so the global switch alone cannot force color on.
func newColor(enabled bool, attributes ...color.Attribute) *color.Color {
	c := color.New(attributes...)
	if enabled {
		c.EnableColor()
	} else {
		c.DisableColor()
	}
	return c
}

// printColored writes one line in the given color, following the resolved color setting.
func printColored(attribute color.Attribute, format string, a ...any) {
	_, _ = newColor(colorEnabled, attribute).Fprintln(color.Output, fmt.Sprintf(format, a...))
}
