package iohelper

import (
	"bufio"
	"io"
	"os"
)

// StdinProcessor encapsulates functions for work with stdin.
type StdinProcessor struct{}

// NewStdinProcessor creates new StdinProcessor ыекгсе and returns a pointer to it.
func NewStdinProcessor() *StdinProcessor {
	return &StdinProcessor{}
}

// ready checks stdin and return bool result with error.
func (*StdinProcessor) ready() (bool, error) {

	stat, err := os.Stdin.Stat()
	if err != nil {
		return false, err
	}

	ready := (stat.Mode() & os.ModeCharDevice) == 0
	return ready, nil
}

// Ready checks stdin and return true, if stdin contains data.
func (p *StdinProcessor) Ready() bool {
	ready, err := p.ready()
	return err == nil && ready
}

// Process calls function for every row from stdin.
// The function is specified as a parameter.
func (p *StdinProcessor) Process(pfunc func(row []byte) error) error {

	var input []byte

	stdinReady, err := p.ready()
	if err != nil {
		return err
	}

	if stdinReady {
		reader := bufio.NewReader(os.Stdin)
		for {
			bytes, hasMoreInLine, err := reader.ReadLine()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			input = append(input, bytes...)
			if !hasMoreInLine {
				if err := pfunc(input); err != nil {
					return err
				}
				input = nil
			}
		}
	}

	return nil
}

// AsyncProcess sends every row from stdin to channel.
func (p *StdinProcessor) AsyncProcess(line chan []byte) error {

	process := func(row []byte) error {
		line <- row
		return nil
	}

	return p.Process(process)
}
