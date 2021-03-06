// Copyright 2019 The SQLFlow Authors. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by goyacc -p sql -o extended_syntax_parser.go sql.y. DO NOT EDIT.

//line sql.y:2
package sql

import __yyfmt__ "fmt"

//line sql.y:2

import (
	"fmt"
	"strings"
	"sync"
)

/* expr defines an expression as a Lisp list.  If len(val)>0,
   it is an atomic expression, in particular, NUMBER, IDENT,
   or STRING, defined by typ and val; otherwise, it is a
   Lisp S-expression. */
type expr struct {
	typ  int
	val  string
	sexp exprlist
}

type exprlist []*expr

/* construct an atomic expr */
func atomic(typ int, val string) *expr {
	return &expr{
		typ: typ,
		val: val,
	}
}

/* construct a funcall expr */
func funcall(name string, oprd exprlist) *expr {
	return &expr{
		sexp: append(exprlist{atomic(IDENT, name)}, oprd...),
	}
}

/* construct a unary expr */
func unary(typ int, op string, od1 *expr) *expr {
	return &expr{
		sexp: append(exprlist{atomic(typ, op)}, od1),
	}
}

/* construct a binary expr */
func binary(typ int, od1 *expr, op string, od2 *expr) *expr {
	return &expr{
		sexp: append(exprlist{atomic(typ, op)}, od1, od2),
	}
}

/* construct a variadic expr */
func variadic(typ int, op string, ods exprlist) *expr {
	return &expr{
		sexp: append(exprlist{atomic(typ, op)}, ods...),
	}
}

type extendedSelect struct {
	extended bool
	train    bool
	analyze  bool
	standardSelect
	trainClause
	predictClause
	explainClause
}

type standardSelect struct {
	fields exprlist
	tables []string
	where  *expr
	limit  string
	origin string
}

type trainClause struct {
	estimator  string
	trainAttrs attrs
	columns    columnClause
	label      string
	save       string
}

/* If no FOR in the COLUMN, the key is "" */
type columnClause map[string]exprlist
type fieldClause exprlist

type attrs map[string]*expr

type predictClause struct {
	predAttrs attrs
	model     string
	// FIXME(tony): rename into to predTable
	into string
}

type explainClause struct {
	explainAttrs attrs
	trainedModel string
	explainer    string
}

var parseResult *extendedSelect

func attrsUnion(as1, as2 attrs) attrs {
	for k, v := range as2 {
		if _, ok := as1[k]; ok {
			log.Panicf("attr %q already specified", as2)
		}
		as1[k] = v
	}
	return as1
}

//line sql.y:116
type sqlSymType struct {
	yys   int
	val   string /* NUMBER, IDENT, STRING, and keywords */
	flds  exprlist
	tbls  []string
	expr  *expr
	expl  exprlist
	atrs  attrs
	eslt  extendedSelect
	slct  standardSelect
	tran  trainClause
	colc  columnClause
	labc  string
	infr  predictClause
	expln explainClause
}

const SELECT = 57346
const FROM = 57347
const WHERE = 57348
const LIMIT = 57349
const TRAIN = 57350
const PREDICT = 57351
const EXPLAIN = 57352
const WITH = 57353
const COLUMN = 57354
const LABEL = 57355
const USING = 57356
const INTO = 57357
const FOR = 57358
const AS = 57359
const TO = 57360
const IDENT = 57361
const NUMBER = 57362
const STRING = 57363
const AND = 57364
const OR = 57365
const GE = 57366
const LE = 57367
const NE = 57368
const NOT = 57369
const POWER = 57370
const UMINUS = 57371

var sqlToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"SELECT",
	"FROM",
	"WHERE",
	"LIMIT",
	"TRAIN",
	"PREDICT",
	"EXPLAIN",
	"WITH",
	"COLUMN",
	"LABEL",
	"USING",
	"INTO",
	"FOR",
	"AS",
	"TO",
	"IDENT",
	"NUMBER",
	"STRING",
	"AND",
	"OR",
	"'>'",
	"'<'",
	"'='",
	"'!'",
	"GE",
	"LE",
	"NE",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"'%'",
	"NOT",
	"POWER",
	"UMINUS",
	"';'",
	"'('",
	"')'",
	"','",
	"'['",
	"']'",
	"'\"'",
	"'\\''",
}
var sqlStatenames = [...]string{}

const sqlEofCode = 1
const sqlErrCode = 2
const sqlInitialStackSize = 16

//line sql.y:350

/* Like Lisp's builtin function cdr. */
func (e *expr) cdr() (r []string) {
	for i := 1; i < len(e.sexp); i++ {
		r = append(r, e.sexp[i].String())
	}
	return r
}

/* Convert exprlist to string slice. */
func (el exprlist) Strings() (r []string) {
	for i := 0; i < len(el); i++ {
		r = append(r, el[i].String())
	}
	return r
}

func (e *expr) String() string {
	if e.typ == 0 { /* a compound expression */
		switch e.sexp[0].typ {
		case '+', '*', '/', '%', '=', '<', '>', '!', LE, GE, AND, OR:
			if len(e.sexp) != 3 {
				log.Panicf("Expecting binary expression, got %.10q", e.sexp)
			}
			return fmt.Sprintf("%s %s %s", e.sexp[1], e.sexp[0].val, e.sexp[2])
		case '-':
			switch len(e.sexp) {
			case 2:
				return fmt.Sprintf(" -%s", e.sexp[1])
			case 3:
				return fmt.Sprintf("%s - %s", e.sexp[1], e.sexp[2])
			default:
				log.Panicf("Expecting either unary or binary -, got %.10q", e.sexp)
			}
		case '(':
			if len(e.sexp) != 2 {
				log.Panicf("Expecting ( ) as unary operator, got %.10q", e.sexp)
			}
			return fmt.Sprintf("(%s)", e.sexp[1])
		case '[':
			return "[" + strings.Join(e.cdr(), ", ") + "]"
		case NOT:
			return fmt.Sprintf("NOT %s", e.sexp[1])
		case IDENT: /* function call */
			return e.sexp[0].val + "(" + strings.Join(e.cdr(), ", ") + ")"
		}
	} else {
		return fmt.Sprintf("%s", e.val)
	}

	log.Panicf("Cannot print an unknown expression")
	return ""
}

func (s standardSelect) String() string {
	if s.origin != "" {
		return s.origin
	}

	r := "SELECT "
	if len(s.fields) == 0 {
		r += "*"
	} else {
		for i := 0; i < len(s.fields); i++ {
			r += s.fields[i].String()
			if i != len(s.fields)-1 {
				r += ", "
			}
		}
	}
	r += "\nFROM " + strings.Join(s.tables, ", ")
	if s.where != nil {
		r += fmt.Sprintf("\nWHERE %s", s.where)
	}
	if len(s.limit) > 0 {
		r += fmt.Sprintf("\nLIMIT %s", s.limit)
	}
	return r
}

// sqlReentrantParser makes sqlParser, generated by goyacc and using a
// global variable parseResult to return the result, reentrant.
type extendedSyntaxParser struct {
	pr sqlParser
}

func newExtendedSyntaxParser() *extendedSyntaxParser {
	return &extendedSyntaxParser{sqlNewParser()}
}

var mu sync.Mutex

func (p *extendedSyntaxParser) Parse(s string) (r *extendedSelect, e error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			e, ok = r.(error)
			if !ok {
				e = fmt.Errorf("%v", r)
			}
		}
	}()

	mu.Lock()
	defer mu.Unlock()

	p.pr.Parse(newLexer(s))
	return parseResult, nil
}

//line yacctab:1
var sqlExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const sqlPrivate = 57344

const sqlLast = 205

var sqlAct = [...]int{

	47, 127, 39, 126, 83, 113, 82, 16, 41, 40,
	42, 110, 38, 61, 109, 111, 115, 116, 143, 140,
	58, 49, 119, 93, 61, 48, 118, 60, 61, 44,
	28, 29, 50, 80, 45, 46, 35, 7, 62, 63,
	64, 65, 66, 117, 141, 141, 114, 75, 59, 25,
	114, 78, 79, 77, 114, 24, 57, 23, 8, 14,
	86, 92, 88, 81, 94, 95, 96, 97, 98, 99,
	100, 101, 102, 103, 104, 105, 106, 107, 41, 40,
	42, 13, 12, 64, 65, 66, 131, 129, 132, 18,
	76, 49, 120, 146, 144, 48, 142, 139, 137, 44,
	37, 128, 50, 19, 45, 46, 135, 134, 41, 40,
	42, 84, 91, 87, 85, 36, 130, 34, 121, 125,
	133, 49, 32, 31, 130, 48, 30, 138, 27, 44,
	123, 116, 50, 122, 45, 46, 136, 55, 124, 51,
	54, 90, 130, 145, 73, 74, 69, 68, 67, 26,
	71, 70, 72, 62, 63, 64, 65, 66, 73, 74,
	69, 68, 67, 108, 71, 70, 72, 62, 63, 64,
	65, 66, 69, 68, 67, 6, 71, 70, 72, 62,
	63, 64, 65, 66, 53, 5, 15, 52, 11, 7,
	20, 21, 22, 4, 3, 43, 10, 9, 56, 33,
	17, 112, 89, 2, 1,
}
var sqlPact = [...]int{

	171, -1000, 19, 43, 42, 20, 70, 182, -1000, 18,
	16, 10, -1000, -1000, -1000, 144, 111, -12, -9, -1000,
	107, 104, 103, -1000, -1000, -1000, 98, -4, 96, 59,
	128, 173, 126, 14, -1000, 89, -1000, -1000, -14, 136,
	-1000, -9, -1000, -1000, 89, 69, 32, -1000, 89, 89,
	-11, 92, 95, 92, 94, 92, 134, 93, 89, -18,
	-1000, 89, 89, 89, 89, 89, 89, 89, 89, 89,
	89, 89, 89, 89, 89, 122, -31, -35, -1000, -1000,
	-1000, -29, 4, -1000, 17, -1000, 12, -1000, 8, -1000,
	72, -1000, 136, -1000, 136, 50, 50, -1000, -1000, -1000,
	7, 7, 7, 7, 7, 7, 148, 148, -1000, -1000,
	-1000, -1000, 118, 123, 92, 68, 67, 89, 88, 87,
	-1000, 121, 79, 68, 78, -1000, 3, -1000, -1000, -9,
	-1000, -1000, -1000, 136, -1000, -1000, 77, -1000, 2, -1000,
	75, 68, -1000, 74, -1000, -1000, -1000,
}
var sqlPgo = [...]int{

	0, 204, 203, 202, 194, 201, 5, 193, 185, 200,
	199, 2, 0, 1, 198, 12, 195, 3, 186, 4,
	6,
}
var sqlR1 = [...]int{

	0, 1, 1, 1, 1, 1, 1, 1, 2, 14,
	14, 3, 3, 4, 4, 4, 7, 7, 8, 8,
	5, 5, 5, 18, 18, 9, 9, 9, 13, 13,
	13, 17, 17, 6, 6, 10, 10, 19, 20, 20,
	12, 12, 15, 15, 16, 16, 11, 11, 11, 11,
	11, 11, 11, 11, 11, 11, 11, 11, 11, 11,
	11, 11, 11, 11, 11, 11, 11, 11, 11,
}
var sqlR2 = [...]int{

	0, 2, 3, 3, 3, 2, 2, 2, 6, 0,
	2, 0, 2, 9, 8, 8, 5, 7, 5, 7,
	2, 4, 5, 5, 1, 1, 1, 3, 1, 1,
	1, 1, 3, 2, 2, 1, 3, 3, 1, 3,
	3, 4, 1, 3, 2, 3, 1, 1, 1, 1,
	3, 3, 3, 1, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 2, 2,
}
var sqlChk = [...]int{

	-1000, -1, -2, -4, -7, -8, 4, 18, 39, -4,
	-7, -8, 39, 39, 39, -18, -12, -9, 19, 33,
	8, 9, 10, 39, 39, 39, 5, 17, 42, 40,
	19, 19, 19, -10, 19, 40, 19, 41, -15, -11,
	20, 19, 21, -16, 40, 45, 46, -12, 36, 32,
	43, 11, 14, 11, 14, 11, -14, 42, 6, -15,
	41, 42, 31, 32, 33, 34, 35, 26, 25, 24,
	29, 28, 30, 22, 23, -11, 21, 21, -11, -11,
	44, -15, -20, -19, 19, 19, -20, 19, -20, -3,
	7, 19, -11, 41, -11, -11, -11, -11, -11, -11,
	-11, -11, -11, -11, -11, -11, -11, -11, 41, 45,
	46, 44, -5, -6, 42, 12, 13, 26, 14, 14,
	20, -6, 15, 12, 15, -19, -17, -13, 33, 19,
	-12, 19, 21, -11, 19, 19, 15, 19, -17, 19,
	16, 42, 19, 16, 19, -13, 19,
}
var sqlDef = [...]int{

	0, -2, 0, 0, 0, 0, 0, 0, 1, 0,
	0, 0, 5, 6, 7, 0, 0, 24, 26, 25,
	0, 0, 0, 2, 3, 4, 0, 0, 0, 0,
	0, 0, 0, 9, 35, 0, 27, 40, 0, 42,
	46, 47, 48, 49, 0, 0, 0, 53, 0, 0,
	0, 0, 0, 0, 0, 0, 11, 0, 0, 0,
	41, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 67, 68,
	44, 0, 0, 38, 0, 16, 0, 18, 0, 8,
	0, 36, 10, 23, 43, 54, 55, 56, 57, 58,
	59, 60, 61, 62, 63, 64, 65, 66, 50, 51,
	52, 45, 0, 0, 0, 0, 0, 0, 0, 0,
	12, 0, 0, 0, 0, 39, 20, 31, 28, 29,
	30, 33, 34, 37, 17, 19, 0, 14, 0, 15,
	0, 0, 13, 0, 21, 32, 22,
}
var sqlTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 27, 45, 3, 3, 35, 3, 46,
	40, 41, 33, 31, 42, 32, 3, 34, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 39,
	25, 26, 24, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 43, 3, 44,
}
var sqlTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 28, 29, 30, 36, 37, 38,
}
var sqlTok3 = [...]int{
	0,
}

var sqlErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	sqlDebug        = 0
	sqlErrorVerbose = false
)

type sqlLexer interface {
	Lex(lval *sqlSymType) int
	Error(s string)
}

type sqlParser interface {
	Parse(sqlLexer) int
	Lookahead() int
}

type sqlParserImpl struct {
	lval  sqlSymType
	stack [sqlInitialStackSize]sqlSymType
	char  int
}

func (p *sqlParserImpl) Lookahead() int {
	return p.char
}

func sqlNewParser() sqlParser {
	return &sqlParserImpl{}
}

const sqlFlag = -1000

func sqlTokname(c int) string {
	if c >= 1 && c-1 < len(sqlToknames) {
		if sqlToknames[c-1] != "" {
			return sqlToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func sqlStatname(s int) string {
	if s >= 0 && s < len(sqlStatenames) {
		if sqlStatenames[s] != "" {
			return sqlStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func sqlErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !sqlErrorVerbose {
		return "syntax error"
	}

	for _, e := range sqlErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + sqlTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := sqlPact[state]
	for tok := TOKSTART; tok-1 < len(sqlToknames); tok++ {
		if n := base + tok; n >= 0 && n < sqlLast && sqlChk[sqlAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if sqlDef[state] == -2 {
		i := 0
		for sqlExca[i] != -1 || sqlExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; sqlExca[i] >= 0; i += 2 {
			tok := sqlExca[i]
			if tok < TOKSTART || sqlExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if sqlExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += sqlTokname(tok)
	}
	return res
}

func sqllex1(lex sqlLexer, lval *sqlSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = sqlTok1[0]
		goto out
	}
	if char < len(sqlTok1) {
		token = sqlTok1[char]
		goto out
	}
	if char >= sqlPrivate {
		if char < sqlPrivate+len(sqlTok2) {
			token = sqlTok2[char-sqlPrivate]
			goto out
		}
	}
	for i := 0; i < len(sqlTok3); i += 2 {
		token = sqlTok3[i+0]
		if token == char {
			token = sqlTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = sqlTok2[1] /* unknown char */
	}
	if sqlDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", sqlTokname(token), uint(char))
	}
	return char, token
}

func sqlParse(sqllex sqlLexer) int {
	return sqlNewParser().Parse(sqllex)
}

func (sqlrcvr *sqlParserImpl) Parse(sqllex sqlLexer) int {
	var sqln int
	var sqlVAL sqlSymType
	var sqlDollar []sqlSymType
	_ = sqlDollar // silence set and not used
	sqlS := sqlrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	sqlstate := 0
	sqlrcvr.char = -1
	sqltoken := -1 // sqlrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		sqlstate = -1
		sqlrcvr.char = -1
		sqltoken = -1
	}()
	sqlp := -1
	goto sqlstack

ret0:
	return 0

ret1:
	return 1

sqlstack:
	/* put a state and value onto the stack */
	if sqlDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", sqlTokname(sqltoken), sqlStatname(sqlstate))
	}

	sqlp++
	if sqlp >= len(sqlS) {
		nyys := make([]sqlSymType, len(sqlS)*2)
		copy(nyys, sqlS)
		sqlS = nyys
	}
	sqlS[sqlp] = sqlVAL
	sqlS[sqlp].yys = sqlstate

sqlnewstate:
	sqln = sqlPact[sqlstate]
	if sqln <= sqlFlag {
		goto sqldefault /* simple state */
	}
	if sqlrcvr.char < 0 {
		sqlrcvr.char, sqltoken = sqllex1(sqllex, &sqlrcvr.lval)
	}
	sqln += sqltoken
	if sqln < 0 || sqln >= sqlLast {
		goto sqldefault
	}
	sqln = sqlAct[sqln]
	if sqlChk[sqln] == sqltoken { /* valid shift */
		sqlrcvr.char = -1
		sqltoken = -1
		sqlVAL = sqlrcvr.lval
		sqlstate = sqln
		if Errflag > 0 {
			Errflag--
		}
		goto sqlstack
	}

sqldefault:
	/* default state action */
	sqln = sqlDef[sqlstate]
	if sqln == -2 {
		if sqlrcvr.char < 0 {
			sqlrcvr.char, sqltoken = sqllex1(sqllex, &sqlrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if sqlExca[xi+0] == -1 && sqlExca[xi+1] == sqlstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			sqln = sqlExca[xi+0]
			if sqln < 0 || sqln == sqltoken {
				break
			}
		}
		sqln = sqlExca[xi+1]
		if sqln < 0 {
			goto ret0
		}
	}
	if sqln == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			sqllex.Error(sqlErrorMessage(sqlstate, sqltoken))
			Nerrs++
			if sqlDebug >= 1 {
				__yyfmt__.Printf("%s", sqlStatname(sqlstate))
				__yyfmt__.Printf(" saw %s\n", sqlTokname(sqltoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for sqlp >= 0 {
				sqln = sqlPact[sqlS[sqlp].yys] + sqlErrCode
				if sqln >= 0 && sqln < sqlLast {
					sqlstate = sqlAct[sqln] /* simulate a shift of "error" */
					if sqlChk[sqlstate] == sqlErrCode {
						goto sqlstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if sqlDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", sqlS[sqlp].yys)
				}
				sqlp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if sqlDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", sqlTokname(sqltoken))
			}
			if sqltoken == sqlEofCode {
				goto ret1
			}
			sqlrcvr.char = -1
			sqltoken = -1
			goto sqlnewstate /* try again in the same state */
		}
	}

	/* reduction by production sqln */
	if sqlDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", sqln, sqlStatname(sqlstate))
	}

	sqlnt := sqln
	sqlpt := sqlp
	_ = sqlpt // guard against "declared and not used"

	sqlp -= sqlR2[sqln]
	// sqlp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if sqlp+1 >= len(sqlS) {
		nyys := make([]sqlSymType, len(sqlS)*2)
		copy(nyys, sqlS)
		sqlS = nyys
	}
	sqlVAL = sqlS[sqlp+1]

	/* consult goto table to find next state */
	sqln = sqlR1[sqln]
	sqlg := sqlPgo[sqln]
	sqlj := sqlg + sqlS[sqlp].yys + 1

	if sqlj >= sqlLast {
		sqlstate = sqlAct[sqlg]
	} else {
		sqlstate = sqlAct[sqlj]
		if sqlChk[sqlstate] != -sqln {
			sqlstate = sqlAct[sqlg]
		}
	}
	// dummy call; replaced with literal code
	switch sqlnt {

	case 1:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:161
		{
			parseResult = &extendedSelect{
				extended:       false,
				standardSelect: sqlDollar[1].slct}
		}
	case 2:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:166
		{
			parseResult = &extendedSelect{
				extended:       true,
				train:          true,
				standardSelect: sqlDollar[1].slct,
				trainClause:    sqlDollar[2].tran}
		}
	case 3:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:173
		{
			parseResult = &extendedSelect{
				extended:       true,
				train:          false,
				standardSelect: sqlDollar[1].slct,
				predictClause:  sqlDollar[2].infr}
		}
	case 4:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:180
		{
			parseResult = &extendedSelect{
				extended:       true,
				train:          false,
				analyze:        true,
				standardSelect: sqlDollar[1].slct,
				explainClause:  sqlDollar[2].expln}
		}
	case 5:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:188
		{ // FIXME(tony): remove above rules that include select clause
			parseResult = &extendedSelect{
				extended:    true,
				train:       true,
				trainClause: sqlDollar[1].tran}
		}
	case 6:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:194
		{
			parseResult = &extendedSelect{
				extended:      true,
				train:         false,
				predictClause: sqlDollar[1].infr}
		}
	case 7:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:200
		{
			parseResult = &extendedSelect{
				extended:      true,
				train:         false,
				analyze:       true,
				explainClause: sqlDollar[1].expln}
		}
	case 8:
		sqlDollar = sqlS[sqlpt-6 : sqlpt+1]
//line sql.y:210
		{
			sqlVAL.slct.fields = sqlDollar[2].expl
			sqlVAL.slct.tables = sqlDollar[4].tbls
			sqlVAL.slct.where = sqlDollar[5].expr
			sqlVAL.slct.limit = sqlDollar[6].val
		}
	case 9:
		sqlDollar = sqlS[sqlpt-0 : sqlpt+1]
//line sql.y:219
		{
		}
	case 10:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:220
		{
			sqlVAL.expr = sqlDollar[2].expr
		}
	case 11:
		sqlDollar = sqlS[sqlpt-0 : sqlpt+1]
//line sql.y:224
		{
		}
	case 12:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:225
		{
			sqlVAL.val = sqlDollar[2].val
		}
	case 13:
		sqlDollar = sqlS[sqlpt-9 : sqlpt+1]
//line sql.y:229
		{
			sqlVAL.tran.estimator = sqlDollar[3].val
			sqlVAL.tran.trainAttrs = sqlDollar[5].atrs
			sqlVAL.tran.columns = sqlDollar[6].colc
			sqlVAL.tran.label = sqlDollar[7].labc
			sqlVAL.tran.save = sqlDollar[9].val
		}
	case 14:
		sqlDollar = sqlS[sqlpt-8 : sqlpt+1]
//line sql.y:236
		{
			sqlVAL.tran.estimator = sqlDollar[3].val
			sqlVAL.tran.trainAttrs = sqlDollar[5].atrs
			sqlVAL.tran.columns = sqlDollar[6].colc
			sqlVAL.tran.save = sqlDollar[8].val
		}
	case 15:
		sqlDollar = sqlS[sqlpt-8 : sqlpt+1]
//line sql.y:242
		{
			sqlVAL.tran.estimator = sqlDollar[3].val
			sqlVAL.tran.trainAttrs = sqlDollar[5].atrs
			sqlVAL.tran.label = sqlDollar[6].labc
			sqlVAL.tran.save = sqlDollar[8].val
		}
	case 16:
		sqlDollar = sqlS[sqlpt-5 : sqlpt+1]
//line sql.y:251
		{
			sqlVAL.infr.into = sqlDollar[3].val
			sqlVAL.infr.model = sqlDollar[5].val
		}
	case 17:
		sqlDollar = sqlS[sqlpt-7 : sqlpt+1]
//line sql.y:252
		{
			sqlVAL.infr.into = sqlDollar[3].val
			sqlVAL.infr.predAttrs = sqlDollar[5].atrs
			sqlVAL.infr.model = sqlDollar[7].val
		}
	case 18:
		sqlDollar = sqlS[sqlpt-5 : sqlpt+1]
//line sql.y:256
		{
			sqlVAL.expln.trainedModel = sqlDollar[3].val
			sqlVAL.expln.explainer = sqlDollar[5].val
		}
	case 19:
		sqlDollar = sqlS[sqlpt-7 : sqlpt+1]
//line sql.y:257
		{
			sqlVAL.expln.trainedModel = sqlDollar[3].val
			sqlVAL.expln.explainAttrs = sqlDollar[5].atrs
			sqlVAL.expln.explainer = sqlDollar[7].val
		}
	case 20:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:261
		{
			sqlVAL.colc = map[string]exprlist{"feature_columns": sqlDollar[2].expl}
		}
	case 21:
		sqlDollar = sqlS[sqlpt-4 : sqlpt+1]
//line sql.y:262
		{
			sqlVAL.colc = map[string]exprlist{sqlDollar[4].val: sqlDollar[2].expl}
		}
	case 22:
		sqlDollar = sqlS[sqlpt-5 : sqlpt+1]
//line sql.y:263
		{
			sqlVAL.colc[sqlDollar[5].val] = sqlDollar[3].expl
		}
	case 23:
		sqlDollar = sqlS[sqlpt-5 : sqlpt+1]
//line sql.y:267
		{
			sqlVAL.expl = exprlist{sqlDollar[1].expr, atomic(IDENT, "AS"), funcall("", sqlDollar[4].expl)}
		}
	case 24:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:270
		{
			sqlVAL.expl = sqlDollar[1].flds
		}
	case 25:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:274
		{
			sqlVAL.flds = append(sqlVAL.flds, atomic(IDENT, "*"))
		}
	case 26:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:275
		{
			sqlVAL.flds = append(sqlVAL.flds, atomic(IDENT, sqlDollar[1].val))
		}
	case 27:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:276
		{
			sqlVAL.flds = append(sqlDollar[1].flds, atomic(IDENT, sqlDollar[3].val))
		}
	case 28:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:280
		{
			sqlVAL.expr = atomic(IDENT, "*")
		}
	case 29:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:281
		{
			sqlVAL.expr = atomic(IDENT, sqlDollar[1].val)
		}
	case 30:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:282
		{
			sqlVAL.expr = sqlDollar[1].expr
		}
	case 31:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:286
		{
			sqlVAL.expl = exprlist{sqlDollar[1].expr}
		}
	case 32:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:287
		{
			sqlVAL.expl = append(sqlDollar[1].expl, sqlDollar[3].expr)
		}
	case 33:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:291
		{
			sqlVAL.labc = sqlDollar[2].val
		}
	case 34:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:292
		{
			sqlVAL.labc = sqlDollar[2].val[1 : len(sqlDollar[2].val)-1]
		}
	case 35:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:296
		{
			sqlVAL.tbls = []string{sqlDollar[1].val}
		}
	case 36:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:297
		{
			sqlVAL.tbls = append(sqlDollar[1].tbls, sqlDollar[3].val)
		}
	case 37:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:301
		{
			sqlVAL.atrs = attrs{sqlDollar[1].val: sqlDollar[3].expr}
		}
	case 38:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:305
		{
			sqlVAL.atrs = sqlDollar[1].atrs
		}
	case 39:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:306
		{
			sqlVAL.atrs = attrsUnion(sqlDollar[1].atrs, sqlDollar[3].atrs)
		}
	case 40:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:310
		{
			sqlVAL.expr = funcall(sqlDollar[1].val, nil)
		}
	case 41:
		sqlDollar = sqlS[sqlpt-4 : sqlpt+1]
//line sql.y:311
		{
			sqlVAL.expr = funcall(sqlDollar[1].val, sqlDollar[3].expl)
		}
	case 42:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:315
		{
			sqlVAL.expl = exprlist{sqlDollar[1].expr}
		}
	case 43:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:316
		{
			sqlVAL.expl = append(sqlDollar[1].expl, sqlDollar[3].expr)
		}
	case 44:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:320
		{
			sqlVAL.expl = nil
		}
	case 45:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:321
		{
			sqlVAL.expl = sqlDollar[2].expl
		}
	case 46:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:325
		{
			sqlVAL.expr = atomic(NUMBER, sqlDollar[1].val)
		}
	case 47:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:326
		{
			sqlVAL.expr = atomic(IDENT, sqlDollar[1].val)
		}
	case 48:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:327
		{
			sqlVAL.expr = atomic(STRING, sqlDollar[1].val)
		}
	case 49:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:328
		{
			sqlVAL.expr = variadic('[', "square", sqlDollar[1].expl)
		}
	case 50:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:329
		{
			sqlVAL.expr = unary('(', "paren", sqlDollar[2].expr)
		}
	case 51:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:330
		{
			sqlVAL.expr = unary('"', "quota", atomic(STRING, sqlDollar[2].val))
		}
	case 52:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:331
		{
			sqlVAL.expr = unary('\'', "quota", atomic(STRING, sqlDollar[2].val))
		}
	case 53:
		sqlDollar = sqlS[sqlpt-1 : sqlpt+1]
//line sql.y:332
		{
			sqlVAL.expr = sqlDollar[1].expr
		}
	case 54:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:333
		{
			sqlVAL.expr = binary('+', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 55:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:334
		{
			sqlVAL.expr = binary('-', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 56:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:335
		{
			sqlVAL.expr = binary('*', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 57:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:336
		{
			sqlVAL.expr = binary('/', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 58:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:337
		{
			sqlVAL.expr = binary('%', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 59:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:338
		{
			sqlVAL.expr = binary('=', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 60:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:339
		{
			sqlVAL.expr = binary('<', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 61:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:340
		{
			sqlVAL.expr = binary('>', sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 62:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:341
		{
			sqlVAL.expr = binary(LE, sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 63:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:342
		{
			sqlVAL.expr = binary(GE, sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 64:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:343
		{
			sqlVAL.expr = binary(NE, sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 65:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:344
		{
			sqlVAL.expr = binary(AND, sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 66:
		sqlDollar = sqlS[sqlpt-3 : sqlpt+1]
//line sql.y:345
		{
			sqlVAL.expr = binary(OR, sqlDollar[1].expr, sqlDollar[2].val, sqlDollar[3].expr)
		}
	case 67:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:346
		{
			sqlVAL.expr = unary(NOT, sqlDollar[1].val, sqlDollar[2].expr)
		}
	case 68:
		sqlDollar = sqlS[sqlpt-2 : sqlpt+1]
//line sql.y:347
		{
			sqlVAL.expr = unary('-', sqlDollar[1].val, sqlDollar[2].expr)
		}
	}
	goto sqlstack /* stack new state and value */
}
