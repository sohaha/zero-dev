package model

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/xeipuuv/gojsonschema"
)

//go:embed model_schema.json
var modelSchema []byte

func getSchema(schema []byte) []byte {
	var j ztype.Map
	err := zjson.Unmarshal(schema, &j)
	if err == nil {
		var removeKey func(ztype.Map) ztype.Map
		removeKey = func(j ztype.Map) ztype.Map {
			for k, v := range j {
				if strings.HasPrefix(k, "x-api") {
					delete(j, k)
					continue
				}
				m, ok := v.(map[string]interface{})
				if ok {
					_ = j.Set(k, removeKey(m))
				}
			}
			return j
		}
		j = removeKey(j)
		b, err := zjson.Marshal(j)
		if err == nil {
			return b
		}
	}
	return schema
}

func GetModelSchema() []byte {
	return getSchema(modelSchema)
}

var jsonschemaLoader, _ = gojsonschema.NewSchema(gojsonschema.NewStringLoader(zstring.Bytes2String(modelSchema)))

func ValidateModelSchema(data []byte) error {
	res, err := jsonschemaLoader.Validate(gojsonschema.NewBytesLoader(data))
	if err != nil {
		return err
	}

	if !res.Valid() {
		return fmt.Errorf("%s", res.Errors())
	}
	return nil
}
