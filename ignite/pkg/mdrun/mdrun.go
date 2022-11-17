package mdrun

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Asserter is responsible for ensuring that cmd is executed properly
// regarding codeBlock
//
//go:generate mockery --srcpkg . --name Asserter --with-expecter
type Asserter interface {
	Assert(cmd string, codeBlock CodeBlock) error
}

// CodeBlock represents a markdown fenced code block.
type CodeBlock struct {
	Lang  string
	Lines []string
}

// Inspect detects all md files in dir, sort them by folder and assert mdrun
// commands.
// fileSets group markdown files per directory.
// Each directory is considered as a group of instructions that should be
// successfully executed to mark the group as sucessful.
// Files are sorted lexicographycally and instructions order is expected to
// be the same.
func Inspect(dir string, r Asserter) error {
	// build file sets
	var (
		currentDir = dir
		fileSets   = make(map[string][]fs.DirEntry)
	)
	err := fs.WalkDir(os.DirFS(dir), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			currentDir = filepath.Join(dir, path)
			return nil
		}
		if filepath.Ext(d.Name()) != ".md" {
			return nil
		}
		fileSets[currentDir] = append(fileSets[currentDir], d)
		return nil
	})
	if err != nil {
		return fmt.Errorf("mdrun walking %s: %w", dir, err)
	}

	// Loop in filesets to parse markdown, find mdrun blocks and assert them.
	for dir, files := range fileSets {
		for i := 0; i < len(files); i++ {
			bz, err := os.ReadFile(filepath.Join(dir, files[i].Name()))
			if err != nil {
				return fmt.Errorf("read file %s: %w", files[i].Name(), err)
			}
			root := NewParser().Parse(text.NewReader(bz))
			err = ast.Walk(root, visitor{bz: bz, r: r}.visit)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// visitor expose a visit method usable in ast.Walk
type visitor struct {
	r  Asserter
	bz []byte
}

func (v visitor) visit(n ast.Node, entering bool) (ast.WalkStatus, error) {
	mdrun, ok := n.(mdrunNode)
	if !ok {
		// skip if n is not a mdrunNode
		return ast.WalkContinue, nil
	}
	// expected next node must be a FencedCodeBlock
	codeBlock, ok := n.NextSibling().(*ast.FencedCodeBlock)
	if !ok {
		return ast.WalkStop, errors.Errorf("mdrun sibling must be a FencedCodeBlock, got %T", n.NextSibling())
	}
	var (
		lang  = string(codeBlock.Language(v.bz))
		lines []string
	)
	for i := 0; i < codeBlock.Lines().Len(); i++ {
		line := codeBlock.Lines().At(i)
		lines = append(lines, string(line.Value(v.bz)))
	}
	err := v.r.Assert(mdrun.content, CodeBlock{Lang: lang, Lines: lines})
	if err != nil {
		return ast.WalkStop, err
	}
	return ast.WalkContinue, nil
}
