package goast

type unit struct {
	Index int
	Type  Type
	Value string
}

func Unit(index int, t Type, value string) *unit {
	return &unit{
		Index: index,
		Type:  t,
		Value: value,
	}
}

func (u unit) isZero() bool {
	return u.Index == 0 && u.Type == Raw && len(u.Value) == 0
}
