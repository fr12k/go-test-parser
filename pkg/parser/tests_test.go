package parser

import (
	"go/scanner"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"

	"github.com/fr12k/go-file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// @markdown
// TestParseDirThisFile this test parse this test file
func TestParseDirThisFile(t *testing.T) {
	filename, dir := getCwd(t)

	p := New()
	result, err := p.ParseDir(dir)
	assert.NoError(t, err)
	assert.Contains(t, result, filename)
	assert.Contains(t, result[filename].Tests, "TestParseDirThisFile")
	assert.Equal(t, "TestParseDirThisFile this test parse this test file\n", result[filename].Tests["TestParseDirThisFile"].Comment)
	assert.Contains(t, result[filename].Tests["TestParseDirThisFile"].Code, "func TestParseDirThisFile")
	assert.Contains(t, result[filename].Tests, "TestParseDir")
	assert.Equal(t, "TestParseDir this test parse this test file\n", result[filename].Tests["TestParseDir"].Comment)
	assert.Contains(t, result[filename].Tests["TestParseDir"].Code, "func TestParseDir")
}

// @markdown TestParseDir this test parse this test file
func TestParseDir(t *testing.T) {
	// Just a test that is used for parsing see above.
}

// TestParseFileError this test the error handling when the source code file couldn't be read
func TestParseDirError(t *testing.T) {
	_, dir := getCwd(t)

	p := New()
	p.openFile = file.OpenFile(file.NewReaderError(io.EOF))

	result, err := p.ParseDir(dir)
	assert.ErrorIs(t, err, io.EOF)

	assert.Nil(t, result)
}

func TestParseFileError(t *testing.T) {
	tempDir := t.TempDir()
	testFileContent := `package main

import "testing"

// @markdown
// TestExample this is a test function
func TestExample(t *testing.T)
	t.Log("example test")
}`
	testFilePath := filepath.Join(tempDir, "example_test.go")
	err := os.WriteFile(testFilePath, []byte(testFileContent), 0644)
	assert.NoError(t, err)

	p := New()

	result, err := p.ParseFile(testFilePath)
	errLst, ok := err.(scanner.ErrorList)
	require.True(t, ok)
	assert.Len(t, errLst, 1)
	assert.Nil(t, result)
}

func TestParseFile(t *testing.T) {
	filename, _ := getCwd(t)

	p := New()
	result, err := p.ParseFile(filename)
	assert.NoError(t, err)
	assert.Contains(t, result, filename)
	assert.Contains(t, result[filename].Tests, "TestParseDirThisFile")
	assert.Equal(t, "TestParseDirThisFile this test parse this test file\n", result[filename].Tests["TestParseDirThisFile"].Comment)
	assert.Contains(t, result[filename].Tests["TestParseDirThisFile"].Code, "func TestParseDirThisFile")
	assert.Contains(t, result[filename].Tests, "TestParseDir")
	assert.Equal(t, "TestParseDir this test parse this test file\n", result[filename].Tests["TestParseDir"].Comment)
	assert.Contains(t, result[filename].Tests["TestParseDir"].Code, "func TestParseDir")
}

func TestParseMissingDirError(t *testing.T) {
	p := New()

	result, err := p.ParseDir(t.Name())
	pathErr, ok := err.(*fs.PathError)
	require.True(t, ok)
	assert.ErrorIs(t, pathErr, syscall.Errno(2))
	assert.Nil(t, result)
}

// test utility

func getCwd(t *testing.T) (filename string, dir string) {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to get current test file directory")
	}

	dir = filepath.Dir(filename)
	t.Log("Test file directory:", dir)
	return
}
