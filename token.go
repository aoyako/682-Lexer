package main

import "fmt"

// Token represents lexer token.
type Token struct {
	Data string
	Name string
}

func buildToken(data string) Token {
	return Token{data, tokenDictionary[data]}
}

// Constructs token by name and value.
// If token value is present in dictionary, use dictionary name
func buildFullToken(name, data string) Token {
	if _, ok := tokenDictionary[data]; ok {
		return Token{data, tokenDictionary[data]}
	}
	return Token{data, name}
}

func (t *Token) String() string {
	return fmt.Sprintf("[%s]: %s", t.Data, t.Name)
}
