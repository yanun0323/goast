// This is the comment
// for explain this package.
package exported

import (
	"context"
	"errors"
)

//go:generate domaingen -v -destination=../../usecase/member.go -name=memberUseCase
type MemberUseCase interface {
	Start(ctx context.Context, req *UpdatePhoneReq) (res *UpdatePhoneResp, err error)
	End(context.Context, *UpdatePhoneReq) (*UpdatePhoneResp, error)
	EndAgain( /* 123 */ context.Context /* 456 */, * /* 789 */ UpdatePhoneReq /* 012 */) /* 345 */ (* /* 678 */ UpdatePhoneResp /* 901 */, error /* 234 */)
	// Exit(ctx context.Context) error
}

var (
	ErrNotFound         = errors.New("not found")
	ErrPermissionDenied = errors.New("permission denied")
)

type UpdatePhoneReq struct {
	Phone       string
	AreaCode    string
	CaptchaCode string
	CreateTime  int64
	UpdateTime  int64
	CreateAt    int64
	UpdateAt    int64
}

type UpdatePhoneResp struct {
	Phone       string
	AreaCode    string
	CaptchaCode string
	CreateTime  string
	UpdateTime  string
	CreateAt    string
	UpdateAt    string
}
