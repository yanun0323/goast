package goast

type File struct {
	Path    string
	Package string
	Imports []string
	Nodes   []*Node
}
