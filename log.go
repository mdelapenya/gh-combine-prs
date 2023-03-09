package main

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

// Logger is an interface for logging
type Logger interface {
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Fprintf(w io.Writer, format string, v ...interface{}) (int, error)
	Fprintln(w io.Writer, v ...interface{}) (int, error)
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type logger struct {
	Verbose bool
}

func newLogger(verbose bool) Logger {
	return logger{
		Verbose: verbose,
	}
}

// Debugf prints a formatted debug message
func (l logger) Debugf(format string, v ...interface{}) {
	if !l.Verbose {
		return
	}
	color.White(format, v...)
}

// Errorf prints a formatted error, prepending ">> Error: " to the message
func (l logger) Errorf(format string, v ...interface{}) {
	color.Red(">> Error: "+format, v...)
}

// Infof prints a formatted info message
func (l logger) Infof(format string, v ...interface{}) {
	if !l.Verbose {
		return
	}
	color.Green(format, v...)
}

// Warnf prints a formatted warn message, prepending ">> Warn: " to the message
func (l logger) Warnf(format string, v ...interface{}) {
	if !l.Verbose {
		return
	}
	color.Yellow(">> Warn: "+format, v...)
}

// Fprintf prints a formatted string to a writer
func (l logger) Fprintf(w io.Writer, format string, v ...interface{}) (int, error) {
	if !l.Verbose {
		return 0, nil
	}
	return fmt.Fprintf(w, format, v...)
}

// Fprintln prints a string to a writer
func (l logger) Fprintln(w io.Writer, v ...interface{}) (int, error) {
	if !l.Verbose {
		return 0, nil
	}
	return fmt.Fprintln(w, v...)
}

// Printf prints a formatted string
func (l logger) Printf(format string, v ...interface{}) {
	if !l.Verbose {
		return
	}
	fmt.Printf(format, v...)
}

// Println prints a string
func (l logger) Println(v ...interface{}) {
	if !l.Verbose {
		return
	}
	fmt.Println(v...)
}
