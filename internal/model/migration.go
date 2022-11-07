package model

import (
	"errors"
	"strings"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/zutil"
	"github.com/zlsgo/zdb"
	"github.com/zlsgo/zdb/builder"
	"github.com/zlsgo/zdb/schema"
)

type Migration struct {
	Model
}

const (
	IDKey        = "id"
	CreatedAtKey = "created_at"
	UpdatedAtKey = "updated_at"
	DeletedAtKey = "deleted_at"
)

func init() {
	zdb.IDKey = IDKey
}

const deleteFieldPrefix = "__del__"

func (m *Migration) Auto() (err error) {
	exist := m.HasTable()
	if !exist {
		zlog.Debug("新建")
		err = m.CreateTable()
		if err != nil {
			return
		}

		zlog.Debug("初始化数据")
		return m.InitValue(true)
	} else {
		zlog.Debug("需要更新表结构")
		err = m.UpdateTable()
	}

	return m.InitValue(false)
}

func (m *Migration) InitValue(all bool) error {
	for _, v := range m.Values {
		data, ok := v.(map[string]interface{})
		if !ok {
			return errors.New("Invalid migration value")
		}
		if !all {
			if _, ok := data[IDKey]; ok {
				continue
			}
		}

		data, err := CheckData(data, m.Columns, activeCreate)
		if err != nil {
			return err
		}
		_, err = m.Model.Insert(data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Model) HasTable() bool {
	table := builder.NewTable(m.Table.Name).Create()

	sql, values, process := table.Has()
	res, err := m.DB.QueryToMaps(sql, values...)
	if err != nil {
		return false
	}

	return process(res)
}

func (m *Migration) UpdateTable() error {
	table := builder.NewTable(m.Table.Name)
	d := table.GetDriver()

	sql, values, process := table.GetColumn()
	res, err := m.DB.QueryToMaps(sql, values...)
	if err != nil {
		return err
	}

	newColumns := zarray.Map(zarray.Filter(m.Columns, func(_ int, _ *Column) bool {
		return true
	}), func(i int, v *Column) string {
		return v.Name
	})

	newColumns = append(newColumns, IDKey)

	currentColumns := process(res)
	oldColumns := zarray.Keys(currentColumns)

	{
		if m.Options.SoftDeletes {
			newColumns = append(newColumns, DeletedAtKey)
		}

		if m.Options.Timestamps {
			if zarray.Contains(oldColumns, CreatedAtKey) {
				newColumns = append(newColumns, CreatedAtKey)
			}
			if zarray.Contains(oldColumns, UpdatedAtKey) {
				newColumns = append(newColumns, UpdatedAtKey)
			}
		}
	}

	updateColumns := zarray.Map(zarray.Filter(m.Columns, func(_ int, n *Column) bool {
		c := currentColumns.Get(n.Name)
		if !c.Exists() {
			return false
		}
		t := strings.ToUpper(d.DataTypeOf(&schema.Field{DataType: schema.DataType(n.Type)}))
		return t != c.Get("type").String()
	}), func(i int, v *Column) string { return v.Name })

	addColumns := zarray.Filter(newColumns, func(_ int, n string) bool {
		return !zarray.Contains(oldColumns, n)
	})

	deleteColumns := zarray.Filter(oldColumns, func(_ int, n string) bool {
		return !zarray.Contains(newColumns, n) && !strings.HasPrefix(n, deleteFieldPrefix)
	})

	for _, v := range deleteColumns {
		// TODO 危险操作，考虑重命名字段
		// sql, values = table.RenameColumn(v, deleteFieldPrefix+v)
		sql, values = table.DropColumn(v)
		_, err := m.DB.Exec(sql, values...)
		if err != nil {
			return err
		}
	}

	if m.Options.Timestamps {
		zlog.Debug(zarray.Contains(oldColumns, CreatedAtKey))
		if !zarray.Contains(oldColumns, CreatedAtKey) {
			sql, values := table.AddColumn(CreatedAtKey, "time", func(f *schema.Field) {
				f.Comment = "更新时间"

			})
			_, err := m.DB.Exec(sql, values...)
			if err != nil {
				return err
			}
		}
		if !zarray.Contains(oldColumns, UpdatedAtKey) {
			sql, values := table.AddColumn(UpdatedAtKey, "time", func(f *schema.Field) {
				f.Comment = "更新时间"
			})
			_, err := m.DB.Exec(sql, values...)
			if err != nil {
				return err
			}
		}
	}

	// TODO 如果有删除字段可以按需恢复
	for _, v := range addColumns {
		c, ok := zarray.Find(m.Columns, func(i int, c *Column) bool {
			return c.Name == v
		})
		if !ok {
			continue
		}

		sql, values := table.AddColumn(v, c.Type, func(f *schema.Field) {
			f.Comment = zutil.IfVal(c.Comment != "", c.Comment, c.Label).(string)
			f.NotNull = !c.Nullable
			f.Size = c.Size
		})
		_, err := m.DB.Exec(sql, values...)
		if err != nil {
			return err
		}
	}

	// TODO 是否需要支持修改字段类型
	for _, v := range updateColumns {
		_ = v
	}

	return nil
}

func (m *Migration) fillField(fields []*schema.Field) []*schema.Field {

	if m.Options.SoftDeletes {
		fields = append(fields, schema.NewField(DeletedAtKey, schema.Int, func(f *schema.Field) {
			f.Size = 9999999999
			f.NotNull = false
			f.Comment = "删除时间"
		}))
	}

	if m.Options.Timestamps {
		fields = append(fields, schema.NewField(CreatedAtKey, schema.Time, func(f *schema.Field) {
			f.Comment = "创建时间"
		}))
		fields = append(fields, schema.NewField(UpdatedAtKey, schema.Time, func(f *schema.Field) {
			f.Comment = "更新时间"
		}))
	}
	return fields
}

func (m *Migration) CreateTable() error {
	table := builder.NewTable(m.Table.Name).Create()

	fields := make([]*schema.Field, 0, len(m.Columns))
	// sideFields := make([]*schema.Field, 0, len(m.Columns))

	fields = append(fields, m.getPrimaryKey())

	for _, v := range m.Columns {
		f := schema.NewField(v.Name, schema.DataType(v.Type), func(f *schema.Field) {
			f.Comment = zutil.IfVal(v.Comment != "", v.Comment, v.Label).(string)
			f.NotNull = !v.Nullable
			f.Size = v.Size
		})

		// if !v.Side {
		fields = append(fields, f)
		// } else {
		// 	sideFields = append(sideFields, f)
		// }
	}

	fields = m.fillField(fields)

	table.Column(fields...)

	sql, values := table.Build()
	_, err := m.DB.Exec(sql, values...)

	// if err == nil && len(sideFields) > 0 {
	// 	err = m.createSideTable(sideFields)
	// }

	return err
}

func (m *Migration) getPrimaryKey() *schema.Field {
	return schema.NewField(IDKey, schema.Uint, func(f *schema.Field) {
		f.Comment = "ID"
		f.PrimaryKey = true
		f.AutoIncrement = true
	})
}

// func (m *Migration) createSideTable(fields []*schema.Field) error {
// 	table := builder.NewTable(m.Table.Name + "__side").Create()

// 	table.Column(m.getPrimaryKey())
// 	table.Column(fields...)

// 	sql, values := table.Build()
// 	_, err := m.DB.Exec(sql, values...)
// 	return err
// }
