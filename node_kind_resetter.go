package goast

type resetterRule struct {
	TriggerKeyword string
	KindStack      []Kind
	AcceptedKind   []Kind
}

// nodeKindResetter re-set the node kind with node's position of scope
type nodeKindResetter struct {
}

func Run(scs []Scope) {
	for _, sc := range scs {
		switch sc.Kind() {
		case ScopeImport:

		case ScopeVariable:

		case ScopeConst:

		case ScopeType:

		case ScopeFunc:

		}
	}
}
