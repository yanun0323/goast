// This is a sample file to test ast.
//
// This is the package comment.
package exported

import (
	"context"
	"errors"

	"github.com/yanun0323/goast/exported/enum"
)

/* comment */
var ( /* comment */
	num1             = /* comment */ 5 /* comment */
	num2             = int64(5 /* comment */)
	num3 enum.Number = 5
	num4 int64       = int64(5)
	num5 int64       = int64(int(5))

	/* comment */
	array1 = [2]int{1 /* comment */, 2 /* comment */} /*
		comment1
		comment2
	*/
	array2 [2]int       = [2]int{_num1, 2}
	array3 [2][2][2]int = [2][2][2]int{}
)

const ( /* comment */
	/* comment */ _num1 = 5
	_string1            = "string"
)

// comment
type /* comment */ Enum /* comment */ int

// comment
/* comment */
const ( /* comment */
	/* comment */ enum1/* comment */ Enum = /* comment */ iota /* comment */ + /*
			comment1
		*/1
		/* comment */
	enum2
	enum3 /* comment */
)

/* comment */
func /* comment */ init() /* comment */ { /* comment */
}

type /* comment */ SampleInterface /* comment */ interface /* comment */ { /* comment */
	Run(ctx context.Context, delegator any) error
	Stop(context.Context) error
}

var _ SampleInterface = (*SampleStruct)(nil)

type /* comment */ SampleStruct /* comment */ struct /* comment */ { /* comment */
	/* comment */ Fn /* comment */ func( /* comment */
		/* comment */ int, /* comment */
		/* comment */ string, /* comment */
		/* comment */) /* comment */ ( /* comment */
		/* comment */ int, /* comment */
		/* comment */ error, /* comment */
		/* comment */) /* comment */
}

// SampleStruct.Run
//
// runs something we don't know
func /* comment */ ( /* comment */ ss /* comment */ *SampleStruct /* comment */) /* comment */ Run /* comment */ ( /* comment */
	/* comment */ ctx /* comment */ context.Context, /*
		comment1
		comment2
	*/
	delegator /* comment */ interface{}, /* comment */
) /* comment */ error /* comment */ { /* comment */
	if delegator == nil {
		return errors.New("nil delegator")
	}

	return nil
}

type Generic[T any, V comparable] struct{}

func /* cm */ ( /* cm */ g /* cm */ *Generic[ /* cm */ T /* cm */ /* cm */, V /* cm */] /* cm */) /* cm */ Run /* cm */ ( /* cm */ obj /* cm */ T /* cm */ /* cm */, c /* cm */ V /* cm */) /* cm */ error /* cm */ {
	return nil
}

/*
SampleStruct.Run

stops something we don't know
*/
func (ss *SampleStruct) Stop(context.Context) error {
	return nil
}
