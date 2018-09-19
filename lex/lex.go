// Public Domain (-) 2018-present, The Peerbase Authors.
// See the Peerbase UNLICENSE file for details.

// Package lex provides an engine for building lexical scanners.
//
// It is based on the model put forward by Rob Pike in his "Lexical Scanning in
// Go" talk: https://talks.golang.org/2011/lex.slide
package lex

// EOF is returned by certain methods of the lexing Engine to signal that it has
// reached the end of the input.
const EOF = -1

// Error represents an emitted error.
type Error struct {
	Col   int
	Line  int
	Pos   int
	Type  ErrorType
	Value string
}

// ErrorType represents the type of an emitted error.
type ErrorType int

// StateFn takes the lexing Engine as a parameter and returns the next StateFn
// that will continue the lexical scanning.
type StateFn func(*Engine) StateFn

// Token represents an emitted token.
type Token struct {
	Col   int
	Line  int
	Pos   int
	Type  TokenType
	Value string
}

// TokenType represents the type of an emitted token.
type TokenType int
