package dialect

import (
	"fmt"
	log "gorm/Log"
)

type Dialect interface {
	DataTypeOf(v any) string
	TableExistSQL(tableName string) (string, interface{})
}

var dialects map[string]Dialect

func init() {
	dialects = make(map[string]Dialect)
	dialects["sqlite3"] = &Sqlite3{}
}

func GetDialect(name string) (dial Dialect, err error) {
	dial, ok := dialects[name]
	if !ok {
		err = fmt.Errorf("dialect %s is not found", name)
		log.Error(err)
		return
	}
	return
}
