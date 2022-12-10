package account

import (
	"zlsapp/internal/parse"

	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/zdb"
)

const UsersModel = "inlay::accounts"

func userModel(db *zdb.DB) error {
	json, _ := zjson.SetBytes([]byte("{}"), "name", ztype.Map{})
	json, _ = zjson.SetBytes(json, "name", "账号模型")
	json, _ = zjson.SetBytes(json, "table", ztype.Map{
		"name":    "account_users",
		"comment": "用户表",
	})

	json, _ = zjson.SetBytes(json, "options", ztype.Map{
		"timestamps": true,
		"crypt_id":   true,
	})

	defAccount := defaultAccount()
	json, _ = zjson.SetBytes(json, "values", defAccount)
	json, _ = zjson.SetBytes(json, "columns", ztype.Maps{
		{
			"label":    "头像",
			"name":     "avatar",
			"nullable": true,
			"type":     "string",
			"size":     1024 * 2,
			"validations": ztype.Maps{
				{
					"method": "regex",
					"args":   "^(data:image/|http://|https://)",
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
			"name":  "nickname",
			"type":  "string",
			"size":  20,
			"label": "昵称",
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
					"args":   250,
				},
			},
		},
		{
			"name":  "status",
			"type":  "int8",
			"size":  9,
			"label": "状态",
			"enum": []parse.ColumnEnum{
				{Value: "0", Label: "待激活"},
				{Value: "1", Label: "正常"},
				{Value: "2", Label: "禁用"},
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
			"name":     "remark",
			"type":     "string",
			"size":     100,
			"default":  "",
			"nullable": true,
			"label":    "备注",
		},
		{
			"name":     "roles",
			"type":     "json",
			"nullable": true,
			"size":     200,
			"label":    "角色",
		},
	})

	m, err := parse.AddModelForJSON(UsersModel, json, func(m *parse.Modeler) (parse.Storageer, error) {
		return parse.NewSQL(db, m.Table.Name), nil
	}, false)

	if err != nil {
		return err
	}

	if !m.Migration().HasTable() {
		zlog.Success("初始化管理账号：")
		for _, v := range defAccount {
			zlog.Printf("        账号：%s 密码：%s\n", v["account"], v["password"])
		}
	}

	err = m.Migration().Auto(true)
	if err != nil {
		return zerror.With(err, "用户模型初始化失败")
	}
	return nil
}

const LogsModel = "inlay::logs"

func logsModel(db *zdb.DB) error {
	json, _ := zjson.SetBytes([]byte("{}"), "name", ztype.Map{})
	json, _ = zjson.SetBytes(json, "name", "日志模型")
	json, _ = zjson.SetBytes(json, "table", ztype.Map{
		"name":    "account_logs",
		"comment": "日志表",
	})

	json, _ = zjson.SetBytes(json, "options", ztype.Map{
		"timestamps": true,
		"crypt_id":   true,
	})

	json, _ = zjson.SetBytes(json, "columns", ztype.Maps{
		{
			"name":  "action",
			"type":  "string",
			"label": "操作",
			"validations": ztype.Maps{
				{
					"method": "minLength",
					"args":   1,
				},
				{
					"method": "maxLength",
					"args":   60,
				},
			},
		},
		{
			"name":  "uid",
			"type":  "string",
			"label": "操作用户",
			"validations": ztype.Maps{
				{
					"method": "minLength",
					"args":   1,
				},
			},
		},
		{
			"name":     "ip",
			"type":     "string",
			"size":     100,
			"default":  "",
			"nullable": true,
			"validations": []ztype.Map{
				{"method": "ip"},
			},
			"label": "IP",
		},
		// {
		// 	"name":     "user_agent",
		// 	"type":     "string",
		// 	"size":     250,
		// 	"default":  "",
		// 	"nullable": true,
		// 	"label":    "user_agent",
		// },
		{
			"name":     "device",
			"type":     "string",
			"size":     10,
			"default":  "",
			"nullable": true,
			"label":    "操作设备",
		},
		{
			"name":     "os",
			"type":     "string",
			"size":     10,
			"default":  "",
			"nullable": true,
			"label":    "操作系统",
		},
		{
			"name":     "os_version",
			"type":     "string",
			"size":     10,
			"default":  "",
			"nullable": true,
			"label":    "系统版本",
		},
		{
			"name":    "status",
			"type":    "int8",
			"size":    9,
			"label":   "状态",
			"default": LogsStatusRead,
			"enum": []parse.ColumnEnum{
				{Value: ztype.ToString(LogsStatusUnread), Label: "未读"},
				{Value: ztype.ToString(LogsStatusRead), Label: "已读"},
			},
		},
		{
			"name":     "remark",
			"type":     "string",
			"size":     100,
			"default":  "",
			"nullable": true,
			"label":    "备注",
		},
	})

	m, err := parse.AddModelForJSON(LogsModel, json, func(m *parse.Modeler) (parse.Storageer, error) {
		return parse.NewSQL(db, m.Table.Name), nil
	}, false)
	if err != nil {
		return err
	}

	return m.Migration().Auto(true)
}

const RolesModel = "inlay::roles"

func roleModel(db *zdb.DB) error {
	m := &parse.Modeler{
		Name: "角色模型",
		Table: parse.Table{
			Name:    "account_roles",
			Comment: "角色表",
		},
		Columns: []*parse.Column{
			{
				Name:   "name",
				Type:   "string",
				Label:  "角色名称",
				Unique: true,
				Size:   20,
			},
			{
				Name:   "key",
				Type:   "string",
				Label:  "角色标识",
				Unique: true,
				Size:   20,
				Validations: []parse.Validations{
					{
						Method:  "regex",
						Args:    "^[a-zA-Z0-9_]+$",
						Message: "角色标识只能包含字母、数字、下划线",
					},
				},
			},

			{
				Name:     "menu",
				Type:     "json",
				Label:    "角色菜单",
				Default:  "[]",
				Nullable: true,
			},
		},
		Options: parse.Options{
			Timestamps: true,
		},
		Values: []map[string]interface{}{
			{
				parse.IDKey: 1,
				"name":      "管理员",
				"key":       "admin",
			},
		},
	}
	err := parse.AddModel(RolesModel, m, func(m *parse.Modeler) (parse.Storageer, error) {
		return parse.NewSQL(db, m.Table.Name), nil
	}, false)
	if err != nil {
		return err
	}

	return m.Migration().Auto(true)
}
