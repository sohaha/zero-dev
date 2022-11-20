package parse

import (
	"errors"
	"strings"

	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/ztype"
)

type beforeProcess func(interface{}) (string, error)

func (m *Modeler) GetBeforeProcess(p []string) (fn []beforeProcess, err error) {
	for _, v := range p {
		switch strings.ToLower(v) {
		default:
			return nil, errors.New("before name not found")
		case "json":
			fn = append(fn, func(s interface{}) (string, error) {
				zlog.Debug(s, ztype.GetType(s))
				switch v := s.(type) {
				case string:
					if zjson.Valid(v) {
						return v, nil
					}
					return " ", nil
				case []interface{}:
				case map[string]interface{}:
				default:
					return " ", nil
				}
				j, err := zjson.Marshal(s)
				return zstring.Bytes2String(j), err
			})
		}
	}
	return
}

func (m *Modeler) valuesBeforeProcess(data ztype.Map) (newData ztype.Map, err error) {
	for k := range m.cryptKeys {
		if _, ok := data[k]; ok {
			data[k], err = m.cryptKeys[k](data.Get(k).String())
			if err != nil {
				return nil, err
			}
		}
	}

	for name, fns := range m.beforeProcess {
		val := data.Get(name)
		if !val.Exists() {
			continue
		}
		var v interface{}
		v = val.Value()
		for _, fn := range fns {
			v, err = fn(v)
			if err != nil {
				return data, err
			}
		}
		_ = data.Set(name, v)
	}
	newData = data
	return
}

type afterProcess func(string) (interface{}, error)

func (m *Modeler) GetAfterProcess(p []string) (fn []afterProcess, err error) {
	for _, v := range p {
		switch strings.ToLower(v) {
		default:
			return nil, errors.New("after name not found")
		case "json":
			fn = append(fn, func(s string) (interface{}, error) {
				j := zjson.Parse(s)
				if !j.Exists() {
					return nil, errors.New("json parse error")
				}
				if j.IsArray() {
					return j.Slice().Value(), nil
				}
				return j.MapString(), nil
			})
		}
	}
	return
}
