package app_test

import (
	"testing"

	"github.com/mikeschinkel/gerardus/app"
)

func Test_stringDiff(t *testing.T) {
	type args struct {
		s1  string
		s2  string
		pad int
	}
	var tests = []struct {
		name string
		args args
		want string
	}{
		{
			name: "S1 is empty",
			args: args{
				s2: "ABC",
			},
			want: "<<2[ABC]2>>",
		},
		{
			name: "S2 is empty",
			args: args{
				s1: "ABC",
			},
			want: "<<1[ABC]1>>",
		},
		{
			name: "S1 and S2 are empty",
		},
		{
			name: "S1 and S2 are completely different",
			args: args{
				s1: "ABC",
				s2: "XYZ",
			},
			want: "<<1[ABC]1>><<2[XYZ]2>>",
		},
		{
			name: "S1 has extra middle chars",
			args: args{
				s1:  "ABCDEF123GHIJKLMNOP",
				s2:  "ABCDEFGHIJKLMNOP",
				pad: 5,
			},
			want: "BCDEF<<1[123]1>>GHIJK",
		},
		//{
		//	name: "Lorem Ipsum",
		//	args: args{
		//		s1:  "In publishing and graphic design, Lorem ipsum is a placeholder text commonly used to demonstrate the visual form of a document or a typeface without relying on meaningful content. Lorem ipsum may be used as a placeholder before final copy is available.",
		//		s2:  "In publishing & graphic design, Lorem ipsum is a commonly used text placeholder to demonstrate a document in its visual form, or a typeface sans meaningful content. Lorem ipsum is often used as a placeholder awaiting final copy.",
		//		pad: 25,
		//	},
		//	want: "",
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.StringDiff(tt.args.s1, tt.args.s2, tt.args.pad); got != tt.want {
				t.Errorf("\nStringDiff(1,2):\n\t got: %v\n\twant: %v\n", got, tt.want)
			}
		})
	}
}
