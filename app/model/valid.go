package model

import (
	"errors"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/sohaha/zlsgo/zvalid"
)

type validations struct {
	Method  string      `json:"method"`
	Message string      `json:"message"`
	Args    interface{} `json:"args"`
}

type activeType uint

const (
	activeCreate activeType = iota + 1
	activeUpdate
)

// filterDate 过滤数据字段
func filterDate(data ztype.Map, fields []string) ztype.Map {
	l := len(fields)
	if l == 0 {
		return data
	}

	n := make(ztype.Map, len(data))
	for k := range data {
		if !zarray.Contains(fields, k) {
			n[k] = data[k]
		}
	}

	return n
}

// CheckData 验证数据
func CheckData(data ztype.Map, columns []*Column, active activeType) (ztype.Map, error) {
	d := make(ztype.Map, len(columns))
	for _, column := range columns {
		name := column.Name
		label := column.Label
		if label == "" {
			label = name
		}

		v, ok := data[name]

		{
			if !ok {
				if column.Default != nil {
					v = column.Default
					ok = true
				}
			}
			if !ok && active != activeUpdate && !column.Nullable {
				return d, errors.New(label + "不能为空")
			}
		}

		if ok {
			typ := column.Type
			switch typ {
			case "bool":
				d[name] = ztype.ToBool(v)
			case "time":
				t, _ := v.(string)
				parse, err := ztime.Parse(t)
				if err != nil {
					return d, errors.New(label + ": 时间格式错误")
				}
				d[name] = DataTime{Time: parse}
			case "int", "uint", "float", "string":
				var (
					val interface{}
					err error
				)
				switch typ {
				case "string":
					val, err = validRule(label, v, column.Validations, column.Size).String()
					if val == "" && !column.Nullable {
						return d, errors.New(label + "不能为空")
					}
				default:
					rule := validRule(label, v, column.Validations, column.Size).IsNumber()
					switch typ {
					case "int":
						val, err = rule.Int()
					case "uint":
						val = ztype.ToUint(rule.Value())
					default:
						val, err = rule.Float64()
					}

				}
				if err != nil {
					return d, err
				}
				d[name] = val
			default:
				d[name] = v
			}
		}
	}
	return d, nil
}

var inlayRules = map[string]func(label string, rule *zvalid.Engine, valid validations) *zvalid.Engine{
	"regex": func(label string, rule *zvalid.Engine, valid validations) *zvalid.Engine {
		rule.Regex(ztype.ToString(valid.Args), valid.Message)
		return rule
	},
	"json": func(label string, rule *zvalid.Engine, valid validations) *zvalid.Engine {
		rule.IsJSON(valid.Message)
		return rule
	},
	"enum": func(label string, rule *zvalid.Engine, valid validations) *zvalid.Engine {
		switch val := valid.Args.(type) {
		case []float64:
			rule.EnumFloat64(val)
		case []string:
			rule.EnumString(val)
		case []int:
			rule.EnumInt(val)
		default:
			rule.Customize(func(rawValue string, err error) (string, error) {
				return "", errors.New(label + "枚举值不在合法范围")
			})
		}
		return rule
	},
	"mobile": func(label string, rule *zvalid.Engine, valid validations) *zvalid.Engine {
		rule.IsMobile(valid.Message)
		return rule
	},
	"email": func(label string, rule *zvalid.Engine, valid validations) *zvalid.Engine {
		rule.IsMail(valid.Message)
		return rule
	},
	"url": func(label string, rule *zvalid.Engine, valid validations) *zvalid.Engine {
		rule.IsURL(valid.Message)
		return rule
	},
	"ip": func(label string, rule *zvalid.Engine, valid validations) *zvalid.Engine {
		rule.IsIP(valid.Message)
		return rule
	},
	"minLength": func(label string, rule *zvalid.Engine, valid validations) *zvalid.Engine {
		rule.MinUTF8Length(ztype.ToInt(valid.Args), valid.Message)
		return rule
	},
	"maxLength": func(label string, rule *zvalid.Engine, valid validations) *zvalid.Engine {
		rule.MaxUTF8Length(ztype.ToInt(valid.Args), valid.Message)
		return rule
	},
	"min": func(label string, rule *zvalid.Engine, valid validations) *zvalid.Engine {
		rule.MinFloat(ztype.ToFloat64(valid.Args), valid.Message)
		return rule
	},
	"max": func(label string, rule *zvalid.Engine, valid validations) *zvalid.Engine {
		rule.MaxFloat(ztype.ToFloat64(valid.Args), valid.Message)
		return rule
	},
}

func validRule(label string, v interface{}, valids []validations, max uint) zvalid.Engine {
	rule := zvalid.New().VerifiAny(v, label)

	for _, valid := range valids {
		r, ok := inlayRules[valid.Method]
		if ok {
			r(label, &rule, valid)
		} else {
			fn, ok := valid.Args.(func(label string, rule zvalid.Engine, valid validations) zvalid.Engine)
			if ok {
				rule = fn(label, rule, valid)
			}
		}
	}

	if max > 0 {
		rule = rule.MaxUTF8Length(int(max))
	}
	return rule
}

func validRule2(label string, v interface{}, valids []validations, max uint) zvalid.Engine {
	rule := zvalid.New().VerifiAny(v, label)

	for _, valid := range valids {
		switch valid.Method {
		case "regex":
			rule = rule.Regex(ztype.ToString(valid.Args), valid.Message)
		case "json":
			rule = rule.IsJSON(valid.Message)
		case "enum":
			switch val := valid.Args.(type) {
			case []float64:
				rule = rule.EnumFloat64(val)
			case []string:
				rule = rule.EnumString(val)
			case []int:
				rule = rule.EnumInt(val)
			default:
				rule = rule.Customize(func(rawValue string, err error) (string, error) {
					return "", errors.New(label + "枚举值不在合法范围")
				})
			}
		case "mobile":
			rule = rule.IsMobile(valid.Message)
		case "mail":
			rule = rule.IsMail(valid.Message)
		case "url":
			rule = rule.IsURL(valid.Message)
		case "ip":
			rule = rule.IsIP(valid.Message)
		case "minLength":
			rule = rule.MinUTF8Length(ztype.ToInt(valid.Args), valid.Message)
		case "maxLength":
			rule = rule.MaxUTF8Length(ztype.ToInt(valid.Args), valid.Message)
		case "min":
			rule = rule.MinFloat(ztype.ToFloat64(valid.Args), valid.Message)
		case "max":
			rule = rule.MaxFloat(ztype.ToFloat64(valid.Args), valid.Message)
		}
	}

	if max > 0 {
		rule = rule.MaxUTF8Length(int(max))
	}
	return rule
}
