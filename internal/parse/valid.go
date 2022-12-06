package parse

import (
	"errors"
	"time"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/sohaha/zlsgo/zvalid"
	"github.com/zlsgo/zdb/schema"
)

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

// VerifiData 验证数据
func VerifiData(data ztype.Map, columns []*Column, active activeType) (ztype.Map, error) {
	d := make(ztype.Map, len(columns))
	for _, column := range columns {
		if active == activeUpdate && column.ReadOnly {
			continue
		}

		name, label := column.Name, column.Label
		if label == "" {
			label = name
		}

		v, ok := data[name]

		{
			if !ok && active != activeUpdate {
				if column.Default != nil {
					v = column.Default
					ok = true
				}
			}
			if !ok && !column.Nullable && active != activeUpdate {
				return d, errors.New(label + "不能为空")
			}
		}

		if ok {
			if v == nil {
				if column.Nullable {
					continue
				} else {
					return d, errors.New(label + "不能为 null")
				}
			}
			typ := column.Type
			switch typ {
			case schema.Bool:
				d[name] = ztype.ToBool(v)
			case schema.Time:
				switch t := v.(type) {
				default:
					return d, errors.New(label + ": 未知时间格式")
				case DataTime:
					d[name] = t
				case time.Time:
					d[name] = DataTime{Time: t}
				case int:
					d[name] = DataTime{Time: ztime.Unix(ztype.ToInt64(v))}
				case string:
					r, err := ztime.Parse(t)
					if err != nil {
						return d, errors.New(label + ": 时间格式错误")
					}
					d[name] = DataTime{Time: r}
				}
			case schema.JSON:
				err := column.GetValidations().VerifiAny(v).Error()
				if err != nil {
					return d, err
				}
				d[name] = v
			default:
				var (
					val interface{}
					err error
				)
				switch typ {
				case schema.String:
					val, err = column.GetValidations().VerifiAny(v).String()
					if val == "" && !column.Nullable {
						return d, errors.New(label + "不能为空")
					}
				default:
					rule := column.GetValidations().VerifiAny(v).IsNumber()
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
			}
		}
	}

	return d, nil
}

func resolverColumnOptions(c *Column) {
	if len(c.Options) > 0 {
		c.Options = zarray.Map(c.Options, func(_ int, v ColumnEnum) ColumnEnum {
			if v.Label == "" {
				v.Label = v.Value
			}
			return v
		})
		c.validRules = c.validRules.EnumString(zarray.Map(c.Options, func(_ int, v ColumnEnum) string {
			return v.Value
		}))
	}
}

func resolverValidRule(c *Column) {
	label := c.GetLabel()
	rule := zvalid.New().SetAlias(label)

	if c.Type == schema.JSON {
		rule = rule.Required().IsJSON(c.Name + "必须是JSON格式")
	}

	if c.Size > 0 {
		switch c.Type {
		case schema.JSON:
		case schema.String:
			rule = rule.MaxUTF8Length(int(c.Size))
		case schema.Int, schema.Int8, schema.Int16, schema.Int32, schema.Int64, schema.Uint, schema.Uint8, schema.Uint16, schema.Uint32, schema.Uint64:
			rule = rule.MaxInt(int(c.Size))
		case schema.Float:
			rule = rule.MaxFloat(float64(c.Size))
		case schema.Time:
			rule.Customize(func(rawValue string, err error) (newValue string, newErr error) {
				if err != nil {
					return "", err
				}
				if ztime.Unix(int64(c.Size)).After(time.Now()) {
					return rawValue, errors.New(label + "时间不能大于指定时间")
				}
				return
			})
		}
	}

	for _, valid := range c.Validations {
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
					ok := zarray.Contains(ztype.ToSlice(val).String(), rawValue)
					if !ok {
						return "", errors.New(label + "枚举值不在合法范围")
					}
					return rawValue, nil
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

	c.validRules = rule
}
