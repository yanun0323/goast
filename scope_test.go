package goast

import (
	"testing"

	"github.com/yanun0323/goast/assert"
)

func TestParseScopeMethod(t *testing.T) {
	a := assert.New(t)

	text := "func (m *memberUseCase) Start(ctx context.Context, req *usecase.UpdatePhoneReq) (res *usecase.UpdatePhoneResp, err error) {}"
	sc, err := ParseScope(0, []byte(text))
	a.NoError(err)
	a.Equal(len(sc), 1)
	// sc[0].Node().DebugPrint(10)
}
