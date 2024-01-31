package goast

import (
	"testing"
)

func TestParsing(t *testing.T) {

}

func assert[T comparable](t *testing.T, expected, got T, msg ...string) {
	if expected == got {
		return
	}

	format := "mismatch result. expected: %+v, but got: %+v"
	if len(msg) != 0 && len(msg[0]) != 0 {
		format = msg[0] + ". " + format
	}

	t.Fatalf(format, expected, got)
}

func Test_extractPackage(t *testing.T) {

	tests := []struct {
		name string
		row  string
		want []*unit
	}{
		{
			name: "simple happy case",
			row:  "package goast /* comment */",
			want: []*unit{
				{
					Index: 0,
					Type:  Keyword,
					Value: "package",
				},
				{
					Index: 1,
					Type:  Raw,
					Value: "goast",
				},
				{
					Index: 2,
					Type:  Comment,
					Value: "/* comment */",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractPackage(tt.row)
			assert(t, len(tt.want), len(got))
			for i := range got {
				assert(t, *got[i], *tt.want[i])
			}
		})
	}
}

func Test_extractImport(t *testing.T) {
	type args struct {
		rows []string
		idx  int
	}
	tests := []struct {
		name  string
		args  args
		want  [][]*unit
		want1 int
	}{
		{
			name: "simple happy case",
			args: args{
				rows: []string{
					"import (",
					"\tlog \"github.com/yanun0323/pkg/logs\"",
					"\t/* foo */ \"github.com/yanun0323/gollection/v2\" /* bar */",
					")",
				},
				idx: 0,
			},
			want: [][]*unit{
				{
					{Index: 0, Type: Keyword, Value: "log"},
					{Index: 1, Type: Raw, Value: "\"github.com/yanun0323/pkg/logs\""},
				},
				{
					{Index: 0, Type: Comment, Value: "/* foo */"},
					{Index: 1, Type: Raw, Value: "\"github.com/yanun0323/gollection/v2\""},
					{Index: 2, Type: Comment, Value: "/* bar */"},
				},
			},
			want1: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := extractImport(tt.args.rows, tt.args.idx)
			assert(t, len(tt.want), len(got))
			for i := range got {
				assert(t, len(tt.want[i]), len(got[i]))
				for j := range got[i] {
					assert(t, *tt.want[i][j], *got[i][j])
				}
			}

			assert(t, tt.want1, got1)
		})
	}
}
