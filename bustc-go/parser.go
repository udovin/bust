package main

type Node interface {
	Position() Position
	EndPosition() Position
}

type Parser struct {
	lexer *Lexer
}

func NewParser(lexer *Lexer) *Parser {
	return &Parser{
		lexer: lexer,
	}
}
