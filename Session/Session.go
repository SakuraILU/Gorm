package session

import (
	"database/sql"
	"fmt"
	clause "gorm/Clause"
	dialect "gorm/Dialect"
	log "gorm/Log"
	schema "gorm/Schema"
	"reflect"
	"strings"
)

type Session struct {
	db       *sql.DB
	dial     dialect.Dialect
	reftable *schema.Schema
	sqlcmds  *strings.Builder
	sqlvals  []any

	clause *clause.Clause
}

func NewSession(db *sql.DB, dial dialect.Dialect) (s *Session) {
	return &Session{
		db:      db,
		dial:    dial,
		sqlcmds: &strings.Builder{},
		sqlvals: make([]any, 0),

		clause: clause.NewClause(),
	}
}

func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Model(v any) *Session {
	// model table only if no table is modeled or a new table is to be modeled
	// otherwise, use cached one
	if s.reftable == nil || (s.reftable.GetName() != reflect.Indirect(reflect.ValueOf(v)).Type().Name()) {
		s.reftable = schema.NewSchema(v, s.dial)
	}

	return s
}

func (s *Session) RefTable() *schema.Schema {
	return s.reftable
}

func (s *Session) CreateTable() (err error) {
	fields := s.reftable.GetFields()
	columns := make([]string, 0)
	for _, field := range fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}

	desc := strings.Join(columns, ",")
	log.Info(desc)
	_, err = s.Raw(fmt.Sprintf("CREATE TABLE %s (%s)", s.reftable.GetName(), desc)).Exec()
	return
}

func (s *Session) DropTable() (err error) {
	_, err = s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.reftable.GetName())).Exec()
	return
}

func (s *Session) HasTable() bool {
	cmd, val := s.dial.TableExistSQL(s.reftable.GetName())
	row := s.Raw(cmd, val).QueryRow()
	var tablename string
	if err := row.Scan(&tablename); err != nil {
		log.Error(err)
		return false
	}
	if tablename != s.reftable.GetName() {
		log.Error("want %s, but get %s", s.reftable.GetName(), tablename)
	}
	return true
}

// v1, v2, v3...
func (s *Session) Insert(values ...any) (n int64, err error) {
	// s.Model(values[0])

	s.clause.Set(clause.INSERT, s.reftable.GetName(), s.reftable.GetFieldNames())

	filednames := s.reftable.GetFieldNames()
	vs := make([]any, 0)
	for _, value := range values {
		refv := reflect.Indirect(reflect.ValueOf(value))
		v := make([]any, 0)
		for _, name := range filednames {
			v = append(v, refv.FieldByName(name).Interface())
		}
		vs = append(vs, v)
	}
	s.clause.Set(clause.VALUES, vs...)

	cmd, vals := s.clause.Build(clause.INSERT, clause.VALUES)

	res, err := s.Raw(cmd, vals...).Exec()
	if err != nil {
		log.Error(err)
		return
	}
	return res.RowsAffected()
}

// &vs[]
func (s *Session) Find(values any) (err error) {
	// (*             values)    [0]
	// Indirect       ValueOf    Elem()
	refvs := reflect.Indirect(reflect.ValueOf(values))
	reftyp := refvs.Type().Elem()
	// s.Model(reflect.New(reftyp).Interface())

	s.clause.Set(clause.SELECT, s.reftable.GetName(), s.reftable.GetFieldNames())
	cmd, vals := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(cmd, vals...).Query()
	if err != nil {
		log.Error(err)
		return
	}

	for rows.Next() {
		v := reflect.New(reftyp).Elem()
		elems := make([]any, 0)
		for _, name := range s.reftable.GetFieldNames() {
			elems = append(elems, v.FieldByName(name).Addr().Interface())
		}
		err = rows.Scan(elems...)
		if err != nil {
			log.Error(err)
			return
		}
		refvs.Set(reflect.Append(refvs, v))
	}
	return
}

func (s *Session) First(value any) (err error) {
	refv := reflect.Indirect(reflect.ValueOf(value))
	reftyp := refv.Type()
	// something weired here
	// use v := reflect.MakeSlice(reflect.SliceOf(reftyp), 0, 0) will cause panic
	// it is not addressable...
	// but use New to create a pointer and then Set MakeSlice to it will be ok
	values := reflect.New(reflect.SliceOf(reftyp)).Elem()
	values.Set(reflect.MakeSlice(reflect.SliceOf(reftyp), 0, 0))
	err = s.Limit(1).Find(values.Addr().Interface())
	if err != nil {
		log.Error(err)
		return
	}
	if values.Len() == 0 {
		err = fmt.Errorf("Not Found")
		log.Error(err)
		return
	}

	refv.Set(values.Index(0))
	return
}

func (s *Session) Update(values ...any) (n int64, err error) {
	kvs, ok := values[0].(map[string]any)
	if !ok {
		if len(values)%2 != 0 {
			err = fmt.Errorf("Update: arguments should be a map[string]any or k, v, k, v,...(length mod 2 == 0)")
			log.Error(err)
			return
		}
		kvs = make(map[string]any)
		for i := 0; i < len(values); i = i + 2 {
			kvs[values[i].(string)] = values[i+1]
		}
	}

	s.clause.Set(clause.UPDATE, s.reftable.GetName(), kvs)
	cmd, vals := s.clause.Build(clause.UPDATE, clause.WHERE)
	res, err := s.Raw(cmd, vals...).Exec()
	if err != nil {
		log.Error(err)
		return
	}
	return res.RowsAffected()
}

func (s *Session) Delete() (n int64, err error) {
	s.clause.Set(clause.DELETE, s.reftable.GetName())
	cmd, vals := s.clause.Build(clause.DELETE, clause.WHERE)
	res, err := s.Raw(cmd, vals...).Exec()
	return res.RowsAffected()
}

func (s *Session) Count() (n int64, err error) {
	s.clause.Set(clause.COUNT, s.reftable.GetName())
	cmd, vals := s.clause.Build(clause.COUNT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	row := s.Raw(cmd, vals).QueryRow()
	if err = row.Scan(&n); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) Where(desc string, values ...any) *Session {
	values = append([]any{desc}, values...)
	s.clause.Set(clause.WHERE, values...)
	return s
}

func (s *Session) OrderBy(key string) *Session {
	s.clause.Set(clause.ORDERBY, key)
	return s
}

func (s *Session) Limit(value any) *Session {
	s.clause.Set(clause.LIMIT, value)
	return s
}

func (s *Session) Raw(sqlcmd string, sqlvals ...any) *Session {
	s.sqlcmds.WriteString(sqlcmd)
	s.sqlvals = append(s.sqlvals, sqlvals...)
	return s
}

func (s *Session) Clear() {
	s.sqlcmds.Reset()
	s.sqlvals = make([]any, 0)
}

func (s *Session) Exec() (r sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sqlcmds.String(), s.sqlvals)
	r, err = s.DB().Exec(s.sqlcmds.String(), s.sqlvals...)
	if err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) Query() (rs *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sqlcmds.String(), s.sqlvals)
	rs, err = s.DB().Query(s.sqlcmds.String(), s.sqlvals...)
	if err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) QueryRow() (r *sql.Row) {
	defer s.Clear()
	log.Info(s.sqlcmds.String(), s.sqlvals)
	// the error of QueryRow will be deferred until the Scan() method
	return s.DB().QueryRow(s.sqlcmds.String(), s.sqlvals...)
}
