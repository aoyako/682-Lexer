package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	errorState = "3"
	finalState = "2"
)

var (
	connectionRegex  = regexp.MustCompile(`[0-9]+->[[:alnum:]\_]+\[label=".*?"\]`)
	declarationRegex = regexp.MustCompile(`[0-9]+\[peripheries=[0-9]\]`)
)

// DependentItem represents dependence on other class
type DependentItem struct {
	Source          *State
	Transition      string
	ReferencedClass string
}

// CompiledClass represents set of connected states
type CompiledClass struct {
	Start      *State
	References []DependentItem
}

// BuildStates builds input diagrams and returns start state
func BuildStates(core string, files []string) (*State, error) {
	classes := make(map[string]*CompiledClass)

	for _, filename := range files {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		class, err := CompileClass(getClassName(filename), string(data))
		if err != nil {
			return nil, err
		}
		classes[getClassName(filename)] = class
	}

	data, err := ioutil.ReadFile(core)
	if err != nil {
		return nil, err
	}

	class, err := CompileClass(getClassName(core), string(data))
	if err != nil {
		return nil, err
	}
	classes[getClassName(core)] = class

	return LinkClasses(getClassName(core), classes)
}

// LinkClasses resolves class dependencies
func LinkClasses(main string, classes map[string]*CompiledClass) (*State, error) {
	for _, class := range classes {
		for _, item := range class.References {
			if _, ok := classes[item.ReferencedClass]; !ok {
				return nil, fmt.Errorf("Cannot find class with name: %s", item.ReferencedClass)
			}
			item.Source.Transitions[item.Transition] = classes[item.ReferencedClass].Start
		}
	}

	return classes[main].Start, nil
}

// CompileClass compiles single class
func CompileClass(className, data string) (*CompiledClass, error) {
	var first int
	var second int
	var condition string
	var err error

	var references []DependentItem
	states := make(map[int]*State)

	var unprocessedCommands []string

	data = strings.ReplaceAll(data, " ", "")
	data = strings.ReplaceAll(data, "\n", "")
	data = strings.ReplaceAll(data, "\t", "")
	data = strings.Trim(data, "}")
	data = data[strings.Index(data, "{")+1:]

	commands := connectionRegex.FindAllString(data, -1)
	unprocessedCommands = declarationRegex.FindAllString(data, -1)
	for _, c := range commands {
		c = strings.TrimSuffix(c, "]")
		details := strings.SplitN(c, "[", 2)

		directionDetails := strings.Split(details[0], "->")

		conditionDetails := strings.Split(details[1], "label=")
		conditionDetails[1] = strings.TrimPrefix(conditionDetails[1], "\"")
		condition = strings.TrimSuffix(conditionDetails[1], "\"")

		first, err = strconv.Atoi(directionDetails[0])
		if err != nil {
			return nil, err
		}
		if _, ok := states[first]; !ok {
			states[first] = &State{
				Name:        className,
				ID:          first,
				Transitions: make(map[string]*State),
			}
			if first == 0 {
				states[first].IsStart = true
			}
		}

		second, err = strconv.Atoi(directionDetails[1])
		if err != nil {
			references = append(references, DependentItem{
				Source:          states[first],
				Transition:      condition,
				ReferencedClass: directionDetails[1],
			})
		} else {
			if _, ok := states[second]; !ok {
				states[second] = &State{
					Name:        className,
					ID:          second,
					Transitions: make(map[string]*State),
				}
				if second == 0 {
					states[second].IsStart = true
				}
			}

			states[first].Transitions[condition] = states[second]
		}
	}

	for _, c := range unprocessedCommands {
		c = strings.TrimSuffix(c, "]")
		details := strings.Split(c, "[")
		target, err := strconv.Atoi(details[0])
		if err != nil {
			return nil, err
		}

		stateType := strings.Split(details[1], "=")[1]

		switch stateType {
		case finalState:
			states[target].IsFinal = true
		case errorState:
			states[target].IsError = true
		}
	}

	for _, v := range states {
		if v.IsStart {
			return &CompiledClass{v, references}, nil
		}
	}

	return nil, fmt.Errorf("Start state not found for %s class", className)
}

func getClassName(filename string) string {
	name := filepath.Base(filename)
	return strings.TrimSuffix(name, filepath.Ext(name))
}
