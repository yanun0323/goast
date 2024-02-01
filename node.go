package goast

type Node struct {
	Values          []*unit /* node values/defines */
	CollectionCount int     /* array size */
	CollectionKey   Type    /* map key */
	CollectionValue Type    /* map value, slice value, array value */

	Parameters      []*Node /* struct/interface/function/method/var/const parameters */
	Returns         []*Node /* function/method returns */
	UnsupportedData *unit   /* raw string data unsupported type */
}

func (n *Node) Format() string {
	// TODO: format to string value with node structure
	return ""
}
