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

func TestParseScopeStruct(t *testing.T) {
	a := assert.New(t)

	text := "type UpdatePhoneReq struct {\n\t\tPhone       *string `json:\"phone\" binding:\"required\"`\n\t\tAreaCode    map[string]*Hello `json:\"area_code\" binding:\"required\"`\n\t\tCaptchaCode []*Hello `json:\"code\" binding:\"required\"`\n\t}\""

	// n, err := extract([]byte(text))
	// a.NoError(err)
	// n.DebugPrint()
	// a.Nil(1)

	sc, err := ParseScope(0, []byte(text))
	a.NoError(err)
	a.Equal(len(sc), 1)

	// sc[0].Node().DebugPrint()
	// a.Nil(1)
}
