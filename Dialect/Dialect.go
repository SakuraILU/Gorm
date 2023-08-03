package dialect

type Dialect interface {
	DataTypeof(v any) string
	TableExistSQL(tableName string) (string, []interface{})
}
