package mdrun

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func NewParser() parser.Parser {
	return parser.NewParser(
		parser.WithBlockParsers(
			append(
				parser.DefaultBlockParsers(),
				util.Prioritized(mdrunParser{}, 9999),
			)...,
		),
		parser.WithInlineParsers(parser.DefaultInlineParsers()...),
		parser.WithParagraphTransformers(parser.DefaultParagraphTransformers()...),
	)
}

// mdrunParser implements golmark/parser.BlockParser to parse mdrun blocks.
type mdrunParser struct{}

type mdrunNode struct {
	*ast.BaseBlock
	content string
}

var mdrunNodeKind = ast.NewNodeKind("MDRun")

func (mdrunNode) Kind() ast.NodeKind { return mdrunNodeKind }
func (n mdrunNode) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

// Trigger returns a list of characters that triggers Parse method of
// this parser.
// If Trigger returns a nil, Open will be called with any lines.
func (p mdrunParser) Trigger() []byte {
	return []byte{'['}
}

func (p mdrunParser) process(reader text.Reader) ([]byte, bool) {
	line, _ := reader.PeekLine()
	if !bytes.HasPrefix(line, []byte("[mdrun]: # ")) {
		return nil, false
	}
	var (
		start   = bytes.LastIndex(line, []byte{'('}) + 1
		end     = bytes.LastIndex(line, []byte{')'})
		content = line[start:end]
	)
	return content, true
}

// Open parses the current line and returns a result of parsing.
//
// Open must not parse beyond the current line.
// If Open has been able to parse the current line, Open must advance a reader
// position by consumed byte length.
//
// If Open has not been able to parse the current line, Open should returns
// (nil, NoChildren). If Open has been able to parse the current line, Open
// should returns a new Block node and returns HasChildren or NoChildren.
func (p mdrunParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	if content, ok := p.process(reader); ok {
		return mdrunNode{
			BaseBlock: &ast.BaseBlock{},
			content:   string(content),
		}, parser.Continue
	}
	return nil, parser.NoChildren
}

// Continue parses the current line and returns a result of parsing.
//
// Continue must not parse beyond the current line.
// If Continue has been able to parse the current line, Continue must advance
// a reader position by consumed byte length.
//
// If Continue has not been able to parse the current line, Continue should
// returns Close. If Continue has been able to parse the current line,
// Continue should returns (Continue | NoChildren) or
// (Continue | HasChildren)
func (p mdrunParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	if _, ok := p.process(reader); ok {
		return parser.Continue
	}
	return parser.Close
}

// Close will be called when the parser returns Close.
func (p mdrunParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	// nothing to do
}

// CanInterruptParagraph returns true if the parser can interrupt paragraphs,
// otherwise false.
func (p mdrunParser) CanInterruptParagraph() bool {
	return true
}

// CanAcceptIndentedLine returns true if the parser can open new node when
// the given line is being indented more than 3 spaces.
func (p mdrunParser) CanAcceptIndentedLine() bool {
	return false
}
