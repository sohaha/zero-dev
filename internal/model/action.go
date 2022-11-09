package model

import (
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/zdb/builder"
)

// ActionUpdate 更新数据
func (m *Model) ActionUpdate(key interface{}, data ztype.Map) error {
	_, err := m.Update(data, func(b *builder.UpdateBuilder) error {
		b.Where(b.EQ(IDKey, key))
		return nil
	})
	return err
}

// ActionCreate 创建数据
func (m *Model) ActionCreate(data ztype.Map) (lastId int64, err error) {
	return m.Insert(data)
}
