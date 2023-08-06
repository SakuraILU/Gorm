package session

import (
	"database/sql"
	clause "gorm/Clause"
	dialect "gorm/Dialect"
	log "gorm/Log"
	schema "gorm/Schema"
	"reflect"
	"strings"
)

type Session struct {
	db       *sql.DB
	tx       *sql.Tx
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

func (s *Session) DB() DB {
	if s.tx != nil {
		return s.tx
	}
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

func (s *Session) Raw(sqlcmd string, sqlvals ...any) *Session {
	s.sqlcmds.WriteString(sqlcmd)
	s.sqlvals = append(s.sqlvals, sqlvals...)
	return s
}

func (s *Session) callMethod(name string, value any) {
	if value == nil {
		value = s.reftable.GetMdoel()
	}

	fn, err := s.reftable.GetMethod(name)
	if err != nil {
		return
	}
	rets := fn.Call([]reflect.Value{reflect.ValueOf(value), reflect.ValueOf(s)})
	if len(rets) > 0 {
		if err, ok := rets[0].Interface().(error); ok {
			log.Error(err)
		}
	}
	return
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

func (s *Session) Begin() (err error) {
	tx, err := s.db.Begin()
	if err != nil {
		log.Error(err)
		return
	}

	s.tx = tx
	return
}

func (s *Session) Commit() (err error) {
	err = s.tx.Commit()
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Commit transaction")
	s.tx = nil
	return
}

func (s *Session) Rollback() (err error) {
	err = s.tx.Rollback()
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Rollback transaction")
	s.tx = nil
	return
}
