package meta

import (
	"reflect"
	"testing"
)

func TestQuery_GetArguments(t *testing.T) {
	tests := []struct {
		name  string
		query *Query
		want  []string
	}{
		{
			name: "test0",
			query: &Query{
				Host:   "host",
				Path:   "path",
				Method: "method",
			},
			want: []string{"host", "path", "method"},
		},
	}
	for _, tt := range tests {
		if got := tt.query.GetArguments(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Query.GetArguments() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
