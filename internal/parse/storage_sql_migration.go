package parse

import (
	"errors"
	"strings"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/sohaha/zlsgo/zutil"
	"github.com/zlsgo/zdb"
	"github.com/zlsgo/zdb/builder"
	"github.com/zlsgo/zdb/schema"
)

type Migration struct {
	Model *Modeler
	DB    *zdb.DB
}

func (m *Migration) Auto(deleteColumn bool) (err error) {
	if m.Model.Table.Name == "" {
		return errors.New("表名不能为空")
	}

	exist := m.HasTable()

	defer func() {
		if err == nil {
			err = m.Indexs()
			if err == nil {
				err = m.InitValue(!exist)
			}
		}
	}()

	if !exist {
		err = m.CreateTable()
		return
	}

	err = m.UpdateTable(deleteColumn)

	return
}

func (m *Migration) InitValue(all bool) error {
	if !all {
		row, err := FindOne(m.Model, ztype.Map{}, func(o *StorageOptions) error {
			o.Fields = []string{"COUNT(*) AS count"}
			return nil
		})
		if err != nil {
			all = row.Get("count").Int() == 0
		}
	}
	for _, data := range m.Model.Values {
		// data, ok := v.(map[string]interface{})
		// if !ok {
		// 	return errors.New("初始化数据格式错误")
		// }
		if !all {
			if _, ok := data[IDKey]; ok {
				continue
			}
		}
		uid := ""
		_, err := Insert(m.Model, data, uid)
		if err != nil {
			return zerror.With(err, "初始化数据失败")
		}
	}

	return nil
}

func (m *Migration) HasTable() bool {
	table := builder.NewTable(m.Model.Table.Name).Create()

	sql, values, process := table.Has()
	res, err := m.DB.QueryToMaps(sql, values...)

	if err != nil {
		return false
	}

	return process(res)
}

func (m *Migration) UpdateTable(deleteColumn bool) error {
	table := builder.NewTable(m.Model.Table.Name)
	// d := table.GetDriver()

	sql, values, process := table.GetColumn()
	res, err := m.DB.QueryToMaps(sql, values...)
	if err != nil {
		return err
	}

	newColumns := zarray.Map(zarray.Filter(m.Model.Columns, func(_ int, _ *Column) bool {
		return true
	}), func(i int, v *Column) string {
		return v.Name
	})

	newColumns = append(newColumns, IDKey)

	currentColumns := process(res)
	oldColumns := zarray.Keys(currentColumns)

	{
		if m.Model.Options.SoftDeletes {
			newColumns = append(newColumns, DeletedAtKey)
		}

		if m.Model.Options.CreatedBy {
			newColumns = append(newColumns, CreatedByKey)
		}

		if m.Model.Options.Timestamps {
			if zarray.Contains(oldColumns, CreatedAtKey) {
				newColumns = append(newColumns, CreatedAtKey)
			}
			if zarray.Contains(oldColumns, UpdatedAtKey) {
				newColumns = append(newColumns, UpdatedAtKey)
			}
		}
	}

	// updateColumns := zarray.Map(zarray.Filter(m.Model.Columns, func(_ int, n *Column) bool {
	// 	c := currentColumns.Get(n.Name)
	// 	if !c.Exists() {
	// 		return false
	// 	}
	// 	nf := schema.NewField(n.Name, schema.DataType(n.Type), func(f *schema.Field) {
	// 		f.Size = n.Size
	// 	})
	// 	t := d.DataTypeOf(nf, true)
	// 	return !strings.EqualFold(t, c.Get("type").String())
	// }), func(i int, v *Column) string { return v.Name })

	addColumns := zarray.Filter(newColumns, func(_ int, n string) bool {
		return !zarray.Contains(oldColumns, n)
	})

	deleteColumns := zarray.Filter(oldColumns, func(_ int, n string) bool {
		return !zarray.Contains(newColumns, n) && !strings.HasPrefix(n, deleteFieldPrefix)
	})

	for _, v := range deleteColumns {
		if deleteColumn {
			sql, values = table.DropColumn(v)
		} else {
			sql, values = table.RenameColumn(v, deleteFieldPrefix+v)
		}

		_, err := m.DB.Exec(sql, values...)
		if err != nil {
			return err
		}
	}

	if m.Model.Options.SoftDeletes {
		if !zarray.Contains(oldColumns, DeletedAtKey) {
			sql, values := table.AddColumn(DeletedAtKey, "int", func(f *schema.Field) {
				f.Comment = "删除时间戳"
			})
			_, err := m.DB.Exec(sql, values...)
			if err != nil {
				return err
			}
		}
	}

	if m.Model.Options.CreatedBy {
		if !zarray.Contains(oldColumns, CreatedByKey) {
			sql, values := table.AddColumn(CreatedByKey, "string", func(f *schema.Field) {
				f.Comment = "创建人 ID"
				f.NotNull = false
				f.Size = 120
			})
			_, err := m.DB.Exec(sql, values...)
			if err != nil {
				return err
			}
		}
	}

	if m.Model.Options.Timestamps {
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

	for _, v := range addColumns {
		c, ok := zarray.Find(m.Model.Columns, func(i int, c *Column) bool {
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

		if !deleteColumn {
			recovery := deleteFieldPrefix + v
			_, ok := zarray.Find(oldColumns, func(i int, n string) bool {
				return n == recovery
			})
			if ok {
				sql, values = table.RenameColumn(recovery, v)
			}
		}

		_, err := m.DB.Exec(sql, values...)
		if err != nil {
			return err
		}
	}

	// TODO 是否需要支持修改字段类型
	// if len(updateColumns) > 0 {
	// 	zlog.Warn("暂不支持修改字段类型：", updateColumns)
	// }

	return nil
}

func (m *Migration) fillField(fields []*schema.Field) []*schema.Field {
	if m.Model.Options.SoftDeletes {
		fields = append(fields, schema.NewField(DeletedAtKey, schema.Int, func(f *schema.Field) {
			f.Size = 9999999999
			f.NotNull = false
			f.Comment = "删除时间"
		}))
	}

	if m.Model.Options.Timestamps {
		fields = append(fields, schema.NewField(CreatedAtKey, schema.Time, func(f *schema.Field) {
			f.Comment = "创建时间"
		}))
		fields = append(fields, schema.NewField(UpdatedAtKey, schema.Time, func(f *schema.Field) {
			f.Comment = "更新时间"
		}))
	}

	if m.Model.Options.CreatedBy {
		fields = append(fields, schema.NewField(CreatedByKey, schema.String, func(f *schema.Field) {
			f.Comment = "创建人 ID"
			f.NotNull = false
			f.Size = 120
		}))
	}
	return fields
}

func (m *Migration) CreateTable() error {
	table := builder.NewTable(m.Model.Table.Name).Create()

	fields := make([]*schema.Field, 0, len(m.Model.Columns))

	fields = append(fields, m.getPrimaryKey())

	for _, v := range m.Model.Columns {
		f := schema.NewField(v.Name, schema.DataType(v.Type), func(f *schema.Field) {
			f.Comment = zutil.IfVal(v.Comment != "", v.Comment, v.Label).(string)
			f.NotNull = !v.Nullable
			f.Size = v.Size
		})
		fields = append(fields, f)
	}

	fields = m.fillField(fields)

	table.Column(fields...)

	if len(fields) == 0 {
		return errors.New("表字段不能为空")
	}

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

func (m *Migration) Indexs() error {
	table := builder.NewTable(m.Model.Table.Name).Create()

	uniques := make(map[string][]string, 0)
	indexs := make(map[string][]string, 0)
	for _, column := range m.Model.Columns {
		unique := ztype.ToString(column.Unique)
		if unique != "" {
			if unique == "true" {
				unique = column.Name
			}
			uniques[unique] = append(uniques[unique], column.Name)
		}

		index := ztype.ToString(column.Index)
		if index != "" {
			if index == "true" {
				index = column.Name
			}
			indexs[index] = append(indexs[index], column.Name)
		}
	}

	for name, v := range uniques {
		name = m.Model.Table.Name + "__unique__" + name
		sql, values, process := table.HasIndex(name)
		res, err := m.DB.QueryToMaps(sql, values...)

		if err == nil && !process(res) {
			sql, values := table.CreateIndex(name, v, "UNIQUE")
			_, err = m.DB.Exec(sql, values...)
			if err != nil {
				return err
			}
		}
	}

	for name, v := range indexs {
		name = m.Model.Table.Name + "__idx__" + name
		sql, values, process := table.HasIndex(name)
		res, err := m.DB.QueryToMaps(sql, values...)
		if err == nil && !process(res) {
			sql, values := table.CreateIndex(name, v, "")
			_, err = m.DB.Exec(sql, values...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
