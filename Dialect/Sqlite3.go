package dialect

import (
	log "gorm/Log"
	"reflect"
	"time"
)

type Sqlite3 struct {
}

func (s *Sqlite3) DataTypeOf(v any) string {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int8:
		return "tinyint"
	case reflect.Int16:
		return "smallint"
	case reflect.Int, reflect.Int32:
		return "integer"
	case reflect.Int64:
		return "bigint"
	case reflect.Uint8:
		return "unsigned tinyint"
	case reflect.Uint16:
		return "unsigned smallint"
	case reflect.Uint, reflect.Uint32:
		return "unsigned integer"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		// datetime--time.Time
		if _, ok := v.(time.Time); ok {
			return "datetime"
		}
	}

	log.Infof("Unsupport data type: %v", v)
	return ""
}

func (s *Sqlite3) TableExistSQL(tableName string) (string, interface{}) {
	return "SELECT name FROM sqlite_master WHERE type='table' and name = ?", tableName
}
