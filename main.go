package main

import (
	"c_compiler/analyzer"
	"c_compiler/lexer"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func ReadSourceCode(filename string) string {
	sourceCode, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(errors.New("failed to read source code from " + filename))
	}
	return string(sourceCode)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type ViewQuaternion struct {
	Idx  int
	Quad analyzer.Quaternion
}

type ViewDate struct {
	SourceCode string
	Tokens     []lexer.Token
	Quads      []ViewQuaternion
}

func cpNewHandler(writer http.ResponseWriter, request *http.Request) {
	sourceCode := request.FormValue("sourceCode")
	option := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile("./code.txt", option, os.FileMode(0600))
	checkErr(err)
	if err := os.Truncate("./code.txt", 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}
	_, err = fmt.Fprintln(file, sourceCode)
	checkErr(err)
	err = file.Close()
	checkErr(err)
	http.Redirect(writer, request, "/compiler", http.StatusFound)
}

func cpViewHandler(writer http.ResponseWriter, _ *http.Request) {
	html, err := template.ParseFiles("view.html")
	checkErr(err)
	sourceCode := ReadSourceCode("./code.txt")
	mLexer := lexer.Lexer{
		SourceCode: sourceCode,
		CurrentIdx: 0,
		EndIdx:     0,
		Row:        1,
		Col:        1,
	}
	mLexer.Exec()
	mAnalyzer := analyzer.Analyzer{
		Tokens: mLexer.Tokens,
	}
	mAnalyzer.Exec()
	mViewData := ViewDate{
		SourceCode: sourceCode,
		Tokens:     mLexer.Tokens,
		Quads:      nil,
	}
	var viewQuads []ViewQuaternion
	for i := 1; mAnalyzer.Quads[i].Op != ""; i++ {
		viewQuads = append(viewQuads, ViewQuaternion{
			Idx: i,
			Quad: analyzer.Quaternion{
				Op:   mAnalyzer.Quads[i].Op,
				Arg1: mAnalyzer.Quads[i].Arg1,
				Arg2: mAnalyzer.Quads[i].Arg2,
				Res:  mAnalyzer.Quads[i].Res,
			},
		})
	}
	mViewData.Quads = viewQuads
	err = html.Execute(writer, mViewData)
	checkErr(err)
}

func main() {
	http.HandleFunc("/compiler", cpViewHandler)
	fmt.Println("compiler")
	http.HandleFunc("/compiler/new", cpNewHandler)
	fmt.Println("compiler/new")
	err := http.ListenAndServe("localhost:8080", nil)
	fmt.Println("finish")
	log.Fatal(err)
}
