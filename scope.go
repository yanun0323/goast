package goast

type ScopeKind int

const (
	ScopeUnknown  ScopeKind = iota + 1
	ScopePackage            // package
	ScopeComment            // comment
	ScopeImport             // import
	ScopeVariable           // var
	ScopeConst              // const
	ScopeType               // type
	ScopeFunc               // func
	ScopeMethod             // method
)

type Scope struct {
	Kind ScopeKind
}
