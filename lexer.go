package main

import (
	"bufio"
	"fmt"
	"io"
)

// Lexer represents lexer worker.
type Lexer struct {
	Automata *Automata
}

// NewLexer constructor.
func NewLexer(state *State) *Lexer {
	return &Lexer{&Automata{
		StartState: state,
	}}
}

// GetAllTokens retireives all tokens. (Starts the lexer).
func (l *Lexer) GetAllTokens() (list []Token, e error) {
	defer func() {
		r := recover()
		if r != nil {
			e = fmt.Errorf("[%v]", r)
		}
	}()

	list = make([]Token, 0)
	end := false
	var token Token

	for !end {
		token, end = l.Automata.NextToken()
		// Ignore empty tokens
		if (token != Token{} && token.Name != emptyTokenName) {
			list = append(list, token)
		}
	}

	return
}

// SetInput sets new source file for lexer
func (l *Lexer) SetInput(input io.Reader) {
	l.Automata.Reader = bufio.NewReader(input)
	l.Automata.Finished = false
	l.Automata.Position = 0
	l.Automata.ErrorToken = false
	l.Automata.LastSymbol = 0
	l.Automata.Cache.Reset()
}
