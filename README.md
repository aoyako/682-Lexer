# Lexer
## Author
Lohachov Mykhailo, ISS-31

## Abstract
Usually, lexical analysis is the first phase in the compilation process.
Lexer takes the source code which consists of characters and converts it into a stream of tokens.

## Project goals
Design a program, that will accept source files and return tokens in the correct order.

## Realisation
During planning, I found convenient describing FA in DOT graph description language. 
With a slight modification I managed to create a builder for a variety of automatas.
To run the program, create description for your automata, define token dictionary, and then run:
```
go run . -c <main_automata> -t <dictionary> [<automata_i> ...]
```
Input path of desired file and press __enter__.

### Automata declaration
Same as usual dot file, but:
- In scope of the project, automata files are called ___classes___.
- ```
  <first_state> -> <second_state> [label="<transition_symbol>"]
  ```
- When in `<first_state>`, automata will check if current symbol matches `<transition_symbol>`. If so, then it will move into `<second_state>`
  - You can put regexp in `<transition_symbol>`. Just add `$` before it, like so: `$[a-Z]`
  - If you want to match an empty symbol, use `$$e`
  - If `<second_state>` is a string, you should define another automata in the file with name `<second_state>` and `<second_state>`
  diagram. Initial state of new automata will be jointed in place of `<second_state>` string.

- ```
  <state_number> [peripheries=<type>]
  ```
  `<state_number>` is a special state.
  Possible types:
    - `2` - Final state
    - `3` - Error state. Note that error state must lead to final state. All processed symbols will construct `__ERROR_TOKEN`

- Initial state in `<main_automata>` will be initial in the final automata.
- Initial state must be 0

### Dictionary declaration
Put newline-separated `<key> <value>` pairs in a dictionary file.
`<key>` is a sequence of accepted symbols, `<value>` - name of a token.
If a dictionary mapping is absent for input, token will be named as automata class.

## Examples
You can find examples in the `lua/` folder
```
go run . -c lua/main.dot -t lua/dictionary.txt lua/equals.dot lua/init.dot lua/number_constant.dot lua/string_constant_1.dot lua/string_constant_2.dot lua/identifier.dot lua/commentary.dot lua/hex_number_constant.dot
```
Then, 
```
lua/program_sample_2.lua
```