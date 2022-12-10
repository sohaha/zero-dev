package parse

import (
	"errors"
	"strings"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/zdb"
	"github.com/zlsgo/zdb/builder"
)

func (s *SQL) m(f string) string {
	if strings.Contains(f, ".") {
		return f
	}
	if strings.Contains(f, " ") {
		return f
	}
	return s.table + "." + f
}
func (s *SQL) parseExprs(d *builder.Cond, filter ztype.Map) (exprs []string, err error) {
	if len(filter) > 0 {
		for k := range filter {
			v := ztype.New(filter[k])
			if k == "" {
				exprs = append(exprs, d.And(v.String()))
				continue
			}
			f := strings.SplitN(zstring.TrimSpace(k), " ", 2)
			l := len(f)
			if l != 2 {
				switch val := v.Value().(type) {
				case []interface{}:
					exprs = append(exprs, d.In(s.m(f[0]), val...))
				default:
					exprs = append(exprs, d.EQ(s.m(f[0]), val))
				}
			} else {
				switch f[1] {
				default:
					err = errors.New("Unknown operator:" + f[1])
					return
				case "=":
					exprs = append(exprs, d.EQ(f[0], v.Value()))
				case ">":
					exprs = append(exprs, d.GT(f[0], v.Value()))
				case ">=":
					exprs = append(exprs, d.GE(f[0], v.Value()))
				case "<":
					exprs = append(exprs, d.LT(f[0], v.Value()))
				case "<=":
					exprs = append(exprs, d.LE(f[0], v.Value()))
				case "!=":
					exprs = append(exprs, d.NE(f[0], v.Value()))
					// case "like":
					// 	exprs = append(exprs, d.Like(f[0], v.String()))
					// case "in":
					// 	exprs = append(exprs, d.In(f[0], v.Value()))
					// case "not in":
					// 	exprs = append(exprs, d.NotIn(f[0], v.Value()))
					// case "between":
					// 	exprs = append(exprs, d.Between(f[0], v.Value()))
				}
			}
		}
	}

	return
}

func (s *SQL) Insert(data ztype.Map) (lastId interface{}, err error) {
	return s.db.Insert(s.table, data)
}

func (s *SQL) Delete(filter ztype.Map, fn ...StorageOptionFn) (int64, error) {
	o := StorageOptions{}
	for _, f := range fn {
		if err := f(&o); err != nil {
			return 0, err
		}
	}
	return s.db.Delete(s.table, func(b *builder.DeleteBuilder) error {
		exprs, err := s.parseExprs(&b.Cond, filter)
		if err != nil {
			return err
		}

		if len(exprs) > 0 {
			b.Where(exprs...)
		}

		if len(o.OrderBy) > 0 {
			for k, v := range o.OrderBy {
				if v == -1 {
					b.OrderBy(k + " DESC")
				}
			}
		}

		return nil
	})
}

func (s *SQL) FindOne(filter ztype.Map, fn ...StorageOptionFn) (ztype.Map, error) {
	rows, err := s.Find(filter, func(so *StorageOptions) error {
		so.Limit = 1
		if len(fn) > 0 {
			return fn[0](so)
		}
		return nil
	})

	if err == nil && rows.Len() > 0 {
		return rows[0], nil
	}

	return ztype.Map{}, err
}

func (s *SQL) Find(filter ztype.Map, fn ...StorageOptionFn) (ztype.Maps, error) {
	o := StorageOptions{}
	for _, f := range fn {
		if err := f(&o); err != nil {
			return nil, err
		}
	}

	items, err := s.db.Find(s.table, func(b *builder.SelectBuilder) error {
		fields := o.Fields
		if len(fields) > 0 {
			b.Select(zarray.Map(o.Fields, func(_ int, v string) string {
				return s.m(v)
			})...)
		}

		exprs, err := s.parseExprs(&b.Cond, filter)
		if err != nil {
			return err
		}

		if len(exprs) > 0 {
			b.Where(exprs...)
		}

		if len(o.Join) > 0 {
			for _, v := range o.Join {
				b.JoinWithOption("", b.As(v.Table, v.As), v.Expr)
			}
		}

		if len(o.OrderBy) > 0 {
			for k, v := range o.OrderBy {
				if v == -1 {
					b.OrderBy(k + " DESC")
				}
			}
		}

		if o.Limit > 0 {
			b.Limit(o.Limit)
		}

		return nil
	})

	if err != nil && err != zdb.ErrNotFound {
		return items, err
	}

	return items, nil
}

func (s *SQL) Pages(page, pagesize int, filter ztype.Map, fn ...StorageOptionFn) (ztype.Maps, PageInfo, error) {
	o := StorageOptions{}
	for _, f := range fn {
		if err := f(&o); err != nil {
			return nil, PageInfo{}, err
		}
	}

	rows, p, err := s.db.Pages(s.table, page, pagesize, func(b *builder.SelectBuilder) error {
		if len(o.Fields) > 0 {
			b.Select(zarray.Map(o.Fields, func(_ int, v string) string {
				return s.m(v)
			})...)
		}

		exprs, err := s.parseExprs(&b.Cond, filter)
		if err != nil {
			return err
		}

		if len(exprs) > 0 {
			b.Where(exprs...)
		}

		if len(o.OrderBy) > 0 {
			for k, v := range o.OrderBy {
				if v == -1 {
					b.OrderBy(k + " DESC")
				}
			}
		}

		if len(o.Join) > 0 {
			for _, v := range o.Join {
				b.JoinWithOption("", b.As(v.Table, v.As), v.Expr)
			}
		}

		if o.Limit > 0 {
			b.Limit(o.Limit)
		}

		return nil
	})
	if err != nil && err != zdb.ErrNotFound {
		return nil, PageInfo{}, err
	}

	return rows, PageInfo{
		p,
	}, nil
}

func (s *SQL) Update(data ztype.Map, filter ztype.Map, fn ...StorageOptionFn) (int64, error) {
	o := StorageOptions{}
	for _, f := range fn {
		if err := f(&o); err != nil {
			return 0, err
		}
	}
	return s.db.Update(s.table, data, func(b *builder.UpdateBuilder) error {
		exprs, err := s.parseExprs(&b.Cond, filter)
		if err != nil {
			return err
		}

		if len(exprs) > 0 {
			b.Where(exprs...)
		}

		if o.Limit > 0 {
			b.Limit(o.Limit)
		}

		if len(o.OrderBy) > 0 {
			for k, v := range o.OrderBy {
				if v == -1 {
					b.OrderBy(k + " DESC")
				}
			}
		}

		return nil
	})
}
