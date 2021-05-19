package main

import (
	"fmt"
	"io"
	"os"
)

var tokenDictionary map[string]string

func readTokenDictionary(filename string) error {
	tokenDictionary = make(map[string]string)
	var tokenValue string
	var tokenName string

	dict, err := os.Open(filename)
	if err != nil {
		return err
	}

	for {
		_, err := fmt.Fscanf(dict, "%s %s\n", &tokenValue, &tokenName)
		if err == io.EOF {
			break
		}
		tokenDictionary[tokenValue] = tokenName
	}

	return nil
}
