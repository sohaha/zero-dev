package meta

import (
	"reflect"
	"testing"
)

func TestResource_Match(t *testing.T) {
	type fields struct {
		Host   string
		Path   string
		Method string
	}
	type args struct {
		query *Query
	}
	tests := []struct {
		args    args
		fields  fields
		name    string
		want    bool
		wantErr bool
	}{
		{
			name: "test0",
			fields: fields{
				Host:   "*",
				Path:   "*",
				Method: "*",
			},
			args: args{
				query: &Query{
					Host:   "host",
					Path:   "path",
					Method: "method",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "test1",
			fields: fields{
				Host:   "host",
				Path:   "path",
				Method: "method",
			},
			args: args{
				query: &Query{
					Host:   "host",
					Path:   "path",
					Method: "method",
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		r := &Resource{
			Host:   tt.fields.Host,
			Path:   tt.fields.Path,
			Method: tt.fields.Method,
		}
		got := r.Match(tt.args.query)
		if got != tt.want {
			t.Errorf("%q. Resource.Match() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestResource_GetArguments(t *testing.T) {
	type fields struct {
		Host   string
		Path   string
		Method string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "test0",
			fields: fields{
				Host:   "host",
				Path:   "path",
				Method: "method",
			},
			want: []string{"host", "path", "method"},
		},
		{
			name: "test1",
			fields: fields{
				Host:   "",
				Path:   "",
				Method: "",
			},
			want: []string{"", "", ""},
		},
	}
	for _, tt := range tests {
		r := &Resource{
			Host:   tt.fields.Host,
			Path:   tt.fields.Path,
			Method: tt.fields.Method,
		}
		if got := r.GetArguments(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Resource.GetArguments() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestResource_IsValid(t *testing.T) {
	type fields struct {
		Host   string
		Path   string
		Method string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "test0",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "test1",
			fields: fields{
				Host:   "host",
				Path:   "path",
				Method: "method",
			},
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				Host:   "host",
				Path:   "",
				Method: "method",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		r := &Resource{
			Host:   tt.fields.Host,
			Path:   tt.fields.Path,
			Method: tt.fields.Method,
		}
		if err := r.IsValid(); (err != nil) != tt.wantErr {
			t.Errorf("%q. Resource.IsValid() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
