package mdrun

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// mdrunParser implements golmark/parser.BlockParser to parse mdrun blocks.
type mdrunParser struct{}

type mdrunNode struct {
	*ast.BaseNode
}

var mdrunNodeKind = ast.NewNodeKind("mdrun")

func (mdrunNode) Dump(source []byte, level int) {}
func (mdrunNode) HasBlankPreviousLines() bool   { return false }
func (mdrunNode) IsRaw() bool                   { return false }
func (mdrunNode) Kind() ast.NodeKind            { return mdrunNodeKind }
func (mdrunNode) Lines() *text.Segments         {}

// Trigger returns a list of characters that triggers Parse method of
// this parser.
// If Trigger returns a nil, Open will be called with any lines.
func (p mdrunParser) Trigger() []byte {
	return []byte("[mdrun]: # ")
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
	return mdrunNode{}, parser.Continue
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
	panic("not implemented") // TODO: Implement
}

// Close will be called when the parser returns Close.
func (p mdrunParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	panic("not implemented") // TODO: Implement
}

// CanInterruptParagraph returns true if the parser can interrupt paragraphs,
// otherwise false.
func (p mdrunParser) CanInterruptParagraph() bool {
	panic("not implemented") // TODO: Implement
}

// CanAcceptIndentedLine returns true if the parser can open new node when
// the given line is being indented more than 3 spaces.
func (p mdrunParser) CanAcceptIndentedLine() bool {
	panic("not implemented") // TODO: Implement
}
