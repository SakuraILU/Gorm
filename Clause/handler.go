package clause

import (
	"fmt"
	log "gorm/Log"
	"strings"
)

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
)

type GenFn func(...any) (string, []any)

var generators map[Type]GenFn = map[Type]GenFn{
	INSERT:  _insert,
	VALUES:  _values,
	SELECT:  _select,
	LIMIT:   _limit,
	WHERE:   _where,
	ORDERBY: _orderBy,
	UPDATE:  _update,
}

// string, []string
func _insert(args ...any) (string, []any) {
	// INSERT INTO User (name, age, career)
	//			  tablename   fieldnames
	//				string		[]string
	tablename, ok := args[0].(string)
	if !ok {
		log.Fatal("table name is not string")
	}
	fieldnames, ok := args[1].([]string)
	if !ok {
		log.Fatal("field names %v are not string", fieldnames)
	}

	return fmt.Sprintf("INSERT INTO %s (%s)", tablename, strings.Join(fieldnames, ", ")), []any{}
}

// []any, []any, ...
func _values(args ...any) (string, []any) {
	// args (?,?,?) (?,?,>), ["Jack", Worker, 35], ["Bob", "Student", 21]
	//									v1					v2
	//									[]any				[]any
	cmds := make([]string, 0)
	vs := make([]any, 0)
	for _, val := range args {
		v := val.([]any)
		vtemplate := make([]string, len(v))
		for i := 0; i < len(v); i++ {
			vtemplate[i] = "?"
		}
		cmds = append(cmds, fmt.Sprintf("(%s)", strings.Join(vtemplate, ",")))
		vs = append(vs, v...)
	}

	return "VALUES " + strings.Join(cmds, ", "), vs
}

// string, []string
func _select(args ...any) (string, []any) {
	// SELECT age, career FROM User
	//		    fieldnames   tablename
	//			 []string     string
	tablename := args[0].(string)
	fieldnames := args[1].([]string)

	return fmt.Sprintf("SELECT %s FROM %s", strings.Join(fieldnames, ","), tablename), []any{}
}

// desc, any, any, ...
func _where(args ...any) (string, []any) {
	// WHERE type = ?, age < ?, Worker, 35
	// 			  desc			  v1    v2
	//			 string		 	 any   any
	desc := args[0].(string)
	vs := args[1:]

	return fmt.Sprintf("WHERE %s", desc), vs
}

// string
func _orderBy(args ...any) (string, []any) {
	// ORDER BY user1.age
	//			   key
	//			  string
	key := args[0].(string)
	return fmt.Sprintf("ORDER BY %s", key), []any{}
}

// any
func _limit(args ...any) (string, []any) {
	// LIMIT 30
	//       v
	//		any
	v := args[0]
	return "LIMIT ?", []any{v}
}

// string map[string]any
func _update(args ...any) (string, []any) {
	// UPDATE tablename SET k1 = ? k2 = ? ..., v1, v2 ...
	//			string		str   str          any any
	tablename := args[0].(string)

	kvs := args[1].(map[string]any)
	ks := make([]string, 0)
	vs := make([]any, 0)
	for k, v := range kvs {
		ks = append(ks, k+" = ?")
		vs = append(vs, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tablename, strings.Join(ks, ", ")), vs
}

func _delete(args ...any) (string, []any) {
	// DELETE FROM tablename
	//			   string
	tablename := args[0].(string)
	return fmt.Sprintf("DELETE FROM %s", tablename), []any{}
}
