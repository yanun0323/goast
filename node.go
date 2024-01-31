package goast

type Node struct {
	Name            string
	Type            Type
	Define          string  /* 'type', 'var', 'const', 'func' */
	Receiver        string  /* method receiver */
	CollectionCount int     /* array size */
	CollectionKey   Type    /* map key */
	CollectionValue Type    /* map value, slice value, array value */
	Values          []*Node /* node values */
	Parameters      []*Node /* function/method parameters */
	Returns         []*Node /* function/method returns */
	UnsupportedData string  /* raw string data unsupported type */
}

func (n *Node) Format() string {
	// TODO: format to string value with node structure
	return ""
}
