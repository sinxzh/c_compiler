package lexer

import (
	"c_compiler/helper"
	"errors"
	"fmt"
	"log"
)

var KeyWords = []string{
	"main", "if", "else", "while",
}

const (
	SynError uint8 = iota // 0

	SynMain  uint8 = iota
	SynInt   uint8 = iota
	SynFloat uint8 = iota
	SynChar  uint8 = iota
	SynIf    uint8 = iota
	SynElse  uint8 = iota // 6
	SynFor   uint8 = iota
	SynWhile uint8 = iota
	SynVoid  uint8 = iota
	SynRet   uint8 = iota

	SynId  uint8 = iota // 11
	SynNum uint8 = iota

	SynAssign uint8 = iota
	SynPlus   uint8 = iota
	SynMinus  uint8 = iota
	SynTimes  uint8 = iota // 16
	SynDivide uint8 = iota
	SynLBranL uint8 = iota
	SynLBranR uint8 = iota
	SynMBranL uint8 = iota
	SynMBranR uint8 = iota // 21
	SynBBranL uint8 = iota
	SynBBranR uint8 = iota

	SynColon     uint8 = iota
	SynComma     uint8 = iota
	SynSemicolon uint8 = iota // 26

	SynLT uint8 = iota
	SynLE uint8 = iota
	SynEq uint8 = iota
	SynNE uint8 = iota
	SynGE uint8 = iota // 31
	SynGT uint8 = iota
)

type Token struct {
	Value string
	Syn   uint8
	Row   int
	Col   int
}

type Lexer struct {
	SourceCode string
	CurrentIdx int
	EndIdx     int
	Row        int
	Col        int
	Tokens     []Token
}

func (l *Lexer) Exec() {
	for l.CurrentIdx < len(l.SourceCode) {
		l.RecFirChar()
	}
}

func (l *Lexer) RecFirChar() {
	switch l.SourceCode[l.CurrentIdx] {
	case '+':
		l.HanLasChar(true, SynPlus)
	case '-':
		l.HanLasChar(true, SynMinus)
	case '*':
		l.HanLasChar(true, SynTimes)
	case '/':
		l.HanLasChar(true, SynDivide)
	case '(':
		l.HanLasChar(true, SynLBranL)
	case ')':
		l.HanLasChar(true, SynLBranR)
	case '[':
		l.HanLasChar(true, SynMBranL)
	case ']':
		l.HanLasChar(true, SynMBranR)
	case '{':
		l.HanLasChar(true, SynBBranL)
	case '}':
		l.HanLasChar(true, SynBBranR)
	case ':':
		l.HanLasChar(true, SynColon)
	case ',':
		l.HanLasChar(true, SynComma)
	case ';':
		l.HanLasChar(true, SynSemicolon)
	default:
		l.HanChar()
	}
}

func (l *Lexer) HanChar() {
	switch l.SourceCode[l.CurrentIdx] {
	case '>':
		if l.CurrentIdx+1 < len(l.SourceCode) && l.SourceCode[l.CurrentIdx+1] == '=' {
			l.HanLasChar(false, SynGE)
		} else {
			l.HanLasChar(true, SynGT)
		}
	case '<':
		if l.CurrentIdx+1 < len(l.SourceCode) && l.SourceCode[l.CurrentIdx+1] == '=' {
			l.HanLasChar(false, SynLE)
		} else {
			l.HanLasChar(true, SynLT)
		}
	case '=':
		if l.CurrentIdx+1 < len(l.SourceCode) && l.SourceCode[l.CurrentIdx+1] == '=' {
			l.HanLasChar(false, SynEq)
		} else {
			l.HanLasChar(true, SynAssign)
		}
	case '!':
		if l.CurrentIdx+1 < len(l.SourceCode) && l.SourceCode[l.CurrentIdx+1] == '=' {
			l.HanLasChar(false, SynNE)
		} else {
			l.HanLasChar(true, SynError)
			l.PrintTokens()
			fmt.Printf("%c %d %d\n", l.SourceCode[l.CurrentIdx], l.Row, l.Col)
			log.Fatal(errors.New("source code has illegal character"))
		}
	default:
		if helper.IsAlpha(rune(l.SourceCode[l.CurrentIdx])) {
			l.HanKeyAndId()
		} else if helper.IsDigit(rune(l.SourceCode[l.CurrentIdx])) {
			l.HanNum()
		} else if helper.IsBlank(rune(l.SourceCode[l.CurrentIdx])) {
			l.HanBlank()
		} else {
			l.HanLasChar(true, SynError)
			l.PrintTokens()
			fmt.Printf("%c %d %d\n", l.SourceCode[l.CurrentIdx], l.Row, l.Col)
			log.Fatal(errors.New("source code has illegal character"))
		}
	}
}

// 单个字符的EndIdx放在HanLasChar中处理，多个字符的EndIdx在各自的处理函数中处理

func (l *Lexer) HanKeyAndId() {
	l.EndIdx = l.CurrentIdx + 1
	for l.EndIdx < len(l.SourceCode) && (helper.IsAlpha(rune(l.SourceCode[l.EndIdx])) || helper.IsDigit(rune(l.SourceCode[l.EndIdx]))) {
		l.EndIdx++
	}

	runesSourceCode := []rune(l.SourceCode)
	for _, keyword := range KeyWords {
		if string(runesSourceCode[l.CurrentIdx:l.EndIdx]) == keyword {
			l.HanKeyChar(keyword)
			return
		}
	}

	l.HanLasChar(false, SynId)
}

func (l *Lexer) HanKeyChar(keyword string) {
	switch keyword {
	case "main":
		l.HanLasChar(false, SynMain)
	case "int":
		l.HanLasChar(false, SynInt)
	case "double":
		l.HanLasChar(false, SynFloat)
	case "char":
		l.HanLasChar(false, SynChar)
	case "if":
		l.HanLasChar(false, SynIf)
	case "else":
		l.HanLasChar(false, SynElse)
	case "for":
		l.HanLasChar(false, SynFor)
	case "while":
		l.HanLasChar(false, SynWhile)
	case "void":
		l.HanLasChar(false, SynVoid)
	case "return":
		l.HanLasChar(false, SynRet)
	}
}

func (l *Lexer) HanNum() {
	l.EndIdx = l.CurrentIdx + 1
	isFloat := false
	for l.EndIdx < len(l.SourceCode) {
		if helper.IsDigit(rune(l.SourceCode[l.EndIdx])) {
			l.EndIdx++
		} else if l.SourceCode[l.EndIdx] == '.' {
			if !isFloat {
				isFloat = true
				l.EndIdx++
			} else {
				l.HanLasChar(false, SynError)
				l.PrintTokens()
				fmt.Println("词法分析错误：数字非法")
				fmt.Printf("%c %d %d\n", l.SourceCode[l.EndIdx], l.Row, l.Col)
				//log.Fatal(errors.New("the number is in wrong format. "))
				log.Fatal()
			}
		} else {
			break
		}
	}

	l.HanLasChar(false, SynNum)
}

// HanBlank 遇到空白字符时单独处理row和col
func (l *Lexer) HanBlank() {
	l.EndIdx = l.CurrentIdx + 1
	for l.EndIdx < len(l.SourceCode) && helper.IsBlank(rune(l.SourceCode[l.EndIdx])) {
		l.EndIdx++
	}
	for i := l.CurrentIdx; i < l.EndIdx; i++ {
		if l.SourceCode[i] == ' ' {
			l.Col++
		} else if l.SourceCode[i] == '\t' {
			l.Col += 4
		} else if l.SourceCode[i] == '\n' { // 经测试，'\r' 与 '\n' 位置相同
			//fmt.Printf("%d %d\n", l.Row, l.Col)
			//fmt.Printf("%c test\n", l.SourceCode[i])
			l.Col = 1
			l.Row++
		} else {
			//fmt.Printf("%d %d\n", l.Row, l.Col)
		}
	}
	l.CurrentIdx = l.EndIdx
}

// HanLasChar 存放单词符号和标志，处理非空白字符的row和col，处理错误
func (l *Lexer) HanLasChar(isSingle bool, syn uint8) {
	if isSingle {
		l.EndIdx = l.CurrentIdx + 1
	}
	// 存放单词符号和标志
	l.Tokens = append(l.Tokens, Token{
		Value: string([]rune(l.SourceCode)[l.CurrentIdx:l.EndIdx]),
		Syn:   syn,
		Row:   l.Row,
		Col:   l.Col,
	})
	// 处理非空白字符的row和col
	l.Col += l.EndIdx - l.CurrentIdx
	l.CurrentIdx = l.EndIdx
}

func (l *Lexer) PrintTokens() {
	for _, token := range l.Tokens {
		fmt.Println(token)
	}
}
