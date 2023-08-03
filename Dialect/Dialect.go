package dialect

type Dialect interface {
	DataTypeOf(v any) string
	TableExistSQL(tableName string) (string, interface{})
}
