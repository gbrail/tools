// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains simple golden tests for various examples.
// Besides validating the results when the implementation changes,
// it provides a way to look at the generated code without having
// to execute the print statements in one's head.

package main

import (
	"strings"
	"testing"
)

// Golden represents a test case.
type Golden struct {
	name   string
	input  string // input; the package clause is provided when running the test.
	output string // exected output.
}

var golden = []Golden{
	{"day", dayIn, dayOut},
	{"offset", offsetIn, offsetOut},
	{"gap", gapIn, gapOut},
	{"num", numIn, numOut},
	{"unum", unumIn, unumOut},
	{"prime", primeIn, primeOut},
}

// Each example starts with "type XXX [u]int", with a single space separating them.

// Simple test: enumeration of type int starting at 0.
const dayIn = `type Day int
const (
	Monday Day = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)
`

const dayOut = `
const _Dayname = "MondayTuesdayWednesdayThursdayFridaySaturdaySunday"

var _Dayindex = [...]uint8{0, 6, 13, 22, 30, 36, 44, 50}

func (i Day) String() string {
	if i < 0 || i >= Day(len(_Dayindex)-1) {
		return fmt.Sprintf("Day(%d)", i)
	}
	return _Dayname[_Dayindex[i]:_Dayindex[i+1]]
}
`

// Enumeration with an offset.
// Also includes a duplicate.
const offsetIn = `type Number int
const (
	_ Number = iota
	One
	Two
	Three
	AnotherOne = One  // Duplicate; note that AnotherOne doesn't appear below.
)
`

const offsetOut = `
const _Numbername = "OneTwoThree"

var _Numberindex = [...]uint8{0, 3, 6, 11}

func (i Number) String() string {
	i--
	if i < 0 || i >= Number(len(_Numberindex)-1) {
		return fmt.Sprintf("Number(%d)", i+1)
	}
	return _Numbername[_Numberindex[i]:_Numberindex[i+1]]
}
`

// Gaps and an offset.
const gapIn = `type Gap int
const (
	Two Gap = 2
	Three Gap = 3
	Five Gap = 5
	Six Gap = 6
	Seven Gap = 7
	Eight Gap = 8
	Nine Gap = 9
	Eleven Gap = 11
)
`

const gapOut = `
const (
	_Gapname0 = "TwoThree"
	_Gapname1 = "FiveSixSevenEightNine"
	_Gapname2 = "Eleven"
)

var (
	_Gapindex0 = [...]uint8{0, 3, 8}
	_Gapindex1 = [...]uint8{0, 4, 7, 12, 17, 21}
	_Gapindex2 = [...]uint8{0, 6}
)

func (i Gap) String() string {
	switch {
	case 2 <= i && i <= 3:
		i -= 2
		return _Gapname0[_Gapindex0[i]:_Gapindex0[i+1]]
	case 5 <= i && i <= 9:
		i -= 5
		return _Gapname1[_Gapindex1[i]:_Gapindex1[i+1]]
	case i == 11:
		return _Gapname2
	default:
		return fmt.Sprintf("Gap(%d)", i)
	}
}
`

// Signed integers spanning zero.
const numIn = `type Num int
const (
	m_2 Num = -2 + iota
	m_1
	m0
	m1
	m2
)
`

const numOut = `
const _Numname = "m_2m_1m0m1m2"

var _Numindex = [...]uint8{0, 3, 6, 8, 10, 12}

func (i Num) String() string {
	i -= -2
	if i < 0 || i >= Num(len(_Numindex)-1) {
		return fmt.Sprintf("Num(%d)", i+-2)
	}
	return _Numname[_Numindex[i]:_Numindex[i+1]]
}
`

// Unsigned integers spanning zero.
const unumIn = `type Unum uint
const (
	m_2 Unum = iota + 253
	m_1
)

const (
	m0 Unum = iota
	m1
	m2
)
`

const unumOut = `
const (
	_Unumname0 = "m0m1m2"
	_Unumname1 = "m_2m_1"
)

var (
	_Unumindex0 = [...]uint8{0, 2, 4, 6}
	_Unumindex1 = [...]uint8{0, 3, 6}
)

func (i Unum) String() string {
	switch {
	case 0 <= i && i <= 2:
		return _Unumname0[_Unumindex0[i]:_Unumindex0[i+1]]
	case 253 <= i && i <= 254:
		i -= 253
		return _Unumname1[_Unumindex1[i]:_Unumindex1[i+1]]
	default:
		return fmt.Sprintf("Unum(%d)", i)
	}
}
`

// Enough gaps to trigger a map implementation of the method.
// Also includes a duplicate to test that it doesn't cause problems
const primeIn = `type Prime int
const (
	p2 Prime = 2
	p3 Prime = 3
	p5 Prime = 5
	p7 Prime = 7
	p77 Prime = 7 // Duplicate; note that p77 doesn't appear below.
	p11 Prime = 11
	p13 Prime = 13
	p17 Prime = 17
	p19 Prime = 19
	p23 Prime = 23
	p29 Prime = 29
	p37 Prime = 31
	p41 Prime = 41
	p43 Prime = 43
)
`

const primeOut = `
const _Primename = "p2p3p5p7p11p13p17p19p23p29p37p41p43"

var _Primemap = map[Prime]string{
	2:  _Primename[0:2],
	3:  _Primename[2:4],
	5:  _Primename[4:6],
	7:  _Primename[6:8],
	11: _Primename[8:11],
	13: _Primename[11:14],
	17: _Primename[14:17],
	19: _Primename[17:20],
	23: _Primename[20:23],
	29: _Primename[23:26],
	31: _Primename[26:29],
	41: _Primename[29:32],
	43: _Primename[32:35],
}

func (i Prime) String() string {
	if str, ok := _Primemap[i]; ok {
		return str
	}
	return fmt.Sprintf("Prime(%d)", i)
}
`

func TestGolden(t *testing.T) {
	for _, test := range golden {
		var g Generator
		input := "package test\n" + test.input
		file := test.name + ".go"
		g.parsePackage(".", []string{file}, input)
		// Extract the name and type of the constant from the first line.
		tokens := strings.SplitN(test.input, " ", 3)
		if len(tokens) != 3 {
			t.Fatalf("%s: need type declaration on first line", test.name)
		}
		g.generate(tokens[1])
		got := string(g.format())
		if got != test.output {
			t.Errorf("%s: got\n====\n%s====\nexpected\n====%s", test.name, got, test.output)
		}
	}
}
