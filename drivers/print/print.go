package print

import (
	"github.com/deployithq/deployit/utils"
	"github.com/fatih/color"
)

type Print struct {
	DebugMode bool
	Cyan      *color.Color
	Yellow    *color.Color
	Red       *color.Color
	White     *color.Color
}

func Init() *Print {

	p := new(Print)

	p.Cyan = color.New(color.FgCyan)
	p.Yellow = color.New(color.FgYellow)
	p.Red = color.New(color.FgRed)
	p.White = color.New(color.FgWhite)

	return p
}

func (p *Print) SetDebug(debug bool) {
	if debug {
		p.DebugMode = true
	}
}

func (p *Print) Info(a ...interface{}) {
	p.Cyan.Print(a, "\n")
}

func (p *Print) WhiteInfo(a ...interface{}) {
	p.White.Print(a, "\n")
}

func (p *Print) Infof(format string, a ...interface{}) {
	format += "\n"
	p.Cyan.Printf(format, a)
}

func (p *Print) Debug(a ...interface{}) {
	if p.DebugMode {
		p.White.Print("[DEBUG, ", utils.FileLine(), "]: ", a, "\n")
	}
}

func (p *Print) Debugf(format string, a ...interface{}) {
	if p.DebugMode {
		format += "\n"
		p.White.Printf(format, "[DEBUG, ", utils.FileLine(), "]: ", a)
	}
}

func (p *Print) Error(a ...interface{}) {
	p.Red.Print("[ERROR, ", utils.FileLine(), "]: ", a, "\n")
}

func (p *Print) Errorf(format string, a ...interface{}) {
	format += "\n"
	p.Red.Printf(format, "[ERROR, ", utils.FileLine(), "]: ", a)
}

func (p *Print) Warn(a ...interface{}) {
	p.Yellow.Print("[WARNING]: ", a, "\n")
}

func (p *Print) Warnf(format string, a ...interface{}) {
	format += "\n"
	p.Yellow.Printf(format, "[WARNING]: ", a)
}
