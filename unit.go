package goast

type unit struct {
	Index int
	Type  Type
	Value string
}

func (u unit) isZero() bool {
	return u.Index == 0 && u.Type == Raw && len(u.Value) == 0
}
