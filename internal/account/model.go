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
	m := &parse.Modeler{
		Name: "日志模型",
		Table: parse.Table{
			Name:    "account_logs",
			Comment: "日志表",
		},
		Options: parse.Options{
			Timestamps: true,
			CryptID:    true,
		},
	}

	m.Columns = []*parse.Column{
		{
			Name:    "action",
			Type:    "string",
			Label:   "操作",
			Default: "",
			Validations: []parse.Validations{
				{
					Method: "minLength",
					Args:   1,
				},
				{
					Method: "maxLength",
					Args:   60,
				},
			},
		},
		{
			Name:  "uid",
			Type:  "string",
			Label: "操作用户",
			Validations: []parse.Validations{
				{
					Method: "minLength",
					Args:   1,
				},
			},
		},
		{
			Name:    "ip",
			Type:    "string",
			Label:   "请求 IP",
			Size:    100,
			Default: "",
			Validations: []parse.Validations{
				{Method: "ip"},
			},
		},
		{
			Name:    "ip_region",
			Type:    "string",
			Label:   "IP 归属地",
			Size:    100,
			Default: "",
		},
		{
			Name:    "method",
			Type:    "string",
			Label:   "请求方法",
			Size:    10,
			Default: "",
		},
		{
			Name:  "path",
			Type:  "string",
			Label: "请求路径",
			Size:  200,
		},
		{
			Name:    "device",
			Type:    "string",
			Label:   "操作设备",
			Size:    10,
			Default: "",
		},
		{
			Name:    "browser",
			Type:    "string",
			Label:   "浏览器",
			Size:    10,
			Default: "",
		}, {
			Name:    "browser_version",
			Type:    "string",
			Label:   "浏览器版本",
			Size:    20,
			Default: "",
		}, {
			Name:    "os",
			Type:    "string",
			Label:   "操作系统",
			Size:    10,
			Default: "",
		},
		{
			Name:    "os_version",
			Type:    "string",
			Label:   "系统版本",
			Size:    10,
			Default: "",
		},
		{
			Name:    "module",
			Type:    "string",
			Label:   "访问模块",
			Size:    20,
			Default: "",
		},
		{
			Name:    "result",
			Type:    "bool",
			Label:   "操作状态",
			Default: true,
		},
		{
			Name:    "status",
			Type:    "int8",
			Label:   "状态",
			Size:    9,
			Default: LogsStatusRead,
			Options: []parse.ColumnEnum{
				{Value: ztype.ToString(LogsStatusUnread), Label: "未读"},
				{Value: ztype.ToString(LogsStatusRead), Label: "已读"},
			},
		},
		{
			Name:    "category",
			Type:    "int8",
			Label:   "类别",
			Default: LogTypeCommon,
			Options: []parse.ColumnEnum{
				{Value: ztype.ToString(LogTypeCommon), Label: "普通日志"},
				{Value: ztype.ToString(LogTypeLogin), Label: "登录日志"},
				{Value: ztype.ToString(LogTypeAction), Label: "操作日志"},
			},
		},
		{
			Name:    "detail",
			Type:    "string",
			Label:   "详情",
			Size:    200,
			Default: "",
		},
	}

	m.Relations = map[string]*parse.ModelRelation{
		"user": {
			Name:    "user",
			Key:     "uid",
			Model:   UsersModel,
			Foreign: parse.IDKey,
			Fields:  []string{"nickname", "account"},
		},
	}
	err := parse.AddModel(LogsModel, m, func(m *parse.Modeler) (parse.Storageer, error) {
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
