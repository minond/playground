package main

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

type tokenid int
type segment int

type parser struct {
	tokens []token
	pos    int
}

type tokenizer struct {
	chars []rune
	pos   int
}

type token struct {
	id   tokenid
	line int
	val  *string
}

type statement interface {
	asm() []string
}

type pushStmt struct {
	seg  segment
	val  int
	line int
}

type popStmt struct {
	seg  segment
	val  int
	line int
}

type errStmt struct {
	token token
	error error
	line  int
}

type addStmt struct{ line int }
type andStmt struct{ line int }
type eqStmt struct{ line int }
type gtStmt struct{ line int }
type ltStmt struct{ line int }
type negStmt struct{ line int }
type notStmt struct{ line int }
type orStmt struct{ line int }
type subStmt struct{ line int }

const (
	argumentMem segment = iota
	constantMem
	localMem
	staticMem
	tempMem
	thatMem
	thisMem
)

const (
	errToken tokenid = iota
	eolToken
	numToken
	pushToken
	popToken
	addToken
	andToken
	argumentToken
	constantToken
	eqToken
	gtToken
	localToken
	ltToken
	negToken
	notToken
	orToken
	staticToken
	subToken
	tempToken
	thatToken
	thisToken
)

const (
	fslashRn = rune('/')
	nilRn    = rune(0)
	nlRn     = rune('\n')
)

var (
	currID = 0

	segmentsMap = map[tokenid]segment{
		argumentToken: argumentMem,
		constantToken: constantMem,
		localToken:    localMem,
		staticToken:   staticMem,
		tempToken:     tempMem,
		thatToken:     thatMem,
		thisToken:     thisMem,
	}

	tokensMap = map[string]tokenid{
		"add":      addToken,
		"and":      andToken,
		"argument": argumentToken,
		"constant": constantToken,
		"eq":       eqToken,
		"gt":       gtToken,
		"local":    localToken,
		"lt":       ltToken,
		"neg":      negToken,
		"not":      notToken,
		"or":       orToken,
		"pop":      popToken,
		"push":     pushToken,
		"static":   staticToken,
		"sub":      subToken,
		"temp":     tempToken,
		"that":     thatToken,
		"this":     thisToken,
	}

	tokensPopMem = []tokenid{
		argumentToken,
		localToken,
		staticToken,
		tempToken,
		thatToken,
		thisToken,
	}

	tokensPushMem = []tokenid{
		argumentToken,
		constantToken,
		localToken,
		staticToken,
		tempToken,
		thatToken,
		thisToken,
	}
)

func (t token) String() string {
	if t.id == eolToken {
		return "EOL"
	} else if t.val == nil || *t.val == "" {
		return fmt.Sprintf("%s on line %d", t.id, t.line)
	} else {
		return fmt.Sprintf(`%s ("%s") on line %d`, t.id, *t.val, t.line)
	}
}

func (id tokenid) String() string {
	switch id {
	case errToken:
		return "error"
	case numToken:
		return "number"
	case pushToken:
		return "push"
	case popToken:
		return "pop"
	case addToken:
		return "add"
	case andToken:
		return "and"
	case argumentToken:
		return "argument"
	case constantToken:
		return "constant"
	case eqToken:
		return "eq"
	case gtToken:
		return "gt"
	case localToken:
		return "local"
	case ltToken:
		return "lt"
	case negToken:
		return "neg"
	case notToken:
		return "not"
	case orToken:
		return "or"
	case staticToken:
		return "static"
	case subToken:
		return "sub"
	case tempToken:
		return "temp"
	case thatToken:
		return "that"
	case thisToken:
		return "this"
	case eolToken:
		return "EOL"
	default:
		panic(fmt.Sprintf("invalid token id: %d", id))
	}
}

func (s segment) String() string {
	switch s {
	case argumentMem:
		return "argument"
	case constantMem:
		return "constant"
	case localMem:
		return "local"
	case staticMem:
		return "static"
	case tempMem:
		return "temp"
	case thatMem:
		return "that"
	case thisMem:
		return "this"
	default:
		panic(fmt.Sprintf("invalid segment id: %d", s))
	}
}

func (s pushStmt) asm() []string {
	header := comment("line %03d: push %s %d", s.line, s.seg, s.val)
	switch s.seg {
	case argumentMem:
		return []string{}
	case constantMem:
		return []string{
			header,
			fmt.Sprintf("@%d", s.val),
			"D=A",
			"@SP",
			"AM=M+1",
			"A=A-1",
			"M=D",
		}
	case localMem:
		return []string{}
	case staticMem:
		return []string{}
	case tempMem:
		return []string{}
	case thatMem:
		return []string{}
	case thisMem:
		return []string{}
	default:
		panic(fmt.Sprintf("Unimplemented push %v", s))
	}
}

func (s popStmt) asm() []string {
	header := comment("line %03d: pop %s %d", s.line, s.seg, s.val)

	switch s.seg {
	case argumentMem:
		return popOp(header, "ARG", s.val)
	case localMem:
		return popOp(header, "LCL", s.val)
	case staticMem:
		return popOp(header, "STATIC", s.val)
	case tempMem:
		return popOp(header, "TEMP", s.val)
	case thatMem:
		return popOp(header, "THAT", s.val)
	case thisMem:
		return popOp(header, "THIS", s.val)
	default:
		panic(fmt.Sprintf("Unimplemented pop %v", s))
	}
}

func (s addStmt) asm() []string {
	return binOp(comment("line %03d: add", s.line), "+")
}

func (s andStmt) asm() []string {
	return binOp(comment("line %03d: and", s.line), "&")
}

func (s eqStmt) asm() []string {
	id := nextID()
	return []string{
		comment("line %03d: eq (%d)", s.line, id),
	}
}

func (s gtStmt) asm() []string {
	id := nextID()
	return []string{
		comment("line %03d: gt (%d)", s.line, id),
	}
}

func (s ltStmt) asm() []string {
	id := nextID()
	return []string{
		comment("line %03d: lt (%d)", s.line, id),
	}
}

func (s negStmt) asm() []string {
	return uniOp(comment("line %03d: neg", s.line), "-")
}

func (s notStmt) asm() []string {
	return uniOp(comment("line %03d: not", s.line), "!")
}

func (s orStmt) asm() []string {
	return binOp(comment("line %03d: or", s.line), "|")
}

func (s subStmt) asm() []string {
	return binOp(comment("line %03d: sub", s.line), "-")
}

func (s errStmt) asm() []string {
	id := nextID()
	return []string{
		comment("Error: %v, %s\n", s.error, s.token),
		fmt.Sprintf("(ERROR.%d)", id),
		fmt.Sprintf("@ERROR.%d", id),
		"0; JMP",
	}
}

func (t tokenizer) run() (tokens []token) {
	line := 1
	isNl := func(r rune) bool {
		return r == nlRn
	}

	for !t.done() {
		if t.curr() == nlRn {
			line++
			t.eat()
			continue
		} else if t.curr() == nilRn {
			t.eat()
			continue
		} else if unicode.IsSpace(t.curr()) {
			t.eat()
			continue
		} else if t.curr() == fslashRn && t.peek() == fslashRn {
			t.eatUntil(isNl)
		} else if unicode.IsDigit(t.curr()) {
			val := string(t.eatUntil(unicode.IsSpace))
			tokens = append(tokens, token{
				id:   numToken,
				val:  &val,
				line: line,
			})
		} else if unicode.IsLetter(t.curr()) {
			str := string(t.eatWhile(unicode.IsLetter))
			val := ""
			id, ok := tokensMap[str]

			if !ok {
				val = str
				id = errToken
			}

			tokens = append(tokens, token{id: id, val: &val, line: line})
		} else {
			t.eat()
		}
	}

	return tokens
}

func (t tokenizer) curr() rune {
	return t.chars[t.pos]
}

func (t tokenizer) peek() rune {
	if t.pos+1 < len(t.chars) {
		return t.chars[t.pos+1]
	}

	return rune(0)
}

func (t tokenizer) done() bool {
	return t.pos >= len(t.chars)
}

func (t *tokenizer) eat() rune {
	next := t.chars[t.pos]
	t.pos++
	return next
}

func (t *tokenizer) eatWhile(f func(rune) bool) (buff []rune) {
	for !t.done() {
		if !f(t.curr()) {
			break
		}

		buff = append(buff, t.eat())
	}

	return buff
}

func (t *tokenizer) eatUntil(f func(rune) bool) (buff []rune) {
	for !t.done() {
		if f(t.curr()) {
			break
		}

		buff = append(buff, t.eat())
	}

	return buff
}

func (p parser) run() (statements []statement, ok bool) {
	ok = true

	for !p.done() {
		switch p.eat().id {
		case pushToken:
			statements = append(statements, p.parsePushPop(true, tokensPushMem))
		case popToken:
			statements = append(statements, p.parsePushPop(false, tokensPopMem))
		case addToken:
			statements = append(statements, addStmt{p.prev().line})
		case andToken:
			statements = append(statements, addStmt{p.prev().line})
		case eqToken:
			statements = append(statements, eqStmt{p.prev().line})
		case gtToken:
			statements = append(statements, gtStmt{p.prev().line})
		case ltToken:
			statements = append(statements, ltStmt{p.prev().line})
		case negToken:
			statements = append(statements, negStmt{p.prev().line})
		case notToken:
			statements = append(statements, notStmt{p.prev().line})
		case orToken:
			statements = append(statements, orStmt{p.prev().line})
		case subToken:
			statements = append(statements, subStmt{p.prev().line})

		case errToken:
			line := p.prev().line
			p.eatLine()
			ok = false
			statements = append(statements, errStmt{
				token: p.prev(),
				error: errors.New("invalid token"),
				line:  line,
			})

		default:
			line := p.prev().line
			p.eatLine()
			ok = false
			statements = append(statements, errStmt{
				token: p.prev(),
				error: errors.New("unexpected token"),
				line:  line,
			})
		}
	}

	return statements, ok
}

func (p *parser) parsePushPop(isPush bool, memTokens []tokenid) statement {
	segTok, err := p.expect(memTokens...)

	if err != nil {
		p.eatLine()
		return errStmt{
			token: p.curr(),
			error: err,
		}
	}

	str, err := p.expect(numToken)

	if err != nil {
		p.eatLine()
		return errStmt{
			token: p.curr(),
			error: err,
		}
	}

	if str.val == nil {
		p.eatLine()
		return errStmt{
			token: str,
			error: errors.New("expecting a number value"),
		}
	}

	num, err := strconv.Atoi(*str.val)

	if err != nil {
		p.eatLine()
		return errStmt{
			token: str,
			error: fmt.Errorf("unable to convert %s to a number", *str.val),
		}
	}

	seg, ok := segmentsMap[segTok.id]

	if !ok {
		p.eatLine()
		return errStmt{
			token: segTok,
			error: fmt.Errorf("expecting %v but found [%s] instead",
				tokensPushMem, segTok.id),
		}
	}

	if isPush {
		return pushStmt{
			seg:  seg,
			val:  num,
			line: segTok.line,
		}
	}

	return popStmt{
		seg:  seg,
		val:  num,
		line: segTok.line,
	}
}

func (p parser) done() bool {
	return p.pos >= len(p.tokens)
}

func (p *parser) eat() token {
	if p.done() {
		return token{eolToken, -1, nil}
	}

	next := p.tokens[p.pos]
	p.pos++
	return next
}

func (p *parser) eatLine() {
	if p.done() {
		return
	}
	line := p.eat().line
	for line != 0 && !p.done() && p.curr().line == line {
		p.eat()
	}
}

func (p parser) curr() token {
	if p.done() {
		return token{eolToken, -1, nil}
	}

	return p.tokens[p.pos]
}

func (p parser) prev() token {
	return p.tokens[p.pos-1]
}

func (p *parser) expect(ids ...tokenid) (token, error) {
	curr := p.curr()
	for _, id := range ids {
		if curr.id == id {
			p.eat()
			return curr, nil
		}
	}

	return token{}, fmt.Errorf("expecting %v but found [%s] instead",
		ids, curr.id)
}

func tokenize(src string) []token {
	toks := tokenizer{chars: []rune(src)}
	return toks.run()
}

func parse(tokens []token) ([]statement, bool) {
	parse := parser{tokens: tokens}
	return parse.run()
}

func compile(stmts []statement) (code []string) {
	for _, stmt := range stmts {
		code = append(code, stmt.asm()...)
	}
	return
}

func nextID() int {
	currID++
	return currID
}

func binOp(header, op string) []string {
	return []string{
		header,
		"@SP",                     // Load the SP
		"AM=M-1",                  // Update the SP = SP-1 and point to SP-1
		"D=M",                     // Store SP-1 value in D
		"AM=A-1",                  // Update the SP = A-1 and point to SP-1
		fmt.Sprintf("M=D%sM", op), // Run operation on D and M, which is SP-2
	}
}

func uniOp(header, op string) []string {
	return []string{
		header,
		"@SP",                    // Load the SP
		"A=M-1",                  // Point to SP-1 but do not update SP
		fmt.Sprintf("M=%sM", op), // Run operation on M, which is SP-1
	}
}

func popOp(header, seg string, offset int) []string {
	return []string{
		header,

		fmt.Sprintf("@%s", seg), // Load the segment
		"D=M", // Store the start of that segment's address in D
		fmt.Sprintf("@%d", offset), // Load the offset
		"D=D+A",                    // Store the start + offset of address in D
		"@R13",
		"M=D",    // Store that address in R13
		"@SP",    // Load the SP
		"AM=M-1", // Update the SP = SP-1 and point to SP-1
		"D=M",    // Store value of SP-1 in D
		"@R13",
		"A=M", // Pointer to address previously stored in R13, which is start + offset
		"M=D", // Set that value in memory to D which is SP-1
	}
}

func comment(str string, args ...interface{}) string {
	return fmt.Sprintf("// "+str, args...)
}

func main() {
	sample := `
push constant 7
push constant 8
add
`

	statements, _ := parse(tokenize(sample))
	for _, line := range compile(statements) {
		fmt.Println(line)
	}
}
