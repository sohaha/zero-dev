package parse

import (
	"zlsapp/common/hashid"
	"zlsapp/conf"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/zlsgo/zdb/schema"
)

func InitModel(alias string, m *Modeler) {
	m.Alias = alias
	salt := m.Options.Salt
	m.Hashid = hashid.New(salt, 8)
	m.readOnlyKeys = make([]string, 0)
	m.cryptKeys = make(map[string]cryptProcess, 0)
	m.afterProcess = make(map[string][]afterProcess, 0)
	m.beforeProcess = make(map[string][]beforeProcess, 0)
	m.Fields = zarray.Map(m.Columns, func(_ int, c *Column) string {
		resolverColumn(m, c)

		resolverValidRule(c)

		resolverColumnOptions(c)

		return c.Name
	})

	m.inlayFields = []string{IDKey}
	if m.Options.Timestamps {
		m.inlayFields = append(m.inlayFields, CreatedAtKey, UpdatedAtKey)
	}

	if m.Options.CreatedBy {
		m.inlayFields = append(m.inlayFields, CreatedByKey)
	}

	// if m.Options.SoftDeletes {
	// 	m.inlayFields = append(m.inlayFields, DeletedAtKey)
	// }

	m.fullFields = append([]string{IDKey}, m.Fields...)
	m.fullFields = zarray.Unique(append(m.fullFields, m.inlayFields...))
	if m.Raw == nil {
		m.Raw, _ = zjson.Marshal(m)
	}

	relations := m.Relations
	if len(relations) > 0 {
		for k := range relations {
			v := relations[k]
			if v.Foreign == "" {
				m.Relations[k].Foreign = IDKey
			}
		}

		newRelations := make(map[string]*ModelRelation, len(relations))
		for k := range relations {
			v := relations[k]
			newRelations[zstring.CamelCaseToSnakeCase(k)] = v
		}
		m.Relations = newRelations
	}

	if m.Options.CreatedBy {
		m.Relations[zstring.SnakeCaseToCamelCase(CreatedByKey, true)] = &ModelRelation{
			Key:     CreatedByKey,
			Model:   conf.UsersModel,
			Foreign: "_id",
			Fields: []string{
				"account",
				"nickname",
			},
		}
	}

	resolverView(m)
	// resolverApi(m)
}

func (m *Modeler) isInlayField(field string) bool {
	if field == IDKey {
		return true
	}

	if !m.Options.Timestamps && !m.Options.CreatedBy {
		return false
	}

	return field == CreatedAtKey || field == UpdatedAtKey || field == CreatedByKey
}

func (m *Modeler) GetFields(exclude ...string) []string {
	f := m.fullFields
	if len(exclude) == 0 {
		return f
	}

	return zarray.Filter(f, func(_ int, v string) bool {
		return !zarray.Contains(exclude, v)
	})
}

func (m *Modeler) GetColumn(name string) (*Column, bool) {
	column, ok := zarray.Find(m.Columns, func(_ int, c *Column) bool {
		return c.Name == name
	})
	if ok {
		return column, true
	}

	if name == IDKey {
		return &Column{
			Name:     IDKey,
			Type:     schema.Int,
			Nullable: false,
			Label:    "ID",
			ReadOnly: true,
		}, true

	}
	if m.Options.Timestamps {
		switch name {
		case CreatedAtKey:
			return &Column{
				Name:  name,
				Type:  schema.Time,
				Label: "创建时间"}, true
		case UpdatedAtKey:
			return &Column{
				Name:  name,
				Type:  schema.Time,
				Label: "更新时间"}, true
		}
	}

	if m.Options.SoftDeletes {
		if name == DeletedAtKey {
			return &Column{
				Name:  name,
				Type:  schema.Int,
				Size:  11,
				Label: "删除时间戳"}, true
		}
	}

	if m.Options.CreatedBy {
		if name == CreatedByKey {
			return &Column{
				Name:     name,
				Type:     schema.String,
				Nullable: true,
				Default:  "",
				Size:     120,
				ReadOnly: true,
				Label:    "创建人 ID"}, true
		}
	}

	return nil, false
}
