package mdrun

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
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
	spew.Config.DisableMethods = true
	for dir, files := range fileSets {
		fmt.Println("DIR", dir, len(files))
		for i := 0; i < len(files); i++ {
			bz, err = os.ReadFile(filepath.Join(dir, files[i].Name()))
			if err != nil {
				return fmt.Errorf("read file %s: %w", files[i].Name(), err)
			}
			n := NewParser().Parse(text.NewReader(bz))
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
	mdrun, ok := n.(mdrunNode)
	if !ok {
		return ast.WalkContinue, nil
	}
	fmt.Println(mdrun.content)
	cmds := strings.Fields(mdrun.content)
	switch cmds[0] {
	case "exec":
		// expect the next node is a code block
		n := n.NextSibling()
		codeBlock, ok := n.(*ast.FencedCodeBlock)
		if !ok {
			return ast.WalkStop, errors.Errorf("expected FencedCodeBlock, got %T", n)
		}
		lang := string(codeBlock.Language(bz))
		for i := 0; i < codeBlock.Lines().Len(); i++ {
			line := codeBlock.Lines().At(i)
			s := string(line.Value(bz))
			fmt.Println("MEXT", lang, s)
		}
	default:
		return ast.WalkStop, errors.Errorf("unknow mdrun commands %q", cmds[0])
	}
	return ast.WalkContinue, nil
}
