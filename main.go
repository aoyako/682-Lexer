package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	source string // file to split

	coreDiagram string
	diagrams    []string

	initialState *State

	dictionary string
	err        error
)

func init() {
	flag.StringVar(&coreDiagram, "c", "", "Main diagram")
	flag.StringVar(&dictionary, "t", "", "Token dictionary")
	flag.Parse()
	diagrams = flag.Args()

	readTokenDictionary(dictionary)
	initialState, err = BuildStates(coreDiagram, diagrams)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	lexer := NewLexer(initialState)

	for {
		AcceptSource(lexer)
	}
}

// AcceptSource asks for a file and prints tokens from it.
func AcceptSource(l *Lexer) {
	fmt.Fscanf(os.Stdin, "%s\n", &source)
	// source := "lua/program_sample_3.lua"
	file, err := os.Open(source)
	if err != nil {
		fmt.Printf("Cannot open file: %s\n", source)
		return
	}
	defer file.Close()

	l.SetInput(file)
	PrintTokens(l.GetAllTokens())
}

// PrintTokens prints tokens with formatting.
func PrintTokens(tokens []Token, err error) {
	for _, t := range tokens {
		fmt.Println(t.String())
	}
	if err != nil {
		fmt.Printf("Received error: %s\n", err.Error())
	}
}
