package model

import (
	"errors"
	"strings"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/zutil"
	"github.com/zlsgo/zdb/builder"
	"github.com/zlsgo/zdb/schema"
)

type Migration struct {
	Model
}

const deleteFieldPrefix = "__del__"

func (m *Migration) InitValue() error {
	for _, v := range m.Values {
		data, ok := v.(map[string]interface{})
		if !ok {
			return errors.New("Invalid migration value")
		}

		_, err := m.Model.Insert(data)
		if err != nil {
			return err
		}
	}

	return nil

	// zlog.Debug(i, err)
	// vof := reflect.ValueOf(m.Values)

	// vof.Len()
	// zlog.Debug(vof.Len())

	// for i := 0; i < vof.Len(); i++ {
	// 	zlog.Debug(vof.Index(i).Interface())
	// }
	// zlog.Debug(ztype.GetType(m.Values))
	// return err
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
	// d := table.GetDriver()

	sql, values, process := table.GetColumn()
	res, err := m.DB.QueryToMaps(sql, values...)
	if err != nil {
		return err
	}

	newColumns := zarray.Map(m.Columns, func(i int, v *Column) string {
		return v.Name
	})

	newColumns = append(newColumns, "id")

	{
		if m.Options.SoftDeletes {
			newColumns = append(newColumns, "deleted_at")
		}

		if m.Options.Timestamps {
			newColumns = append(newColumns, "created_at")
			newColumns = append(newColumns, "updated_at")
		}
	}

	oldColumns := zarray.Keys(process(res))

	addColumns := zarray.Filter(newColumns, func(_ int, n string) bool {
		return !zarray.Contains(oldColumns, n)
	})
	deleteColumns := zarray.Filter(oldColumns, func(_ int, n string) bool {
		return !zarray.Contains(newColumns, n) && !strings.HasPrefix(n, deleteFieldPrefix)
	})

	for _, v := range deleteColumns {
		// sql, values = table.RenameColumn(v, deleteFieldPrefix+v)
		// TODO 危险操作，考虑重命名字段
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

	return nil
}

func (m *Migration) CreateTable() error {
	table := builder.NewTable(m.Table.Name).Create()

	fields := make([]*schema.Field, 0, len(m.Columns))

	fields = append(fields, schema.NewField("id", schema.Uint, func(f *schema.Field) {
		f.Comment = "ID"
		f.Size = 64
		f.PrimaryKey = true
		f.AutoIncrement = true
	}))

	for _, v := range m.Columns {
		f := schema.NewField(v.Name, v.Type, func(f *schema.Field) {
			f.Comment = zutil.IfVal(v.Comment != "", v.Comment, v.Label).(string)
			f.NotNull = !v.Nullable
		})
		fields = append(fields, f)
	}

	if m.Options.SoftDeletes {
		fields = append(fields, schema.NewField("deleted_at", "int", func(f *schema.Field) {
			f.NotNull = false
			f.Comment = "删除时间"
		}))
	}

	if m.Options.Timestamps {
		fields = append(fields, schema.NewField("created_at", schema.Time, func(f *schema.Field) {
			f.Comment = "创建时间"
		}))
		fields = append(fields, schema.NewField("updated_at", schema.Time, func(f *schema.Field) {
			f.Comment = "更新时间"
		}))
	}

	table.Column(fields...)

	sql, values := table.Build()
	zlog.Debug(sql, values)
	_, err := m.DB.Exec(sql, values...)
	return err
}
