package path

import (
	"testing"
)

func TestHasWildcardPrefix(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test0",
			args: args{
				pattern: "*",
			},
			want: true,
		},
		{
			name: "test1",
			args: args{
				pattern: "jack*",
			},
			want: false,
		},
		{
			name: "test2",
			args: args{
				pattern: `\*tom`,
			},
			want: false,
		},
		{
			name: "test3",
			args: args{
				pattern: "/test",
			},
			want: false,
		},
		{
			name: "test4",
			args: args{
				pattern: "[t]est",
			},
			want: true,
		},
		{
			name: "test5",
			args: args{
				pattern: "{t,j}est",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		if got := HasWildcardPrefix(tt.args.pattern); got != tt.want {
			t.Errorf("%q. HasWildcardPrefix() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestTrimWildcard(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name            string
		args            args
		wantTrimmed     string
		wantHasWildcard bool
	}{
		{
			name: "test0",
			args: args{
				pattern: "*test",
			},
			wantTrimmed:     "",
			wantHasWildcard: true,
		},
		{
			name: "test1",
			args: args{
				pattern: "test*",
			},
			wantTrimmed:     "test",
			wantHasWildcard: true,
		},
		{
			name: "test2",
			args: args{
				pattern: "te*st",
			},
			wantTrimmed:     "te",
			wantHasWildcard: true,
		},
		{
			name: "test3",
			args: args{
				pattern: "test",
			},
			wantTrimmed:     "test",
			wantHasWildcard: false,
		},
		{
			name: "test4",
			args: args{
				pattern: `test\[]`,
			},
			wantTrimmed:     `test[]`,
			wantHasWildcard: false,
		},
	}
	for _, tt := range tests {
		gotTrimmed, gotHasWildcard := TrimWildcard(tt.args.pattern)
		if gotTrimmed != tt.wantTrimmed {
			t.Errorf("%q. TrimWildcard() gotTrimmed = %v, want %v", tt.name, gotTrimmed, tt.wantTrimmed)
		}
		if gotHasWildcard != tt.wantHasWildcard {
			t.Errorf("%q. TrimWildcard() gotHasWildcard = %v, want %v", tt.name, gotHasWildcard, tt.wantHasWildcard)
		}
	}
}

func TestMatch(t *testing.T) {
	type Result struct {
		Matched bool
	}
	var TestMatchEqual = func(wanted Result, pattern, s string) {
		matched := Match(pattern, s)
		if matched != wanted.Matched {
			t.Errorf("Match(%s, %s) = %v, wanted %v", pattern, s, matched, wanted.Matched)
		}
	}
	TestMatchEqual(Result{true}, `*`, ``)
	TestMatchEqual(Result{false}, `*`, `/`)
	TestMatchEqual(Result{true}, `/*`, `//`)
	TestMatchEqual(Result{true}, `*/`, `debug/`)
	TestMatchEqual(Result{true}, `/*`, `/debug`)
	TestMatchEqual(Result{true}, `/*`, `/debug/`)
	TestMatchEqual(Result{true}, `/*`, `/debug/`)
	TestMatchEqual(Result{true}, `/*`, `/debug/pprof`)
	TestMatchEqual(Result{true}, `/*/`, `/debug/`)
	TestMatchEqual(Result{true}, `/*/*`, `/debug/pprof`)
	TestMatchEqual(Result{true}, `debug/*/`, `debug/test/`)
	TestMatchEqual(Result{true}, `aa/*`, `aa/`) // Wrong
	TestMatchEqual(Result{true}, `**`, ``)
	TestMatchEqual(Result{true}, `/**`, `/debug`)
	TestMatchEqual(Result{true}, `/**`, `/debug/pprof/profile`)
	TestMatchEqual(Result{true}, `/**`, `/debug/pprof/profile/`)
	TestMatchEqual(Result{false}, `/**/profile`, `/debug/pprof/profile/`) // Wrong
	TestMatchEqual(Result{true}, `/**/profile`, `/debug/pprof/profile`)
	TestMatchEqual(Result{true}, `/*/*/profile`, `/debug/pprof/profile`)
	TestMatchEqual(Result{true}, `/**/*`, `/debug/pprof/profile`)
	TestMatchEqual(Result{true}, `/**/pprof/*`, `/debug/pprof/profile`)
	TestMatchEqual(Result{true}, `/**/pprof/*/`, `/debug/pprof/profile/`)
	TestMatchEqual(Result{true}, `/{debug,test}/profile`, `/debug/profile`)
	TestMatchEqual(Result{false}, `/{debug,test}/profile`, `/debug/profile/`)

	TestMatchEqual(Result{true}, `dashboard*.xxxx.com`, `dashboard.xxxx.com`)
	TestMatchEqual(Result{true}, `dashboard{-sit,-prod}.xxxx.com`, `dashboard-sit.xxxx.com`)
	TestMatchEqual(Result{false}, `dashboard{-sit,-prod}.xxxx.com`, `dashboard-si.xxxx.com`)
	TestMatchEqual(Result{true}, `/{config/*,instance}`, `/config/delete`)
	TestMatchEqual(Result{true}, `**`, `/config`)
}
