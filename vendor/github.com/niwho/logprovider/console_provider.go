package logprovier

import (
	"fmt"
	"os"
	"runtime"
)

type Brush func(string) string

func NewBrush(color string) Brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

var colors = []Brush{
	NewBrush("1;35"), // panic    magenta
	NewBrush("1;35"), // Fatal    magenta
	NewBrush("1;31"), // Error    red
	NewBrush("1;33"), // Warn     yellow
	NewBrush("1;36"), // Info     cyan
	NewBrush("1;34"), // Debug    blue
}

type ConsoleProvider struct {
	level int
}

func NewConsoleProvider() *ConsoleProvider {
	return &ConsoleProvider{}
}

func (cp *ConsoleProvider) Init() error {
	return nil
}

func (cp *ConsoleProvider) SetLevel(l int) {
	cp.level = l
}

func (cp *ConsoleProvider) WriteMsg(msg string, level int) error {
	if level < cp.level {
		return nil
	}
	if goos := runtime.GOOS; goos == "windows" {
		fmt.Fprint(os.Stdout, msg)
	}
	fmt.Fprint(os.Stdout, colors[level](msg))
	return nil
}

func (cp *ConsoleProvider) Flush() error {
	return nil
}

func (cp *ConsoleProvider) Close() error {
	return nil
}
