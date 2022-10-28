package util

import "testing"

func TestContains(t *testing.T) {
	type args struct {
		arr []string
		s   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test0",
			args: args{
				arr: []string{"a", "b"},
				s:   "b",
			},
			want: true,
		},
		{
			name: "test1",
			args: args{
				arr: []string{"a"},
				s:   "b",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		if got := Contains(tt.args.arr, tt.args.s); got != tt.want {
			t.Errorf("%q. Contains() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
