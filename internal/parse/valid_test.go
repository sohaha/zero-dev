package parse

import (
	"reflect"
	"sync"
	"testing"

	"github.com/sohaha/zlsgo/ztype"
)

var testColumns = []*Column{
	{
		Label:    "用户名",
		Name:     "username",
		Type:     "string",
		Size:     10,
		Nullable: false,
		Validations: []Validations{
			{
				Method: "minLength",
				Args:   "5",
			},
		},
	},
	{
		Name:  "age",
		Label: "年龄",
		Type:  "int",
		Validations: []Validations{
			{
				Method: "min",
				Args:   "18",
			},
			{
				Method: "max",
				Args:   "200",
			},
		},
	},
	{
		Name:     "gender",
		Type:     "float",
		Nullable: true,
		Validations: []Validations{
			{
				Method: "enum",
				Args:   []float64{1.0, 1.1},
			},
		},
		Label: "性别",
	},
	{
		Name:     "login_ip",
		Type:     "string",
		Nullable: true,
		Validations: []Validations{
			{
				Method: "ip",
			},
		},
		Label: "登录IP",
	},
}

func TestCheckData(t *testing.T) {
	type (
		args struct {
			data    ztype.Map
			columns []*Column
			active  activeType
		}
		test struct {
			want    ztype.Map
			name    string
			args    args
			wantErr bool
		}
	)

	columns := testColumns

	var tests []test

	{
		tests = append(tests, test{
			name: "空数据",
			args: args{
				data:    nil,
				columns: columns,
				active:  activeCreate,
			},
			want:    nil,
			wantErr: true,
		})

		tests = append(tests, test{
			name: "正常",
			args: args{
				data: map[string]interface{}{
					"username": "admin",
					"age":      "18",
				},
				columns: columns,
				active:  activeCreate,
			},
			want: map[string]interface{}{
				"username": "admin",
				"age":      18,
			},
		})
		tests = append(tests, test{
			name: "用户名最大长度",
			args: args{
				data: map[string]interface{}{
					"username": "admin1234567890",
				},
				columns: columns,
				active:  activeCreate,
			},
			wantErr: true,
		})
		tests = append(tests, test{
			name: "用户名最小长度",
			args: args{
				data: map[string]interface{}{
					"username": "a",
				},
				columns: columns,
				active:  activeCreate,
			},
			wantErr: true,
		})
	}

	{
		tests = append(tests, test{
			name: "年龄非数字值",
			args: args{
				data: map[string]interface{}{
					"username": "admin",
					"age":      "xxx",
				},
				columns: columns,
				active:  activeCreate,
			},
			want: map[string]interface{}{
				"username": "admin",
			},
			wantErr: true,
		})
		tests = append(tests, test{
			name: "年龄空值",
			args: args{
				data: map[string]interface{}{
					"username": "admin",
				},
				columns: columns,
				active:  activeCreate,
			},
			want: map[string]interface{}{
				"username": "admin",
			},
			wantErr: true,
		})
		tests = append(tests, test{
			name: "年龄空值-更新",
			args: args{
				data: map[string]interface{}{
					"username": "admin",
				},
				columns: columns,
				active:  activeUpdate,
			},
			want: map[string]interface{}{
				"username": "admin",
			},
			wantErr: false,
		})
		tests = append(tests, test{
			name: "年龄零值-更新",
			args: args{
				data: map[string]interface{}{
					"username": "admin",
					"age":      0,
				},
				columns: columns,
				active:  activeUpdate,
			},
			wantErr: true,
		})
		tests = append(tests, test{
			name: "年龄最大值",
			args: args{
				data: map[string]interface{}{
					"username": "admin",
					"age":      1000,
				},
				columns: columns,
				active:  activeUpdate,
			},
			wantErr: true,
		})
	}

	{
		tests = append(tests, test{
			name: "性别枚举值",
			args: args{
				data: map[string]interface{}{
					"username": "admin",
					"age":      "18",
					"test":     "xxx",
					"gender":   "1.1",
				},
				columns: columns,
				active:  activeCreate,
			},
			want: map[string]interface{}{
				"username": "admin",
				"age":      18,
				"gender":   1.1,
			},
		})
		tests = append(tests, test{
			name: "性别非枚举值",
			args: args{
				data: map[string]interface{}{
					"username": "admin",
					"age":      "18",
					"test":     "xxx",
					"gender":   111,
				},
				columns: columns,
				active:  activeCreate,
			},
			wantErr: true,
		})
	}

	{
		tests = append(tests, test{
			name: "登录IP正确",
			args: args{
				data: map[string]interface{}{
					"username": "admin",
					"age":      18,
					"login_ip": "192.168.3.3",
				},
				columns: columns,
			},
			want: map[string]interface{}{
				"username": "admin",
				"age":      18,
				"login_ip": "192.168.3.3",
			},
		})
		tests = append(tests, test{
			name: "登录IP错误",
			args: args{
				data: map[string]interface{}{
					"username": "admin",
					"age":      18,
					"login_ip": "This is IP",
				},
				columns: columns,
			},
			wantErr: true,
		})
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			got, err := VerifiData(tt.args.data, tt.args.columns, tt.args.active)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("error = %v", err)
				}
			}

			if tt.want != nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("data = %+v, want %+v", got, tt.want)
			}
		})
	}
	wg.Wait()
}
