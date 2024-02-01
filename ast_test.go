package goast

import (
	"testing"
)

func TestParsing(t *testing.T) {

}

func equal[T comparable](t *testing.T, expected, got T, msg ...string) {
	if expected == got {
		return
	}

	format := "mismatch result. expected equal to: %+v, but got: %+v"
	if len(msg) != 0 && len(msg[0]) != 0 {
		format = msg[0] + ". " + format
	}

	t.Fatalf(format, expected, got)
}

func notEqual[T comparable](t *testing.T, expected, got T, msg ...string) {
	if expected != got {
		return
	}

	format := "mismatch result. expected not equal to: %+v, but got: %+v"
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
				Unit(0, Keyword, "package"),
				Unit(1, Raw, "goast"),
				Unit(2, Comment, "/* comment */"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractPackage(tt.row)
			equal(t, len(tt.want), len(got))
			for i := range got {
				equal(t, *got[i], *tt.want[i])
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
					Unit(0, Keyword, "log"),
					Unit(1, Raw, "\"github.com/yanun0323/pkg/logs\""),
				},
				{
					Unit(0, Comment, "/* foo */"),
					Unit(1, Raw, "\"github.com/yanun0323/gollection/v2\""),
					Unit(2, Comment, "/* bar */"),
				},
			},
			want1: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := extractImport(tt.args.rows, tt.args.idx)
			equal(t, len(tt.want), len(got))
			for i := range got {
				equal(t, len(tt.want[i]), len(got[i]))
				for j := range got[i] {
					equal(t, *tt.want[i][j], *got[i][j])
				}
			}

			equal(t, tt.want1, got1)
		})
	}
}

func Test_extractComment(t *testing.T) {
	type args struct {
		rows []string
		idx  int
	}
	tests := []struct {
		name  string
		args  args
		want  unit
		want1 int
	}{
		{
			name: "one row comment with double slashes happy case",
			args: args{
				rows: []string{"// one row comment with double slashes", ""},
				idx:  0,
			},
			want:  *Unit(0, Comment, "// one row comment with double slashes"),
			want1: 0,
		},
		{
			name: "two rows comments with double slashes happy case",
			args: args{
				rows: []string{"// two rows comments", "// with double slashes", ""},
				idx:  0,
			},
			want:  *Unit(0, Comment, "// two rows comments\n// with double slashes"),
			want1: 1,
		},
		{
			name: "one row comment with single slash happy case",
			args: args{
				rows: []string{"/* one row comment with single slash */", ""},
				idx:  0,
			},
			want:  *Unit(0, Comment, "/* one row comment with single slash */"),
			want1: 0,
		},
		{
			name: "four rows comments with single slash happy case",
			args: args{
				rows: []string{"/* multiline comments", "with", "single slash", "*/", ""},
				idx:  0,
			},
			want:  *Unit(0, Comment, "/* multiline comments\nwith\nsingle slash\n*/"),
			want1: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := extractComment(tt.args.rows, tt.args.idx)
			equal(t, tt.want, *got)
			equal(t, tt.want1, got1)
		})
	}
}

func Test_extractConst(t *testing.T) {
	type args struct {
		rows []string
		idx  int
	}
	tests := []struct {
		name  string
		args  args
		want  *Node
		want1 int
	}{
		{
			name: "one row happy case",
			args: args{
				rows: []string{
					"const _name = \"hello\"",
				},
				idx: 0,
			},
			want: &Node{
				Values: []*unit{
					Unit(0, Keyword, "const"),
					Unit(1, Raw, "_name"),
					Unit(2, Keyword, "="),
					Unit(3, Raw, "\"hello\""),
				},
			},
			want1: 0,
		},
		{
			name: "multi rows happy case",
			args: args{
				rows: []string{
					"const (",
					"\t_name = \"hello\"",
					"\t_age = 5",
					"\t_foo = map[int]int{",
					"\t\t1 : 11",
					"\t\t2 : 22",
					"\t\t3 : 33",
					"\t}",
					"\t_multiString = `",
					"one",
					"two",
					"three",
					"\t`",
					")",
				},
				idx: 0,
			},
			want: &Node{
				Values: []*unit{
					Unit(0, Keyword, "const"),
				},
				Parameters: []*Node{
					{
						Values: []*unit{Unit(0, Keyword, "name"), Unit(1, Keyword, "="), Unit(2, Raw, "\"hello\"")},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := extractConst(tt.args.rows, tt.args.idx)
			equal(t, len(tt.want.Values), len(got.Values))
			for i := range tt.want.Values {
				equal(t, *tt.want.Values[i], *got.Values[i])
			}
			equal(t, len(tt.want.Parameters), len(got.Parameters))
			for i := range tt.want.Parameters {
				equal(t, len(tt.want.Parameters[i].Values), len(got.Parameters[i].Values))
				for j := range tt.want.Parameters[i].Values {
					equal(t, *tt.want.Parameters[i].Values[j], *got.Parameters[i].Values[j])
				}
				equal(t, len(tt.want.Parameters[i].Parameters), len(got.Parameters[i].Parameters))
			}
			equal(t, tt.want1, got1)
		})
	}
}
