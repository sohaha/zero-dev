package parse

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/zvalid"
	"github.com/zlsgo/zdb"
	"github.com/zlsgo/zdb/schema"
)

type (
	Model struct {
		Schema        string `json:"$schema"`
		Raw           []byte
		Storage       Storageer
		Name          string               `json:"name"`
		Path          string               `json:"-"`
		Table         Table                `json:"table"`
		Columns       []*Column            `json:"columns"`
		Views         map[string]*View     `json:"views"`
		Relations     map[string]*relation `json:"relations"`
		Values        []interface{}        `json:"values"`
		fields        []string
		inlayFields   []string
		fullFields    []string
		readOnlyKeys  []string
		cryptKeys     map[string]cryptProcess
		beforeProcess map[string][]beforeProcess
		afterProcess  map[string][]afterProcess
		Options       struct {
			Api              interface{} `json:"api"`
			ApiPath          string      `json:"api_path"`
			CryptID          bool        `json:"crypt_id"`
			DisabledMigrator bool        `json:"disabled_migrator"`
			SoftDeletes      bool        `json:"soft_deletes"`
			Timestamps       bool        `json:"timestamps"`
		} `json:"options"`
	}

	Table struct {
		Name    string `json:"name"`
		Comment string `json:"comment"`
	}

	relation struct {
		Name    string   `json:"name"`
		Type    string   `json:"type"`
		Model   string   `json:"model"`
		Foreign string   `json:"foreign"`
		Key     string   `json:"key"`
		Fields  []string `json:"fields"`
	}

	Column struct {
		Default     interface{}     `json:"default"`
		Unique      interface{}     `json:"unique"`
		Index       interface{}     `json:"index"`
		Crypt       string          `json:"crypt"`
		Name        string          `json:"name"`
		Comment     string          `json:"comment"`
		Label       string          `json:"label"`
		Type        schema.DataType `json:"type"`
		Validations []validations   `json:"validations"`
		Options     []ColumnEnum    `json:"options"`
		Before      []string        `json:"before"`
		After       []string        `json:"after"`
		validRules  zvalid.Engine   `json:"-"`
		Size        uint64          `json:"size"`
		ReadOnly    bool            `json:"readonly"`
		Nullable    bool            `json:"nullable"`
	}

	validations struct {
		Args    interface{} `json:"args"`
		Method  string      `json:"method"`
		Message string      `json:"message"`
	}

	ColumnEnum struct {
		Value string `json:"value"`
		Label string `json:"label"`
	}
)

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

type DataTime struct {
	time.Time
}

func (t *DataTime) UnmarshalJSON(data []byte) error {
	if len(data) == 2 {
		*t = DataTime{Time: time.Time{}}
		return nil
	}
	now, err := ztime.Parse(zstring.Bytes2String(data))
	if err != nil {
		return err
	}
	*t = DataTime{Time: now}
	return nil
}

func (t DataTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	return zstring.String2Bytes(ztime.FormatTime(t.Time, "\"Y-m-d H:i:s\"")), nil
}

func (t DataTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.IsZero() || t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t DataTime) String() string {
	if t.Time.IsZero() {
		return "0000-00-00 00:00:00"
	}
	return ztime.FormatTime(t.Time)
}

// Scan value of time.Time
func (t *DataTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = DataTime{Time: value}
		return nil
	}
	if b, ok := v.([]byte); ok {
		parse, err := ztime.Parse(zstring.Bytes2String(b))
		if err != nil {
			return err
		}
		*t = DataTime{Time: parse}
		return nil
	}

	return fmt.Errorf("can not convert %v to timestamp", v)
}