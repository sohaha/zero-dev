package model

import (
	"errors"
	"strings"

	"github.com/sohaha/zlsgo/zarray"
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

func (m *Migration) InitValue() error {
	for _, v := range m.Values {
		data, ok := v.(map[string]interface{})
		if !ok {
			return errors.New("Invalid migration value")
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

	{
		if m.Options.SoftDeletes {
			newColumns = append(newColumns, DeletedAtKey)
		}

		if m.Options.Timestamps {
			newColumns = append(newColumns, CreatedAtKey)
			newColumns = append(newColumns, UpdatedAtKey)
		}
	}

	currentColumns := process(res)
	oldColumns := zarray.Keys(currentColumns)

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

	// TODO 如果有删除字段可以按需恢复
	for _, v := range addColumns {
		c, ok := zarray.Find(m.Columns, func(i int, c *Column) bool {
			return c.Name == v
		})
		if !ok {
			continue
		}

		sql, values := table.AddColumn(v, c.Type)
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

func (m *Migration) CreateTable() error {
	table := builder.NewTable(m.Table.Name).Create()

	fields := make([]*schema.Field, 0, len(m.Columns))
	// sideFields := make([]*schema.Field, 0, len(m.Columns))

	fields = append(fields, m.getPrimaryKey())

	for _, v := range m.Columns {
		f := schema.NewField(v.Name, v.Type, func(f *schema.Field) {
			f.Comment = zutil.IfVal(v.Comment != "", v.Comment, v.Label).(string)
			f.NotNull = !v.Nullable
		})
		// if !v.Side {
		fields = append(fields, f)
		// } else {
		// 	sideFields = append(sideFields, f)
		// }
	}

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
