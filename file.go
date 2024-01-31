package goast

type File struct {
	Path    string
	Package []*unit
	Imports [][]*unit
	Nodes   []*Node
}
