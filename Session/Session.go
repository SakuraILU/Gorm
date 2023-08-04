package session

import (
	"database/sql"
	"fmt"
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
}

func NewSession(db *sql.DB, dial dialect.Dialect) (s *Session) {
	return &Session{
		db:      db,
		dial:    dial,
		sqlcmds: &strings.Builder{},
		sqlvals: make([]any, 0),
	}
}

func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Model(v any) *Session {
	// model table only if no table is modeled or a new table is to be modeled
	// otherwise, use cached one
	if s.reftable == nil || (s.reftable.GetName() != reflect.TypeOf(v).Name()) {
		s.reftable = schema.NewSchema(v, s.dial)
	}

	return s
}

func (s *Session) RefTable() *schema.Schema {
	return s.reftable
}

func (s *Session) CreateTable() (err error) {
	fields := s.reftable.GetFields()
	columns := make([]string, len(fields))
	for _, field := range fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}

	desc := strings.Join(columns, ",")
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
