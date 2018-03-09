package main

import (
	"fmt"
	"strings"
	"unicode"
)

type tokenId string

type scanner struct {
	pos   int
	chars []rune
}

type parser struct {
	pos    int
	tokens []token
}

type token struct {
	id     tokenId
	lexeme string
	pos    int
}

type statement interface {
	error() error
	binary() string
	String() string
}

type label struct {
	err error
	val token
}

type computation struct {
	bit *token
	reg *token
	op  *token
	lhs *computation
	rhs *computation
}

type cinstruction struct {
	err  error
	dest *token
	comp computation
	jump *token
}

type ainstruction struct {
	err error
	val token
}

const (
	andToken    tokenId = "and"
	atToken     tokenId = "at"
	bitToken    tokenId = "bit"
	cparenToken tokenId = "cparen"
	eofToken    tokenId = "eof"
	eolToken    tokenId = "eol"
	eqToken     tokenId = "eq"
	errToken    tokenId = "err"
	idToken     tokenId = "id"
	minusToken  tokenId = "minus"
	notToken    tokenId = "not"
	numToken    tokenId = "num"
	oparenToken tokenId = "oparen"
	orToken     tokenId = "or"
	plusToken   tokenId = "plus"
	scolonToken tokenId = "scolon"

	eofRn    = rune(0)
	fslashRn = rune('/')
	nlRn     = rune('\n')
)

var (
	runeToks = map[rune]tokenId{
		rune('!'): notToken,
		rune('&'): andToken,
		rune('('): oparenToken,
		rune(')'): cparenToken,
		rune('+'): plusToken,
		rune('-'): minusToken,
		rune('0'): bitToken,
		rune('1'): bitToken,
		rune(';'): scolonToken,
		rune('='): eqToken,
		rune('@'): atToken,
		rune('|'): orToken,
	}

	idFn = or(unicode.IsDigit, unicode.IsLetter, is(rune('_')))
	nlFn = is(nlRn)
)

func main() {
	source := `

// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/06/max/Max.asm

// Computes R2 = max(R0, R1)  (R0,R1,R2 refer to RAM[0],RAM[1],RAM[2])

   @R0
   D=M              // D = first number
   @R1
   D=D-M            // D = first number - second number
   @OUTPUT_FIRST
   D;JGT            // if D>0 (first is greater) goto output_first
   @R1
   D=M              // D = second number
   @OUTPUT_D
   0;JMP            // goto output_d
(OUTPUT_FIRST)
   @R0
   D=M              // D = first number
(OUTPUT_D)
   @R2
   M=D              // M[2] = D (greatest number)
(INFINITE_LOOP)
   @INFINITE_LOOP
   0;JMP            // infinite loop
D=D+A
@0


`

	stmt := parse(scan(source))

	for _, s := range stmt {
		fmt.Println(s)
	}
}

/******************************************************************************
 *
 * Parser
 *
 *****************************************************************************/

func parse(tokens []token) []statement {
	p := parser{tokens: tokens}
	return p.parse()
}

func (l label) String() string {
	if l.err != nil {
		return fmt.Sprintf("(%s) // ERROR: %s", l.val.lexeme, l.err)
	} else {
		return fmt.Sprintf("(%s)", l.val.lexeme)
	}
}

func (l label) binary() string {
	return ""
}

func (l label) error() error {
	return l.err
}

func (i cinstruction) String() string {
	if i.err != nil {
		return fmt.Sprintf("    XXX // ERROR: %s", i.err)
	} else {
		return fmt.Sprintf("    XXX")
	}
}

func (i cinstruction) binary() string {
	return ""
}

func (i cinstruction) error() error {
	return i.err
}

func (i ainstruction) String() string {
	if i.err != nil {
		return fmt.Sprintf("    @%s // ERROR: %s", i.val.lexeme, i.err)
	} else {
		return fmt.Sprintf("    @%s", i.val.lexeme)
	}
}

func (i ainstruction) binary() string {
	return ""
}

func (i ainstruction) error() error {
	return i.err
}

func (p parser) parse() []statement {
	var stmt []statement

	for !p.done() {
		curr := p.curr()

		if curr.id == atToken {
			stmt = append(stmt, p.skipCheck(p.ainstruction()))
		} else if curr.id == oparenToken {
			stmt = append(stmt, p.skipCheck(p.label()))
		} else {
			stmt = append(stmt, p.skipCheck(p.cinstruction()))
		}
	}

	return stmt
}

func (p *parser) cinstruction() cinstruction {
	eq := p.peek()
	ci := cinstruction{}

	if eq.id == eqToken {
		dest, err := p.expect(idToken)
		// Eat the equals sign
		p.eat()

		if err != nil {
			ci.err = err
			return ci
		}

		ci.dest = &dest
	}

	// XXX parse computation
	for !p.done() && p.curr().id != eolToken && p.curr().id != eofToken && p.curr().id != scolonToken {
		p.eat()
	}

	if p.is(scolonToken) {
		// Eat semicolon
		p.eat()
		jump, err := p.expect(idToken)

		if err != nil {
			ci.err = err
			return ci
		}

		ci.jump = &jump
	}

	_, err := p.expect(eolToken, eofToken)

	if err != nil {
		ci.err = err
		return ci
	}

	return ci
}

func (p *parser) label() label {
	// Eat the opening paren
	p.eat()
	val, err := p.expect(idToken)

	if err != nil {
		return label{err: err}
	}

	_, err = p.expect(cparenToken)

	if err != nil {
		return label{err: err}
	}

	_, err = p.expect(eolToken, eofToken)

	if err != nil {
		return label{err: err}
	}

	return label{err: nil, val: val}
}

func (p *parser) ainstruction() ainstruction {
	// Eat the at-sign
	p.eat()
	val, err := p.expect(idToken, numToken, bitToken)

	if err != nil {
		return ainstruction{err: err}
	}

	_, err = p.expect(eolToken, eofToken)

	if err != nil {
		return ainstruction{err: err}
	}

	return ainstruction{err: nil, val: val}
}

func (p *parser) skipCheck(stmt statement) statement {
	if stmt.error() != nil {
		// Error on last parse, eat until end of current line
		for !p.done() && !p.is(eolToken, eofToken) {
			p.eat()
		}
	}

	return stmt
}

func (p *parser) is(ids ...tokenId) bool {
	curr := p.curr()

	for _, id := range ids {
		if id == curr.id {
			return true
		}
	}

	return false
}

func (p *parser) expect(ids ...tokenId) (token, error) {
	curr := p.curr()

	for _, id := range ids {
		if id == curr.id {
			p.eat()
			return curr, nil
		}
	}

	return token{}, fmt.Errorf("Expecting one of %v but found [%s] instead.",
		ids, curr.id)
}

func (p parser) peek() token {
	if p.pos+1 >= len(p.tokens) {
		return p.tokens[p.pos]
	} else {
		return p.tokens[p.pos+1]
	}
}

func (p parser) curr() token {
	return p.tokens[p.pos]
}

func (p *parser) eat() {
	p.pos += 1
}

func (p parser) done() bool {
	return p.pos >= len(p.tokens)
}

/******************************************************************************
 *
 * Scanner
 *
 *****************************************************************************/

func scan(source string) []token {
	s := scanner{chars: []rune(strings.TrimSpace(source))}
	s.chars = append(s.chars, rune(0))
	return s.scan()
}

func (s scanner) scan() []token {
	var tokens []token

	for !s.done() {
		curr := s.curr()

		if curr == nlRn {
			tokens = append(tokens, tok(eolToken, "<eol>", s.pos))
			s.eat()
		} else if curr == eofRn {
			tokens = append(tokens, tok(eofToken, "<eof>", s.pos))
			s.eat()
		} else if id, ok := runeToks[curr]; ok {
			tokens = append(tokens, tok(id, string(curr), s.pos))
			s.eat()
		} else if curr == fslashRn && s.peek() == fslashRn {
			s.takeUntil(nlFn)
			tokens = append(tokens, tok(eolToken, "<eol>", s.pos))
		} else if lexeme := s.takeWhile(idFn); len(lexeme) > 0 {
			tokens = append(tokens, tok(idToken, lexeme, s.pos))
		} else if unicode.IsSpace(curr) {
			s.takeWhile(unicode.IsSpace)
		} else {
			tokens = append(tokens, tok(errToken, string(curr), s.pos))
			s.eat()
		}
	}

	return tokens
}

func (s scanner) done() bool {
	return s.pos >= len(s.chars)
}

func (s scanner) curr() rune {
	if s.done() {
		return rune(0)
	} else {
		return s.chars[s.pos]
	}
}

func (s scanner) peek() rune {
	if s.pos+1 >= len(s.chars) {
		return rune(0)
	} else {
		return s.chars[s.pos+1]
	}
}

func (s *scanner) eat() {
	s.pos += 1
}

func (s *scanner) takeWhile(f func(rune) bool) string {
	var buff []rune

	for !s.done() {
		if !f(s.curr()) {
			break
		}

		buff = append(buff, s.curr())
		s.eat()
	}

	return string(buff)
}

func (s *scanner) takeUntil(f func(rune) bool) string {
	var buff []rune

	for !s.done() {
		if f(s.curr()) {
			s.eat()
			break
		}

		buff = append(buff, s.curr())
		s.eat()
	}

	return string(buff)
}

func or(fs ...func(rune) bool) func(rune) bool {
	return func(r rune) bool {
		for _, f := range fs {
			if f(r) {
				return true
			}
		}

		return false
	}
}

func is(v rune) func(rune) bool {
	return func(r rune) bool {
		return r == v
	}
}

func tok(id tokenId, lexeme string, pos int) token {
	return token{
		id:     id,
		lexeme: lexeme,
		pos:    pos,
	}
}
