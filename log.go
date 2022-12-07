package main

import (
	"fmt"
	"io"
)

// Logger is an interface for logging
type Logger interface {
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
