package mdrun

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type Runner struct {
	dir     string
	content []byte
}

func Run(dir string) error {
	var (
		currentDir = dir
		// fileSets will group markdown files per directory.
		// Each directory is considered as a group of instructions that should
		// be successfully executed to mark the group as sucessful.
		// Files are sorted lexicographycally and instructions order is expected to
		// be the same.
		fileSets = make(map[string][]fs.DirEntry)
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
	parser := parser.NewParser(
		parser.WithBlockParsers(
			append(
				parser.DefaultBlockParsers(),
				util.Prioritized(mdrunParser{}, 9999),
			)...,
		),
		parser.WithInlineParsers(parser.DefaultInlineParsers()...),
		parser.WithParagraphTransformers(parser.DefaultParagraphTransformers()...),
	)
	spew.Config.DisableMethods = true
	for dir, files := range fileSets {
		fmt.Println("DIR", dir, len(files))
		for i := 0; i < len(files); i++ {
			bz, err = os.ReadFile(filepath.Join(dir, files[i].Name()))
			if err != nil {
				return fmt.Errorf("read file %s: %w", files[i].Name(), err)
			}
			n := parser.Parse(text.NewReader(bz))
			err = ast.Walk(n, visitMD)
			if err != nil {
				return err
			}

		}
	}
	return nil
}

var bz []byte

func visitMD(n ast.Node, entering bool) (ast.WalkStatus, error) {
	fmt.Println(n.Kind(), string(n.Text(bz)))
	if n.Type() == ast.TypeBlock {
		n.Dump(bz, 0)
	}
	if n.Kind() != ast.KindTextBlock {
		return ast.WalkContinue, nil
	}
	return ast.WalkContinue, nil
}
