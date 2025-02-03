package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/fr12k/go-test-parser/pkg/parser"
)

func main() {
	dir := os.Args[1]
	md := &bytes.Buffer{}
	fmt.Fprintln(md, "## Test Cases")

	p := parser.New()

	testFiles, err := p.ParseDir(dir)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for filename, tests := range testFiles {
		fmt.Fprintln(md, "") // Ensure a blank line after the comment
		fmt.Fprintf(md, "## %s\n\n", filename)
		for fncName, test := range tests.Tests {
			fmt.Fprintf(md, "### %s\n", fncName)
			fmt.Fprintln(md, test.Comment)
			fmt.Fprintln(md, "```go")
			fmt.Fprintln(md, test.Code)

			fmt.Fprintln(md, "\n```")
		}
	}

	fmt.Println(md.String())
}
