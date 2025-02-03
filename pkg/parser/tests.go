package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/fr12k/go-file"
)

const exportKey = "@markdown"

type Parser struct {
	tests    map[string]TestFile
	openFile file.OpenFunc
}
type TestFile struct {
	Name  string
	Tests map[string]Test
}

type Test struct {
	Name    string
	Comment string
	Code    string
}

func New() *Parser {
	return &Parser{
		tests:    map[string]TestFile{},
		openFile: file.Open(),
	}
}

func (p *Parser) ParseDir(dir string) (map[string]TestFile, error) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), "_test.go") {
			return p.processTestFile(path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return p.tests, nil
}

func (p *Parser) ParseFile(file string) (map[string]TestFile, error) {
	err := p.processTestFile(file)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return p.tests, nil
}
func (p *Parser) processTestFile(filename string) error {
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, filename, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing:", filename, err)
		return err
	}

	// Extract test functions
	for _, decl := range node.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && strings.HasPrefix(fn.Name.Name, "Test") {
			if fn.Doc != nil {
				// Convert the comment block to a string
				commentText := fn.Doc.Text()
				if !strings.Contains(commentText, exportKey) {
					continue
				}

				commentText = strings.Replace(commentText, exportKey+"\n", "", 1)
				if strings.Contains(commentText, exportKey) {
					commentText = strings.Replace(commentText, exportKey+" ", "", 1)
				}

				// Print the full function with comments & blank lines preserved
				code, err := p.printFullFunction(fs, node, fn)
				if err != nil {
					fmt.Println("Error reading code block:", filename, err)
					return err
				}

				t := Test{
					Name:    fn.Name.Name,
					Comment: commentText,
					Code:    code,
				}
				if _, ok := p.tests[filename]; !ok {
					p.tests[filename] = TestFile{
						Tests: map[string]Test{},
					}
				}
				p.tests[filename].Tests[t.Name] = t

			}
		}
	}
	return nil
}

func (p *Parser) printFullFunction(fs *token.FileSet, node *ast.File, fn *ast.FuncDecl) (string, error) {
	out := &bytes.Buffer{}

	// Extract the exact portion of the source that contains the function
	start := fs.Position(fn.Pos()).Offset
	end := fs.Position(fn.End()).Offset

	source, err := p.openFile(fs.Position(node.Package).Filename).Read()
	if err != nil {
		fmt.Println("Error reading file:", err)
		return "", err
	}

	// Extract and write only the function part
	out.Write(source[start:end])
	return out.String(), nil
}
