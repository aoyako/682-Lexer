package main

import (
	"fmt"
	"regexp"
)

const (
	emptySymbol = "$$e"
)

// NoTransitions error happens when there are no
// possible transitions from state.
type NoTransitions struct {
	About string
}

func (n NoTransitions) Error() string {
	return "no possible transitions: " + n.About
}

// State represents automata state.
type State struct {
	Name        string // Class name
	ID          int    // ID in class
	Transitions map[string]*State
	IsFinal     bool
	IsStart     bool
	IsError     bool
}

func (s *State) String() string {
	return fmt.Sprintf("%s-%d", s.Name, s.ID)
}

// Move tries to accept next symbol.
// In case of success returns new state, true, nil.
// In case of empty symbol transition returns new state, false, nil.
func (s *State) Move(input rune) (*State, bool, error) {
	// EOF check
	if input == 0 {
		if _, ok := s.Transitions[emptySymbol]; ok {
			return s.Transitions[emptySymbol], false, nil
		}
		return nil, false, NoTransitions{"from " + s.String()}
	}

	for k, s := range s.Transitions {
		if match(k, input) {
			return s, true, nil
		}
	}

	if _, ok := s.Transitions[emptySymbol]; ok {
		return s.Transitions[emptySymbol], false, nil
	}

	return nil, false, NoTransitions{"from " + s.String()}
}

// match returns true if input matches check string.
func match(check string, input rune) bool {
	if len([]rune(check)) < 1 {
		panic("found bad transition condition by " + check + " input " + string(input))
	}

	// Character comparison
	if len([]rune(check)) == 1 {
		return []rune(check)[0] == input
	}

	// Regexp comparison
	if check[0] == '$' && check[1] != '$' {
		matched, _ := regexp.MatchString(check[1:], string(input))
		return matched
	}

	return false
}
