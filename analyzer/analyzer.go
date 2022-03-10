package analyzer

import (
	"c_compiler/helper"
	"c_compiler/lexer"
	"errors"
	"fmt"
	"log"
	"strconv"
)

type Analyzer struct {
	Tokens      []lexer.Token
	CurToken    lexer.Token
	CurTokenIdx int
	Quads       [helper.MaxQuadSize]Quaternion
	NexQuadIdx  int
	TmpIdx      int
}

type Quaternion struct {
	Op   string
	Arg1 string
	Arg2 string
	Res  string
}

func (a *Analyzer) Exec() {
	a.CheckMain()
	var chain int
	a.NexQuadIdx = 1
	a.CheckStatementBlock(&chain)
	//a.PrintQuads()
}

func (a *Analyzer) ScanToken() {
	a.CurTokenIdx++
	if a.CurTokenIdx >= len(a.Tokens) {
		return
	}
	a.CurToken = a.Tokens[a.CurTokenIdx]
}

func (a *Analyzer) Match(syn uint8) {
	if a.CurTokenIdx >= len(a.Tokens) {
		fmt.Println("索引错误")
		log.Fatal()
	}
	if a.Tokens[a.CurTokenIdx].Syn != syn {
		fmt.Println("语法分析错误：单词不匹配")
		curToken := a.Tokens[a.CurTokenIdx]
		fmt.Printf("需要%d，发现%s\n", syn, curToken.Value)
		fmt.Printf("%d %d\n", curToken.Row, curToken.Col)
		log.Fatal()
	}
	a.ScanToken()
}

func (a *Analyzer) CheckMain() {
	a.Match(lexer.SynMain)
	a.Match(lexer.SynLBranL)
	a.Match(lexer.SynLBranR)
}

func (a *Analyzer) CheckStatementBlock(chain *int) {
	a.Match(lexer.SynBBranL)
	a.CheckStatementSeq(chain)
	a.Match(lexer.SynBBranR)
}

func (a *Analyzer) CheckStatementSeq(chain *int) {
	a.CheckStatement(chain)
	for a.CurToken.Syn == lexer.SynId || a.CurToken.Syn == lexer.SynIf || a.CurToken.Syn == lexer.SynWhile {
		a.FillResOfLink(*chain, a.NexQuadIdx)
		a.CheckStatement(chain)
	}
	a.FillResOfLink(*chain, a.NexQuadIdx)
}

func (a *Analyzer) CheckStatement(chain *int) {
	quadIdx1, quadIdx2, chainTmp := 0, 0, 0
	switch a.CurToken.Syn {
	case lexer.SynId:
		res := a.CurToken.Value
		a.ScanToken()
		a.Match(lexer.SynAssign)
		arg := a.CheckExpression()
		a.Match(lexer.SynSemicolon)
		a.GenQuaternion("=", arg, "", res)
		*chain = 0
	case lexer.SynIf:
		a.Match(lexer.SynIf)
		a.Match(lexer.SynLBranL)
		a.CheckCondition(&quadIdx1, &quadIdx2)
		a.FillResOfLink(quadIdx1, a.NexQuadIdx)
		a.Match(lexer.SynLBranR)
		a.CheckStatementBlock(&chainTmp)
		*chain = a.MergeLink(chainTmp, quadIdx2)
	case lexer.SynWhile:
		a.Match(lexer.SynWhile)
		res := strconv.Itoa(a.NexQuadIdx)
		a.Match(lexer.SynLBranL)
		a.CheckCondition(&quadIdx1, &quadIdx2)
		a.FillResOfLink(quadIdx1, a.NexQuadIdx)
		a.Match(lexer.SynLBranR)
		a.CheckStatementBlock(&chainTmp)
		quadTmp, _ := strconv.Atoi(res)
		a.FillResOfLink(chainTmp, quadTmp)
		a.GenQuaternion("j", "", "", res)
		*chain = quadIdx2
	default:
		log.Fatal(errors.New("can't identify such statement"))
	}
}

func (a *Analyzer) CheckExpression() string {
	arg1 := a.CheckTerm()
	res := arg1
	for a.CurToken.Syn == lexer.SynPlus || a.CurToken.Syn == lexer.SynMinus {
		op := a.CurToken.Value
		a.ScanToken()
		arg2 := a.CheckTerm()
		res = a.GenTmp()
		a.GenQuaternion(op, arg1, arg2, res)
		arg1 = res
	}

	return res
}

func (a *Analyzer) CheckTerm() string {
	arg1 := a.CheckFactor()
	res := arg1

	for a.CurToken.Syn == lexer.SynTimes || a.CurToken.Syn == lexer.SynDivide {
		op := a.CurToken.Value
		a.ScanToken()
		arg2 := a.CheckFactor()
		res = a.GenTmp()
		a.GenQuaternion(op, arg1, arg2, res)
		arg1 = res
	}

	return res
}

func (a *Analyzer) CheckFactor() string {
	var res string
	if a.CurToken.Syn == lexer.SynId || a.CurToken.Syn == lexer.SynNum {
		res = a.CurToken.Value
		a.ScanToken()
	} else {
		a.Match(lexer.SynLBranL)
		res = a.CheckExpression()
		a.Match(lexer.SynLBranR)
	}
	return res
}

func (a *Analyzer) CheckCondition(quadIdx1 *int, quadIdx2 *int) {
	arg1 := a.CheckExpression()
	if a.CurToken.Syn >= lexer.SynLT && a.CurToken.Syn <= lexer.SynGT {
		op := a.CurToken.Value
		a.ScanToken()
		arg2 := a.CheckExpression()
		*quadIdx1 = a.NexQuadIdx
		*quadIdx2 = a.NexQuadIdx + 1
		a.GenQuaternion("j" + op, arg1, arg2, "0")
		a.GenQuaternion("j", "", "", "0")
	} else {
		fmt.Println(a.CurToken.Value)
		fmt.Printf("%d %d\n", a.Tokens[a.CurTokenIdx].Row, a.Tokens[a.CurTokenIdx].Col)
		log.Fatal(errors.New("relational operators' error"))
	}
}

func (a *Analyzer) GenQuaternion(op string, arg1 string, arg2 string, res string) {
	a.Quads[a.NexQuadIdx] = Quaternion{
		Op:   op,
		Arg1: arg1,
		Arg2: arg2,
		Res:  res,
	}
	a.NexQuadIdx++
}

func (a *Analyzer) GenTmp() string {
	a.TmpIdx++
	res := "T" + strconv.Itoa(a.TmpIdx)
	return res
}

func (a *Analyzer) MergeLink(p1 int, p2 int) int {
	res := p2
	if res == 0 {
		res = p1
	}
	tmpRes, _ := strconv.Atoi(a.Quads[p2].Res)
	for tmpRes != 0 {
		p2, _ = strconv.Atoi(a.Quads[p2].Res)
		a.Quads[p2].Res = strconv.Itoa(p1)
		tmpRes, _ = strconv.Atoi(a.Quads[p2].Res)
	}
	return res
}

func (a *Analyzer) FillResOfLink(p int, res int) {
	for p != 0 {
		tmpRes, _ := strconv.Atoi(a.Quads[p].Res)
		a.Quads[p].Res = strconv.Itoa(res)
		p = tmpRes
	}
}

func (a *Analyzer) PrintQuads() {
	for i := 1; a.Quads[i].Op != ""; i++ {
		fmt.Println(a.Quads[i])
	}
}
