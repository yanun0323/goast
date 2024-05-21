package goast

/*
Combo stands for any elements inside a pair of parentheses
*/
type Combo struct {
	sep  string /* '\n' or ',' */
	elem [][]*Element
}

func NewCombo(s string, sep string /* '\n' or ',' */) Combo {
	// TODO: create another func to parse all combo/context/... recursively

	return Combo{}
}
