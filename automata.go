package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const (
	errorTokenName = "__ERROR_TOKEN"
	emptyTokenName = "_EMPTY_TOKEN"
)

// Automata represents DFA.
type Automata struct {
	StartState *State
	Reader     *bufio.Reader
	Cache      strings.Builder // Cache stores symbols read between start and final state
	Finished   bool
	Position   int
	ErrorToken bool // true when error state was visited
	LastSymbol rune // Stores last symbol from input if recent transition was on empty symbol
}

// NextToken returns next token and indicator if file processing is finished
func (a *Automata) NextToken() (Token, bool) {
	currentState := a.StartState

	defer func() {
		r := recover()
		if r != nil {
			panic(fmt.Sprintf("Position %d - cannot process, current cache: %s, state: %s-%d",
				a.Position, a.Cache.String(), currentState.Name, currentState.ID))
		}
	}()

	a.Finished = false
	a.ErrorToken = false
	var ok bool

	for {
		var input rune
		if a.LastSymbol == 0 {
			input, _, err = a.Reader.ReadRune()
			a.LastSymbol = input
			a.Position++
		} else {
			input = a.LastSymbol
		}

		if err != nil {
			if err == io.EOF {
				a.Finished = true
				currentState, _, err = currentState.Move(0)
			} else {
				panic("Error reading next symbol " + err.Error())
			}
		} else {
			currentState, ok, err = currentState.Move(input)
			if err != nil {
				panic("Error accepting next symbol " + err.Error())
			}

			if ok {
				a.Cache.WriteRune(input)
				a.LastSymbol = 0
			}
		}

		if currentState.IsError {
			a.ErrorToken = true
		}

		if currentState.IsFinal {
			tokenData := a.Cache.String()
			a.Cache.Reset()

			if a.ErrorToken {
				a.ErrorToken = false
				return buildFullToken(errorTokenName, tokenData), a.Finished
			}

			if strings.TrimSpace(tokenData) == "" {
				return buildFullToken(emptyTokenName, tokenData), a.Finished
			}

			return buildFullToken(currentState.Name, tokenData), a.Finished
		}
	}
}
