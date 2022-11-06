package account

import (
	"zlsapp/app/model"

	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/zdb"
)

func migration(di zdi.Invoker) (m *model.Model, err error) {
	_, diErr := di.Invoke(func(db *zdb.DB) {
		json, _ := zjson.SetBytes([]byte("{}"), "name", ztype.Map{})
		json, _ = zjson.SetBytes(json, "table", ztype.Map{
			"name":    "account_user",
			"comment": "用户表",
		})

		json, _ = zjson.SetBytes(json, "options", ztype.Map{
			"timestamps": true,
			"crypt_id":   true,
		})

		json, _ = zjson.SetBytes(json, "values", ztype.Maps{
			{
				"account":  "admin",
				"password": "admin",
				"status":   "1",
				"roles":    "admin",
			},
		})

		zlog.Success("初始化用户")
		zlog.Printf("        账号: %s\n        密码: %s\n", "admin", "admin")

		json, _ = zjson.SetBytes(json, "columns", ztype.Maps{
			{
				"label":    "头像",
				"name":     "avatar",
				"nullable": true,
				"type":     "string",
				"validations": ztype.Maps{
					{
						"method": "url",
					},
				},
			},
			{
				"name":  "account",
				"type":  "string",
				"label": "账号",
				"validations": ztype.Maps{
					{
						"method": "minLength",
						"args":   3,
					},
					{
						"method": "maxLength",
						"args":   10,
					},
				},
			},
			{
				"name":  "password",
				"type":  "string",
				"label": "密码",
				"crypt": "PASSWORD",
				"validations": ztype.Maps{
					{
						"method": "minLength",
						"args":   3,
					},
					{
						"method": "maxLength",
						"args":   20,
					},
				},
			},
			{
				"name":  "status",
				"type":  "int",
				"size":  10,
				"label": "状态",
				"validations": ztype.Maps{
					{
						"method": "enum",
						"args":   []string{"0", "1", "2"},
					},
				},
			},
			{
				"name":     "salt",
				"type":     "string",
				"size":     8,
				"nullable": true,
				"label":    "盐",
			},
			{
				"name":     "login_time",
				"type":     "time",
				"nullable": true,
				"label":    "登录时间",
			},
			{
				"name":     "roles",
				"type":     "string",
				"nullable": true,
				"size":     200,
				"label":    "角色",
			},
		})

		m, err = model.Add(db, "account", json)
		zerror.Panic(err)

		zerror.Panic(m.Migration().Auto())

	})

	if diErr != nil {
		return nil, diErr
	}

	return
}
